package codec

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

func brushSideTestBytes() []byte {
	return []byte{
		0xee, 0xff, 0x00, 0x80, 0xff, 0xff, 0x00, 0x80,
		0xcd, 0xab, 0xff, 0x7f, 0x34, 0x12, 0xff, 0x7f,
	}
}

func brushSideTestRecords(thin bool) []lumps.BrushSide {
	if thin {
		return []lumps.BrushSide{
			{PlaneNum: 0xffee, TexInfo: -32768, DispInfo: -1, Bevel: 0, Thin: 128, HasThin: true},
			{PlaneNum: 0xabcd, TexInfo: 32767, DispInfo: 0x1234, Bevel: 255, Thin: 127, HasThin: true},
		}
	}
	return []lumps.BrushSide{
		{PlaneNum: 0xffee, TexInfo: -32768, DispInfo: -1, Bevel: -32768},
		{PlaneNum: 0xabcd, TexInfo: 32767, DispInfo: 0x1234, Bevel: 32767},
	}
}

func TestBrushSidesIndependentDialectBytes(t *testing.T) {
	for _, profile := range []string{"sdk2013", "csgo", "l4d2"} {
		context := lightTestContext(profile, 0)
		raw, want := brushSideTestBytes(), brushSideTestRecords(profile == "csgo")
		got, err := DecodeBrushSides(raw, context)
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("%s: got=%+v err=%v want=%+v", profile, got, err, want)
		}
		encoded, err := EncodeBrushSides(want, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("%s: independent encode differs: %v", profile, err)
		}
		got[1].TexInfo = -7
		mutated := brushSideTestBytes()
		mutated[10], mutated[11] = 0xf9, 0xff
		encoded, err = EncodeBrushSides(got, context)
		if err != nil || !bytes.Equal(encoded, mutated) {
			t.Fatalf("%s: mutation changed unrelated fields: %v", profile, err)
		}
		context.ByteOrder = nil
		if got, err := DecodeBrushSides(raw, context); err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("%s: canonical nil order: %v", profile, err)
		}
	}
	sdk19 := lightTestContext("sdk2013", 0)
	sdk19.BSPVersion = 19
	if got, err := DecodeBrushSides(brushSideTestBytes(), sdk19); err != nil || !reflect.DeepEqual(got, brushSideTestRecords(false)) {
		t.Fatalf("SDK19 short bevel failed: %v", err)
	}
}

func TestBrushSidesMalformedAndContexts(t *testing.T) {
	raw := brushSideTestBytes()
	for _, profile := range []string{"sdk2013", "csgo", "l4d2"} {
		context := lightTestContext(profile, 0)
		for n := 1; n < len(raw); n++ {
			if n == 8 {
				continue
			}
			if got, err := DecodeBrushSides(raw[:n], context); err == nil || got != nil {
				t.Fatalf("%s accepted incomplete length %d", profile, n)
			}
		}
		for n := 1; n < 8; n++ {
			if got, err := DecodeBrushSides(append(append([]byte(nil), raw...), make([]byte, n)...), context); err == nil || got != nil {
				t.Fatalf("%s accepted trailing length %d", profile, n)
			}
		}
		if got, err := DecodeBrushSides(nil, context); err != nil || len(got) != 0 {
			t.Fatalf("%s empty payload failed: %v", profile, err)
		}
	}
	invalid := []format.Context{
		{}, lightTestContext("unknown", 0), lightTestContext("sdk2013", 1), lightTestContext("csgo", 1), lightTestContext("l4d2", -1),
		{Profile: "csgo", BSPVersion: 20}, {Profile: "sdk2013", BSPVersion: 21}, {Profile: "l4d2", BSPVersion: 19},
		{Profile: "csgo", BSPVersion: 21, ByteOrder: binary.BigEndian},
	}
	for _, context := range invalid {
		if got, err := DecodeBrushSides(nil, context); err == nil || got != nil {
			t.Fatalf("decode accepted invalid context %+v", context)
		}
		if got, err := EncodeBrushSides(nil, context); err == nil || got != nil {
			t.Fatalf("encode accepted invalid context %+v", context)
		}
	}
}

func TestBrushSidesRejectNonrepresentableSetters(t *testing.T) {
	csgo := lightTestContext("csgo", 0)
	for _, bevel := range []int16{-32768, -1, 256, 32767} {
		records := brushSideTestRecords(true)
		records[1].Bevel = bevel
		if raw, err := EncodeBrushSides(records, csgo); err == nil || raw != nil {
			t.Fatalf("accepted non-byte bevel %d", bevel)
		}
	}
	for _, profile := range []string{"sdk2013", "csgo", "l4d2"} {
		context := lightTestContext(profile, 0)
		records := brushSideTestRecords(profile == "csgo")
		records[1].HasThin = !records[1].HasThin
		if raw, err := EncodeBrushSides(records, context); err == nil || raw != nil {
			t.Fatalf("%s silently converted tail dialect", profile)
		}
		if profile != "csgo" {
			records = brushSideTestRecords(false)
			records[1].Thin = 1
			if raw, err := EncodeBrushSides(records, context); err == nil || raw != nil {
				t.Fatalf("%s discarded thin", profile)
			}
		}
	}
	// Presence is explicit even when Thin is zero; byte values are not booleans.
	records := brushSideTestRecords(true)
	records[0].Bevel, records[0].Thin = 128, 0
	records[1].Thin = 255
	want := brushSideTestBytes()
	want[6], want[7], want[15] = 128, 0, 255
	raw, err := EncodeBrushSides(records, csgo)
	if err != nil || !bytes.Equal(raw, want) {
		t.Fatalf("byte boundaries/explicit zero thin not retained: %v", err)
	}
	if got, err := DecodeBrushSides(want, csgo); err != nil || !reflect.DeepEqual(got, records) {
		t.Fatalf("byte boundaries did not reread: %+v %v", got, err)
	}
}
