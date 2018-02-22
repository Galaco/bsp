package versions

import (
	"github.com/galaco/bsp/lumps"
)


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
		lumps.Occlusion{}, //occlusion
		lumps.Leaf{},
		lumps.Unimplemented{}, //faceids
		lumps.Edge{},
		lumps.Surfedge{},
		lumps.Model{},
		lumps.Unimplemented{}, //worldlights
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
		lumps.Unimplemented{}, //originalfaces
		lumps.Unimplemented{}, //physdisp
		lumps.Unimplemented{}, //physcollide
		lumps.Unimplemented{}, //vertnormals
		lumps.Unimplemented{}, //vertnormalindices
		lumps.Unimplemented{}, //disp lightmap alphas
		lumps.Unimplemented{}, //disp verts
		lumps.Unimplemented{}, //disp lightmap sample positions
		lumps.Game{},
		lumps.Unimplemented{}, //LeafWaterData FIXME
		lumps.Unimplemented{}, //primitives
		lumps.Unimplemented{}, //primverts
		lumps.Unimplemented{}, //primindices
		lumps.Unimplemented{}, //pakfile
		lumps.Unimplemented{}, //clipportalverts
		lumps.Unimplemented{}, //cubemap FIXME
		lumps.TexdataStringData{},
		lumps.TexDataStringTable{},
		lumps.Overlay{},       //overlays
		lumps.Unimplemented{}, //leafmindisttowater
		lumps.Unimplemented{}, //face marco texture info
		lumps.Unimplemented{}, //disp tris
		lumps.Unimplemented{}, //physcollidesurface | prop blob
		lumps.Unimplemented{}, //wateroverlays
		lumps.Unimplemented{}, //lightmappages | leaf ambient index hdr
		lumps.LeafAmbientIndex{}, //lightmappageinfos | leaf ambient index
		lumps.Unimplemented{}, //lighting hdr
		lumps.WorldLightHDR{}, //worldlights hdr
		lumps.Unimplemented{}, //leaf ambient lighting hdr
		lumps.Unimplemented{}, //leaf ambient lighting
		lumps.Unimplemented{}, //xzippakfile
		lumps.Unimplemented{}, //faces hdr
		lumps.Unimplemented{}, //map flags
		lumps.OverlayFade{},   //overlay fades
		lumps.Unimplemented{}, //overlay system levels
		lumps.Unimplemented{}, //physlevel
		lumps.Unimplemented{}, //disp multiblend
	}
}