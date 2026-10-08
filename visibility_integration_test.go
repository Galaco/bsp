package bsp_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/galaco/bsp"
	"github.com/galaco/bsp/lumps"
)

func publicVisibilityWire(pairs [][2]int32, data []byte) []byte {
	raw := make([]byte, 4+8*len(pairs)+len(data))
	binary.LittleEndian.PutUint32(raw, uint32(len(pairs)))
	for cluster, pair := range pairs {
		for kind, offset := range pair {
			binary.LittleEndian.PutUint32(raw[4+8*cluster+4*kind:], uint32(offset))
		}
	}
	copy(raw[4+8*len(pairs):], data)
	return raw
}

func requireVisibilityIDs(t *testing.T, file *bsp.File, cluster int, kind lumps.VisibilityKind, expected []int, wantPresent bool) {
	t.Helper()
	value, err := file.Visibility()
	if err != nil {
		t.Fatal(err)
	}
	ids, present, err := value.VisibleClusters(cluster, kind)
	if err != nil || present != wantPresent || len(ids) != len(expected) || len(ids) != 0 && !reflect.DeepEqual(ids, expected) {
		t.Fatalf("cluster %d kind %d: IDs=%v presence=%v error=%v want=%v/%v", cluster, kind, ids, present, err, expected, wantPresent)
	}
}

func TestVisibilityPublicEditsExportAndOwnership(t *testing.T) {
	for _, profile := range []bsp.Profile{bsp.ProfileSDK2013, bsp.ProfileCSGO, bsp.ProfileL4D2, bsp.ProfileReactiveDrop} {
		t.Run(string(profile), func(t *testing.T) {
			pairs := make([][2]int32, 9)
			for i := range pairs {
				pairs[i] = [2]int32{-1, -1}
			}
			pairs[0] = [2]int32{76, 78}
			pairs[1] = pairs[0]
			data := []byte{1, 0x81, 0x80, 0x80}
			original := publicVisibilityWire(pairs, data)
			file := readOrdinaryFixture(t, profile, map[int]ordinaryWire{4: {0, original}})
			requireVisibilityIDs(t, file, 0, lumps.PVS, []int{0, 8}, true)
			requireVisibilityIDs(t, file, 1, lumps.PAS, []int{7}, true)
			requireVisibilityIDs(t, file, 8, lumps.PVS, nil, false)
			value, err := file.Visibility()
			if err != nil {
				t.Fatal(err)
			}
			// Change count, offsets, and stream; preserve shared offset semantics.
			value.Offsets = value.Offsets[:8]
			for i := range value.Offsets {
				value.Offsets[i] = [2]int32{-1, -1}
			}
			value.Offsets[0], value.Offsets[1] = [2]int32{68, 69}, [2]int32{68, -1}
			value.Data = []byte{0x82, 4}
			stored, err := file.RawLump(bsp.LumpVisibility)
			if err != nil || !bytes.Equal(stored, original) {
				t.Fatalf("snapshot edits changed bytes: %v", err)
			}
			if err := file.SetVisibility(value); err != nil {
				t.Fatal(err)
			}
			want := publicVisibilityWire(value.Offsets, value.Data)
			value.Offsets[0][0], value.Data[0] = -1, 0
			stored, err = file.RawLump(bsp.LumpVisibility)
			if err != nil || !bytes.Equal(stored, want) {
				t.Fatalf("setter retained caller slices or emitted wrong count/table/data: %v", err)
			}
			file = rereadOrdinary(t, file)
			requireVisibilityIDs(t, file, 0, lumps.PVS, []int{1, 7}, true)
			requireVisibilityIDs(t, file, 0, lumps.PAS, []int{2}, true)
			requireVisibilityIDs(t, file, 1, lumps.PVS, []int{1, 7}, true)
			requireVisibilityIDs(t, file, 1, lumps.PAS, nil, false)
			stored, err = file.RawLump(bsp.LumpVisibility)
			if err != nil || !bytes.Equal(stored, want) {
				t.Fatalf("export/reimport changed edited wire: %v", err)
			}
		})
	}
}

