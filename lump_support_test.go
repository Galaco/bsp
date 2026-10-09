package bsp_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
	"testing"

	"github.com/galaco/bsp"
	"github.com/galaco/bsp/lumps"
)

func supportFixture(t *testing.T, profile bsp.Profile, version int32) *bsp.File {
	t.Helper()
	raw := make([]byte, bsp.HeaderSize)
	copy(raw, "VBSP")
	binary.LittleEndian.PutUint32(raw[4:], uint32(version))
	file, err := bsp.Read(bytes.NewReader(raw), bsp.ReadOptions{Profile: profile})
	if err != nil {
		t.Fatal(err)
	}
	return file
}

func TestSupportRawCapabilitiesEveryNumericID(t *testing.T) {
	for _, tc := range []struct {
		profile bsp.Profile
		version int32
	}{
		{bsp.ProfileSDK2013, 19}, {bsp.ProfileSDK2013, 20},
		{bsp.ProfileCSGO, 21}, {bsp.ProfileL4D2, 21}, {bsp.ProfileReactiveDrop, 21},
	} {
		file := supportFixture(t, tc.profile, tc.version)
		for id := bsp.LumpID(0); id < bsp.LumpCount; id++ {
			// Deliberately opaque unknown revisions and non-record payloads:
			// metadata inspection must not decode them or deny raw replacement.
			want := []byte{byte(id), 0x80, 0, byte(tc.version), 0xff}
			desc := bsp.HeaderLump{Version: 9876, Offset: 4096, FourCC: [4]byte{1, 2, 3, 4}}
			if err := file.SetRawLump(id, desc, want); err != nil {
				t.Fatal(err)
			}
			info, err := file.Support(id)
			if err != nil || info.ID != id || info.Name == "" || !info.RawReadable || !info.RawEditable || !info.Preserved {
				t.Fatalf("%s ID%d raw capabilities=%+v error=%v", tc.profile, id, info, err)
			}
			if info.Status == bsp.SupportTyped || info.Status == bsp.SupportByteStream || info.Readable || info.Editable || info.RecordWidth != 0 {
				t.Fatalf("%s ID%d unknown revision advertised decoded support: %+v", tc.profile, id, info)
			}
			got, err := file.RawLump(id)
			if err != nil || !bytes.Equal(got, want) {
				t.Fatalf("%s ID%d raw=%x error=%v", tc.profile, id, got, err)
			}
			gotDesc, err := file.Descriptor(id)
			desc.Length = int32(len(want))
			if err != nil || gotDesc != desc {
				t.Fatalf("%s ID%d descriptor=%+v want=%+v error=%v", tc.profile, id, gotDesc, desc, err)
			}
		}
	}
}

func TestSupportInvalidIDs(t *testing.T) {
	file := supportFixture(t, bsp.ProfileSDK2013, 20)
	for _, id := range []bsp.LumpID{-1, 64, 1 << 30} {
		if _, err := file.Support(id); !errors.Is(err, bsp.ErrInvalidLumpID) {
			t.Fatalf("ID%d error=%v", id, err)
		}
	}
}

func TestSupportProfileAliasesAndGates(t *testing.T) {
	for _, profile := range []bsp.Profile{bsp.ProfileSDK2013, bsp.ProfileCSGO, bsp.ProfileL4D2, bsp.ProfileReactiveDrop} {
		version := int32(21)
		if profile == bsp.ProfileSDK2013 {
			version = 20
		}
		file := supportFixture(t, profile, version)
		for _, tc := range []struct {
			id         bsp.LumpID
			legacy, rd string
			width      int
		}{
			{22, "portals", "propcollision", 8}, {23, "clusters", "prophulls", 16},
			{24, "portalverts", "prophullverts", 12}, {25, "clusterportals", "proptris", 8},
			{49, "physcollidesurface", "propblob", 0},
		} {
			info, err := file.Support(tc.id)
			if err != nil {
				t.Fatal(err)
			}
			if profile == bsp.ProfileReactiveDrop {
				status := bsp.SupportTyped
				if tc.id == 49 {
					status = bsp.SupportByteStream
				}
				if info.Name != tc.rd || info.Status != status || info.RecordWidth != tc.width || !info.Readable || !info.Editable || len(info.Aliases) == 0 || info.Aliases[0] != tc.legacy {
					t.Fatalf("RD ID%d=%+v", tc.id, info)
				}
			} else if info.Name != tc.legacy || info.Status != bsp.SupportUnsupportedProfile || info.Readable || info.Editable || info.RecordWidth != 0 {
				t.Fatalf("%s ID%d=%+v", profile, tc.id, info)
			}
		}
		for _, id := range []bsp.LumpID{50, 61, 63} {
			info, err := file.Support(id)
			if err != nil {
				t.Fatal(err)
			}
			allowed := id == 50 || (id == 61 && profile != bsp.ProfileSDK2013) || (id == 63 && (profile == bsp.ProfileCSGO || profile == bsp.ProfileReactiveDrop))
			if allowed != (info.Status == bsp.SupportTyped && info.Readable && info.Editable) {
				t.Fatalf("%s ID%d allowed=%v support=%+v", profile, id, allowed, info)
			}
		}
	}
}

