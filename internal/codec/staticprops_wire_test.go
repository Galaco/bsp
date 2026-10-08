package codec

import (
	"bytes"
	"encoding/binary"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
	"github.com/go-gl/mathgl/mgl32"
)

// Each width/tail is authored here independently of production layout helpers.
var staticPropWireCases = []struct {
	name    string
	profile string
	version int32
	layout  lumps.StaticPropLayout
	width   int
}{
	{"common4", "sdk2013", 4, lumps.StaticPropCommon4, 56},
	{"common5", "sdk2013", 5, lumps.StaticPropCommon5, 60},
	{"common6", "sdk2013", 6, lumps.StaticPropCommon6, 64},
	{"common7", "l4d2", 7, lumps.StaticPropCommon7, 68},
	{"common8", "l4d2", 8, lumps.StaticPropCommon8, 68},
	{"common9", "l4d2", 9, lumps.StaticPropCommon9, 72},
	{"common10", "csgo", 10, lumps.StaticPropCommon10, 76},
	{"candidate11", "csgo", 11, lumps.StaticPropCommon11, 80},
	{"sdk10", "sdk2013", 10, lumps.StaticPropSDK2013V10, 72},
}

func staticPropTestContext(profile string, version int32) format.Context {
	bspVersion := int32(21)
	if profile == "sdk2013" {
		bspVersion = 20
	}
	return format.Context{Profile: profile, BSPVersion: bspVersion, LumpVersion: version, ByteOrder: binary.LittleEndian}
}

func staticPropLiteralRecords(layout lumps.StaticPropLayout) []byte {
	first := []byte{
		0, 0, 0x80, 0x3f, 0, 0, 0, 0xc0, 0, 0, 0x60, 0x40, // origin 1,-2,3.5
		0, 0, 0x80, 0x40, 0, 0, 0xa0, 0x40, 0, 0, 0xc0, 0x40, // angles 4,5,6
		0, 0, 0, 0, 1, 0, 0x82, 0xf1, // model0 span0+1, solid, flags
		0xf9, 0xff, 0xff, 0xff, // skin -7
		0, 0, 0, 0x41, 0, 0, 0x10, 0x41, // fades 8,9
		0, 0, 0x20, 0x41, 0, 0, 0x30, 0x41, 0, 0, 0x40, 0x41, // lighting 10,11,12
	}
	second := []byte{
		0, 0, 0xc8, 0x42, 0, 0, 0x48, 0xc3, 0, 0x20, 0x96, 0x43, // origin 100,-200,300.25
		0, 0, 0x60, 0x41, 0, 0, 0x70, 0x41, 0, 0, 0x80, 0x41, // angles 14,15,16
		1, 0, 1, 0, 2, 0, 0xff, 0x8a, // model1 span1+2, solid, flags
		0, 0, 0, 0x80, // skin MinInt32
		0, 0, 0x90, 0x41, 0, 0, 0x98, 0x41, // fades 18,19
		0, 0, 0xa0, 0x41, 0, 0, 0xa8, 0x41, 0, 0, 0xb0, 0x41, // lighting 20,21,22
	}
	if layout == lumps.StaticPropSDK2013V10 {
		first[31], second[31] = 0xa5, 0x5a
	}
	switch layout {
	case lumps.StaticPropCommon4:
	case lumps.StaticPropCommon5:
		first = append(first, 0, 0, 0xc0, 0x3f)   // forced 1.5
		second = append(second, 0, 0, 0x20, 0x40) // forced 2.5
	case lumps.StaticPropCommon6:
		first = append(first, 0, 0, 0xc0, 0x3f, 0x34, 0x12, 0xcd, 0xab)
		second = append(second, 0, 0, 0x20, 0x40, 0xff, 0xff, 1, 0)
	case lumps.StaticPropCommon7:
		first = append(first, 0, 0, 0xc0, 0x3f, 0x34, 0x12, 0xcd, 0xab, 1, 2, 3, 4)
		second = append(second, 0, 0, 0x20, 0x40, 0xff, 0xff, 1, 0, 0xff, 0x80, 0x40, 0x20)
	case lumps.StaticPropCommon8:
		first = append(first, 0, 0, 0xc0, 0x3f, 5, 6, 7, 8, 1, 2, 3, 4)
		second = append(second, 0, 0, 0x20, 0x40, 0xff, 0xfe, 0xfd, 0xfc, 0xff, 0x80, 0x40, 0x20)
	case lumps.StaticPropCommon9:
		first = append(first, 0, 0, 0xc0, 0x3f, 5, 6, 7, 8, 1, 2, 3, 4, 0x83, 0xde, 0xad, 0xbe)
		second = append(second, 0, 0, 0x20, 0x40, 0xff, 0xfe, 0xfd, 0xfc, 0xff, 0x80, 0x40, 0x20, 0xff, 0x12, 0x34, 0x56)
	case lumps.StaticPropCommon10:
		first = append(first, 0, 0, 0xc0, 0x3f, 5, 6, 7, 8, 1, 2, 3, 4, 0x83, 0xde, 0xad, 0xbe, 0x78, 0x56, 0x34, 0xf2)
		second = append(second, 0, 0, 0x20, 0x40, 0xff, 0xfe, 0xfd, 0xfc, 0xff, 0x80, 0x40, 0x20, 0xff, 0x12, 0x34, 0x56, 0xef, 0xcd, 0xab, 0x89)
	case lumps.StaticPropCommon11:
		first = append(first, 0, 0, 0xc0, 0x3f, 5, 6, 7, 8, 1, 2, 3, 4, 0x83, 0xde, 0xad, 0xbe, 0x78, 0x56, 0x34, 0xf2, 0, 0, 0x20, 0x40)
		second = append(second, 0, 0, 0x20, 0x40, 0xff, 0xfe, 0xfd, 0xfc, 0xff, 0x80, 0x40, 0x20, 0xff, 0x12, 0x34, 0x56, 0xef, 0xcd, 0xab, 0x89, 0, 0, 0x80, 0x3e)
	case lumps.StaticPropSDK2013V10:
		first = append(first, 0, 0, 0xc0, 0x3f, 0x34, 0x12, 0xcd, 0xab, 0x78, 0x56, 0x34, 0xf2, 0x40, 0, 0x80, 0)
		second = append(second, 0, 0, 0x20, 0x40, 0xff, 0xff, 1, 0, 0xef, 0xcd, 0xab, 0x89, 0xff, 0xff, 1, 0)
	}
	return append(first, second...)
}

