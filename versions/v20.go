package versions

import (
	"github.com/galaco/bsp/lumps"
	"errors"
)

func Getv20Lump(index int) (lumps.ILump,error) {
	var ret lumps.ILump
	var err error
	switch index {
	case 0:
		ret = &lumps.EntData{}
	case 1:
		ret = &lumps.Planes{}
	case 2:
		ret = &lumps.TexData{}
	case 3:
		ret = &lumps.Vertex{}
	case 4:
		ret = &lumps.Visibility{}
	case 5:
		ret = &lumps.Node{}
	case 6:
		ret = &lumps.TexInfo{}
	case 7:
		ret = &lumps.Face{}
	case 8:
		ret = &lumps.Lighting{}
	case 9:
		ret = &lumps.Occlusion{}
	case 10:
		ret = &lumps.Leaf{}
	case 11:
		ret = &lumps.FaceId{}
	case 12:
		ret = &lumps.Edge{}
	case 13:
		ret = &lumps.Surfedge{}
	case 14:
		ret = &lumps.Model{}
	case 15:
		ret = &lumps.WorldLight{}
	case 16:
		ret = &lumps.LeafFace{}
	case 17:
		ret = &lumps.LeafBrush{}
	case 18:
		ret = &lumps.Brush{}
	case 19:
		ret = &lumps.BrushSide{}
	case 20:
		ret = &lumps.Area{}
	case 21:
		ret = &lumps.AreaPortal{}
	case 22:
		ret = &lumps.Unimplemented{} //portals | unused0 | propcollision
	case 23:
		ret = &lumps.Unimplemented{} //clusters | unused1 | prophulls
	case 24:
		ret = &lumps.Unimplemented{} //portalverts | unused2 | prophullverts
	case 25:
		ret = &lumps.Unimplemented{} //clusterportals | unused3 | proptris
	case 26:
		ret = &lumps.Unimplemented{}
	case 27:
		ret = &lumps.Face{}
	case 28:
		ret = &lumps.PhysDisp{}
	case 29:
		ret = &lumps.Unimplemented{} //physcollide - IN PROGRESS
	case 30:
		ret = &lumps.VertNormal{}
	case 31:
		ret = &lumps.VertNormalIndice{}
	case 32:
		ret = &lumps.Unimplemented{} //disp lightmap alphas - IS STRIPPED ANYWAY?
	case 33:
		ret = &lumps.DispVert{}
	case 34:
		ret = &lumps.DispLightmapSamplePosition{}
	case 35:
		ret = &lumps.Game{}
	case 36:
		ret = &lumps.LeafWaterData{}
	case 37:
		ret = &lumps.Unimplemented{} //primitives FIXME - Appears to be 4bytes unaccounted for at end of lump?
	case 38:
		ret = &lumps.PrimVert{}
	case 39:
		ret = &lumps.PrimIndice{}
	case 40:
		ret = &lumps.Unimplemented{} //pakfile
	case 41:
		ret = &lumps.ClipPortalVerts{}
	case 42:
		ret = &lumps.Cubemap{}
	case 43:
		ret = &lumps.TexdataStringData{}
	case 44:
		ret = &lumps.TexDataStringTable{}
	case 45:
		ret = &lumps.Overlay{}
	case 46:
		ret = &lumps.LeafMinDistToWater{}
	case 47:
		ret = &lumps.FaceMacroTextureInfo{}
	case 48:
		ret = &lumps.DispTris{}
	case 49:
		ret = &lumps.Unimplemented{} //physcollidesurface | prop blob
	case 50:
		ret = &lumps.Unimplemented{} //wateroverlays
	case 51:
		ret = &lumps.LeafAmbientIndexHDR{}
	case 52:
		ret = &lumps.LeafAmbientIndex{}
	case 53:
		ret = &lumps.Unimplemented{} //lighting hdr
	case 54:
		ret = &lumps.WorldLightHDR{} //worldlights hdr
	case 55:
		ret = &lumps.LeafAmbientLightingHDR{}
	case 56:
		ret = &lumps.LeafAmbientLighting{} //leaf ambient lighting
	case 57:
		ret = &lumps.Unimplemented{} //xzippakfile
	case 58:
		ret = &lumps.FaceHDR{}
	case 59:
		ret = &lumps.MapFlags{}
	case 60:
		ret = &lumps.OverlayFade{}
	case 61:
		ret = &lumps.Unimplemented{} //overlay system levels
	case 62:
		ret = &lumps.Unimplemented{} //physlevel
	case 63:
		ret = &lumps.Unimplemented{} //disp multiblend
	default:
		err = errors.New("invalid lump id")
	}

	return ret,err
}
