package bsp

import (
	"fmt"

	"github.com/galaco/bsp/internal/codec"
	"github.com/galaco/bsp/lumps"
	"github.com/go-gl/mathgl/mgl32"
)

// WaterOverlays returns complete 1120-byte records with checked active faces.
func (bsp *Bsp) WaterOverlays() ([]lumps.WaterOverlay, error) {
	decoded, err := bsp.Lump(LumpWaterOverlays)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.WaterOverlays).GetData(), nil
}

// SetWaterOverlays commits complete records after active reference validation.
func (bsp *Bsp) SetWaterOverlays(values []lumps.WaterOverlay) error {
	if err := bsp.checkTypedRevision(LumpWaterOverlays); err != nil {
		return err
	}
	if err := bsp.validateWaterOverlayReferences(values); err != nil {
		return err
	}
	raw, err := codec.EncodeWaterOverlays(values, bsp.codecContext(LumpWaterOverlays))
	return bsp.commitOrdinary(LumpWaterOverlays, raw, err)
}
func (bsp *Bsp) validateWaterOverlayReferences(values []lumps.WaterOverlay) error {
	if len(values) == 0 {
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
	for i := range values {
		record := &values[i]
		if record.FaceCount() > 256 {
			return &FormatError{"water overlay references", fmt.Errorf("record %d face count %d exceeds 256", i, record.FaceCount())}
		}
		if record.TexInfo < 0 || int(record.TexInfo) >= textures {
			return &FormatError{"water overlay references", fmt.Errorf("record %d texinfo %d outside %d records", i, record.TexInfo, textures)}
		}
		for j := 0; j < int(record.FaceCount()); j++ {
			face := record.Faces[j]
			if face < 0 || int64(face) >= int64(faces) {
				return &FormatError{"water overlay references", fmt.Errorf("record %d active face %d index %d outside %d records", i, j, face, faces)}
			}
		}
	}
	return nil
}

// OverlaySystemLevels returns explicit CPU/GPU byte codes; an empty optional
// channel stays absent rather than synthesizing one default per overlay.
func (bsp *Bsp) OverlaySystemLevels() ([]lumps.OverlaySystemLevel, error) {
	decoded, err := bsp.Lump(LumpOverlaySystemLevels)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.OverlaySystemLevels).GetData(), nil
}

// SetOverlaySystemLevels requires one explicit record per overlay, or absence.
func (bsp *Bsp) SetOverlaySystemLevels(values []lumps.OverlaySystemLevel) error {
	if err := bsp.checkTypedRevision(LumpOverlaySystemLevels); err != nil {
		return err
	}
	if err := bsp.validateOverlaySystemLevelReferences(values); err != nil {
		return err
	}
	raw, err := codec.EncodeOverlaySystemLevels(values, bsp.codecContext(LumpOverlaySystemLevels))
	return bsp.commitOrdinary(LumpOverlaySystemLevels, raw, err)
}
func (bsp *Bsp) validateOverlaySystemLevelReferences(values []lumps.OverlaySystemLevel) error {
	if len(values) == 0 {
		return nil
	}
	if err := bsp.checkTypedRevision(LumpOverlays); err != nil {
		return err
	}
	overlays, err := bsp.logicalRecordCount(LumpOverlays, 352)
	if err != nil {
		return err
	}
	if len(values) != overlays {
		return &FormatError{"overlay system levels", fmt.Errorf("got %d records for %d overlays", len(values), overlays)}
	}
	return nil
}

// Existing explicit levels constrain overlay count edits, but absent optional
// levels do not make otherwise valid overlays unreadable.
func (bsp *Bsp) validateOverlaySystemLevelCount(overlays int) error {
	if bsp.profile != ProfileCSGO && bsp.profile != ProfileL4D2 && bsp.profile != ProfileReactiveDrop {
		return nil
	}
	if len(bsp.lumps[LumpOverlaySystemLevels].raw) == 0 {
		return nil
	}
	if err := bsp.checkTypedRevision(LumpOverlaySystemLevels); err != nil {
		return err
	}
	count, err := bsp.logicalRecordCount(LumpOverlaySystemLevels, 4)
	if err != nil {
		return err
	}
	if count != overlays {
		return &FormatError{"overlay system levels", fmt.Errorf("%d explicit records for %d overlays", count, overlays)}
	}
	return nil
}