func staticPropTestBytes(layout lumps.StaticPropLayout) []byte {
	raw := make([]byte, 274)
	raw[0] = 2 // models
	copy(raw[4:132], "models/first.mdl")
	copy(raw[132:260], "models/second.mdl")
	raw[131], raw[259] = 0xfe, 0xfd // original bytes after NUL
	raw[260] = 3                    // leaves
	raw[266], raw[268] = 2, 4       // leaf IDs 0,2,4
	raw[270] = 2                    // props
	return append(raw, staticPropLiteralRecords(layout)...)
}

func staticPropTestValue(layout lumps.StaticPropLayout) lumps.StaticProps {
	value := lumps.StaticProps{
		Layout: layout, Models: make([]lumps.StaticPropModel, 2), LeafIndices: []uint16{0, 2, 4},
		Props: []lumps.StaticProp{
			{Origin: mgl32.Vec3{1, -2, 3.5}, Angles: mgl32.Vec3{4, 5, 6}, ModelIndex: 0, FirstLeaf: 0, LeafCount: 1, Solid: 0x82, Flags: 0xf1, Skin: -7, FadeMinDist: 8, FadeMaxDist: 9, LightingOrigin: mgl32.Vec3{10, 11, 12}},
			{Origin: mgl32.Vec3{100, -200, 300.25}, Angles: mgl32.Vec3{14, 15, 16}, ModelIndex: 1, FirstLeaf: 1, LeafCount: 2, Solid: 0xff, Flags: 0x8a, Skin: math.MinInt32, FadeMinDist: 18, FadeMaxDist: 19, LightingOrigin: mgl32.Vec3{20, 21, 22}},
		},
	}
	copy(value.Models[0][:], "models/first.mdl")
	copy(value.Models[1][:], "models/second.mdl")
	value.Models[0][127], value.Models[1][127] = 0xfe, 0xfd
	if layout != lumps.StaticPropCommon4 {
		value.Props[0].ForcedFadeScale, value.Props[1].ForcedFadeScale = 1.5, 2.5
	}
	if layout == lumps.StaticPropCommon6 || layout == lumps.StaticPropCommon7 || layout == lumps.StaticPropSDK2013V10 {
		value.Props[0].MinDXLevel, value.Props[0].MaxDXLevel = 0x1234, 0xabcd
		value.Props[1].MinDXLevel, value.Props[1].MaxDXLevel = 0xffff, 1
	}
	if layout == lumps.StaticPropSDK2013V10 {
		value.Props[0].Flags, value.Props[1].Flags = 0xf2345678, 0x89abcdef
		value.Props[0].BasePadding, value.Props[1].BasePadding = 0xa5, 0x5a
		value.Props[0].LightmapResX, value.Props[0].LightmapResY = 64, 128
		value.Props[1].LightmapResX, value.Props[1].LightmapResY = 65535, 1
	} else {
		if layout >= lumps.StaticPropCommon7 {
			value.Props[0].DiffuseModulation = [4]byte{1, 2, 3, 4}
			value.Props[1].DiffuseModulation = [4]byte{255, 128, 64, 32}
		}
		if layout >= lumps.StaticPropCommon8 {
			value.Props[0].MinCPULevel, value.Props[0].MaxCPULevel, value.Props[0].MinGPULevel, value.Props[0].MaxGPULevel = 5, 6, 7, 8
			value.Props[1].MinCPULevel, value.Props[1].MaxCPULevel, value.Props[1].MinGPULevel, value.Props[1].MaxGPULevel = 255, 254, 253, 252
		}
		if layout >= lumps.StaticPropCommon9 {
			value.Props[0].DisableXBox360, value.Props[1].DisableXBox360 = 0x83, 0xff
			value.Props[0].Padding, value.Props[1].Padding = [3]byte{0xde, 0xad, 0xbe}, [3]byte{0x12, 0x34, 0x56}
		}
		if layout >= lumps.StaticPropCommon10 {
			value.Props[0].ExtendedFlags, value.Props[1].ExtendedFlags = 0xf2345678, 0x89abcdef
		}
		if layout == lumps.StaticPropCommon11 {
			value.Props[0].UniformScale, value.Props[1].UniformScale = 2.5, 0.25
		}
	}
	return value
}

