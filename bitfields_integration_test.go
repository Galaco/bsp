package bsp_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"github.com/galaco/bsp"
	facevalue "github.com/galaco/bsp/primitives/face"
)

func TestLeafPackedEditsPersistInBothDescriptorLayouts(t *testing.T) {
	for _, pc := range []struct {
		profile bsp.Profile
		version int32
	}{{bsp.ProfileSDK2013, 19}, {bsp.ProfileSDK2013, 20}, {bsp.ProfileCSGO, 21}, {bsp.ProfileL4D2, 21}, {bsp.ProfileReactiveDrop, 21}} {
		for _, revision := range []int32{0, 1} {
			stride := 56
			if revision == 1 {
				stride = 32
			}
			payload := make([]byte, 2*stride)
			payload[6], payload[7] = 0xff, 0xff
			payload[stride+6], payload[stride+7] = 0x80, 0x80
			for i := 0; i < 2; i++ {
				payload[i*stride+stride-2], payload[i*stride+stride-1] = byte(0xa1+i), 0xb2
			}
			file, err := bsp.ReadFromStream(bytes.NewReader(profileContainer(pc.profile, pc.version, nil)), bsp.ReadOptions{Profile: pc.profile})
			if err != nil {
				t.Fatal(err)
			}
			mustDesc := bsp.HeaderLump{Version: revision}
			if err := file.SetRawLump(bsp.LumpLeafs, mustDesc, payload); err != nil {
				t.Fatal(err)
			}
			leaves, err := file.Leaves()
			if err != nil {
				t.Fatal(err)
			}
			if leaves[0].Area() != 511 || leaves[0].Flags() != 127 || leaves[1].Area() != 128 || leaves[1].Flags() != 64 {
				t.Fatalf("unsigned fields=%+v", leaves)
			}
			leaves[0].SetArea(512)
			leaves[1].SetFlags(1)
			if err := file.SetLeaves(leaves); err != nil {
				t.Fatal(err)
			}
			binary.LittleEndian.PutUint16(payload[6:], 0xfe00)
			binary.LittleEndian.PutUint16(payload[stride+6:], 0x0280)
			next := rereadOrdinary(t, file)
			raw, _ := next.RawLump(bsp.LumpLeafs)
			got, err := next.Leaves()
			if err != nil || !bytes.Equal(raw, payload) || got[0].Area() != 0 || got[0].Flags() != 127 || got[1].Area() != 128 || got[1].Flags() != 1 {
				t.Fatalf("leaf%d %s bitfield reread=%+v %v", revision, pc.profile, got, err)
			}
		}
	}
}

func facePackedPayloads() map[bsp.LumpId][]byte {
	faces := make([]byte, 112)
	for i := 0; i < 2; i++ {
		binary.LittleEndian.PutUint16(faces[i*56+8:], 3)
		binary.LittleEndian.PutUint16(faces[i*56+48:], 0x8001)
		binary.LittleEndian.PutUint32(faces[i*56+52:], 0xf1234567+uint32(i))
	}
	primitive := make([]byte, 10)
	binary.LittleEndian.PutUint16(primitive[4:], 3)
	indices := []byte{0, 0, 1, 0, 2, 0}
	return map[bsp.LumpId][]byte{bsp.LumpFaces: append([]byte(nil), faces...), bsp.LumpOriginalFaces: append([]byte(nil), faces...), bsp.LumpFacesHDR: append([]byte(nil), faces...), bsp.LumpPrimitives: primitive, bsp.LumpPrimIndices: indices}
}

func TestFaceShadowBitNeverEntersPrimitiveCountAndEditsPersist(t *testing.T) {
	for _, profile := range []bsp.Profile{bsp.ProfileSDK2013, bsp.ProfileCSGO, bsp.ProfileL4D2, bsp.ProfileReactiveDrop} {
		for _, revision := range []int32{0, 1} {
			payloads := facePackedPayloads()
			file := spatialFixture(t, profile, payloads)
			for _, id := range []bsp.LumpId{bsp.LumpFaces, bsp.LumpOriginalFaces, bsp.LumpFacesHDR} {
				desc, _ := file.Descriptor(id)
				desc.Version = revision
				if err := file.SetRawLump(id, desc, payloads[id]); err != nil {
					t.Fatal(err)
				}
			}
			// One actual primitive, not32769; the face owns exactly three winding indices.
			primitives, err := file.Primitives()
			if err != nil || len(primitives) != 1 {
				t.Fatalf("flagged primitive consumer=%+v %v", primitives, err)
			}
			for _, api := range []struct {
				id  bsp.LumpId
				get func() ([]facevalue.Face, error)
				set func([]facevalue.Face) error
			}{{bsp.LumpFaces, file.Faces, file.SetFaces}, {bsp.LumpOriginalFaces, file.OriginalFaces, file.SetOriginalFaces}, {bsp.LumpFacesHDR, file.FacesHDR, file.SetFacesHDR}} {
				faces, err := api.get()
				if err != nil || faces[0].PrimitiveCount() != 1 || faces[0].DynamicShadowsEnabled() {
					t.Fatalf("face%d packed fields=%+v %v", api.id, faces, err)
				}
				before, _ := file.RawLump(api.id)
				desc, _ := file.Descriptor(api.id)
				faces[1].SetPrimitiveCount(2)
				if err := api.set(faces); !errors.Is(err, bsp.ErrMalformed) {
					t.Fatalf("masked invalid range error=%v", err)
				}
				raw, _ := file.RawLump(api.id)
				after, _ := file.Descriptor(api.id)
				if !bytes.Equal(raw, before) || after != desc {
					t.Fatal("failed face range partially committed")
				}
				faces[1].SetPrimitiveCount(1)
				faces[1].SetDynamicShadowsEnabled(true)
				if err := api.set(faces); err != nil {
					t.Fatal(err)
				}
				binary.LittleEndian.PutUint16(payloads[api.id][56+48:], 1)
			}
			next := rereadOrdinary(t, file)
			for _, api := range []struct {
				id  bsp.LumpId
				get func() ([]facevalue.Face, error)
			}{{bsp.LumpFaces, next.Faces}, {bsp.LumpOriginalFaces, next.OriginalFaces}, {bsp.LumpFacesHDR, next.FacesHDR}} {
				raw, _ := next.RawLump(api.id)
				faces, err := api.get()
				if err != nil || !bytes.Equal(raw, payloads[api.id]) || faces[1].PrimitiveCount() != 1 || !faces[1].DynamicShadowsEnabled() || faces[0].NumPrims != 0x8001 {
					t.Fatalf("face%d persisted fields=%+v %v", api.id, faces, err)
				}
			}
		}
	}
}

func TestFaceWindingEditRejectsDanglingLocalIndices(t *testing.T) {
	file := spatialFixture(t, bsp.ProfileCSGO, facePackedPayloads())
	faces, err := file.Faces()
	if err != nil {
		t.Fatal(err)
	}
	before, _ := file.RawLump(bsp.LumpFaces)
	faces[1].NumEdges = 2
	if err := file.SetFaces(faces); !errors.Is(err, bsp.ErrMalformed) {
		t.Fatalf("dangling winding index=%v", err)
	}
	raw, _ := file.RawLump(bsp.LumpFaces)
	if !bytes.Equal(raw, before) {
		t.Fatal("dangling face edit committed")
	}
}
