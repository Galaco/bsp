package bsp

import (
	"github.com/galaco/bsp/internal/codec"
	"github.com/galaco/bsp/lumps"
	"github.com/go-gl/mathgl/mgl32"
)

// TextureData returns an owned snapshot of the decoded records.
func (file *File) TextureData() ([]lumps.TextureData, error) {
	decoded, err := file.decodeLump(LumpTexData)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.TextureData).GetData(), nil
}

// SetTextureData atomically replaces the payload, preserving its selected revision.
func (file *File) SetTextureData(records []lumps.TextureData) error {
	if err := file.checkTypedRevision(LumpTexData); err != nil {
		return err
	}
	raw, err := codec.EncodeTextureData(records)
	return file.commitOrdinary(LumpTexData, raw, err)
}

// Vertices returns an owned snapshot of the decoded records.
func (file *File) Vertices() ([]mgl32.Vec3, error) {
	decoded, err := file.decodeLump(LumpVertexes)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.Vertices).GetData(), nil
}

// SetVertices atomically replaces the payload, preserving its selected revision.
func (file *File) SetVertices(records []mgl32.Vec3) error {
	if err := file.checkTypedRevision(LumpVertexes); err != nil {
		return err
	}
	raw, err := codec.EncodeVertices(records)
	return file.commitOrdinary(LumpVertexes, raw, err)
}

// Nodes returns an owned snapshot of the decoded records.
func (file *File) Nodes() ([]lumps.Node, error) {
	decoded, err := file.decodeLump(LumpNodes)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.Nodes).GetData(), nil
}

// SetNodes atomically replaces the payload, preserving its selected revision.
func (file *File) SetNodes(records []lumps.Node) error {
	if err := file.checkTypedRevision(LumpNodes); err != nil {
		return err
	}
	raw, err := codec.EncodeNodes(records)
	return file.commitOrdinary(LumpNodes, raw, err)
}

// TextureInfo returns an owned snapshot of the decoded records.
func (file *File) TextureInfo() ([]lumps.TextureInfo, error) {
	decoded, err := file.decodeLump(LumpTexInfo)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.TextureInfo).GetData(), nil
}

// SetTextureInfo atomically replaces the payload, preserving its selected revision.
func (file *File) SetTextureInfo(records []lumps.TextureInfo) error {
	if err := file.checkTypedRevision(LumpTexInfo); err != nil {
		return err
	}
	raw, err := codec.EncodeTextureInfo(records)
	return file.commitOrdinary(LumpTexInfo, raw, err)
}

// Lighting returns the independent LDR channel without HDR fallback.
func (file *File) Lighting() ([]lumps.ColorRGBExponent32, error) {
	decoded, err := file.decodeLump(LumpLighting)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.Lighting).GetData(), nil
}

// SetLighting atomically replaces the payload, preserving its selected revision.
func (file *File) SetLighting(records []lumps.ColorRGBExponent32) error {
	if err := file.checkTypedRevision(LumpLighting); err != nil {
		return err
	}
	raw, err := codec.EncodeLighting(records)
	return file.commitOrdinary(LumpLighting, raw, err)
}

// LightingHDR returns the independent HDR channel without LDR fallback.
func (file *File) LightingHDR() ([]lumps.ColorRGBExponent32, error) {
	decoded, err := file.decodeLump(LumpLightingHDR)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.Lighting).GetData(), nil
}

// SetLightingHDR atomically replaces the payload, preserving its selected revision.
func (file *File) SetLightingHDR(records []lumps.ColorRGBExponent32) error {
	if err := file.checkTypedRevision(LumpLightingHDR); err != nil {
		return err
	}
	raw, err := codec.EncodeLighting(records)
	return file.commitOrdinary(LumpLightingHDR, raw, err)
}

// FaceIDs returns an owned snapshot of the decoded records.
func (file *File) FaceIDs() ([]lumps.FaceID, error) {
	decoded, err := file.decodeLump(LumpFaceIDs)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.FaceIDs).GetData(), nil
}

// SetFaceIDs atomically replaces the payload, preserving its selected revision.
func (file *File) SetFaceIDs(records []lumps.FaceID) error {
	if err := file.checkTypedRevision(LumpFaceIDs); err != nil {
		return err
	}
	raw, err := codec.EncodeFaceIDs(records)
	return file.commitOrdinary(LumpFaceIDs, raw, err)
}