func TestStaticPropsIndependentMultipleRecordLayouts(t *testing.T) {
	for _, tc := range staticPropWireCases {
		t.Run(tc.name, func(t *testing.T) {
			context := staticPropTestContext(tc.profile, tc.version)
			raw, want := staticPropTestBytes(tc.layout), staticPropTestValue(tc.layout)
			if len(raw) != 274+2*tc.width {
				t.Fatalf("independent fixture has wrong width")
			}
			if width, err := StaticPropRecordSize(context); err != nil || width != tc.width {
				t.Fatalf("width=%d err=%v", width, err)
			}
			got, err := DecodeStaticProps(raw, context)
			if err != nil || !reflect.DeepEqual(got, want) {
				t.Fatalf("decoded=%+v want=%+v err=%v", got, want, err)
			}
			encoded, err := EncodeStaticProps(want, context)
			if err != nil || !bytes.Equal(encoded, raw) {
				t.Fatalf("independent encoding differs: %v", err)
			}
			if err := ValidateStaticPropReferences(got, 5); err != nil {
				t.Fatal(err)
			}
			got.Props[1].Skin = 0x12345678
			got.Models[1][127] = 0x9a
			got.LeafIndices[2] = 3
			mutated := staticPropTestBytes(tc.layout)
			mutated[259], mutated[268] = 0x9a, 3
			copy(mutated[274+tc.width+32:], []byte{0x78, 0x56, 0x34, 0x12})
			encoded, err = EncodeStaticProps(got, context)
			if err != nil || !bytes.Equal(encoded, mutated) || !bytes.Equal(raw, staticPropTestBytes(tc.layout)) {
				t.Fatalf("mutation changed unrelated fields or aliased input: %v", err)
			}
			context.ByteOrder = nil
			if _, err := DecodeStaticProps(raw, context); err != nil {
				t.Fatalf("canonical nil order rejected: %v", err)
			}
		})
	}
}

