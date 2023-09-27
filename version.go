package bsp

import (
	"fmt"

	"github.com/galaco/bsp/lump"
)

// LumpResolverByBSPVersion returns an empty bsp lump for the specified bsp version and lump id
// If a version is not 19,20,21 then a generic lump that holds
// raw bytes only ([]byte) is returned.
func LumpResolverByBSPVersion(id, version int) (l Lump, err error) {
	if id < 0 || id > 63 {
		return nil, fmt.Errorf("invalid lump id: %d provided", id)
	}

	switch version {
	case 19:
		// @TODO: Implement v21.
		l, err = getv20Lump(id)
	case 20:
		l, err = getv20Lump(id)
	case 21:
		// @TODO: Implement v21.
		l, err = getv20Lump(id)
	default:
		l, err = &lump.RawBytes{}, nil
	}

	l.SetVersion(int32(version))

	return l, nil
}

// Getv20Lump returns the corresponding v20 lump for provided index.
func getv20Lump(index int) (Lump, error) {
	switch index {
	case 0:
		return &lump.EntData{}, nil
	case 1:
		return &lump.Planes{}, nil
	case 2:
		return &lump.TexData{}, nil
	case 3:
		return &lump.Vertex{}, nil
	case 4:
		return &lump.Visibility{}, nil
	case 5:
		return &lump.Node{}, nil
	case 6:
		return &lump.TexInfo{}, nil
	case 7:
		return &lump.Face{}, nil
	case 8:
		return &lump.Lighting{}, nil
	case 9:
		return &lump.Occlusion{}, nil
	case 10:
		return &lump.Leaf{}, nil
	case 11:
		return &lump.FaceId{}, nil
	case 12:
		return &lump.Edge{}, nil
	case 13:
		return &lump.Surfedge{}, nil
	case 14:
		return &lump.Model{}, nil
	case 15:
		return &lump.WorldLight{}, nil
	case 16:
		return &lump.LeafFace{}, nil
	case 17:
		return &lump.LeafBrush{}, nil
	case 18:
		return &lump.Brush{}, nil
	case 19:
		return &lump.BrushSide{}, nil
	case 20:
		return &lump.Area{}, nil
	case 21:
		return &lump.AreaPortal{}, nil
	case 22:
		return &lump.RawBytes{}, nil //portals | unused0 | propcollision
	case 23:
		return &lump.RawBytes{}, nil //clusters | unused1 | prophulls
	case 24:
		return &lump.RawBytes{}, nil //portalverts | unused2 | prophullverts
	case 25:
		return &lump.RawBytes{}, nil //clusterportals | unused3 | proptris
	case 26:
		return &lump.DispInfo{}, nil
	case 27:
		return &lump.Face{}, nil
	case 28:
		return &lump.PhysDisp{}, nil
	case 29:
		return &lump.RawBytes{}, nil //physcollide - IN PROGRESS
	case 30:
		return &lump.VertNormal{}, nil
	case 31:
		return &lump.VertNormalIndice{}, nil
	case 32:
		return &lump.RawBytes{}, nil //disp lightmap alphas - IS STRIPPED ANYWAY?
	case 33:
		return &lump.DispVert{}, nil
	case 34:
		return &lump.DispLightmapSamplePosition{}, nil
	case 35:
		return &lump.Game{}, nil
	case 36:
		return &lump.LeafWaterData{}, nil
	case 37:
		return &lump.RawBytes{}, nil //primitives FIXME - Appears to be 4bytes unaccounted for at end of lump?
	case 38:
		return &lump.PrimVert{}, nil
	case 39:
		return &lump.PrimIndice{}, nil
	case 40:
		return &lump.Pakfile{}, nil //pakfile
	case 41:
		return &lump.ClipPortalVerts{}, nil
	case 42:
		return &lump.Cubemap{}, nil
	case 43:
		return &lump.TexDataStringData{}, nil
	case 44:
		return &lump.TexDataStringTable{}, nil
	case 45:
		return &lump.Overlay{}, nil
	case 46:
		return &lump.LeafMinDistToWater{}, nil
	case 47:
		return &lump.FaceMacroTextureInfo{}, nil
	case 48:
		return &lump.DispTris{}, nil
	case 49:
		return &lump.RawBytes{}, nil //physcollidesurface | prop blob
	case 50:
		return &lump.RawBytes{}, nil //wateroverlays
	case 51:
		return &lump.LeafAmbientIndexHDR{}, nil
	case 52:
		return &lump.LeafAmbientIndex{}, nil
	case 53:
		return &lump.Lighting{}, nil //lighting hdr
	case 54:
		return &lump.WorldLightHDR{}, nil //worldlights hdr
	case 55:
		return &lump.LeafAmbientLightingHDR{}, nil
	case 56:
		return &lump.LeafAmbientLighting{}, nil //leaf ambient lighting
	case 57:
		return &lump.RawBytes{}, nil //xzippakfile
	case 58:
		return &lump.FaceHDR{}, nil
	case 59:
		return &lump.MapFlags{}, nil
	case 60:
		return &lump.OverlayFade{}, nil
	case 61:
		return &lump.RawBytes{}, nil //overlay system levels
	case 62:
		return &lump.RawBytes{}, nil //physlevel
	case 63:
		return &lump.RawBytes{}, nil //disp multiblend
	default:
		return nil, fmt.Errorf("invalid lump id")
	}
}