// Edges returns an owned snapshot of the decoded records.
func (file *File) Edges() ([][2]uint16, error) {
	decoded, err := file.decodeLump(LumpEdges)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.Edges).GetData(), nil
}

// SetEdges atomically replaces the payload, preserving its selected revision.
func (file *File) SetEdges(records [][2]uint16) error {
	if err := file.checkTypedRevision(LumpEdges); err != nil {
		return err
	}
	raw, err := codec.EncodeEdges(records)
	return file.commitOrdinary(LumpEdges, raw, err)
}

// SurfaceEdges returns an owned snapshot of the decoded records.
func (file *File) SurfaceEdges() ([]int32, error) {
	decoded, err := file.decodeLump(LumpSurfEdges)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.SurfaceEdges).GetData(), nil
}

// SetSurfaceEdges atomically replaces the payload, preserving its selected revision.
func (file *File) SetSurfaceEdges(records []int32) error {
	if err := file.checkTypedRevision(LumpSurfEdges); err != nil {
		return err
	}
	raw, err := codec.EncodeSurfaceEdges(records)
	return file.commitOrdinary(LumpSurfEdges, raw, err)
}

// Models returns an owned snapshot of the decoded records.
func (file *File) Models() ([]lumps.Model, error) {
	decoded, err := file.decodeLump(LumpModels)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.Models).GetData(), nil
}

// SetModels atomically replaces the payload, preserving its selected revision.
func (file *File) SetModels(records []lumps.Model) error {
	if err := file.checkTypedRevision(LumpModels); err != nil {
		return err
	}
	raw, err := codec.EncodeModels(records)
	return file.commitOrdinary(LumpModels, raw, err)
}

// LeafFaces returns an owned snapshot of the decoded records.
func (file *File) LeafFaces() ([]uint16, error) {
	decoded, err := file.decodeLump(LumpLeafFaces)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.LeafFaces).GetData(), nil
}

// SetLeafFaces atomically replaces the payload, preserving its selected revision.
func (file *File) SetLeafFaces(records []uint16) error {
	if err := file.checkTypedRevision(LumpLeafFaces); err != nil {
		return err
	}
	raw, err := codec.EncodeLeafFaces(records)
	return file.commitOrdinary(LumpLeafFaces, raw, err)
}

// LeafBrushes returns an owned snapshot of the decoded records.
func (file *File) LeafBrushes() ([]uint16, error) {
	decoded, err := file.decodeLump(LumpLeafBrushes)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.LeafBrushes).GetData(), nil
}

// SetLeafBrushes atomically replaces the payload, preserving its selected revision.
func (file *File) SetLeafBrushes(records []uint16) error {
	if err := file.checkTypedRevision(LumpLeafBrushes); err != nil {
		return err
	}
	raw, err := codec.EncodeLeafBrushes(records)
	return file.commitOrdinary(LumpLeafBrushes, raw, err)
}

// Brushes returns an owned snapshot of the decoded records.
func (file *File) Brushes() ([]lumps.Brush, error) {
	decoded, err := file.decodeLump(LumpBrushes)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.Brushes).GetData(), nil
}

// SetBrushes atomically replaces the payload, preserving its selected revision.
func (file *File) SetBrushes(records []lumps.Brush) error {
	if err := file.checkTypedRevision(LumpBrushes); err != nil {
		return err
	}
	raw, err := codec.EncodeBrushes(records)
	return file.commitOrdinary(LumpBrushes, raw, err)
}

// Areas returns an owned snapshot of the decoded records.
func (file *File) Areas() ([]lumps.Area, error) {
	decoded, err := file.decodeLump(LumpAreas)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.Areas).GetData(), nil
}

// SetAreas atomically replaces the payload, preserving its selected revision.
func (file *File) SetAreas(records []lumps.Area) error {
	if err := file.checkTypedRevision(LumpAreas); err != nil {
		return err
	}
	raw, err := codec.EncodeAreas(records)
	return file.commitOrdinary(LumpAreas, raw, err)
}

// AreaPortals returns an owned snapshot of the decoded records.
func (file *File) AreaPortals() ([]lumps.AreaPortal, error) {
	decoded, err := file.decodeLump(LumpAreaPortals)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.AreaPortals).GetData(), nil
}

// SetAreaPortals atomically replaces the payload, preserving its selected revision.
func (file *File) SetAreaPortals(records []lumps.AreaPortal) error {
	if err := file.checkTypedRevision(LumpAreaPortals); err != nil {
		return err
	}
	raw, err := codec.EncodeAreaPortals(records)
	return file.commitOrdinary(LumpAreaPortals, raw, err)
}

