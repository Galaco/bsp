package bsp_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/galaco/bsp"
)

func containerBytes(id int, payload []byte) []byte {
	raw := make([]byte, bsp.HeaderSize+len(payload))
	binary.LittleEndian.PutUint32(raw, 0x50534256)
	binary.LittleEndian.PutUint32(raw[4:], 20)
	if id >= 0 {
		binary.LittleEndian.PutUint32(raw[8+16*id:], bsp.HeaderSize)
		binary.LittleEndian.PutUint32(raw[12+16*id:], uint32(len(payload)))
		copy(raw[bsp.HeaderSize:], payload)
	}
	return raw
}

func TestReadValidation(t *testing.T) {
	cases := []struct {
		name   string
		change func([]byte) []byte
		want   error
	}{
		{"magic", func(b []byte) []byte { b[0] = 0; return b }, bsp.ErrMalformed},
		{"byte order", func(b []byte) []byte { copy(b, "PSBV"); return b }, bsp.ErrUnsupportedFormat},
		{"version", func(b []byte) []byte { binary.LittleEndian.PutUint32(b[4:], 987); return b }, bsp.ErrUnsupportedFormat},
		{"short header", func(b []byte) []byte { return b[:1035] }, io.ErrUnexpectedEOF},
		{"negative offset", func(b []byte) []byte { binary.LittleEndian.PutUint32(b[8:], 0xffffffff); return b }, bsp.ErrMalformed},
		{"negative length", func(b []byte) []byte { binary.LittleEndian.PutUint32(b[12:], 0xffffffff); return b }, bsp.ErrMalformed},
		{"header range", func(b []byte) []byte { binary.LittleEndian.PutUint32(b[8:], 1035); return b }, bsp.ErrMalformed},
		{"huge length", func(b []byte) []byte { binary.LittleEndian.PutUint32(b[12:], 0x7fffffff); return b }, io.ErrUnexpectedEOF},
		{"truncated payload", func(b []byte) []byte { return b[:len(b)-1] }, io.ErrUnexpectedEOF},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := bsp.Read(bytes.NewReader(tc.change(containerBytes(0, []byte("hello")))), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
			if !errors.Is(err, tc.want) {
				t.Fatalf("error=%v, want %v", err, tc.want)
			}
		})
	}
	file, err := bsp.Read(bytes.NewReader(containerBytes(0, []byte("hello"))), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := file.RawLump(0)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(raw, []byte("hello")) {
		t.Fatalf("payload=%q", raw)
	}
}

func TestLazyDecodeErrors(t *testing.T) {
	file, err := bsp.Read(bytes.NewReader(containerBytes(4, []byte{0})), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		data, err := file.Visibility()
		if len(data.Offsets) != 0 || len(data.Data) != 0 || !errors.Is(err, bsp.ErrMalformed) || !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("data=%v error=%v", data, err)
		}
	}
	if raw, err := file.RawLump(4); err != nil || !bytes.Equal(raw, []byte{0}) {
		t.Fatalf("raw passthrough=%x error=%v", raw, err)
	}
}