func TestStaticPropsExplicitContextMatrix(t *testing.T) {
	matrix := map[string]map[int32]lumps.StaticPropLayout{
		"sdk2013": {4: 4, 5: 5, 6: 6, 10: lumps.StaticPropSDK2013V10},
		"l4d2":    {4: 4, 5: 5, 6: 6, 7: 7, 8: 8, 9: 9},
		"csgo":    {4: 4, 5: 5, 6: 6, 7: 7, 8: 8, 9: 9, 10: 10, 11: 11},
	}
	for profile, accepted := range matrix {
		for _, bspVersion := range []int32{19, 20, 21} {
			for _, version := range []int32{-1, 0, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 99} {
				context := staticPropTestContext(profile, version)
				context.BSPVersion = bspVersion
				want, supported := accepted[version]
				supported = supported && ((profile == "sdk2013" && bspVersion != 21) || (profile != "sdk2013" && bspVersion == 21))
				layout, err := StaticPropLayoutFor(context)
				if supported {
					if err != nil || layout != want {
						t.Fatalf("%+v selected %d: %v", context, layout, err)
					}
					// Empty records cannot be used to infer SDK/common v10.
					got, err := DecodeStaticProps(make([]byte, 12), context)
					if err != nil || got.Layout != want || len(got.Props) != 0 {
						t.Fatalf("empty explicit layout failed: %v", err)
					}
					encoded, err := EncodeStaticProps(lumps.StaticProps{Layout: want}, context)
					if err != nil || !bytes.Equal(encoded, make([]byte, 12)) {
						t.Fatalf("empty encode failed: %v", err)
					}
				} else {
					if err == nil {
						t.Fatalf("accepted %+v", context)
					}
					if _, err := DecodeStaticProps(make([]byte, 12), context); err == nil {
						t.Fatalf("decoded unsupported %+v", context)
					}
					if raw, err := EncodeStaticProps(lumps.StaticProps{Layout: want}, context); err == nil || raw != nil {
						t.Fatalf("encoded unsupported %+v", context)
					}
				}
			}
		}
	}
	for _, context := range []format.Context{{}, {Profile: "unknown", BSPVersion: 21, LumpVersion: 10}, {Profile: "csgo", BSPVersion: 21, LumpVersion: 10, ByteOrder: binary.BigEndian}} {
		if _, err := StaticPropRecordSize(context); err == nil {
			t.Fatalf("accepted invalid context %+v", context)
		}
	}
	// Width is a consistency check, never permission to select another dialect.
	if _, err := DecodeStaticProps(staticPropTestBytes(lumps.StaticPropCommon10), staticPropTestContext("sdk2013", 10)); err == nil {
		t.Fatal("SDK context accepted common v10")
	}
	if _, err := DecodeStaticProps(staticPropTestBytes(lumps.StaticPropSDK2013V10), staticPropTestContext("csgo", 10)); err == nil {
		t.Fatal("CSGO context accepted SDK v10")
	}
}

