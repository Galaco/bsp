package bsp_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/galaco/bsp"
	"testing"
)

func profileContainer(profile bsp.Profile, version int32, payload []byte) []byte {
	raw := containerBytes(0, payload)
	binary.LittleEndian.PutUint32(raw[4:], uint32(version))
	binary.LittleEndian.PutUint32(raw[1032:], 1234)
	if profile == bsp.ProfileL4D2 {
		for i := 0; i < bsp.LumpCount; i++ {
			entry := raw[8+i*16 : 8+(i+1)*16]
			offset, length, revision := binary.LittleEndian.Uint32(entry), binary.LittleEndian.Uint32(entry[4:]), binary.LittleEndian.Uint32(entry[8:])
			binary.LittleEndian.PutUint32(entry, revision)
			binary.LittleEndian.PutUint32(entry[4:], offset)
			binary.LittleEndian.PutUint32(entry[8:], length)
		}
	}
	return raw
}

func TestProfileHeaderNormalization(t *testing.T) {
	var common bsp.Header
	for _, profile := range []bsp.Profile{bsp.ProfileCSGO, bsp.ProfileL4D2} {
		raw := profileContainer(profile, 21, []byte("entity"))
		copy(raw[8+63*16+12:], []byte{0xab, 0xcd, 0xef, 1})
		file, err := bsp.ReadFromStream(bytes.NewReader(raw), bsp.ReadOptions{Profile: profile})
		if err != nil {
			t.Fatal(err)
		}
		text, err := file.Entities()
		if err != nil || text != "entity" {
			t.Fatalf("%s entity=%q error=%v", profile, text, err)
		}
		header := file.Header()
		if profile == bsp.ProfileCSGO {
			common = header
		} else if header != common {
			t.Fatalf("%s normalized header=%+v want %+v", profile, header, common)
		}
	}
}

func TestDescriptorRevisionSelectsLeafLayout(t *testing.T) {
	for _, tc := range []struct {
		profile                  bsp.Profile
		fileVersion, leafVersion int32
		stride                   int
	}{{bsp.ProfileSDK2013, 19, 1, 32}, {bsp.ProfileCSGO, 21, 0, 56}} {
		payload := make([]byte, tc.stride*2)
		for i := 0; i < 2; i++ {
			record := payload[i*tc.stride : (i+1)*tc.stride]
			binary.LittleEndian.PutUint32(record, uint32(100+i))
			binary.LittleEndian.PutUint16(record[4:], uint16(5+i))
			binary.LittleEndian.PutUint16(record[28:], uint16(19+i))
			if tc.leafVersion == 0 {
				record[30] = byte(31 + i)
			}
		}
		raw := containerBytes(10, payload)
		binary.LittleEndian.PutUint32(raw[4:], uint32(tc.fileVersion))
		binary.LittleEndian.PutUint32(raw[8+10*16+8:], uint32(tc.leafVersion))
		file, err := bsp.ReadFromStream(bytes.NewReader(raw), bsp.ReadOptions{Profile: tc.profile})
		if err != nil {
			t.Fatal(err)
		}
		records, err := file.Leaves()
		if err != nil {
			t.Fatal(err)
		}
		if len(records) != 2 || records[1].Contents != 101 || records[1].Cluster != 6 || records[1].LeafWaterDataID != 20 {
			t.Fatalf("records=%+v", records)
		}
		if tc.leafVersion == 0 && records[1].LightSample.Color[0].R != 32 {
			t.Fatalf("light sample=%+v", records[1].LightSample)
		}
		if err := file.SetLeaves(records); err != nil {
			t.Fatal(err)
		}
		encoded, err := file.RawLump(bsp.LumpLeafs)
		if err != nil || !bytes.Equal(encoded, payload) {
			t.Fatalf("encoded=%x error=%v", encoded, err)
		}
	}
}

func TestUnknownRevisionIsOpaqueAndCannotBeTypedEdited(t *testing.T) {
	raw := containerBytes(0, []byte("future"))
	binary.LittleEndian.PutUint32(raw[16:], 987)
	file, err := bsp.ReadFromStream(bytes.NewReader(raw), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.Entities(); !errors.Is(err, bsp.ErrUnsupportedLump) {
		t.Fatal(err)
	}
	if err := file.SetEntities("replacement"); !errors.Is(err, bsp.ErrUnsupportedLump) {
		t.Fatal(err)
	}
	stored, err := bsp.NewWriter(file).WriteLump(0)
	if err != nil || string(stored) != "future" {
		t.Fatalf("opaque=%q error=%v", stored, err)
	}
}

func TestUnsupportedProfileStopsAtHeader(t *testing.T) {
	for _, tc := range []struct {
		profile bsp.Profile
		version int32
	}{{"", 20}, {"unknown", 20}, {bsp.ProfileSDK2013, 21}, {bsp.ProfileCSGO, 20}, {bsp.ProfileL4D2, 20}, {bsp.ProfileReactiveDrop, 19}, {bsp.ProfileReactiveDrop, 20}} {
		stream := bytes.NewReader(profileContainer(bsp.ProfileCSGO, tc.version, []byte("unread payload")))
		if _, err := bsp.ReadFromStream(stream, bsp.ReadOptions{Profile: tc.profile}); !errors.Is(err, bsp.ErrUnsupportedFormat) {
			t.Fatalf("profile=%s version=%d error=%v", tc.profile, tc.version, err)
		}
		if stream.Len() != len("unread payload") {
			t.Fatalf("unsupported profile consumed payload: remaining=%d", stream.Len())
		}
	}
}
