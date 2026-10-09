package bsp_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"testing"

	"github.com/galaco/bsp"
	"github.com/galaco/bsp/lumps"
)

func mapFlagsProfiles() []struct {
	profile bsp.Profile
	version int32
} {
	return []struct {
		profile bsp.Profile
		version int32
	}{
		{bsp.ProfileSDK2013, 19},
		{bsp.ProfileSDK2013, 20},
		{bsp.ProfileCSGO, 21},
		{bsp.ProfileL4D2, 21},
		{bsp.ProfileReactiveDrop, 21},
	}
}

func readMapFlagsFixture(t *testing.T, profile bsp.Profile, version int32, wire ordinaryWire) *bsp.File {
	t.Helper()
	file := readOrdinaryFixture(t, profile, map[int]ordinaryWire{int(bsp.LumpMapFlags): wire})
	if file.Header().Version == version {
		return file
	}
	raw, err := file.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	binary.LittleEndian.PutUint32(raw[4:], uint32(version))
	file, err = bsp.Read(bytes.NewReader(raw), bsp.ReadOptions{Profile: profile})
	if err != nil {
		t.Fatal(err)
	}
	return file
}

func requireMapFlagsState(t *testing.T, file *bsp.File, value uint32, present bool, wantRaw []byte) {
	t.Helper()
	got, gotPresent, err := file.MapFlags()
	if err != nil || got.LevelFlags != value || gotPresent != present {
		t.Fatalf("flags=%+v present=%v error=%v; want=%08x present=%v", got, gotPresent, err, value, present)
	}
	raw, err := file.RawLump(bsp.LumpMapFlags)
	if err != nil || !bytes.Equal(raw, wantRaw) {
		t.Fatalf("exported flags=%x error=%v; want=%x", raw, err, wantRaw)
	}
	descriptor, err := file.Descriptor(bsp.LumpMapFlags)
	if err != nil || descriptor.Version != 0 || descriptor.Length != int32(len(wantRaw)) || descriptor.FourCC != [4]byte{} {
		t.Fatalf("flags descriptor=%+v error=%v", descriptor, err)
	}
}

func TestMapFlagsPresenceSnapshotAndRemovalAcrossProfiles(t *testing.T) {
	for _, profile := range mapFlagsProfiles() {
		t.Run(fmt.Sprintf("%s-%d", profile.profile, profile.version), func(t *testing.T) {
			file := readMapFlagsFixture(t, profile.profile, profile.version, ordinaryWire{})
			requireMapFlagsState(t, file, 0, false, nil)
			file = rereadOrdinary(t, file)
			requireMapFlagsState(t, file, 0, false, nil)
			if err := file.SetMapFlags(lumps.MapFlags{}); err != nil {
				t.Fatal(err)
			}
			requireMapFlagsState(t, file, 0, true, []byte{0, 0, 0, 0})
			file = rereadOrdinary(t, file)
			requireMapFlagsState(t, file, 0, true, []byte{0, 0, 0, 0})
			snapshot, present, err := file.MapFlags()
			if err != nil || !present {
				t.Fatalf("snapshot present=%v error=%v", present, err)
			}
			snapshot.LevelFlags = 0x81234567
			requireMapFlagsState(t, file, 0, true, []byte{0, 0, 0, 0})
			if err := file.SetMapFlags(snapshot); err != nil {
				t.Fatal(err)
			}
			snapshot.LevelFlags = 0
			requireMapFlagsState(t, file, 0x81234567, true, []byte{0x67, 0x45, 0x23, 0x81})
			file = rereadOrdinary(t, file)
			requireMapFlagsState(t, file, 0x81234567, true, []byte{0x67, 0x45, 0x23, 0x81})
			if err := file.RemoveMapFlags(); err != nil {
				t.Fatal(err)
			}
			requireMapFlagsState(t, file, 0, false, nil)
			file = rereadOrdinary(t, file)
			requireMapFlagsState(t, file, 0, false, nil)
			if err := file.RemoveMapFlags(); err != nil {
				t.Fatal(err)
			}
			requireMapFlagsState(t, file, 0, false, nil)
			if err := file.SetMapFlags(lumps.MapFlags{}); err != nil {
				t.Fatal(err)
			}
			file = rereadOrdinary(t, file)
			requireMapFlagsState(t, file, 0, true, []byte{0, 0, 0, 0})
		})
	}
}