func TestStaticPropsMalformedBudgetsAndSlots(t *testing.T) {
	for _, tc := range staticPropWireCases {
		context := staticPropTestContext(tc.profile, tc.version)
		raw := staticPropTestBytes(tc.layout)
		for n := 0; n < len(raw); n++ {
			if _, err := DecodeStaticProps(raw[:n], context); err == nil {
				t.Fatalf("%s accepted truncation %d", tc.name, n)
			}
		}
		for _, countOffset := range []int{0, 260, 270} {
			for _, bits := range []uint32{0x80000000, 0xffffffff, 0x7fffffff, 0x40000000} {
				bad := append([]byte(nil), raw...)
				binary.LittleEndian.PutUint32(bad[countOffset:], bits)
				if _, err := DecodeStaticProps(bad, context); err == nil {
					t.Fatalf("%s accepted count %#x at %d", tc.name, bits, countOffset)
				}
			}
		}
		for _, extra := range []int{1, 2, 4, tc.width} {
			if _, err := DecodeStaticProps(append(append([]byte(nil), raw...), make([]byte, extra)...), context); err == nil {
				t.Fatalf("%s accepted trailing %d bytes", tc.name, extra)
			}
		}
		zeroProps := append([]byte(nil), raw[:274]...)
		zeroProps[270] = 0
		if _, err := DecodeStaticProps(zeroProps, context); err != nil {
			t.Fatalf("%s rejected dictionaries with zero props: %v", tc.name, err)
		}
		if _, err := DecodeStaticProps(append(zeroProps, 0), context); err == nil {
			t.Fatalf("%s accepted zero-prop trailing byte", tc.name)
		}
		for _, modelStart := range []int{4, 132} {
			bad := append([]byte(nil), raw...)
			copy(bad[modelStart:modelStart+128], bytes.Repeat([]byte{'x'}, 128))
			if _, err := DecodeStaticProps(bad, context); err == nil {
				t.Fatalf("%s accepted unterminated slot", tc.name)
			}
		}
	}
	for _, counts := range [][4]int{{math.MaxInt32, 0, 0, 80}, {0, math.MaxInt32, 0, 80}, {0, 0, math.MaxInt32, 80}, {-1, 0, 0, 80}, {0, -1, 0, 80}, {0, 0, -1, 80}, {0, 0, 0, 0}} {
		if _, err := staticPropPayloadSize(counts[0], counts[1], counts[2], counts[3]); err == nil {
			t.Fatalf("accepted unrepresentable size %v", counts)
		}
	}
	limit := (math.MaxInt32 - 12) / 128
	if size, err := staticPropPayloadSize(limit, 0, 0, 80); err != nil || size != 12+limit*128 {
		t.Fatalf("bounded size failed: %d %v", size, err)
	}
	if _, err := staticPropPayloadSize(limit+1, 0, 0, 80); err == nil {
		t.Fatal("accepted just-over signed payload limit")
	}
}

func TestStaticPropModelNamePolicy(t *testing.T) {
	var model lumps.StaticPropModel
	copy(model[:], "first\x00retained")
	if name, err := model.Name(); err != nil || name != "first" {
		t.Fatalf("first NUL: %q %v", name, err)
	}
	for _, name := range []string{"", strings.Repeat("x", 127), strings.Repeat("é", 63)} {
		if err := model.SetName(name); err != nil {
			t.Fatal(err)
		}
		if got, err := model.Name(); err != nil || got != name {
			t.Fatalf("name mismatch %q %v", got, err)
		}
		if !bytes.Equal(model[len(name):], make([]byte, 128-len(name))) {
			t.Fatal("edited tail not zeroed")
		}
	}
	for _, name := range []string{strings.Repeat("x", 128), strings.Repeat("é", 64), "embedded\x00nul"} {
		before := model
		if err := model.SetName(name); err == nil || model != before {
			t.Fatal("invalid name accepted or mutation was not atomic")
		}
	}
	copy(model[:], bytes.Repeat([]byte{'x'}, 128))
	if _, err := model.Name(); err == nil {
		t.Fatal("unterminated Name accepted")
	}
	value := staticPropTestValue(lumps.StaticPropCommon10)
	value.Models[1] = model
	if raw, err := EncodeStaticProps(value, staticPropTestContext("csgo", 10)); err == nil || raw != nil {
		t.Fatal("unterminated encode accepted")
	}
}