// Collision arrays expose verified scalar wire records only. Public sources
// do not establish model-dictionary linkage or blob triangle index units.
func (bsp *Bsp) PropCollisions() ([]lumps.PropCollision, error) {
	decoded, err := bsp.Lump(LumpPropCollision)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.PropCollisions).GetData(), nil
}

// SetPropCollisions commits scalar records without guessing dictionary linkage.
func (bsp *Bsp) SetPropCollisions(values []lumps.PropCollision) error {
	if err := bsp.checkTypedRevision(LumpPropCollision); err != nil {
		return err
	}
	raw, err := codec.EncodePropCollisions(values, bsp.codecContext(LumpPropCollision))
	return bsp.commitOrdinary(LumpPropCollision, raw, err)
}

// PropHulls returns owned RD hull scalars, including all contents bits.
func (bsp *Bsp) PropHulls() ([]lumps.PropHull, error) {
	decoded, err := bsp.Lump(LumpPropHulls)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.PropHulls).GetData(), nil
}

// SetPropHulls commits complete wire fields without interpreting surface lookup.
func (bsp *Bsp) SetPropHulls(values []lumps.PropHull) error {
	if err := bsp.checkTypedRevision(LumpPropHulls); err != nil {
		return err
	}
	raw, err := codec.EncodePropHulls(values, bsp.codecContext(LumpPropHulls))
	return bsp.commitOrdinary(LumpPropHulls, raw, err)
}

// PropHullVertices returns owned RD float3 values with original float bits.
func (bsp *Bsp) PropHullVertices() ([]mgl32.Vec3, error) {
	decoded, err := bsp.Lump(LumpPropHullVerts)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.PropHullVertices).GetData(), nil
}

// SetPropHullVertices commits the verified RD vertex array.
func (bsp *Bsp) SetPropHullVertices(values []mgl32.Vec3) error {
	if err := bsp.checkTypedRevision(LumpPropHullVerts); err != nil {
		return err
	}
	raw, err := codec.EncodePropHullVertices(values, bsp.codecContext(LumpPropHullVerts))
	return bsp.commitOrdinary(LumpPropHullVerts, raw, err)
}

// PropHullTriangles returns signed index span scalars; index units remain opaque.
func (bsp *Bsp) PropHullTriangles() ([]lumps.PropHullTriangles, error) {
	decoded, err := bsp.Lump(LumpPropTris)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.PropHullTriangles).GetData(), nil
}

// SetPropHullTriangles commits scalars without inventing blob-relative bounds.
func (bsp *Bsp) SetPropHullTriangles(values []lumps.PropHullTriangles) error {
	if err := bsp.checkTypedRevision(LumpPropTris); err != nil {
		return err
	}
	raw, err := codec.EncodePropHullTriangles(values, bsp.codecContext(LumpPropTris))
	return bsp.commitOrdinary(LumpPropTris, raw, err)
}

// PropBlob exposes RD's defined byte vector without interpreting its interior.
func (bsp *Bsp) PropBlob() ([]byte, error) {
	if err := bsp.checkTypedRevision(LumpPropBlob); err != nil {
		return nil, err
	}
	return bsp.DecodedLump(LumpPropBlob)
}

// SetPropBlob commits owned uncompressed bytes without inspecting its interior.
func (bsp *Bsp) SetPropBlob(raw []byte) error {
	if err := bsp.checkTypedRevision(LumpPropBlob); err != nil {
		return err
	}
	descriptor := bsp.header.Lumps[LumpPropBlob]
	descriptor.FourCC = [4]byte{}
	return bsp.SetRawLump(LumpPropBlob, descriptor, raw)
}
