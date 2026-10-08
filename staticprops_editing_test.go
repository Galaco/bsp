package bsp_test

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/galaco/bsp"
	"github.com/galaco/bsp/lumps"
)

const editingSprpID uint32 = 0x73707270

type editingGameWire struct {
	id             uint32
	flags, version uint16
	logical        int32
	data           []byte
}

// Author all container/directory fields directly. No production static-prop or
// game encoder participates in the input fixture or the expected payload.
func readStaticEditingFixture(t *testing.T, profile bsp.Profile, entries []editingGameWire) *bsp.File {
	t.Helper()
	game := make([]byte, 4+16*(len(entries)+1))
	binary.LittleEndian.PutUint32(game, uint32(len(entries)+1))
	game = append(game, []byte("HEAD")...)
	for i, entry := range entries {
		for len(game)%4 != 0 {
			game = append(game, 0x9b)
		}
		length := entry.logical
		if entry.flags&1 == 0 {
			length = int32(len(entry.data))
		}
		putGameEntry(game, i, entry.id, entry.flags, entry.version, int32(bsp.HeaderSize+len(game)), length)
		game = append(game, entry.data...)
		game = append(game, []byte("TGAP")...)
	}
	putGameEntry(game, len(entries), 0, 0, 0, int32(bsp.HeaderSize+len(game)), 0)
	raw := make([]byte, bsp.HeaderSize)
	copy(raw, "VBSP")
	version := uint32(21)
	if profile == bsp.ProfileSDK2013 {
		version = 20
	}
	binary.LittleEndian.PutUint32(raw[4:], version)
	binary.LittleEndian.PutUint32(raw[8+35*16:], bsp.HeaderSize)
	binary.LittleEndian.PutUint32(raw[12+35*16:], uint32(len(game)))
	raw = append(raw, game...)
	leafOffset := len(raw)
	// Three independently authored version-1 leaves, with no water record.
	leaves := make([]byte, 3*32)
	for i := 0; i < 3; i++ {
		binary.LittleEndian.PutUint16(leaves[i*32+28:], 0xffff)
	}
	raw = append(raw, leaves...)
	binary.LittleEndian.PutUint32(raw[8+10*16:], uint32(leafOffset))
	binary.LittleEndian.PutUint32(raw[12+10*16:], uint32(len(leaves)))
	binary.LittleEndian.PutUint32(raw[16+10*16:], 1)
	file, err := bsp.Read(bytes.NewReader(raw), bsp.ReadOptions{Profile: profile})
	if err != nil {
		t.Fatal(err)
	}
	return file
}

func staticEditingEntries(t *testing.T, version uint16, payload []byte) []editingGameWire {
	t.Helper()
	// The sibling frame is the existing independently compressed plane fixture;
	// neither its arbitrary ID nor its version gives it typed static-prop meaning.
	frame := compressedPlaneFixture(t)
	return []editingGameWire{
		{0xdecafbad, 0x20, 77, 0, []byte{0xa1, 0, 0xb2, 0xc3}},
		{editingSprpID, 0x8000, version, 0, payload},
		{0xdecafbad, 0x4001, 99, 40, frame},
	}
}