func TestStaticPropsReferenceValidation(t *testing.T) {
	context := staticPropTestContext("csgo", 10)
	for _, mutation := range []func(*lumps.StaticProp){
		func(p *lumps.StaticProp) { p.ModelIndex = 2 },
		func(p *lumps.StaticProp) { p.ModelIndex = 65535 },
		func(p *lumps.StaticProp) { p.FirstLeaf = 4; p.LeafCount = 0 },
		func(p *lumps.StaticProp) { p.FirstLeaf = 3; p.LeafCount = 1 },
		func(p *lumps.StaticProp) { p.FirstLeaf = 65535; p.LeafCount = 65535 },
	} {
		value := staticPropTestValue(lumps.StaticPropCommon10)
		mutation(&value.Props[1])
		if err := ValidateStaticPropReferences(value, 5); err == nil {
			t.Fatal("bad self reference accepted")
		}
		if raw, err := EncodeStaticProps(value, context); err == nil || raw != nil {
			t.Fatal("bad self reference encoded")
		}
		bad := staticPropTestBytes(lumps.StaticPropCommon10)
		binary.LittleEndian.PutUint16(bad[274+76+24:], value.Props[1].ModelIndex)
		binary.LittleEndian.PutUint16(bad[274+76+26:], value.Props[1].FirstLeaf)
		binary.LittleEndian.PutUint16(bad[274+76+28:], value.Props[1].LeafCount)
		if _, err := DecodeStaticProps(bad, context); err == nil {
			t.Fatal("bad self reference decoded")
		}
	}
	value := staticPropTestValue(lumps.StaticPropCommon10)
	value.Props[1].FirstLeaf, value.Props[1].LeafCount = 3, 0
	if _, err := EncodeStaticProps(value, context); err != nil {
		t.Fatalf("empty end span rejected: %v", err)
	}
	for _, leafCount := range []int{-1, 0, 4} {
		if err := ValidateStaticPropReferences(value, leafCount); err == nil {
			t.Fatalf("accepted exclusive leaf bound %d", leafCount)
		}
	}
	value.Props = nil // unused leaf IDs are still validated
	value.LeafIndices[2] = 65535
	if err := ValidateStaticPropReferences(value, 65535); err == nil {
		t.Fatal("exclusive 65535 leaf bound accepted")
	}
	if err := ValidateStaticPropReferences(value, 65536); err != nil {
		t.Fatalf("uint16 max leaf ID rejected: %v", err)
	}
	value.Models = make([]lumps.StaticPropModel, 65536)
	value.LeafIndices = make([]uint16, 65536)
	value.LeafIndices[65535] = 65535
	value.Props = []lumps.StaticProp{{ModelIndex: 65535, FirstLeaf: 65535, LeafCount: 1}}
	if err := ValidateStaticPropReferences(value, 65536); err != nil {
		t.Fatalf("maximum representable indices rejected: %v", err)
	}
	if raw, err := EncodeStaticProps(value, context); err != nil || len(raw) != 12+65536*128+65536*2+76 {
		t.Fatalf("maximum indices encoding failed: %v", err)
	}
	if err := ValidateStaticPropReferences(lumps.StaticProps{}, 0); err != nil {
		t.Fatal(err)
	}
}

func TestStaticPropsRejectEveryUnavailableField(t *testing.T) {
	type fieldCase struct {
		name    string
		layouts []lumps.StaticPropLayout
		set     func(*lumps.StaticProp)
	}
	fields := []fieldCase{
		{"full flags", []lumps.StaticPropLayout{110}, func(p *lumps.StaticProp) { p.Flags = 256 }},
		{"base padding", []lumps.StaticPropLayout{110}, func(p *lumps.StaticProp) { p.BasePadding = 1 }},
		{"forced fade", []lumps.StaticPropLayout{5, 6, 7, 8, 9, 10, 11, 110}, func(p *lumps.StaticProp) { p.ForcedFadeScale = math.Float32frombits(0x80000000) }},
		{"min DX", []lumps.StaticPropLayout{6, 7, 110}, func(p *lumps.StaticProp) { p.MinDXLevel = 1 }},
		{"max DX", []lumps.StaticPropLayout{6, 7, 110}, func(p *lumps.StaticProp) { p.MaxDXLevel = 1 }},
		{"min CPU", []lumps.StaticPropLayout{8, 9, 10, 11}, func(p *lumps.StaticProp) { p.MinCPULevel = 1 }},
		{"max CPU", []lumps.StaticPropLayout{8, 9, 10, 11}, func(p *lumps.StaticProp) { p.MaxCPULevel = 1 }},
		{"min GPU", []lumps.StaticPropLayout{8, 9, 10, 11}, func(p *lumps.StaticProp) { p.MinGPULevel = 1 }},
		{"max GPU", []lumps.StaticPropLayout{8, 9, 10, 11}, func(p *lumps.StaticProp) { p.MaxGPULevel = 1 }},
		{"diffuse", []lumps.StaticPropLayout{7, 8, 9, 10, 11}, func(p *lumps.StaticProp) { p.DiffuseModulation[3] = 1 }},
		{"disable Xbox", []lumps.StaticPropLayout{9, 10, 11}, func(p *lumps.StaticProp) { p.DisableXBox360 = 1 }},
		{"padding", []lumps.StaticPropLayout{9, 10, 11}, func(p *lumps.StaticProp) { p.Padding[2] = 1 }},
		{"extended flags", []lumps.StaticPropLayout{10, 11}, func(p *lumps.StaticProp) { p.ExtendedFlags = 1 }},
		{"uniform scale", []lumps.StaticPropLayout{11}, func(p *lumps.StaticProp) { p.UniformScale = math.Float32frombits(0x80000000) }},
		{"lightmap X", []lumps.StaticPropLayout{110}, func(p *lumps.StaticProp) { p.LightmapResX = 1 }},
		{"lightmap Y", []lumps.StaticPropLayout{110}, func(p *lumps.StaticProp) { p.LightmapResY = 1 }},
	}
	for _, tc := range staticPropWireCases {
		context := staticPropTestContext(tc.profile, tc.version)
		for _, field := range fields {
			available := false
			for _, layout := range field.layouts {
				available = available || layout == tc.layout
			}
			value := staticPropTestValue(tc.layout)
			field.set(&value.Props[1])
			raw, err := EncodeStaticProps(value, context)
			if available {
				if err != nil {
					t.Fatalf("%s rejected available %s: %v", tc.name, field.name, err)
				}
			} else if err == nil || raw != nil {
				t.Fatalf("%s accepted unavailable %s", tc.name, field.name)
			}
		}
		value := staticPropTestValue(tc.layout)
		value.Layout = 0
		if raw, err := EncodeStaticProps(value, context); err == nil || raw != nil {
			t.Fatalf("%s accepted mismatched layout", tc.name)
		}
	}
}