func TestMapFlagsMalformedGetterPreservesRawAndDescriptor(t *testing.T) {
	for _, profile := range mapFlagsProfiles() {
		for _, length := range []int{1, 2, 3, 5, 6, 7, 8, 16} {
			t.Run(fmt.Sprintf("%s-%d-length%d", profile.profile, profile.version, length), func(t *testing.T) {
				raw := bytes.Repeat([]byte{0xa5}, length)
				file := readMapFlagsFixture(t, profile.profile, profile.version, ordinaryWire{data: raw})
				before, err := file.Descriptor(bsp.LumpMapFlags)
				if err != nil {
					t.Fatal(err)
				}
				for i := 0; i < 2; i++ {
					value, present, err := file.MapFlags()
					var formatError *bsp.FormatError
					if !errors.Is(err, bsp.ErrMalformed) || !errors.As(err, &formatError) || present || value != (lumps.MapFlags{}) {
						t.Fatalf("malformed flags=%+v present=%v error=%v", value, present, err)
					}
					stored, err := file.RawLump(bsp.LumpMapFlags)
					if err != nil || !bytes.Equal(stored, raw) {
						t.Fatalf("failed getter changed raw=%x error=%v", stored, err)
					}
					after, err := file.Descriptor(bsp.LumpMapFlags)
					if err != nil || after != before {
						t.Fatalf("failed getter changed descriptor=%+v error=%v", after, err)
					}
				}
				file = rereadOrdinary(t, file)
				if stored, err := file.RawLump(bsp.LumpMapFlags); err != nil || !bytes.Equal(stored, raw) {
					t.Fatalf("opaque malformed export changed raw=%x error=%v", stored, err)
				}
				if _, _, err := file.MapFlags(); !errors.Is(err, bsp.ErrMalformed) {
					t.Fatalf("reread malformed flags error=%v", err)
				}
				if err := file.SetMapFlags(lumps.MapFlags{}); err != nil {
					t.Fatal(err)
				}
				requireMapFlagsState(t, file, 0, true, []byte{0, 0, 0, 0})
			})
		}
	}
}

func TestMapFlagsUnsupportedSetterAndRemovalAreAtomic(t *testing.T) {
	for _, profile := range mapFlagsProfiles() {
		for _, revision := range []int32{-1, 1, 2} {
			file := readMapFlagsFixture(t, profile.profile, profile.version, ordinaryWire{})
			before := bsp.HeaderLump{Version: revision, FourCC: [4]byte{4, 0, 0, 0}}
			raw := []byte{0x67, 0x45, 0x23, 0x81}
			if err := file.SetRawLump(bsp.LumpMapFlags, before, raw); err != nil {
				t.Fatal(err)
			}
			before.Length = int32(len(raw))
			for _, change := range []func() error{
				func() error { return file.SetMapFlags(lumps.MapFlags{}) },
				func() error { return file.SetMapFlags(lumps.MapFlags{LevelFlags: 0xffffffff}) },
				file.RemoveMapFlags,
			} {
				if err := change(); !errors.Is(err, bsp.ErrUnsupportedLump) {
					t.Fatalf("%s-%d revision%d: failed edit error=%v", profile.profile, profile.version, revision, err)
				}
				stored, err := file.RawLump(bsp.LumpMapFlags)
				if err != nil || !bytes.Equal(stored, raw) {
					t.Fatalf("failed edit changed raw=%x error=%v", stored, err)
				}
				after, err := file.Descriptor(bsp.LumpMapFlags)
				if err != nil || after != before {
					t.Fatalf("failed edit changed descriptor=%+v error=%v", after, err)
				}
			}
		}
	}
}

func TestMapFlagsTypedReplacementClearsCompression(t *testing.T) {
	for _, profile := range mapFlagsProfiles() {
		for _, remove := range []bool{false, true} {
			file := readMapFlagsFixture(t, profile.profile, profile.version, ordinaryWire{})
			// A compressed opaque payload need not decode before typed replacement.
			if err := file.SetRawLump(bsp.LumpMapFlags, bsp.HeaderLump{FourCC: [4]byte{40, 0, 0, 0}}, compressedPlaneFixture(t)); err != nil {
				t.Fatal(err)
			}
			if remove {
				if err := file.RemoveMapFlags(); err != nil {
					t.Fatal(err)
				}
				requireMapFlagsState(t, file, 0, false, nil)
				file = rereadOrdinary(t, file)
				requireMapFlagsState(t, file, 0, false, nil)
			} else {
				if err := file.SetMapFlags(lumps.MapFlags{LevelFlags: 0x80000000}); err != nil {
					t.Fatal(err)
				}
				requireMapFlagsState(t, file, 0x80000000, true, []byte{0, 0, 0, 0x80})
				file = rereadOrdinary(t, file)
				requireMapFlagsState(t, file, 0x80000000, true, []byte{0, 0, 0, 0x80})
			}
		}
	}
}