// Common10/11 and SDK10 share the dictionary framing but not their record
// tails. Nonzero bytes after each name's first NUL must remain opaque slot data.
func staticEditingPayload(sdk bool, version uint16) []byte {
	width := 76
	if sdk {
		width = 72
	} else if version == 11 {
		width = 80
	}
	raw := make([]byte, 274+2*width)
	binary.LittleEndian.PutUint32(raw, 2)
	copy(raw[4:132], "models/first.mdl")
	copy(raw[132:260], "models/second.mdl")
	raw[80], raw[131], raw[205], raw[259] = 0xa1, 0xf1, 0xb2, 0xf2
	binary.LittleEndian.PutUint32(raw[260:], 3)
	binary.LittleEndian.PutUint16(raw[264:], 0)
	binary.LittleEndian.PutUint16(raw[266:], 2)
	binary.LittleEndian.PutUint16(raw[268:], 1)
	binary.LittleEndian.PutUint32(raw[270:], 2)
	for i := 0; i < 2; i++ {
		p := raw[274+i*width:]
		for axis := 0; axis < 3; axis++ {
			binary.LittleEndian.PutUint32(p[axis*4:], math.Float32bits(float32(100*i+axis)+0.25))
			binary.LittleEndian.PutUint32(p[12+axis*4:], math.Float32bits(float32(10*i+axis)+0.5))
			binary.LittleEndian.PutUint32(p[44+axis*4:], math.Float32bits(float32(20*i+axis)+0.75))
		}
		binary.LittleEndian.PutUint16(p[24:], uint16(i))
		if i == 0 {
			binary.LittleEndian.PutUint16(p[26:], 0)
			binary.LittleEndian.PutUint16(p[28:], 2)
		} else {
			binary.LittleEndian.PutUint16(p[26:], 2)
			binary.LittleEndian.PutUint16(p[28:], 1)
		}
		p[30] = byte(4 + i)
		p[31] = byte(0xd1 + i)
		binary.LittleEndian.PutUint32(p[32:], uint32(0xfffffff9+uint32(i))) // skin -7/-6
		binary.LittleEndian.PutUint32(p[36:], math.Float32bits(float32(30+i)))
		binary.LittleEndian.PutUint32(p[40:], math.Float32bits(float32(300+i)))
		binary.LittleEndian.PutUint32(p[56:], math.Float32bits(float32(i)+1.25))
		if sdk {
			p[31] = byte(0x6e + i) // SDK padding is not the primary flags byte.
			binary.LittleEndian.PutUint16(p[60:], uint16(70+i))
			binary.LittleEndian.PutUint16(p[62:], uint16(95+i))
			binary.LittleEndian.PutUint32(p[64:], 0x81234567+uint32(i))
			binary.LittleEndian.PutUint16(p[68:], uint16(17+i))
			binary.LittleEndian.PutUint16(p[70:], uint16(31+i))
		} else {
			copy(p[60:64], []byte{1 + byte(i), 3 + byte(i), 5 + byte(i), 7 + byte(i)})
			copy(p[64:68], []byte{0x11 + byte(i), 0x22 + byte(i), 0x88 + byte(i), 0xf0 + byte(i)})
			copy(p[68:72], []byte{0x81 + byte(i), 0xa1 + byte(i), 0xb1 + byte(i), 0xc1 + byte(i)})
			binary.LittleEndian.PutUint32(p[72:], 0x80000001+uint32(i))
			if version == 11 {
				scale := float32(2.5)
				if i == 1 {
					scale = 0.75
				}
				binary.LittleEndian.PutUint32(p[76:], math.Float32bits(scale))
			}
		}
	}
	return raw
}

func requireStaticEditingSiblings(t *testing.T, file *bsp.File, original []editingGameWire) {
	t.Helper()
	entries, err := file.GameLumps()
	if err != nil || len(entries) != len(original)+1 {
		t.Fatalf("directory=%+v, error=%v", entries, err)
	}
	for i, want := range original {
		got := entries[i]
		if got.ID != want.id || got.Flags != want.flags || got.Version != want.version {
			t.Fatalf("entry %d metadata=%+v, want id=%x flags=%x version=%d", i, got, want.id, want.flags, want.version)
		}
		if i == 1 {
			continue
		}
		if want.flags&1 != 0 {
			// This last sibling's independent frame and four-byte gap end on
			// an aligned terminal boundary; its entire stored span is known.
			wantStored := append(append([]byte(nil), want.data...), []byte("TGAP")...)
			if got.LogicalLength != want.logical || !bytes.Equal(got.StoredData, wantStored) {
				t.Fatalf("compressed sibling %d changed: %+v", i, got)
			}
		} else if got.LogicalLength != int32(len(want.data)) || !bytes.Equal(got.StoredData, want.data) {
			t.Fatalf("opaque sibling %d changed: %+v", i, got)
		}
	}
	terminal := entries[len(entries)-1]
	if terminal.ID != 0 || terminal.Flags != 0 || terminal.Version != 0 || terminal.LogicalLength != 0 || len(terminal.StoredData) != 0 {
		t.Fatalf("terminal boundary changed: %+v", terminal)
	}
	stored, err := file.RawLump(35)
	if err != nil || !bytes.Contains(stored, []byte("HEAD")) || bytes.Count(stored, []byte("TGAP")) != len(original) {
		t.Fatalf("directory gaps changed: %x, error=%v", stored, err)
	}
}

func requireStaticEditingPayload(t *testing.T, file *bsp.File, want []byte) {
	t.Helper()
	got, err := file.DecodedGameLump(1)
	if err != nil || !bytes.Equal(got, want) {
		t.Fatalf("sprp payload differs: got=%x want=%x error=%v", got, want, err)
	}
}

