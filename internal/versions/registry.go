package versions

import (
	"fmt"
	"github.com/galaco/bsp/internal/codec"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

// New selects a codec using the actual lump revision and engine dialect.
func New(index int, context format.Context) (lumps.ILump, error) {
	if err := ValidateRevision(index, context); err != nil {
		return nil, err
	}
	var ret lumps.ILump
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
		ret = new(codec.Visibility)
	case 5:
		ret = new(lumps.Node)
	case 6:
		ret = new(lumps.TexInfo)
	case 7:
		ret = new(codec.Faces)
	case 8:
		ret = new(lumps.Lighting)
	case 9:
		ret = new(codec.Occlusion)
	case 10:
		ret = new(codec.Leaf)
	case 11:
		ret = new(lumps.FaceId)
	case 12:
		ret = new(lumps.Edge)
	case 13:
		ret = new(lumps.Surfedge)
	case 14:
		ret = new(lumps.Model)
	case 15:
		ret = new(codec.WorldLight)
	case 16:
		ret = new(lumps.LeafFace)
	case 17:
		ret = new(lumps.LeafBrush)
	case 18:
		ret = new(lumps.Brush)
	case 19:
		ret = new(codec.BrushSide)
	case 20:
		ret = new(lumps.Area)
	case 21:
		ret = new(lumps.AreaPortal)
	case 22:
		ret = new(codec.PropCollisions)
	case 23:
		ret = new(codec.PropHulls)
	case 24:
		ret = new(codec.PropHullVertices)
	case 25:
		ret = new(codec.PropHullTriangles)
	case 26:
		ret = new(codec.Displacements)
	case 27:
		ret = new(codec.Faces)
	case 28, 29:
		return nil, fmt.Errorf("opaque lump %d", index)
	case 30:
		ret = new(lumps.VertNormal)
	case 31:
		ret = new(lumps.VertNormalIndice)
	case 32:
		return nil, fmt.Errorf("opaque lump %d", index)
	case 33:
		ret = new(lumps.DispVert)
	case 34, 35:
		return nil, fmt.Errorf("opaque lump %d", index)
	case 36:
		ret = new(lumps.LeafWaterData)
	case 37:
		ret = new(codec.Primitives)
	case 38:
		ret = new(lumps.PrimVert)
	case 39:
		ret = new(lumps.PrimIndice)
	case 40:
		ret = new(lumps.Pakfile) //pakfile
	case 41:
		ret = new(lumps.ClipPortalVerts)
	case 42:
		ret = new(codec.Cubemaps)
	case 43:
		ret = new(lumps.TexDataStringData)
	case 44:
		ret = new(lumps.TexDataStringTable)
	case 45:
		ret = new(codec.Overlays)
	case 46:
		ret = new(lumps.LeafMinDistToWater)
	case 47:
		ret = new(lumps.FaceMacroTextureInfo)
	case 48:
		ret = new(lumps.DispTris)
	case 49:
		return nil, fmt.Errorf("byte-stream lump %d; use attached PropBlob accessor", index)
	case 50:
		ret = new(codec.WaterOverlays)
	case 51:
		ret = new(codec.LeafAmbientIndexHDR)
	case 52:
		ret = new(codec.LeafAmbientIndex)
	case 53:
		ret = new(lumps.Lighting) //lighting hdr
	case 54:
		ret = new(codec.WorldLightHDR)
	case 55:
		ret = new(codec.LeafAmbientLightingHDR)
	case 56:
		ret = new(codec.LeafAmbientLighting)
	case 57:
		return nil, fmt.Errorf("opaque lump %d", index)
	case 58:
		ret = new(codec.Faces)
	case 59:
		ret = new(codec.MapFlags)
	case 60:
		ret = new(lumps.OverlayFade)
	case 61:
		ret = new(codec.OverlaySystemLevels)
	case 62:
		return nil, fmt.Errorf("opaque lump %d", index)
	case 63:
		ret = new(codec.DisplacementMultiblends)
	default:
		return nil, fmt.Errorf("invalid lump %d", index)
	}
	ret.SetContext(context)
	return ret, nil
}
