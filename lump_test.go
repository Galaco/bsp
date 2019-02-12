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
		&lumps.TexdataStringData{},
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