func TestStaticPropsPublicReadEditExport(t *testing.T) {
	for _, tc := range []struct {
		name    string
		profile bsp.Profile
		version uint16
		layout  lumps.StaticPropLayout
		sdk     bool
	}{
		{"csgo10", bsp.ProfileCSGO, 10, lumps.StaticPropCommon10, false},
		{"csgo11", bsp.ProfileCSGO, 11, lumps.StaticPropCommon11, false},
		{"sdk10", bsp.ProfileSDK2013, 10, lumps.StaticPropSDK2013V10, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			original := staticEditingPayload(tc.sdk, tc.version)
			entries := staticEditingEntries(t, tc.version, original)
			file := readStaticEditingFixture(t, tc.profile, entries)
			value, err := file.StaticProps()
			if err != nil {
				t.Fatal(err)
			}
			if value.Layout != tc.layout || len(value.Models) != 2 || len(value.LeafIndices) != 3 || len(value.Props) != 2 {
				t.Fatalf("snapshot framing=%+v", value)
			}
			if !reflect.DeepEqual(value.LeafIndices, []uint16{0, 2, 1}) {
				t.Fatalf("leaf-index ordering=%v", value.LeafIndices)
			}
			for i := range value.Models {
				if !bytes.Equal(value.Models[i][:], original[4+i*128:132+i*128]) {
					t.Fatalf("model %d lost slot tail", i)
				}
			}
			name, err := value.Models[1].Name()
			if err != nil || name != "models/second.mdl" {
				t.Fatalf("second model=%q, error=%v", name, err)
			}
			for i, p := range value.Props {
				if p.Origin[0] != float32(100*i)+0.25 || p.Origin[1] != float32(100*i+1)+0.25 || p.Origin[2] != float32(100*i+2)+0.25 || p.Skin != int32(-7+i) || p.ModelIndex != uint16(i) || p.ForcedFadeScale != float32(i)+1.25 {
					t.Fatalf("record %d alignment/base fields=%+v", i, p)
				}
				for axis := range p.Angles {
					if p.Angles[axis] != float32(10*i+axis)+0.5 || p.LightingOrigin[axis] != float32(20*i+axis)+0.75 {
						t.Fatalf("record %d vector axis %d=%+v", i, axis, p)
					}
				}
				if p.Solid != byte(4+i) || p.FadeMinDist != float32(30+i) || p.FadeMaxDist != float32(300+i) || p.FirstLeaf != []uint16{0, 2}[i] || p.LeafCount != []uint16{2, 1}[i] {
					t.Fatalf("record %d scalar/reference fields=%+v", i, p)
				}
				if tc.sdk {
					if p.Flags != 0x81234567+uint32(i) || p.BasePadding != byte(0x6e+i) || p.LightmapResX != uint16(17+i) || p.LightmapResY != uint16(31+i) || p.MinDXLevel != uint16(70+i) || p.MaxDXLevel != uint16(95+i) {
						t.Fatalf("SDK tail %d=%+v", i, p)
					}
				} else {
					if p.Flags != uint32(0xd1+i) || p.DiffuseModulation != [4]byte{0x11 + byte(i), 0x22 + byte(i), 0x88 + byte(i), 0xf0 + byte(i)} || p.DisableXBox360 != byte(0x81+i) || p.Padding != [3]byte{0xa1 + byte(i), 0xb1 + byte(i), 0xc1 + byte(i)} || p.ExtendedFlags != 0x80000001+uint32(i) {
						t.Fatalf("common tail %d=%+v", i, p)
					}
					if p.MinCPULevel != byte(1+i) || p.MaxCPULevel != byte(3+i) || p.MinGPULevel != byte(5+i) || p.MaxGPULevel != byte(7+i) {
						t.Fatalf("CPU/GPU levels %d=%+v", i, p)
					}
					if tc.version == 11 && p.UniformScale != []float32{2.5, 0.75}[i] {
						t.Fatalf("scale %d=%v", i, p.UniformScale)
					}
				}
			}
			// Committing an unchanged snapshot preserves every original byte,
			// including arbitrary name-slot tails and nonboolean byte flags.
			if err := file.SetStaticProps(value); err != nil {
				t.Fatal(err)
			}
			requireStaticEditingPayload(t, file, original)
			requireStaticEditingSiblings(t, file, entries)

			value.Models[0][76] = 0x99
			value.LeafIndices[1] = 0
			value.Props[0].Origin[0] = 123.25
			fresh, err := file.StaticProps()
			if err != nil || fresh.Models[0][76] != 0xa1 || fresh.LeafIndices[1] != 2 || fresh.Props[0].Origin[0] != 0.25 {
				t.Fatalf("snapshot mutation autocommitted: %+v, error=%v", fresh, err)
			}
			uncommitted := rereadOrdinary(t, file)
			requireStaticEditingPayload(t, uncommitted, original)
			if err := value.Models[0].SetName("models/edited.mdl"); err != nil {
				t.Fatal(err)
			}
			want := append([]byte(nil), original...)
			for i := 4; i < 132; i++ {
				want[i] = 0
			}
			copy(want[4:132], "models/edited.mdl")
			binary.LittleEndian.PutUint16(want[266:], 0)
			binary.LittleEndian.PutUint32(want[274:], math.Float32bits(123.25))
			if tc.sdk {
				value.Props[0].Flags = 0xfedcba98
				value.Props[0].LightmapResX, value.Props[0].LightmapResY = 123, 456
				binary.LittleEndian.PutUint32(want[274+64:], 0xfedcba98)
				binary.LittleEndian.PutUint16(want[274+68:], 123)
				binary.LittleEndian.PutUint16(want[274+70:], 456)
			} else {
				value.Props[0].Flags = 0xe7
				value.Props[0].DiffuseModulation = [4]byte{7, 19, 137, 249}
				value.Props[0].ExtendedFlags = 0xf2345678
				want[274+31] = 0xe7
				copy(want[274+64:274+68], []byte{7, 19, 137, 249})
				binary.LittleEndian.PutUint32(want[274+72:], 0xf2345678)
				if tc.version == 11 {
					value.Props[0].UniformScale = 3.5
					binary.LittleEndian.PutUint32(want[274+76:], math.Float32bits(3.5))
				}
			}
			if err := file.SetStaticProps(value); err != nil {
				t.Fatal(err)
			}
			// Subsequent caller edits must not alias committed encoder output.
			value.Models[0][0] = 'X'
			value.LeafIndices[0] = 2
			value.Props[0].Skin = 99
			file = rereadOrdinary(t, file)
			requireStaticEditingPayload(t, file, want)
			requireStaticEditingSiblings(t, file, entries)
			persisted, err := file.StaticProps()
			if err != nil || persisted.Props[0].Origin[0] != 123.25 || persisted.Props[1].Origin[0] != 100.25 || persisted.Props[0].Skin != -7 || persisted.LeafIndices[0] != 0 {
				t.Fatalf("persisted snapshot=%+v, error=%v", persisted, err)
			}
			if tc.sdk {
				if persisted.Props[0].Flags != 0xfedcba98 || persisted.Props[0].LightmapResX != 123 || persisted.Props[0].LightmapResY != 456 || persisted.Props[0].BasePadding != 0x6e {
					t.Fatalf("edited SDK consumer fields=%+v", persisted.Props[0])
				}
			} else if persisted.Props[0].Flags != 0xe7 || persisted.Props[0].ExtendedFlags != 0xf2345678 || persisted.Props[0].DiffuseModulation != [4]byte{7, 19, 137, 249} {
				t.Fatalf("edited common consumer fields=%+v", persisted.Props[0])
			}
			if !tc.sdk && tc.version == 11 && persisted.Props[0].UniformScale != 3.5 {
				t.Fatalf("edited scale=%v", persisted.Props[0].UniformScale)
			}
		})
	}
}

