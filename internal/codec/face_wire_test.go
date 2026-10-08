package codec

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"

	"github.com/galaco/bsp/internal/format"
)

func faceWireLiteral() []byte {
	raw := make([]byte, 112)
	for i := 0; i < 2; i++ {
		wire := raw[i*56 : (i+1)*56]
		binary.LittleEndian.PutUint16(wire, uint16(19+i))
		wire[2], wire[3] = byte(0xe1+i), byte(0xd2+i)
		binary.LittleEndian.PutUint32(wire[4:], 0xfffffff0+uint32(i))
		for j := 0; j < 4; j++ {
			binary.LittleEndian.PutUint16(wire[8+j*2:], uint16(0x8a00+i*10+j))
		}
		copy(wire[16:20], []byte{0x80, 0x91, 0xa2, byte(i)})
		binary.LittleEndian.PutUint32(wire[20:], 0xffffffff-uint32(i))
		binary.LittleEndian.PutUint32(wire[24:], uint32(0x7fc12345))
		if i == 1 {
			binary.LittleEndian.PutUint32(wire[24:], 0x80000000)
		}
		for j := 0; j < 4; j++ {
			binary.LittleEndian.PutUint32(wire[28+j*4:], 0xffff0000+uint32(i*10+j))
		}
		binary.LittleEndian.PutUint32(wire[44:], 0xfffffffe-uint32(i))
		binary.LittleEndian.PutUint16(wire[48:], uint16(0x8001+i))
		binary.LittleEndian.PutUint16(wire[50:], uint16(7+i))
		binary.LittleEndian.PutUint32(wire[52:], 0xfedcba90+uint32(i))
	}
	return raw
}

func TestFaceCompleteWireFieldsAndPackedEdits(t *testing.T) {
	for _, pc := range []struct {
		profile string
		version int32
	}{{"sdk2013", 19}, {"sdk2013", 20}, {"csgo", 21}, {"l4d2", 21}, {"reactivedrop", 21}} {
		for _, revision := range []int32{0, 1} {
			context := format.Context{Profile: pc.profile, BSPVersion: pc.version, LumpVersion: revision}
			raw := faceWireLiteral()
			records, err := DecodeFaces(raw, context)
			if err != nil {
				t.Fatal(err)
			}
			r := &records[1]
			if len(records) != 2 || r.Planenum != 20 || r.Side != 0xe2 || r.OnNode != 0xd3 || r.FirstEdge != -15 || uint16(r.NumEdges) != 0x8a0a || uint16(r.TexInfo) != 0x8a0b || uint16(r.DispInfo) != 0x8a0c || uint16(r.SurfaceFogVolumeID) != 0x8a0d || r.Styles != [4]byte{0x80, 0x91, 0xa2, 1} || r.Lightofs != -2 || math.Float32bits(r.Area) != 0x80000000 || uint32(r.LightmapTextureMinsInLuxels[1]) != 0xffff000b || uint32(r.LightmapTextureSizeInLuxels[1]) != 0xffff000d || r.OrigFace != -3 || r.PrimitiveCount() != 2 || r.DynamicShadowsEnabled() || r.FirstPrimID != 8 || r.SmoothingGroups != 0xfedcba91 || math.Float32bits(records[0].Area) != 0x7fc12345 {
				t.Fatalf("complete face=%+v", r)
			}
			encoded, err := EncodeFaces(records, context)
			if err != nil || !bytes.Equal(encoded, raw) {
				t.Fatalf("untouched packed/float bytes lost: %v", err)
			}
			r.SetPrimitiveCount(32767)
			r.SetDynamicShadowsEnabled(true)
			binary.LittleEndian.PutUint16(raw[56+48:], 0x7fff)
			r.LightmapTextureSizeInLuxels[1] = -3
			binary.LittleEndian.PutUint32(raw[56+40:], 0xfffffffd)
			encoded, err = EncodeFaces(records, context)
			if err != nil || !bytes.Equal(encoded, raw) {
				t.Fatalf("late packed/edit bytes changed: %v", err)
			}
		}
	}
}

func TestFacePartialRecordsAndUnsupportedContexts(t *testing.T) {
	context := format.Context{Profile: "sdk2013", BSPVersion: 20}
	raw := faceWireLiteral()
	for length := 1; length < 112; length++ {
		if length == 56 {
			continue
		}
		if _, err := DecodeFaces(raw[:length], context); err == nil {
			t.Fatalf("accepted%d-byte partial faces", length)
		}
	}
	for _, bad := range []format.Context{{Profile: "sdk2013", BSPVersion: 21}, {Profile: "reactivedrop", BSPVersion: 20}, {Profile: "unknown", BSPVersion: 21}, {Profile: "csgo", BSPVersion: 21, LumpVersion: 2}, {Profile: "csgo", BSPVersion: 21, ByteOrder: binary.BigEndian}} {
		if _, err := DecodeFaces(raw, bad); err == nil {
			t.Fatalf("accepted context%+v", bad)
		}
	}
	adapter := new(Faces)
	adapter.SetContext(context)
	if err := adapter.Decode(raw); err != nil {
		t.Fatal(err)
	}
	if err := adapter.Decode(raw[:111]); err == nil {
		t.Fatal("partial second record accepted")
	}
	after, err := adapter.Encode()
	if err != nil || adapter.Length() != 112 || !bytes.Equal(after, raw) {
		t.Fatalf("failed face decode published partial state: %v", err)
	}
}