func TestVisibilityAbsentAndZeroCount(t *testing.T) {
	file := readOrdinaryFixture(t, bsp.ProfileSDK2013, nil)
	if _, err := file.Visibility(); !errors.Is(err, bsp.ErrLumpNotFound) {
		t.Fatalf("absence error=%v", err)
	}
	file = rereadOrdinary(t, file)
	raw, err := file.RawLump(bsp.LumpVisibility)
	if err != nil || len(raw) != 0 {
		t.Fatalf("absent export synthesized bytes: %x %v", raw, err)
	}
	if err := file.SetVisibility(lumps.Visibility{}); err != nil {
		t.Fatal(err)
	}
	file = rereadOrdinary(t, file)
	raw, err = file.RawLump(bsp.LumpVisibility)
	if err != nil || !bytes.Equal(raw, []byte{0, 0, 0, 0}) {
		t.Fatalf("zero-count wire=%x error=%v", raw, err)
	}
	value, err := file.Visibility()
	if err != nil || len(value.Offsets) != 0 || len(value.Data) != 0 {
		t.Fatalf("present zero-count value=%+v error=%v", value, err)
	}
}

func TestVisibilityFailedSettersPreserveBytesAndDescriptor(t *testing.T) {
	file := readOrdinaryFixture(t, bsp.ProfileSDK2013, map[int]ordinaryWire{4: {0, publicVisibilityWire([][2]int32{{12, 12}}, []byte{1})}})
	// Compression metadata is intentionally opaque. Failed preflight must not
	// clear it or touch bytes; a successful typed replacement must clear it.
	descriptor, err := file.Descriptor(bsp.LumpVisibility)
	if err != nil {
		t.Fatal(err)
	}
	descriptor.FourCC = [4]byte{13, 0, 0, 0}
	before, err := file.RawLump(bsp.LumpVisibility)
	if err != nil {
		t.Fatal(err)
	}
	if err := file.SetRawLump(bsp.LumpVisibility, descriptor, before); err != nil {
		t.Fatal(err)
	}
	for _, value := range []lumps.Visibility{
		{Offsets: [][2]int32{{12, -2}}, Data: []byte{1}},
		{Offsets: [][2]int32{{12, 13}}, Data: []byte{1}},
		{Offsets: [][2]int32{{12, 12}}, Data: []byte{0, 0}},
		{Offsets: make([][2]int32, 65537)},
	} {
		err := file.SetVisibility(value)
		var formatError *bsp.FormatError
		if !errors.Is(err, bsp.ErrMalformed) || !errors.As(err, &formatError) {
			t.Fatalf("setter lacks malformed classification: %v", err)
		}
		after, err := file.RawLump(bsp.LumpVisibility)
		afterDescriptor, descriptorErr := file.Descriptor(bsp.LumpVisibility)
		if err != nil || descriptorErr != nil || !bytes.Equal(after, before) || afterDescriptor != descriptor {
			t.Fatal("failed setter changed bytes or descriptor")
		}
	}
	if err := file.SetVisibility(lumps.Visibility{Offsets: [][2]int32{{12, -1}}, Data: []byte{2}}); err != nil {
		t.Fatal(err)
	}
	afterDescriptor, err := file.Descriptor(bsp.LumpVisibility)
	if err != nil || afterDescriptor.Version != descriptor.Version || afterDescriptor.FourCC != [4]byte{} {
		t.Fatalf("successful setter metadata=%+v %v", afterDescriptor, err)
	}
	requireVisibilityIDs(t, file, 0, lumps.PVS, nil, true)
}