func TestStaticPropsEmptyV10KeepsExplicitDialect(t *testing.T) {
	for _, tc := range []struct {
		profile bsp.Profile
		layout  lumps.StaticPropLayout
	}{
		{bsp.ProfileCSGO, lumps.StaticPropCommon10},
		{bsp.ProfileSDK2013, lumps.StaticPropSDK2013V10},
	} {
		t.Run(string(tc.profile), func(t *testing.T) {
			payload := make([]byte, 12)
			file := readStaticEditingFixture(t, tc.profile, staticEditingEntries(t, 10, payload))
			value, err := file.StaticProps()
			if err != nil || value.Layout != tc.layout || len(value.Models) != 0 || len(value.LeafIndices) != 0 || len(value.Props) != 0 {
				t.Fatalf("empty v10 dialect=%+v, error=%v", value, err)
			}
			if err := file.SetStaticProps(value); err != nil {
				t.Fatal(err)
			}
			requireStaticEditingPayload(t, rereadOrdinary(t, file), payload)
		})
	}
}

func TestStaticPropNamesByteLimitAndAtomicFailures(t *testing.T) {
	var slot lumps.StaticPropModel
	copy(slot[:], "old")
	slot[127] = 0xa5
	for _, invalid := range []string{strings.Repeat("x", 128), strings.Repeat("é", 64), "bad\x00name"} {
		before := slot
		if err := slot.SetName(invalid); err == nil || slot != before {
			t.Fatalf("invalid name accepted or changed slot: %q, error=%v", invalid, err)
		}
	}
	for _, valid := range []string{"", strings.Repeat("x", 127), strings.Repeat("é", 63)} {
		if err := slot.SetName(valid); err != nil {
			t.Fatal(err)
		}
		name, err := slot.Name()
		if err != nil || name != valid || slot[len(valid)] != 0 {
			t.Fatalf("name=%q, want=%q, error=%v", name, valid, err)
		}
		for _, b := range slot[len(valid):] {
			if b != 0 {
				t.Fatal("edited name retained stale padding")
			}
		}
		file := readStaticEditingFixture(t, bsp.ProfileCSGO, staticEditingEntries(t, 10, staticEditingPayload(false, 10)))
		value, err := file.StaticProps()
		if err != nil {
			t.Fatal(err)
		}
		value.Models[0] = slot
		if err := file.SetStaticProps(value); err != nil {
			t.Fatal(err)
		}
		value, err = rereadOrdinary(t, file).StaticProps()
		if err != nil || value.Models[0] != slot {
			t.Fatalf("boundary-length name did not persist: %+v, error=%v", value.Models, err)
		}
	}
	for i := range slot {
		slot[i] = 'x'
	}
	if _, err := slot.Name(); err == nil {
		t.Fatal("unterminated slot accepted")
	}
}

