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
