package bsp

import (
	"github.com/galaco/bsp/lumps"
	"reflect"
	"testing"
)

func TestGetReferenceLumpByIndex(t *testing.T) {
	for i := 0; i < 64; i++ {
		lump, _ := getReferenceLumpByIndex(i, 20)
		if reflect.TypeOf(lump) != reflect.TypeOf(getExpectedLump(i)) {
			t.Errorf("Lump type does not match expected for identifer: %d\n", i)
		}
	}

	_, err := getReferenceLumpByIndex(65, 20)
	if err == nil {
		t.Error("invalid lump provided, but no error returned")
	}
}

func TestLump_Contents(t *testing.T) {
	sut := new(Lump)
	l := new(lumps.Generic)
	l.SetLength(654)
	sut.SetContents(l)
	if sut.Contents() != l {
		t.Error("unexpected lump data returned")
	}
}

func TestLump_Length(t *testing.T) {
	sut := Lump{}
	sut.length = 32
	if sut.Length() != 32 {
		t.Error("incorrect length returned for lump")
	}
}

func TestLump_RawContents(t *testing.T) {
	sut := new(Lump)
	data := []byte{0,1,3,4,5,6}
	sut.SetRawContents(data)
	if len(sut.RawContents()) != len(data) {
		t.Error("unexpected lump data returned")
	}
}

func TestLump_SetContents(t *testing.T) {
	sut := new(Lump)
	l := new(lumps.Generic)
	l.SetLength(654)
	sut.SetContents(l)
	if sut.Contents() != l {
		t.Error("unexpected lump data returned")
	}
}

func TestLump_SetId(t *testing.T) {
	sut := Lump{}
	sut.SetId(LumpPakfile)

	if sut.id != LumpPakfile {
		t.Error("incorrect lump id")
	}
}

func TestLump_SetRawContents(t *testing.T) {
	sut := Lump{}
	data := []byte{0, 1, 4, 3, 2}
	sut.SetRawContents(data)
	for idx, b := range sut.RawContents() {
		if data[idx] != b {
			t.Error("raw lump data mismatch")
		}
	}
}

func getExpectedLump(index int) lumps.ILump {
	lMap := [64]lumps.ILump{
		&lumps.EntData{},
		&lumps.Planes{},
		&lumps.TexData{},
		&lumps.Vertex{},
		&lumps.Visibility{},
		&lumps.Node{},
		&lumps.TexInfo{},
		&lumps.Face{},
		&lumps.Lighting{},
		&lumps.Occlusion{},
		&lumps.Leaf{},
		&lumps.FaceId{},
		&lumps.Edge{},
		&lumps.Surfedge{},
		&lumps.Model{},
		&lumps.WorldLight{},
		&lumps.LeafFace{},
		&lumps.LeafBrush{},
		&lumps.Brush{},
		&lumps.BrushSide{},
		&lumps.Area{},
		&lumps.AreaPortal{},
		&lumps.Unimplemented{},
		&lumps.Unimplemented{},
		&lumps.Unimplemented{},
		&lumps.Unimplemented{},
		&lumps.DispInfo{},
		&lumps.Face{},
		&lumps.PhysDisp{},
		&lumps.Unimplemented{},
		&lumps.VertNormal{},
		&lumps.VertNormalIndice{},
		&lumps.Unimplemented{},
		&lumps.DispVert{},
		&lumps.DispLightmapSamplePosition{},
		&lumps.Game{},
		&lumps.LeafWaterData{},
		&lumps.Unimplemented{},
		&lumps.PrimVert{},
		&lumps.PrimIndice{},
		&lumps.Pakfile{},
		&lumps.ClipPortalVerts{},
		&lumps.Cubemap{},
		&lumps.TexDataStringData{},
		&lumps.TexDataStringTable{},
		&lumps.Overlay{},
		&lumps.LeafMinDistToWater{},
		&lumps.FaceMacroTextureInfo{},
		&lumps.DispTris{},
		&lumps.Unimplemented{},
		&lumps.Unimplemented{},
		&lumps.LeafAmbientIndexHDR{},
		&lumps.LeafAmbientIndex{},
		&lumps.Unimplemented{},
		&lumps.WorldLightHDR{},
		&lumps.LeafAmbientLightingHDR{},
		&lumps.LeafAmbientLighting{},
		&lumps.Unimplemented{},
		&lumps.FaceHDR{},
		&lumps.MapFlags{},
		&lumps.OverlayFade{},
		&lumps.Unimplemented{},
		&lumps.Unimplemented{},
		&lumps.Unimplemented{}, //disp multiblend
	}

	return lMap[index]
}