func requireStaticSetterAtomic(t *testing.T, file *bsp.File, value lumps.StaticProps) {
	t.Helper()
	before, err := file.RawLump(35)
	if err != nil {
		t.Fatal(err)
	}
	descriptor, err := file.Descriptor(35)
	if err != nil {
		t.Fatal(err)
	}
	entries, err := file.GameLumps()
	if err != nil {
		t.Fatal(err)
	}
	if err := file.SetStaticProps(value); err == nil {
		t.Fatal("invalid static-prop snapshot accepted")
	}
	after, err := file.RawLump(35)
	if err != nil || !bytes.Equal(after, before) {
		t.Fatalf("failed setter changed directory bytes: %v", err)
	}
	afterDescriptor, err := file.Descriptor(35)
	if err != nil || afterDescriptor != descriptor {
		t.Fatalf("failed setter changed descriptor: %+v, error=%v", afterDescriptor, err)
	}
	afterEntries, err := file.GameLumps()
	if err != nil || !reflect.DeepEqual(afterEntries, entries) {
		t.Fatalf("failed setter changed entries: %+v, error=%v", afterEntries, err)
	}
}

func TestStaticPropsReferencesAndSetterAtomicity(t *testing.T) {
	for _, sdk := range []bool{false, true} {
		profile := bsp.ProfileCSGO
		if sdk {
			profile = bsp.ProfileSDK2013
		}
		file := readStaticEditingFixture(t, profile, staticEditingEntries(t, 10, staticEditingPayload(sdk, 10)))
		mutations := []struct {
			name   string
			change func(*lumps.StaticProps)
		}{
			{"model exclusive bound", func(v *lumps.StaticProps) { v.Props[0].ModelIndex = 2 }},
			{"model removed", func(v *lumps.StaticProps) { v.Models = v.Models[:1] }},
			{"leaf ID exclusive bound", func(v *lumps.StaticProps) { v.LeafIndices[0] = 3 }},
			{"unused leaf ID still validated", func(v *lumps.StaticProps) { v.LeafIndices = append(v.LeafIndices, 3) }},
			{"leaf span exclusive bound", func(v *lumps.StaticProps) { v.Props[0].FirstLeaf = 3; v.Props[0].LeafCount = 1 }},
			{"leaf span overflow", func(v *lumps.StaticProps) { v.Props[0].FirstLeaf = 65535; v.Props[0].LeafCount = 2 }},
			{"empty span beyond end", func(v *lumps.StaticProps) { v.Props[0].FirstLeaf = 4; v.Props[0].LeafCount = 0 }},
			{"leaf list shortened", func(v *lumps.StaticProps) { v.LeafIndices = v.LeafIndices[:2] }},
			{"wrong layout", func(v *lumps.StaticProps) { v.Layout = lumps.StaticPropCommon9 }},
			{"unterminated model", func(v *lumps.StaticProps) {
				for i := range v.Models[0] {
					v.Models[0][i] = 'x'
				}
			}},
			{"absent scale negative zero", func(v *lumps.StaticProps) { v.Props[0].UniformScale = math.Float32frombits(0x80000000) }},
		}
		if sdk {
			mutations = append(mutations, struct {
				name   string
				change func(*lumps.StaticProps)
			}{"absent RGBA", func(v *lumps.StaticProps) { v.Props[0].DiffuseModulation[3] = 1 }})
		} else {
			mutations = append(mutations,
				struct {
					name   string
					change func(*lumps.StaticProps)
				}{"wide common flags", func(v *lumps.StaticProps) { v.Props[0].Flags = 256 }},
				struct {
					name   string
					change func(*lumps.StaticProps)
				}{"absent SDK lightmap", func(v *lumps.StaticProps) { v.Props[0].LightmapResY = 1 }})
		}
		for _, mutation := range mutations {
			t.Run(string(profile)+"/"+mutation.name, func(t *testing.T) {
				value, err := file.StaticProps()
				if err != nil {
					t.Fatal(err)
				}
				mutation.change(&value)
				requireStaticSetterAtomic(t, file, value)
			})
		}
		// Valid coordinated dictionary/reference edits derive all three counts
		// from the new slices and accept a zero-length span exactly at the end.
		value, err := file.StaticProps()
		if err != nil {
			t.Fatal(err)
		}
		var model lumps.StaticPropModel
		if err := model.SetName("models/third.mdl"); err != nil {
			t.Fatal(err)
		}
		value.Models = append(value.Models, model)
		value.LeafIndices = append(value.LeafIndices, 2)
		prop := value.Props[1]
		prop.ModelIndex, prop.FirstLeaf, prop.LeafCount = 2, 4, 0
		value.Props = append(value.Props, prop)
		if err := file.SetStaticProps(value); err != nil {
			t.Fatal(err)
		}
		file = rereadOrdinary(t, file)
		got, err := file.StaticProps()
		if err != nil || !reflect.DeepEqual(got, value) {
			t.Fatalf("coordinated edit=%+v, want=%+v, error=%v", got, value, err)
		}
		raw, err := file.DecodedGameLump(1)
		if err != nil || binary.LittleEndian.Uint32(raw) != 3 || binary.LittleEndian.Uint32(raw[388:]) != 4 || binary.LittleEndian.Uint32(raw[400:]) != 3 {
			t.Fatalf("derived counts=%x, error=%v", raw, err)
		}
	}
}

