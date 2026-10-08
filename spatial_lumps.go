package bsp

import (
	"encoding/binary"
	"fmt"

	"github.com/galaco/bsp/internal/codec"
	"github.com/galaco/bsp/lumps"
	facevalue "github.com/galaco/bsp/primitives/face"
)

// Overlays returns complete owned records, checking active faces and materials.
func (bsp *Bsp) Overlays() ([]lumps.Overlay, error) {
	decoded, err := bsp.Lump(LumpOverlays)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.Overlays).GetData(), nil
}

// SetOverlays preserves inactive slots and commits only representable references.
func (bsp *Bsp) SetOverlays(records []lumps.Overlay) error {
	if err := bsp.checkTypedRevision(LumpOverlays); err != nil {
		return err
	}
	if err := bsp.validateOverlayReferences(records); err != nil {
		return err
	}
	raw, err := codec.EncodeOverlays(records, bsp.codecContext(LumpOverlays))
	return bsp.commitOrdinary(LumpOverlays, raw, err)
}

// Cubemaps returns snapped integer origins and unchanged wire size/padding codes.
func (bsp *Bsp) Cubemaps() ([]lumps.Cubemap, error) {
	decoded, err := bsp.Lump(LumpCubemaps)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.Cubemaps).GetData(), nil
}

// SetCubemaps atomically commits signed coordinates and preserved padding.
func (bsp *Bsp) SetCubemaps(records []lumps.Cubemap) error {
	if err := bsp.checkTypedRevision(LumpCubemaps); err != nil {
		return err
	}
	raw, err := codec.EncodeCubemaps(records, bsp.codecContext(LumpCubemaps))
	return bsp.commitOrdinary(LumpCubemaps, raw, err)
}

// Primitives returns ten-byte records with checked index and vertex spans.
func (bsp *Bsp) Primitives() ([]lumps.Primitive, error) {
	decoded, err := bsp.Lump(LumpPrimitives)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.Primitives).GetData(), nil
}

// SetPrimitives preserves descriptor revision and validates companion geometry.
func (bsp *Bsp) SetPrimitives(records []lumps.Primitive) error {
	if err := bsp.checkTypedRevision(LumpPrimitives); err != nil {
		return err
	}
	for i := range records {
		if records[i].Type > 1 {
			return fmt.Errorf("%w: primitive %d type %d", ErrUnsupportedLump, i, records[i].Type)
		}
	}
	if err := bsp.validatePrimitiveReferences(records); err != nil {
		return err
	}
	raw, err := codec.EncodePrimitives(records, bsp.codecContext(LumpPrimitives))
	return bsp.commitOrdinary(LumpPrimitives, raw, err)
}

func (bsp *Bsp) validateOverlayReferences(records []lumps.Overlay) error {
	if err := bsp.validateOverlaySystemLevelCount(len(records)); err != nil {
		return err
	}
	if len(records) == 0 {
		return nil
	}
	for _, id := range []LumpId{LumpFaces, LumpTexInfo} {
		if err := bsp.checkTypedRevision(id); err != nil {
			return err
		}
	}
	faces, err := bsp.logicalRecordCount(LumpFaces, 56)
	if err != nil {
		return err
	}
	textures, err := bsp.logicalRecordCount(LumpTexInfo, 72)
	if err != nil {
		return err
	}
	for i := range records {
		record := &records[i]
		if record.FaceCount() > 64 {
			return &FormatError{"overlay references", fmt.Errorf("overlay %d face count %d exceeds64", i, record.FaceCount())}
		}
		// SDK Overlay_EmitOverlayFace always creates a material texinfo. Unlike
		// ordinary face texinfo, no negative overlay sentinel is documented.
		if record.TexInfo < 0 || int(record.TexInfo) >= textures {
			return &FormatError{"overlay references", fmt.Errorf("overlay %d texinfo %d outside %d records", i, record.TexInfo, textures)}
		}
		for j := 0; j < int(record.FaceCount()); j++ {
			face := record.Faces[j]
			if face < 0 || int64(face) >= int64(faces) {
				return &FormatError{"overlay references", fmt.Errorf("overlay %d active face %d references %d outside %d records", i, j, face, faces)}
			}
		}
	}
	return nil
}

