package bsp

import (
	"reflect"
	"testing"

	"github.com/galaco/bsp/lumps"
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
		&lumps.RawBytes{},
		&lumps.RawBytes{},
		&lumps.RawBytes{},
		&lumps.RawBytes{},
		&lumps.DispInfo{},
		&lumps.Face{},
		&lumps.PhysDisp{},
		&lumps.RawBytes{},
		&lumps.VertNormal{},
		&lumps.VertNormalIndice{},
		&lumps.RawBytes{},
		&lumps.DispVert{},
		&lumps.DispLightmapSamplePosition{},
		&lumps.Game{},
		&lumps.LeafWaterData{},
		&lumps.RawBytes{},
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
		&lumps.RawBytes{},
		&lumps.RawBytes{},
		&lumps.LeafAmbientIndexHDR{},
		&lumps.LeafAmbientIndex{},
		&lumps.Lighting{},
		&lumps.WorldLightHDR{},
		&lumps.LeafAmbientLightingHDR{},
		&lumps.LeafAmbientLighting{},
		&lumps.RawBytes{},
		&lumps.FaceHDR{},
		&lumps.MapFlags{},
		&lumps.OverlayFade{},
		&lumps.RawBytes{},
		&lumps.RawBytes{},
		&lumps.RawBytes{}, //disp multiblend
	}

	return lMap[index]
}