func TestStaticPropsPublicRejectsMalformedCountsNamesAndReferences(t *testing.T) {
	original := staticEditingPayload(false, 10)
	cases := []struct {
		name   string
		change func([]byte) []byte
	}{
		{"negative models", func(b []byte) []byte { binary.LittleEndian.PutUint32(b, 0xffffffff); return b }},
		{"model byte budget", func(b []byte) []byte { binary.LittleEndian.PutUint32(b, 0x7fffffff); return b }},
		{"negative leaves", func(b []byte) []byte { binary.LittleEndian.PutUint32(b[260:], 0xffffffff); return b }},
		{"leaf byte budget", func(b []byte) []byte { binary.LittleEndian.PutUint32(b[260:], 0x7fffffff); return b }},
		{"negative props", func(b []byte) []byte { binary.LittleEndian.PutUint32(b[270:], 0xffffffff); return b }},
		{"prop byte budget", func(b []byte) []byte { binary.LittleEndian.PutUint32(b[270:], 0x7fffffff); return b }},
		{"count smaller than payload", func(b []byte) []byte { binary.LittleEndian.PutUint32(b[270:], 1); return b }},
		{"truncated record", func(b []byte) []byte { return b[:len(b)-1] }},
		{"trailing record bytes", func(b []byte) []byte { return append(b, 0x77) }},
		{"unterminated slot", func(b []byte) []byte {
			for i := 4; i < 132; i++ {
				b[i] = 'x'
			}
			return b
		}},
		{"model reference", func(b []byte) []byte { binary.LittleEndian.PutUint16(b[274+24:], 2); return b }},
		{"leaf reference", func(b []byte) []byte { binary.LittleEndian.PutUint16(b[264:], 3); return b }},
		{"leaf span", func(b []byte) []byte { binary.LittleEndian.PutUint16(b[274+28:], 4); return b }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload := tc.change(append([]byte(nil), original...))
			file := readStaticEditingFixture(t, bsp.ProfileCSGO, staticEditingEntries(t, 10, payload))
			if _, err := file.StaticProps(); err == nil {
				t.Fatal("malformed public sprp accepted")
			}
			file = rereadOrdinary(t, file)
			requireStaticEditingPayload(t, file, payload)
		})
	}
}