func TestVisibilityMalformedRawIsPreservedAndTypedErrorsClassified(t *testing.T) {
	for _, raw := range [][]byte{
		{255, 255, 255, 255},
		{255, 255, 255, 127},
		{1, 0, 0, 0},
		publicVisibilityWire([][2]int32{{12, -2}}, []byte{1}),
		publicVisibilityWire([][2]int32{{12, 12}}, []byte{0, 0}),
	} {
		file := readOrdinaryFixture(t, bsp.ProfileSDK2013, map[int]ordinaryWire{4: {0, raw}})
		_, err := file.Visibility()
		var formatError *bsp.FormatError
		if !errors.Is(err, bsp.ErrMalformed) || !errors.As(err, &formatError) {
			t.Fatalf("typed malformed error=%v", err)
		}
		file = rereadOrdinary(t, file)
		stored, err := file.RawLump(bsp.LumpVisibility)
		if err != nil || !bytes.Equal(stored, raw) {
			t.Fatalf("malformed raw was not preserved: %v", err)
		}
	}
}

func TestVisibilityUnsupportedRevisionIsAtomicAndOpaque(t *testing.T) {
	raw := publicVisibilityWire([][2]int32{{12, 12}}, []byte{1})
	file := readOrdinaryFixture(t, bsp.ProfileSDK2013, map[int]ordinaryWire{4: {1, raw}})
	descriptor, err := file.Descriptor(bsp.LumpVisibility)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.Visibility(); !errors.Is(err, bsp.ErrUnsupportedLump) {
		t.Fatalf("unsupported getter error=%v", err)
	}
	if err := file.SetVisibility(lumps.Visibility{}); !errors.Is(err, bsp.ErrUnsupportedLump) {
		t.Fatalf("unsupported setter error=%v", err)
	}
	afterDescriptor, err := file.Descriptor(bsp.LumpVisibility)
	if err != nil || afterDescriptor != descriptor {
		t.Fatalf("failed revision changed descriptor: %v", err)
	}
	file = rereadOrdinary(t, file)
	stored, err := file.RawLump(bsp.LumpVisibility)
	if err != nil || !bytes.Equal(stored, raw) {
		t.Fatalf("unsupported wire changed: %v", err)
	}
}

func TestVisibilityPublicIndependentCountsAndMaximumID(t *testing.T) {
	for _, count := range []int{0, 1, 8, 9, 65536} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			pairs := make([][2]int32, count)
			var data []byte
			var want []int
			switch count {
			case 1:
				data, want = []byte{0x81}, []int{0}
			case 8:
				data, want = []byte{0x81}, []int{0, 7}
			case 9:
				data, want = []byte{0x81, 0xff}, []int{0, 7, 8}
			case 65536:
				for remaining := 8191; remaining > 0; {
					run := remaining
					if run > 255 {
						run = 255
					}
					data = append(data, 0, byte(run))
					remaining -= run
				}
				data = append(data, 0x80)
				want = []int{65535}
			}
			for i := range pairs {
				pairs[i] = [2]int32{int32(4 + 8*count), int32(4 + 8*count)}
			}
			raw := publicVisibilityWire(pairs, data)
			file := readOrdinaryFixture(t, bsp.ProfileSDK2013, map[int]ordinaryWire{4: {0, raw}})
			value, err := file.Visibility()
			if err != nil || len(value.Offsets) != count {
				t.Fatalf("count=%d error=%v", len(value.Offsets), err)
			}
			if count > 0 {
				requireVisibilityIDs(t, file, count-1, lumps.PVS, want, true)
				requireVisibilityIDs(t, file, count-1, lumps.PAS, want, true)
			}
			if err := file.SetVisibility(value); err != nil {
				t.Fatal(err)
			}
			file = rereadOrdinary(t, file)
			stored, err := file.RawLump(bsp.LumpVisibility)
			if err != nil || !bytes.Equal(stored, raw) {
				t.Fatalf("count %d wire changed after setter/export/reread: %v", count, err)
			}
			if count > 0 {
				requireVisibilityIDs(t, file, count-1, lumps.PVS, want, true)
				requireVisibilityIDs(t, file, count-1, lumps.PAS, want, true)
			}
		})
	}
}
