package versions

import (
	"errors"

	"github.com/galaco/bsp/lumps"
)

// Getv20Lump returns the corresponding v20 lump for provided index
func Getv20Lump(index int) (lumps.ILump, error) {
	var ret lumps.ILump
	var err error
	switch index {
	case 0:
		ret = new(lumps.EntData)
	case 1:
		ret = new(lumps.Planes)
	case 2:
		ret = new(lumps.TexData)
	case 3:
		ret = new(lumps.Vertex)
	case 4:
		ret = new(lumps.Visibility)
	case 5:
		ret = new(lumps.Node)
	case 6:
		ret = new(lumps.TexInfo)
	case 7:
		ret = new(lumps.Face)
	case 8:
		ret = new(lumps.Lighting)
	case 9:
		ret = new(lumps.Occlusion)
	case 10:
		ret = new(lumps.Leaf)
	case 11:
		ret = new(lumps.FaceId)
	case 12:
		ret = new(lumps.Edge)
	case 13:
		ret = new(lumps.Surfedge)
	case 14:
		ret = new(lumps.Model)
	case 15:
		ret = new(lumps.WorldLight)
	case 16:
		ret = new(lumps.LeafFace)
	case 17:
		ret = new(lumps.LeafBrush)
	case 18:
		ret = new(lumps.Brush)
	case 19:
		ret = new(lumps.BrushSide)
	case 20:
		ret = new(lumps.Area)
	case 21:
		ret = new(lumps.AreaPortal)
	case 22:
		ret = new(lumps.RawBytes) //portals | unused0 | propcollision
	case 23:
		ret = new(lumps.RawBytes) //clusters | unused1 | prophulls
	case 24:
		ret = new(lumps.RawBytes) //portalverts | unused2 | prophullverts
	case 25:
		ret = new(lumps.RawBytes) //clusterportals | unused3 | proptris
	case 26:
		ret = new(lumps.DispInfo)
	case 27:
		ret = new(lumps.Face)
	case 28:
		ret = new(lumps.PhysDisp)
	case 29:
		ret = new(lumps.RawBytes) //physcollide - IN PROGRESS
	case 30:
		ret = new(lumps.VertNormal)
	case 31:
		ret = new(lumps.VertNormalIndice)
	case 32:
		ret = new(lumps.RawBytes) //disp lightmap alphas - IS STRIPPED ANYWAY?
	case 33:
		ret = new(lumps.DispVert)
	case 34:
		ret = new(lumps.DispLightmapSamplePosition)
	case 35:
		ret = new(lumps.Game)
	case 36:
		ret = new(lumps.LeafWaterData)
	case 37:
		ret = new(lumps.RawBytes) //primitives FIXME - Appears to be 4bytes unaccounted for at end of lump?
	case 38:
		ret = new(lumps.PrimVert)
	case 39:
		ret = new(lumps.PrimIndice)
	case 40:
		ret = new(lumps.Pakfile) //pakfile
	case 41:
		ret = new(lumps.ClipPortalVerts)
	case 42:
		ret = new(lumps.Cubemap)
	case 43:
		ret = new(lumps.TexDataStringData)
	case 44:
		ret = new(lumps.TexDataStringTable)
	case 45:
		ret = new(lumps.Overlay)
	case 46:
		ret = new(lumps.LeafMinDistToWater)
	case 47:
		ret = new(lumps.FaceMacroTextureInfo)
	case 48:
		ret = new(lumps.DispTris)
	case 49:
		ret = new(lumps.RawBytes) //physcollidesurface | prop blob
	case 50:
		ret = new(lumps.RawBytes) //wateroverlays
	case 51:
		ret = new(lumps.LeafAmbientIndexHDR)
	case 52:
		ret = new(lumps.LeafAmbientIndex)
	case 53:
		ret = new(lumps.Lighting) //lighting hdr
	case 54:
		ret = new(lumps.WorldLightHDR) //worldlights hdr
	case 55:
		ret = new(lumps.LeafAmbientLightingHDR)
	case 56:
		ret = new(lumps.LeafAmbientLighting) //leaf ambient lighting
	case 57:
		ret = new(lumps.RawBytes) //xzippakfile
	case 58:
		ret = new(lumps.FaceHDR)
	case 59:
		ret = new(lumps.MapFlags)
	case 60:
		ret = new(lumps.OverlayFade)
	case 61:
		ret = new(lumps.RawBytes) //overlay system levels
	case 62:
		ret = new(lumps.RawBytes) //physlevel
	case 63:
		ret = new(lumps.RawBytes) //disp multiblend
	default:
		err = errors.New("invalid lump id")
	}

	return ret, err
}