func TestStaticPropsUnknownMissingAndDuplicateIDs(t *testing.T) {
	for _, tc := range []struct {
		name    string
		entries []editingGameWire
		want    error
	}{
		{"unknown99", staticEditingEntries(t, 99, []byte{0xff, 0, 0x99, 0x80, 0x7f}), bsp.ErrUnsupportedLump},
		{"missing", []editingGameWire{{0xdeadbeef, 0x8000, 10, 0, []byte{1, 2, 3}}}, bsp.ErrLumpNotFound},
		{"duplicate", []editingGameWire{{editingSprpID, 0x20, 10, 0, staticEditingPayload(false, 10)}, {editingSprpID, 0x40, 10, 0, staticEditingPayload(false, 10)}}, bsp.ErrAmbiguousGameLumpID},
	} {
		t.Run(tc.name, func(t *testing.T) {
			file := readStaticEditingFixture(t, bsp.ProfileCSGO, tc.entries)
			before, err := file.RawLump(35)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := file.StaticProps(); !errors.Is(err, tc.want) {
				t.Fatalf("getter error=%v, want=%v", err, tc.want)
			}
			if err := file.SetStaticProps(lumps.StaticProps{Layout: lumps.StaticPropCommon10}); !errors.Is(err, tc.want) {
				t.Fatalf("setter error=%v, want=%v", err, tc.want)
			}
			after, err := file.RawLump(35)
			if err != nil || !bytes.Equal(after, before) {
				t.Fatalf("failed access changed opaque raw bytes: %v", err)
			}
			file = rereadOrdinary(t, file)
			got, err := file.GameLumps()
			if err != nil || len(got) != len(tc.entries)+1 {
				t.Fatalf("reread directory=%+v, error=%v", got, err)
			}
			for i, want := range tc.entries {
				wantStored := want.data
				if want.flags&1 != 0 {
					wantStored = append(append([]byte(nil), want.data...), []byte("TGAP")...)
				}
				if got[i].ID != want.id || got[i].Flags != want.flags || got[i].Version != want.version || !bytes.Equal(got[i].StoredData, wantStored) {
					t.Fatalf("opaque entry %d changed: %+v", i, got[i])
				}
			}
		})
	}
}

func TestStaticPropsNativeCompressedPublicCommit(t *testing.T) {
	// Synthetic compression-path regression, not evidence of an engine dialect.
	// Independent CPython/liblzma raw LZMA1, dictionary4096, lc3/lp0/pb2,
	// EOS, framed as Source LZMA. Plaintext is twelve zero bytes: three counts.
	frame, err := hex.DecodeString("4c5a4d410c0000000c0000005d0010000000006c5f4e7bfffff0004000")
	if err != nil {
		t.Fatal(err)
	}
	entries := staticEditingEntries(t, 10, frame)
	entries[1].flags, entries[1].logical = 0x8001, 12
	file := readStaticEditingFixture(t, bsp.ProfileCSGO, entries)
	value, err := file.StaticProps()
	if err != nil || value.Layout != lumps.StaticPropCommon10 || len(value.Models) != 0 || len(value.LeafIndices) != 0 || len(value.Props) != 0 {
		t.Fatalf("compressed static props=%+v, error=%v", value, err)
	}
	before, err := file.GameLump(1)
	if err != nil || before.Flags != 0x8001 || !bytes.HasPrefix(before.StoredData, frame) {
		t.Fatalf("typed read changed native compressed storage: %+v, error=%v", before, err)
	}
	invalid := value
	invalid.Layout = lumps.StaticPropSDK2013V10
	requireStaticSetterAtomic(t, file, invalid)
	file = rereadOrdinary(t, file)
	untouched, err := file.GameLump(1)
	if err != nil || untouched.Flags != before.Flags || !bytes.Equal(untouched.StoredData, before.StoredData) {
		t.Fatalf("failed compressed setter or export changed stored frame: %+v, error=%v", untouched, err)
	}
	var model lumps.StaticPropModel
	if err := model.SetName("models/new.mdl"); err != nil {
		t.Fatal(err)
	}
	value.Models = []lumps.StaticPropModel{model}
	value.LeafIndices = []uint16{2}
	value.Props = []lumps.StaticProp{{ModelIndex: 0, FirstLeaf: 0, LeafCount: 1, Flags: 0x80, ExtendedFlags: 0x80000000, DiffuseModulation: [4]byte{1, 2, 3, 255}}}
	if err := file.SetStaticProps(value); err != nil {
		t.Fatal(err)
	}
	entries[1].flags = 0x8000
	file = rereadOrdinary(t, file)
	got, err := file.StaticProps()
	if err != nil || !reflect.DeepEqual(got, value) {
		t.Fatalf("compressed edit persistence=%+v, error=%v", got, err)
	}
	requireStaticEditingSiblings(t, file, entries)
	entry, err := file.GameLump(1)
	if err != nil || entry.Flags != 0x8000 || entry.Version != 10 || entry.LogicalLength != int32(4+128+4+2+4+76) || bytes.HasPrefix(entry.StoredData, []byte("LZMA")) {
		t.Fatalf("edited native descriptor=%+v, error=%v", entry, err)
	}
}

