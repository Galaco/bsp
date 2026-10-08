package codec

import (
	"bytes"
	"encoding/binary"
	"math"
	"reflect"
	"testing"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
	"github.com/go-gl/mathgl/mgl32"
)

func lightTestContext(profile string, revision int32) format.Context {
	version := int32(21)
	if profile == "sdk2013" {
		version = 20
	}
	return format.Context{Profile: profile, BSPVersion: version, LumpVersion: revision, ByteOrder: binary.LittleEndian}
}

// Fixture words are authored independently of the production record/encoder.
func worldLightTestBytes(revision int32) []byte {
	width, cluster := 88, 36
	if revision == 1 {
		width, cluster = 100, 48
	}
	raw := make([]byte, width*2)
	vectors := [2][9]uint32{
		{0x3f800000, 0x40000000, 0x40400000, 0x40800000, 0x40a00000, 0x40c00000, 0x40e00000, 0x41000000, 0x41100000},
		{0xbf800000, 0xc0000000, 0xc0400000, 0xc0800000, 0xc0a00000, 0xc0c00000, 0xc0e00000, 0xc1000000, 0xc1100000},
	}
	tails := [2][13]uint32{
		{0xffffffff, 2, 0x10203040, 0x41200000, 0x41300000, 0x41400000, 0x41500000, 0x41600000, 0x41700000, 0x41800000, 0x55667788, 0xfffffcf7, 0x40000000},
		{0x7fffffff, 5, 0x80000000, 0xc1200000, 0xc1300000, 0xc1400000, 0xc1500000, 0xc1600000, 0xc1700000, 0xc1800000, 0x87654321, 0x12345678, 0xffffffff},
	}
	for i := 0; i < 2; i++ {
		wire := raw[i*width : (i+1)*width]
		for j, bits := range vectors[i] {
			binary.LittleEndian.PutUint32(wire[j*4:], bits)
		}
		if revision == 1 {
			shadow := [2][3]uint32{{0x41880000, 0x41900000, 0x41980000}, {0xc1880000, 0xc1900000, 0xc1980000}}
			for j, bits := range shadow[i] {
				binary.LittleEndian.PutUint32(wire[36+j*4:], bits)
			}
		}
		for j, bits := range tails[i] {
			binary.LittleEndian.PutUint32(wire[cluster+j*4:], bits)
		}
	}
	return raw
}

func worldLightTestRecords(revision int32) []lumps.WorldLight {
	records := []lumps.WorldLight{
		{Origin: mgl32.Vec3{1, 2, 3}, Intensity: mgl32.Vec3{4, 5, 6}, Normal: mgl32.Vec3{7, 8, 9}, Cluster: -1, Type: lumps.EmitSpotLight, Style: 0x10203040, Stopdot: 10, Stopdot2: 11, Exponent: 12, Radius: 13, ConstantAttenuation: 14, LinearAttenuation: 15, QuadraticAttenuation: 16, Flags: 0x55667788, TexInfo: -777, Owner: 0x40000000},
		{Origin: mgl32.Vec3{-1, -2, -3}, Intensity: mgl32.Vec3{-4, -5, -6}, Normal: mgl32.Vec3{-7, -8, -9}, Cluster: 0x7fffffff, Type: lumps.EmitSkyAmbient, Style: -2147483648, Stopdot: -10, Stopdot2: -11, Exponent: -12, Radius: -13, ConstantAttenuation: -14, LinearAttenuation: -15, QuadraticAttenuation: -16, Flags: -2023406815, TexInfo: 0x12345678, Owner: -1},
	}
	if revision == 1 {
		records[0].ShadowCastOffset = mgl32.Vec3{17, 18, 19}
		records[1].ShadowCastOffset = mgl32.Vec3{-17, -18, -19}
	}
	return records
}

func TestWorldLightsIndependentVersions(t *testing.T) {
	contexts := []format.Context{lightTestContext("sdk2013", 0), lightTestContext("csgo", 0), lightTestContext("l4d2", 0), lightTestContext("csgo", 1)}
	sdk19 := lightTestContext("sdk2013", 0)
	sdk19.BSPVersion = 19
	contexts = append(contexts, sdk19)
	for _, context := range contexts {
		raw, want := worldLightTestBytes(context.LumpVersion), worldLightTestRecords(context.LumpVersion)
		got, err := DecodeWorldLights(raw, context)
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("%+v: records=%+v err=%v want=%+v", context, got, err, want)
		}
		encoded, err := EncodeWorldLights(want, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("%+v: independent encoding differs, err=%v", context, err)
		}
		got[1].Owner = 77
		mutated := append([]byte(nil), raw...)
		binary.LittleEndian.PutUint32(mutated[len(mutated)-4:], 77)
		encoded, err = EncodeWorldLights(got, context)
		if err != nil || !bytes.Equal(encoded, mutated) {
			t.Fatalf("%+v: supported mutation changed other fields: %v", context, err)
		}
		context.ByteOrder = nil
		if got, err := DecodeWorldLights(raw, context); err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("canonical nil byte order failed: %v", err)
		}
	}
}