func TestInvalidLumpID(t *testing.T) {
	file, err := bsp.Read(bytes.NewReader(containerBytes(-1, nil)), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []bsp.LumpID{-1, 64} {
		if _, err := file.Descriptor(id); !errors.Is(err, bsp.ErrInvalidLumpID) {
			t.Fatal(err)
		}
		if _, err := file.RawLump(id); !errors.Is(err, bsp.ErrInvalidLumpID) {
			t.Fatal(err)
		}
		if err := file.SetRawLump(id, bsp.HeaderLump{}, nil); !errors.Is(err, bsp.ErrInvalidLumpID) {
			t.Fatal(err)
		}
	}
}

func TestStreamOwnership(t *testing.T) {
	path := filepath.Join(t.TempDir(), "owned.bsp")
	if err := os.WriteFile(path, containerBytes(-1, nil), 0600); err != nil {
		t.Fatal(err)
	}
	stream, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()
	if _, err := bsp.Read(stream, bsp.ReadOptions{Profile: bsp.ProfileSDK2013}); err != nil {
		t.Fatal(err)
	}
	if _, err := stream.Seek(0, io.SeekStart); err != nil {
		t.Fatalf("caller stream was closed: %v", err)
	}
	var magic [4]byte
	if _, err := io.ReadFull(stream, magic[:]); err != nil || string(magic[:]) != "VBSP" {
		t.Fatalf("stream magic=%q error=%v", magic, err)
	}
}

func FuzzReadHeaderRanges(f *testing.F) {
	profiles := []bsp.Profile{bsp.ProfileSDK2013, bsp.ProfileCSGO, bsp.ProfileL4D2, bsp.ProfileReactiveDrop}
	for selection, profile := range profiles {
		version := int32(21)
		if profile == bsp.ProfileSDK2013 {
			version = 20
		}
		raw := profileContainer(profile, version, []byte{0, 0xa5, 0xff})
		f.Add(raw, uint8(selection))
		f.Add(raw[:bsp.HeaderSize-1], uint8(selection))
		// Independently authored signed range boundaries in both descriptor orders.
		offsetWord, lengthWord := 8+63*16, 12+63*16
		if profile == bsp.ProfileL4D2 {
			offsetWord, lengthWord = offsetWord+4, lengthWord+4
		}
		for _, span := range [][2]uint32{
			{bsp.HeaderSize, 3},     // shared payload with lump 0
			{bsp.HeaderSize + 1, 2}, // unaligned but bounded
			{bsp.HeaderSize - 1, 1},
			{bsp.HeaderSize + 3, 1},
			{0x7fffffff, 0},
			{0xffffffff, 0},
			{bsp.HeaderSize, 0x7fffffff},
			{bsp.HeaderSize, 0xffffffff},
		} {
			changed := append([]byte(nil), raw...)
			binary.LittleEndian.PutUint32(changed[offsetWord:], span[0])
			binary.LittleEndian.PutUint32(changed[lengthWord:], span[1])
			f.Add(changed, uint8(selection))
		}
	}
	f.Add(profileContainer(bsp.ProfileSDK2013, 19, nil), uint8(0))
	f.Add(profileContainer(bsp.ProfileSDK2013, 21, nil), uint8(0))
	f.Add([]byte("PSBV"), uint8(0))
	f.Add([]byte{}, uint8(0))
	f.Fuzz(func(t *testing.T, raw []byte, selection uint8) {
		// A bounded reader also bounds stored payloads, overlapping views and
		// their canonical exports. No typed accessor triggers decompression.
		if len(raw) > 4096 {
			return
		}
		profile := profiles[int(selection)%len(profiles)]
		options := bsp.ReadOptions{Profile: profile, MaxDecodedLumpBytes: 4096, MaxLZMADictionaryBytes: 4096}
		before := append([]byte(nil), raw...)
		file, err := bsp.Read(bytes.NewReader(raw), options)
		if !bytes.Equal(raw, before) {
			t.Fatal("read mutated input")
		}
		if err != nil {
			if file != nil {
				t.Fatal("read returned a partial file")
			}
			return
		}
		if len(raw) < bsp.HeaderSize || string(raw[:4]) != "VBSP" {
			t.Fatal("read accepted an incomplete or invalid header")
		}
		header := file.Header()
		version := int32(binary.LittleEndian.Uint32(raw[4:8]))
		if header.ID != 0x50534256 || header.Version != version || header.Revision != int32(binary.LittleEndian.Uint32(raw[1032:1036])) {
			t.Fatal("header words were not preserved")
		}
		if profile == bsp.ProfileSDK2013 && version != 19 && version != 20 || profile != bsp.ProfileSDK2013 && version != 21 {
			t.Fatal("read accepted a version outside the explicit profile")
		}
		source := snapshotStreamSource(t, file)
		for id, descriptor := range header.Lumps {
			entry := raw[8+id*16 : 8+(id+1)*16]
			words := [3]int32{
				int32(binary.LittleEndian.Uint32(entry[:4])),
				int32(binary.LittleEndian.Uint32(entry[4:8])),
				int32(binary.LittleEndian.Uint32(entry[8:12])),
			}
			if profile == bsp.ProfileL4D2 {
				words = [3]int32{words[1], words[2], words[0]}
			}
			var fourCC [4]byte
			copy(fourCC[:], entry[12:16])
			if descriptor != (bsp.HeaderLump{Offset: words[0], Length: words[1], Version: words[2], FourCC: fourCC}) {
				t.Fatalf("lump %d descriptor words changed", id)
			}
			if descriptor.Offset < 0 || descriptor.Length < 0 {
				t.Fatalf("lump %d accepted a negative range", id)
			}
			var want []byte
			if descriptor.Length != 0 {
				end := int64(descriptor.Offset) + int64(descriptor.Length)
				if descriptor.Offset < bsp.HeaderSize || end > int64(len(raw)) {
					t.Fatalf("lump %d accepted an out-of-bounds range", id)
				}
				want = raw[int(descriptor.Offset):int(end)]
			}
			if !bytes.Equal(source.raw[id], want) {
				t.Fatalf("lump %d payload differs from its wire range", id)
			}
		}
		exported, exportErr := file.MarshalBinary()
		requireStreamSourceUnchanged(t, file, source)
		var streamed bytes.Buffer
		n, streamErr := file.WriteTo(&streamed)
		requireStreamSourceUnchanged(t, file, source)
		if exportErr != nil {
			// Read is intentionally lazy: malformed game directories can pass
			// container validation, but export must fail without partial output.
			if len(source.raw[bsp.LumpGame]) == 0 {
				t.Fatalf("valid container without a game directory failed export: %v", exportErr)
			}
			if exported != nil || streamErr == nil || n != 0 || streamed.Len() != 0 {
				t.Fatal("failed export published partial output")
			}
			return
		}
		if streamErr != nil || n != int64(len(exported)) || !bytes.Equal(streamed.Bytes(), exported) {
			t.Fatalf("streamed export differs: bytes=%d error=%v", n, streamErr)
		}
		reloaded, err := bsp.Read(bytes.NewReader(exported), options)
		if err != nil {
			t.Fatal(err)
		}
		next := snapshotStreamSource(t, reloaded)
		if next.header.ID != header.ID || next.header.Version != header.Version || next.header.Revision != header.Revision {
			t.Fatal("export changed header identity or revision")
		}
		for id, descriptor := range next.header.Lumps {
			original := header.Lumps[id]
			if descriptor.Version != original.Version || descriptor.FourCC != original.FourCC || descriptor.Length != original.Length {
				t.Fatalf("export changed lump %d metadata", id)
			}
			if id != int(bsp.LumpGame) && !bytes.Equal(next.raw[id], source.raw[id]) {
				t.Fatalf("export changed stored lump %d", id)
			}
			if descriptor.Offset%4 != 0 || descriptor.Length == 0 && descriptor.Offset != 0 {
				t.Fatalf("export did not canonicalize lump %d alignment", id)
			}
		}
		repeated, err := reloaded.MarshalBinary()
		if err != nil || !bytes.Equal(repeated, exported) || !bytes.Equal(raw, before) {
			t.Fatalf("canonical export was not stable or changed input: %v", err)
		}
	})
}