func TestBundledMapPublicStaticPropsAlignment(t *testing.T) {
	// Read only: the user's bundled fixture must never be rewritten by a test.
	input, err := os.Open("ar_baggage.bsp.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	stream, err := gzip.NewReader(input)
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()
	file, err := bsp.Read(stream, bsp.ReadOptions{Profile: bsp.ProfileCSGO})
	if err != nil {
		t.Fatal(err)
	}
	value, err := file.StaticProps()
	if err != nil {
		t.Fatal(err)
	}
	if value.Layout != lumps.StaticPropCommon10 || len(value.Props) != 488 || len(value.Models) != 89 || len(value.LeafIndices) != 1128 {
		t.Fatalf("bundled layout=%d props=%d models=%d leaves=%d", value.Layout, len(value.Props), len(value.Models), len(value.LeafIndices))
	}
	// Known independent consumer values. Compare the actual float32 bits, not a
	// tolerance that could hide reading the wrong field or record width.
	// Direct fixture inspection records Z bits 0x438fa873, represented exactly
	// by float32(287.316); X/Y bits are c4040000 and 438e0000/c38e0000.
	want := [2][3]float32{{-528, 284, 287.316}, {-528, -284, 287.316}}
	for i := range want {
		for axis := range want[i] {
			if math.Float32bits(value.Props[i].Origin[axis]) != math.Float32bits(want[i][axis]) {
				t.Fatalf("prop %d origin axis %d=%v (%08x), want %v (%08x)", i, axis, value.Props[i].Origin[axis], math.Float32bits(value.Props[i].Origin[axis]), want[i][axis], math.Float32bits(want[i][axis]))
			}
		}
	}
}

// Externally retained authored maps avoid redistributing packed third-party
// assets. SHA256 and expected values are independently corroborated by srctools
// cd7e841; see task08 for source/build provenance and the evidence boundary.
func TestExternalGenuineCSGOStaticPropsV11(t *testing.T) {
	directory := os.Getenv("BSP_CSGO_V11_FIXTURE_DIR")
	if directory == "" {
		t.Skip("set BSP_CSGO_V11_FIXTURE_DIR to the pinned authored fixtures in task08")
	}
	for _, tc := range []struct {
		name, hash     string
		models, leaves int
		props          int
		model          string
		origin         [3]float32
		scale          float32
	}{
		{"holo_city.bsp", "c5c2227c8982b59be185e652382d74f24e4d8ccf6a2b41709b053f3b89bd5319", 1, 1, 1, "models/combine_room/combine_monitor003a.mdl", [3]float32{648, -394, -1963}, 0.1},
		{"ar_campusx3.bsp", "5871bedb53c0ac3971dd0fd1626648c08f7266208ca3b60aa5eff53287230c17", 42, 104, 94, "models/props/de_nuke/hr_nuke/nuke_office_desk/nuke_conference_table.mdl", [3]float32{-750, 1551, 200}, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			raw, err := os.ReadFile(filepath.Join(directory, tc.name))
			if err != nil {
				t.Fatal(err)
			}
			if got := fmt.Sprintf("%x", sha256.Sum256(raw)); got != tc.hash {
				t.Fatalf("fixture hash=%s, want %s", got, tc.hash)
			}
			file, err := bsp.Read(bytes.NewReader(raw), bsp.ReadOptions{Profile: bsp.ProfileCSGO})
			if err != nil {
				t.Fatal(err)
			}
			value, err := file.StaticProps()
			if err != nil {
				t.Fatal(err)
			}
			if value.Layout != lumps.StaticPropCommon11 || len(value.Models) != tc.models || len(value.LeafIndices) != tc.leaves || len(value.Props) != tc.props {
				t.Fatalf("genuine v11 framing=%d/%d/%d/%d", value.Layout, len(value.Models), len(value.LeafIndices), len(value.Props))
			}
			prop := &value.Props[0]
			name, err := value.Models[prop.ModelIndex].Name()
			if err != nil || name != tc.model || prop.Origin != tc.origin || prop.UniformScale != tc.scale || prop.ExtendedFlags != 0 {
				t.Fatalf("genuine first prop: model=%q, record=%+v, error=%v", name, prop, err)
			}
			prop.UniformScale, prop.ExtendedFlags = 2.5, 0x80000010
			if tc.props > 1 {
				second := &value.Props[1]
				if second.Origin != [3]float32{-697, 1582, 200} || second.UniformScale != 1 || second.ExtendedFlags != 0 {
					t.Fatalf("genuine second-record alignment=%+v", second)
				}
				second.UniformScale, second.ExtendedFlags = 0.75, 0x40000020
			}
			if err := file.SetStaticProps(value); err != nil {
				t.Fatal(err)
			}
			after, err := rereadOrdinary(t, file).StaticProps()
			if err != nil {
				t.Fatal(err)
			}
			if after.Props[0].UniformScale != 2.5 || after.Props[0].ExtendedFlags != 0x80000010 {
				t.Fatalf("first v11 edited tail=%+v", after.Props[0])
			}
			if tc.props > 1 && (after.Props[1].UniformScale != 0.75 || after.Props[1].ExtendedFlags != 0x40000020) {
				t.Fatalf("second v11 edited tail=%+v", after.Props[1])
			}
		})
	}
}
