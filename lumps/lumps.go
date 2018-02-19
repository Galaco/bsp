package lumps

import (
	"log"
)

/**
	Lump interface.
	Organise Lump data in a cleaner and more accessible manner
 */
type ILump interface {
	FromBytes([]byte, int32) ILump

	GetData() interface{}

	ToBytes() []byte
}

/**
	Helper info for a lump
 */
type LumpInfo struct {
	length int32
}
func (info LumpInfo) GetLength() int32 {
	return info.length
}
func (info LumpInfo) SetLength(length int32) {
	info.length = length
}





/**
	Return an instance of a Lump for a given offset.
 */
func GetLumpForIndex(index int) ILump {
	if index < 0 || index > 63 {
		log.Fatal("Invalid lump index provided.")
	}

	return lumpMap[index]
}


var lumpMap = [64]ILump{
	EntData{},
	Planes{},
	TexData{},
	Vertex{},
	Visibility{},
	Node{},
	TexInfo{},
	Face{},
	Unimplemented{}, //lighting
	Unimplemented{}, //occlusion
	Leaf{},
	Unimplemented{}, //faceids
	Edge{},
	Surfedge{},
	Model{},
	Unimplemented{}, //worldlights
	LeafFace{},
	LeafBrush{},
	Brush{},
	BrushSide{},
	Area{},
	AreaPortal{},
	Unimplemented{}, //portals | unused0 Z propcollision
	Unimplemented{}, //clusters | unused1 | prophulls
	Unimplemented{}, //portalverts | unused2 | prophullverts
	Unimplemented{}, //clusterportals | unused3 | proptris
	Unimplemented{}, //dispinfo
	Unimplemented{}, //originalfaces
	Unimplemented{}, //physdisp
	Unimplemented{}, //physcollide
	Unimplemented{}, //vertnormals
	Unimplemented{}, //vertnormalindices
	Unimplemented{}, //disp lightmap alphas
	Unimplemented{}, //disp verts
	Unimplemented{}, //disp lightmap sample positions
	Game{},
	LeafWaterData{},
	Unimplemented{}, //primitives
	Unimplemented{}, //primverts
	Unimplemented{}, //primindices
	Unimplemented{}, //pakfile
	Unimplemented{}, //clipportalverts
	Cubemap{},
	TexdataStringData{},
	TexDataStringTable{},
	Overlay{}, //overlays
	Unimplemented{}, //leafmindisttowater
	Unimplemented{}, //face marco texture info
	Unimplemented{}, //disp tris
	Unimplemented{}, //physcollidesurface | prop blob
	Unimplemented{}, //wateroverlays
	Unimplemented{}, //lightmappages | leaf ambient index hdr
	Unimplemented{}, //lightmappageinfos | leaf ambient index
	Unimplemented{}, //lighting hdr
	Unimplemented{}, //worldlights hdr
	Unimplemented{}, //leaf ambient lighting hdr
	Unimplemented{}, //leaf ambient lighting
	Unimplemented{}, //xzippakfile
	Unimplemented{}, //faces hdr
	Unimplemented{}, //map flags
	OverlayFade{}, //overlay fades
	Unimplemented{}, //overlay system levels
	Unimplemented{}, //physlevel
	Unimplemented{}, //disp multiblend
}