func (bsp *Bsp) validatePrimitiveReferences(records []lumps.Primitive) error {
	if len(records) == 0 {
		return nil
	}
	for _, id := range []LumpId{LumpPrimVerts, LumpPrimIndices} {
		if err := bsp.checkTypedRevision(id); err != nil {
			return err
		}
	}
	vertices, err := bsp.logicalRecordCount(LumpPrimVerts, 12)
	if err != nil {
		return err
	}
	indices, err := bsp.decodedLumpView(LumpPrimIndices)
	if err != nil {
		return err
	}
	if len(indices)%2 != 0 {
		return &FormatError{"primitive indices", fmt.Errorf("length %d is not divisible by2", len(indices))}
	}
	faceLocal := false
	for i := range records {
		record := &records[i]
		first, count := int(record.FirstIndex), int(record.IndexCount)
		start, verts := int(record.FirstVert), int(record.VertCount)
		if first > len(indices)/2 || count > len(indices)/2-first {
			return &FormatError{"primitive references", fmt.Errorf("primitive %d index span %d+%d outside %d", i, first, count, len(indices)/2)}
		}
		if start > vertices || verts > vertices-start {
			return &FormatError{"primitive references", fmt.Errorf("primitive %d vertex span %d+%d outside %d", i, start, verts, vertices)}
		}
		if verts == 0 {
			faceLocal = faceLocal || count > 0
			continue
		}
		// Source vbsp/faces.cpp AddWindingToPrimverts emits absolute ID38
		// indices, not offsets relative to FirstVert.
		for j := first; j < first+count; j++ {
			index := int(binary.LittleEndian.Uint16(indices[j*2:]))
			if index < start || index >= start+verts {
				return &FormatError{"primitive references", fmt.Errorf("primitive %d index %d references %d outside [%d,%d)", i, j-first, index, start, start+verts)}
			}
		}
	}
	if !faceLocal {
		return nil
	}
	if err := bsp.checkTypedRevision(LumpFaces); err != nil {
		return err
	}
	faces, err := bsp.decodedLumpView(LumpFaces)
	if err != nil {
		return err
	}
	if len(faces)%56 != 0 {
		return &FormatError{"primitive owning faces", fmt.Errorf("length %d is not divisible by56", len(faces))}
	}
	// Zero-VertCount T-junction primitives index the owning face's winding.
	// Check each actual owner; unreferenced records remain unreferenced, not
	// interpreted as indices into an absent ID38 array. No HDR fallback.
	for i := 0; i < len(faces)/56; i++ {
		face := faces[i*56 : (i+1)*56]
		first := int(binary.LittleEndian.Uint16(face[50:]))
		owner := facevalue.Face{NumPrims: binary.LittleEndian.Uint16(face[48:])}
		count := int(owner.PrimitiveCount())
		if count == 0 {
			continue
		}
		if first > len(records) || count > len(records)-first {
			return &FormatError{"primitive owning faces", fmt.Errorf("face %d primitive span %d+%d outside %d", i, first, count, len(records))}
		}
		edges := int(int16(binary.LittleEndian.Uint16(face[8:])))
		for j := first; j < first+count; j++ {
			record := &records[j]
			if record.VertCount != 0 {
				continue
			}
			for k := int(record.FirstIndex); k < int(record.FirstIndex)+int(record.IndexCount); k++ {
				index := int(binary.LittleEndian.Uint16(indices[k*2:]))
				if index >= edges {
					return &FormatError{"primitive face-local indices", fmt.Errorf("face %d primitive %d index %d outside %d winding vertices", i, j, index, edges)}
				}
			}
		}
	}
	return nil
}