// VertexNormals returns an owned snapshot of the decoded records.
func (file *File) VertexNormals() ([]lumps.VertexNormal, error) {
	decoded, err := file.decodeLump(LumpVertNormals)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.VertexNormals).GetData(), nil
}

// SetVertexNormals atomically replaces the payload, preserving its selected revision.
func (file *File) SetVertexNormals(records []lumps.VertexNormal) error {
	if err := file.checkTypedRevision(LumpVertNormals); err != nil {
		return err
	}
	raw, err := codec.EncodeVertexNormals(records)
	return file.commitOrdinary(LumpVertNormals, raw, err)
}

// VertexNormalIndices returns an owned snapshot of the decoded records.
func (file *File) VertexNormalIndices() ([]uint16, error) {
	decoded, err := file.decodeLump(LumpVertNormalIndices)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.VertexNormalIndices).GetData(), nil
}

// SetVertexNormalIndices atomically replaces the payload, preserving its selected revision.
func (file *File) SetVertexNormalIndices(records []uint16) error {
	if err := file.checkTypedRevision(LumpVertNormalIndices); err != nil {
		return err
	}
	raw, err := codec.EncodeVertexNormalIndices(records)
	return file.commitOrdinary(LumpVertNormalIndices, raw, err)
}

// DisplacementVertices returns an owned snapshot of the decoded records.
func (file *File) DisplacementVertices() ([]lumps.DisplacementVertex, error) {
	decoded, err := file.decodeLump(LumpDispVerts)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.DisplacementVertices).GetData(), nil
}

// SetDisplacementVertices atomically replaces the payload, preserving its selected revision.
func (file *File) SetDisplacementVertices(records []lumps.DisplacementVertex) error {
	if err := file.checkTypedRevision(LumpDispVerts); err != nil {
		return err
	}
	raw, err := codec.EncodeDisplacementVertices(records)
	return file.commitOrdinary(LumpDispVerts, raw, err)
}

// LeafWaterData returns an owned snapshot of the decoded records.
func (file *File) LeafWaterData() ([]lumps.LeafWaterData, error) {
	decoded, err := file.decodeLump(LumpLeafWaterData)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.LeafWaterData).GetData(), nil
}

// SetLeafWaterData atomically replaces the payload, preserving its selected revision.
func (file *File) SetLeafWaterData(records []lumps.LeafWaterData) error {
	if err := file.checkTypedRevision(LumpLeafWaterData); err != nil {
		return err
	}
	raw, err := codec.EncodeLeafWaterData(records)
	return file.commitOrdinary(LumpLeafWaterData, raw, err)
}

// PrimitiveVertices returns an owned snapshot of the decoded records.
func (file *File) PrimitiveVertices() ([]lumps.PrimitiveVertex, error) {
	decoded, err := file.decodeLump(LumpPrimVerts)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.PrimitiveVertices).GetData(), nil
}

// SetPrimitiveVertices atomically replaces the payload, preserving its selected revision.
func (file *File) SetPrimitiveVertices(records []lumps.PrimitiveVertex) error {
	if err := file.checkTypedRevision(LumpPrimVerts); err != nil {
		return err
	}
	raw, err := codec.EncodePrimitiveVertices(records)
	return file.commitOrdinary(LumpPrimVerts, raw, err)
}

// PrimitiveIndices returns an owned snapshot of the decoded records.
func (file *File) PrimitiveIndices() ([]uint16, error) {
	decoded, err := file.decodeLump(LumpPrimIndices)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.PrimitiveIndices).GetData(), nil
}

// SetPrimitiveIndices atomically replaces the payload, preserving its selected revision.
func (file *File) SetPrimitiveIndices(records []uint16) error {
	if err := file.checkTypedRevision(LumpPrimIndices); err != nil {
		return err
	}
	raw, err := codec.EncodePrimitiveIndices(records)
	return file.commitOrdinary(LumpPrimIndices, raw, err)
}

// ClipPortalVertices returns an owned snapshot of the decoded records.
func (file *File) ClipPortalVertices() ([]mgl32.Vec3, error) {
	decoded, err := file.decodeLump(LumpClipPortalVerts)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.ClipPortalVertices).GetData(), nil
}