func TestSupportRevisionWidthsAndConsumerLimits(t *testing.T) {
	for _, tc := range []struct {
		profile  bsp.Profile
		id       bsp.LumpID
		revision int32
		status   bsp.LumpSupportStatus
		width    int
	}{
		{bsp.ProfileCSGO, 10, 0, bsp.SupportTyped, 56},
		{bsp.ProfileCSGO, 10, 1, bsp.SupportTyped, 32},
		{bsp.ProfileCSGO, 15, 1, bsp.SupportTyped, 100},
		{bsp.ProfileReactiveDrop, 54, 1, bsp.SupportTyped, 100},
		{bsp.ProfileL4D2, 15, 1, bsp.SupportOpaque, 0},
		{bsp.ProfileReactiveDrop, 9, 1, bsp.SupportOpaque, 0},
		{bsp.ProfileReactiveDrop, 9, 2, bsp.SupportTyped, 0},
		{bsp.ProfileCSGO, 4, 0, bsp.SupportTyped, 0},
		{bsp.ProfileCSGO, 40, 0, bsp.SupportByteStream, 0},
		{bsp.ProfileCSGO, 59, 0, bsp.SupportTyped, 4},
		{bsp.ProfileCSGO, 32, 0, bsp.SupportOpaque, 0},
		{bsp.ProfileCSGO, 28, 0, bsp.SupportOpaque, 0},
	} {
		file := supportFixture(t, tc.profile, 21)
		if err := file.SetRawLump(tc.id, bsp.HeaderLump{Version: tc.revision}, nil); err != nil {
			t.Fatal(err)
		}
		info, err := file.Support(tc.id)
		if err != nil || info.Status != tc.status || info.RecordWidth != tc.width || info.Limitations == "" || info.Evidence == "" {
			t.Fatalf("%s ID%d revision%d=%+v error=%v", tc.profile, tc.id, tc.revision, info, err)
		}
	}
}

func TestSupportGameDirectoryStaticPropSelector(t *testing.T) {
	file := supportFixture(t, bsp.ProfileReactiveDrop, 21)
	// A counted empty sprp payload proves selection, not a genuine fixture.
	// Genuine shipped v9 fixture consumption is covered by root integration.
	for _, revision := range []uint16{9, 99} {
		if err := file.SetGameLumps([]lumps.GameLump{{ID: lumps.StaticPropID, Version: revision, StoredData: make([]byte, 12)}}); err != nil {
			t.Fatal(err)
		}
		info, err := file.Support(35)
		if err != nil || info.Status != bsp.SupportByteStream || info.RecordWidth != 0 || !info.Readable || !info.Editable {
			t.Fatalf("directory=%+v error=%v", info, err)
		}
		want := bsp.SupportOpaque
		if revision == 9 {
			want = bsp.SupportTyped
		}
		if info.StaticPropStatus != want {
			t.Fatalf("sprp revision%d=%+v", revision, info)
		}
		found := false
		for _, version := range info.StaticPropVersions {
			found = found || version == 9
		}
		if !found {
			t.Fatalf("RD accepted sprp revisions=%v", info.StaticPropVersions)
		}
	}
	if err := file.SetGameLumps([]lumps.GameLump{{ID: lumps.StaticPropID, Version: 9, StoredData: make([]byte, 12)}, {ID: lumps.StaticPropID, Version: 9, StoredData: make([]byte, 12)}}); err != nil {
		t.Fatal(err)
	}
	info, err := file.Support(35)
	if err != nil || info.StaticPropStatus != bsp.SupportOpaque || info.Status != bsp.SupportByteStream {
		t.Fatalf("ambiguous sprp=%+v error=%v", info, err)
	}
}

func TestSupportMetadataSlicesAreIndependent(t *testing.T) {
	file := supportFixture(t, bsp.ProfileReactiveDrop, 21)
	info, err := file.Support(22)
	if err != nil {
		t.Fatal(err)
	}
	info.Aliases[0] = "changed"
	info.Versions[0] = 99
	info.Accessors[0] = "changed"
	again, err := file.Support(22)
	if err != nil || !reflect.DeepEqual(again.Aliases, []string{"portals", "unused0"}) || !reflect.DeepEqual(again.Versions, []int32{0}) || again.Accessors[0] != "PropCollisions" {
		t.Fatalf("mutated global metadata=%+v error=%v", again, err)
	}
}