func TestStaticPropsPreserveAllFloatBits(t *testing.T) {
	patterns := []uint32{0x80000000, 0x7fc12345, 0x7f812345, 0xff800000, 1}
	for _, tc := range staticPropWireCases {
		context := staticPropTestContext(tc.profile, tc.version)
		offsets := []int{0, 4, 8, 12, 16, 20, 36, 40, 44, 48, 52}
		if tc.layout != lumps.StaticPropCommon4 {
			offsets = append(offsets, 56)
		}
		if tc.layout == lumps.StaticPropCommon11 {
			offsets = append(offsets, 76)
		}
		for _, bits := range patterns {
			raw := staticPropTestBytes(tc.layout)
			for _, offset := range offsets {
				binary.LittleEndian.PutUint32(raw[274+offset:], bits)
				binary.LittleEndian.PutUint32(raw[274+tc.width+offset:], bits)
			}
			value, err := DecodeStaticProps(raw, context)
			if err != nil {
				t.Fatal(err)
			}
			encoded, err := EncodeStaticProps(value, context)
			if err != nil || !bytes.Equal(encoded, raw) {
				t.Fatalf("%s changed float bits %#x: %v", tc.name, bits, err)
			}
		}
	}
}

func TestStaticPropsCountMismatchAndNameBoundaries(t *testing.T) {
	for _, tc := range staticPropWireCases {
		context := staticPropTestContext(tc.profile, tc.version)
		for _, change := range [][2]int{{0, 0}, {0, 1}, {0, 3}, {260, 0}, {260, 2}, {260, 4}, {270, 0}, {270, 1}, {270, 3}} {
			raw := staticPropTestBytes(tc.layout)
			binary.LittleEndian.PutUint32(raw[change[0]:], uint32(change[1]))
			if _, err := DecodeStaticProps(raw, context); err == nil {
				t.Fatalf("%s accepted inconsistent count %v", tc.name, change)
			}
		}
		for _, nameLength := range []int{0, 127} {
			raw := staticPropTestBytes(tc.layout)
			copy(raw[4:132], bytes.Repeat([]byte{'x'}, 128))
			raw[4+nameLength] = 0
			value, err := DecodeStaticProps(raw, context)
			if err != nil {
				t.Fatalf("%s rejected boundary name length %d: %v", tc.name, nameLength, err)
			}
			name, err := value.Models[0].Name()
			if err != nil || name != strings.Repeat("x", nameLength) {
				t.Fatalf("%s incorrect boundary name: %q %v", tc.name, name, err)
			}
			encoded, err := EncodeStaticProps(value, context)
			if err != nil || !bytes.Equal(encoded, raw) {
				t.Fatalf("%s lost original name-slot tail: %v", tc.name, err)
			}
		}
	}
}