func TestWorldLightsUnknownEnumsAndFloatBits(t *testing.T) {
	for _, revision := range []int32{0, 1} {
		context := lightTestContext("csgo", revision)
		raw := worldLightTestBytes(revision)
		width, enumOffset := 88, 40
		if revision == 1 {
			width, enumOffset = 100, 52
		}
		binary.LittleEndian.PutUint32(raw[enumOffset:], 0x12345678)
		binary.LittleEndian.PutUint32(raw[width+enumOffset:], 0x80000001)
		binary.LittleEndian.PutUint32(raw, 0x7fc12345)
		binary.LittleEndian.PutUint32(raw[width+4:], 0x80000000)
		records, err := DecodeWorldLights(raw, context)
		if err != nil || records[0].Type != lumps.EmitType(0x12345678) || records[1].Type != lumps.EmitType(-2147483647) {
			t.Fatalf("unknown signed32 enum not preserved: %+v %v", records, err)
		}
		encoded, err := EncodeWorldLights(records, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("enum/NaN/negative-zero bits not preserved: %v", err)
		}
	}
}

func TestWorldLightsMalformedAndContext(t *testing.T) {
	for _, revision := range []int32{0, 1} {
		context := lightTestContext("csgo", revision)
		raw := worldLightTestBytes(revision)
		width := len(raw) / 2
		for n := 1; n < len(raw); n++ {
			if n == width {
				continue // One complete record is independently valid.
			}
			if got, err := DecodeWorldLights(raw[:n], context); err == nil || got != nil {
				t.Fatalf("revision %d accepted incomplete length %d", revision, n)
			}
		}
		for n := 1; n < width; n++ {
			if got, err := DecodeWorldLights(append(append([]byte(nil), raw...), make([]byte, n)...), context); err == nil || got != nil {
				t.Fatalf("revision %d accepted trailing length %d", revision, n)
			}
		}
		if got, err := DecodeWorldLights(nil, context); err != nil || len(got) != 0 {
			t.Fatalf("empty payload: %v", err)
		}
	}
	invalid := []format.Context{
		{}, lightTestContext("unknown", 0), lightTestContext("csgo", -1), lightTestContext("csgo", 2), lightTestContext("sdk2013", 1), lightTestContext("l4d2", 1),
		{Profile: "csgo", BSPVersion: 20}, {Profile: "sdk2013", BSPVersion: 21}, {Profile: "l4d2", BSPVersion: 20},
		{Profile: "csgo", BSPVersion: 21, ByteOrder: binary.BigEndian},
	}
	for _, context := range invalid {
		if got, err := DecodeWorldLights(nil, context); err == nil || got != nil {
			t.Fatalf("decode accepted invalid context %+v", context)
		}
		if got, err := EncodeWorldLights(nil, context); err == nil || got != nil {
			t.Fatalf("encode accepted invalid context %+v", context)
		}
	}
}

func TestWorldLightsRejectUnrepresentableShadow(t *testing.T) {
	context := lightTestContext("csgo", 0)
	for axis := 0; axis < 3; axis++ {
		for _, bits := range []uint32{0x3f800000, 0xbf800000, 0x80000000, 0x7fc12345} {
			records := worldLightTestRecords(0)
			records[1].ShadowCastOffset[axis] = math.Float32frombits(bits)
			if raw, err := EncodeWorldLights(records, context); err == nil || raw != nil {
				t.Fatalf("discarded shadow bits %08x on axis %d", bits, axis)
			}
		}
	}
}

func TestLightSidePayloadSizeBoundaries(t *testing.T) {
	for _, width := range []int{8, 88, 100} {
		limit := int((int64(1)<<31 - 1) / int64(width))
		if size, err := lightSidePayloadSize(limit, width); err != nil || size != limit*width {
			t.Fatalf("valid signed32 boundary width %d: %d %v", width, size, err)
		}
		if _, err := lightSidePayloadSize(limit+1, width); err == nil {
			t.Fatalf("signed32 overflow width %d accepted", width)
		}
		if _, err := lightSidePayloadSize(int(^uint(0)>>1), width); err == nil {
			t.Fatalf("native multiplication overflow width %d accepted", width)
		}
	}
}
