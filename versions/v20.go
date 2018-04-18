package versions

import (
	"github.com/galaco/bsp/lumps"
)

// Get mapping of lump types for v20 bsps
func GetVersion20Mapping() [64]lumps.ILump {
	return [64]lumps.ILump{
		lumps.EntData{},
		lumps.Planes{},
		lumps.TexData{},
		lumps.Vertex{},
		lumps.Visibility{},
		lumps.Node{},
		lumps.TexInfo{},
		lumps.Face{},
		lumps.Lighting{},
		lumps.Occlusion{},
		lumps.Leaf{},
		lumps.FaceId{},
		lumps.Edge{},
		lumps.Surfedge{},
		lumps.Model{},
		lumps.WorldLight{},
		lumps.LeafFace{},
		lumps.LeafBrush{},
		lumps.Brush{},
		lumps.BrushSide{},
		lumps.Area{},
		lumps.AreaPortal{},
		lumps.Unimplemented{}, //portals | unused0 | propcollision
		lumps.Unimplemented{}, //clusters | unused1 | prophulls
		lumps.Unimplemented{}, //portalverts | unused2 | prophullverts
		lumps.Unimplemented{}, //clusterportals | unused3 | proptris
		lumps.Unimplemented{}, //dispinfo
		lumps.Face{}, //originalfaces
		lumps.Unimplemented{}, //physdisp NOTE: bsplib actually implements this as a byte* anyway
		lumps.Unimplemented{}, //physcollide
		lumps.Unimplemented{}, //vertnormals
		lumps.Unimplemented{}, //vertnormalindices
		lumps.Unimplemented{}, //disp lightmap alphas
		lumps.DispVert{},
		lumps.Unimplemented{}, //disp lightmap sample positions NOTE: bsplib implements this as an unsigned char vector
		lumps.Game{},
		lumps.LeafWaterData{}, //LeafWaterData FIXME
		lumps.Unimplemented{}, //primitives FIXME
		lumps.PrimVert{}, //primverts
		lumps.Unimplemented{}, //primindices
		lumps.Unimplemented{}, //pakfile
		lumps.Unimplemented{}, //clipportalverts
		lumps.Cubemap{},
		lumps.TexdataStringData{},
		lumps.TexDataStringTable{},
		lumps.Overlay{},
		lumps.LeafMinDistToWater{}, //leafmindisttowater
		lumps.FaceMacroTextureInfo{}, //face macro texture info
		lumps.DispTris{},
		lumps.Unimplemented{}, //physcollidesurface | prop blob
		lumps.Unimplemented{}, //wateroverlays
		lumps.LeafAmbientIndexHDR{},
		lumps.LeafAmbientIndex{},
		lumps.Unimplemented{}, //lighting hdr
		lumps.WorldLightHDR{}, //worldlights hdr
		lumps.LeafAmbientLightingHDR{},
		lumps.LeafAmbientLighting{}, //leaf ambient lighting
		lumps.Unimplemented{}, //xzippakfile
		lumps.FaceHDR{}, //faces hdr
		lumps.MapFlags{}, //map flags
		lumps.OverlayFade{},
		lumps.Unimplemented{}, //overlay system levels
		lumps.Unimplemented{}, //physlevel
		lumps.Unimplemented{}, //disp multiblend
	}
}