// SetClipPortalVertices atomically replaces the payload, preserving its selected revision.
func (file *File) SetClipPortalVertices(records []mgl32.Vec3) error {
	if err := file.checkTypedRevision(LumpClipPortalVerts); err != nil {
		return err
	}
	raw, err := codec.EncodeClipPortalVertices(records)
	return file.commitOrdinary(LumpClipPortalVerts, raw, err)
}

// TextureStringData returns the exact string bytes, including embedded and trailing NULs.
func (file *File) TextureStringData() (string, error) {
	decoded, err := file.decodeLump(LumpTexDataStringData)
	if err != nil {
		return "", err
	}
	return decoded.(*codec.TextureStringData).GetData(), nil
}

// SetTextureStringData atomically replaces the payload, preserving its selected revision.
func (file *File) SetTextureStringData(records string) error {
	if err := file.checkTypedRevision(LumpTexDataStringData); err != nil {
		return err
	}
	raw, err := codec.EncodeTextureStringData(records)
	return file.commitOrdinary(LumpTexDataStringData, raw, err)
}

// TextureStringTable returns an owned snapshot of the decoded records.
func (file *File) TextureStringTable() ([]int32, error) {
	decoded, err := file.decodeLump(LumpTexDataStringTable)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.TextureStringTable).GetData(), nil
}

// SetTextureStringTable atomically replaces the payload, preserving its selected revision.
func (file *File) SetTextureStringTable(records []int32) error {
	if err := file.checkTypedRevision(LumpTexDataStringTable); err != nil {
		return err
	}
	raw, err := codec.EncodeTextureStringTable(records)
	return file.commitOrdinary(LumpTexDataStringTable, raw, err)
}

// LeafMinDistancesToWater returns an owned snapshot of the decoded records.
func (file *File) LeafMinDistancesToWater() ([]uint16, error) {
	decoded, err := file.decodeLump(LumpLeafMinDistToWater)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.LeafMinDistancesToWater).GetData(), nil
}

// SetLeafMinDistancesToWater atomically replaces the payload, preserving its selected revision.
func (file *File) SetLeafMinDistancesToWater(records []uint16) error {
	if err := file.checkTypedRevision(LumpLeafMinDistToWater); err != nil {
		return err
	}
	raw, err := codec.EncodeLeafMinDistancesToWater(records)
	return file.commitOrdinary(LumpLeafMinDistToWater, raw, err)
}

// FaceMacroTextureInfo returns an owned snapshot of the decoded records.
func (file *File) FaceMacroTextureInfo() ([]lumps.FaceMacroTextureInfo, error) {
	decoded, err := file.decodeLump(LumpFaceMacroTextureInfo)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.FaceMacroTextureInfo).GetData(), nil
}

// SetFaceMacroTextureInfo atomically replaces the payload, preserving its selected revision.
func (file *File) SetFaceMacroTextureInfo(records []lumps.FaceMacroTextureInfo) error {
	if err := file.checkTypedRevision(LumpFaceMacroTextureInfo); err != nil {
		return err
	}
	raw, err := codec.EncodeFaceMacroTextureInfo(records)
	return file.commitOrdinary(LumpFaceMacroTextureInfo, raw, err)
}

// DisplacementTriangles returns an owned snapshot of the decoded records.
func (file *File) DisplacementTriangles() ([]lumps.DisplacementTriangle, error) {
	decoded, err := file.decodeLump(LumpDispTris)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.DisplacementTriangles).GetData(), nil
}

// SetDisplacementTriangles atomically replaces the payload, preserving its selected revision.
func (file *File) SetDisplacementTriangles(records []lumps.DisplacementTriangle) error {
	if err := file.checkTypedRevision(LumpDispTris); err != nil {
		return err
	}
	raw, err := codec.EncodeDisplacementTriangles(records)
	return file.commitOrdinary(LumpDispTris, raw, err)
}

// OverlayFades returns an owned snapshot of the decoded records.
func (file *File) OverlayFades() ([]lumps.OverlayFade, error) {
	decoded, err := file.decodeLump(LumpOverlayFades)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.OverlayFades).GetData(), nil
}

// SetOverlayFades atomically replaces the payload, preserving its selected revision.
func (file *File) SetOverlayFades(records []lumps.OverlayFade) error {
	if err := file.checkTypedRevision(LumpOverlayFades); err != nil {
		return err
	}
	raw, err := codec.EncodeOverlayFades(records)
	return file.commitOrdinary(LumpOverlayFades, raw, err)
}
