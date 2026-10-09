package bsp

import (
	"fmt"

	"github.com/galaco/bsp/internal/codec"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

func (bsp *File) commitOrdinary(id LumpID, raw []byte, err error) error {
	if err != nil {
		return &FormatError{fmt.Sprintf("encode lump %d", id), err}
	}
	descriptor := bsp.header.Lumps[id]
	descriptor.FourCC = [4]byte{}
	bsp.storeOwned(id, descriptor, raw)
	return nil
}

// Leaves returns owned records in the descriptor-selected inline/compact layout.
func (bsp *File) Leaves() ([]lumps.Leaf, error) {
	decoded, err := bsp.decodeLump(LumpLeafs)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.Leaf).GetData(), nil
}

// SetLeaves preserves the selected revision and rejects incompatible lighting.
func (bsp *File) SetLeaves(records []lumps.Leaf) error {
	if err := bsp.checkTypedRevision(LumpLeafs); err != nil {
		return err
	}
	context := bsp.codecContext(LumpLeafs)
	raw, err := codec.EncodeLeaves(records, context)
	if err != nil {
		return bsp.commitOrdinary(LumpLeafs, nil, err)
	}
	if context.LeafVersion == 1 {
		for _, channel := range [][2]LumpID{{LumpLeafAmbientIndex, LumpLeafAmbientLighting}, {LumpLeafAmbientIndexHDR, LumpLeafAmbientLightingHDR}} {
			sampleContext := bsp.codecContext(channel[1])
			if err := bsp.validateAmbientChannel(channel[0], channel[1], nil, nil, sampleContext, len(records)); err != nil {
				return err
			}
		}
	}
	return bsp.commitOrdinary(LumpLeafs, raw, nil)
}

// WorldLights returns the LDR channel without implicit HDR fallback.
func (bsp *File) WorldLights() ([]lumps.WorldLight, error) {
	decoded, err := bsp.decodeLump(LumpWorldLights)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.WorldLight).GetData(), nil
}

// WorldLightsHDR returns the separate HDR worldlight channel.
func (bsp *File) WorldLightsHDR() ([]lumps.WorldLight, error) {
	decoded, err := bsp.decodeLump(LumpWorldLightsHDR)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.WorldLightHDR).GetData(), nil
}

// SetWorldLights atomically commits the selected LDR layout.
func (bsp *File) SetWorldLights(records []lumps.WorldLight) error {
	return bsp.setWorldLights(LumpWorldLights, records)
}

// SetWorldLightsHDR atomically commits the selected HDR layout.
func (bsp *File) SetWorldLightsHDR(records []lumps.WorldLight) error {
	return bsp.setWorldLights(LumpWorldLightsHDR, records)
}
func (bsp *File) setWorldLights(id LumpID, records []lumps.WorldLight) error {
	if err := bsp.checkTypedRevision(id); err != nil {
		return err
	}
	raw, err := codec.EncodeWorldLights(records, bsp.codecContext(id))
	return bsp.commitOrdinary(id, raw, err)
}

// BrushSides returns profile-selected eight-byte brush side records.
func (bsp *File) BrushSides() ([]lumps.BrushSide, error) {
	decoded, err := bsp.decodeLump(LumpBrushSides)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.BrushSide).GetData(), nil
}

// SetBrushSides rejects thin fields unsupported by the retained profile.
func (bsp *File) SetBrushSides(records []lumps.BrushSide) error {
	if err := bsp.checkTypedRevision(LumpBrushSides); err != nil {
		return err
	}
	raw, err := codec.EncodeBrushSides(records, bsp.codecContext(LumpBrushSides))
	return bsp.commitOrdinary(LumpBrushSides, raw, err)
}

// AmbientIndices returns the LDR index channel; zero-count entries reference leaves.
func (bsp *File) AmbientIndices() ([]lumps.AmbientIndex, error) {
	decoded, err := bsp.decodeLump(LumpLeafAmbientIndex)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.LeafAmbientIndex).GetData(), nil
}

// AmbientIndicesHDR returns the separate HDR ambient indices.
func (bsp *File) AmbientIndicesHDR() ([]lumps.AmbientIndex, error) {
	decoded, err := bsp.decodeLump(LumpLeafAmbientIndexHDR)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.LeafAmbientIndexHDR).GetData(), nil
}

// AmbientSamples returns standalone LDR samples, not leaf-v0 inline lighting.
func (bsp *File) AmbientSamples() ([]lumps.AmbientSample, error) {
	decoded, err := bsp.decodeLump(LumpLeafAmbientLighting)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.LeafAmbientLighting).GetData(), nil
}

// AmbientSamplesHDR returns standalone HDR samples without LDR fallback.
func (bsp *File) AmbientSamplesHDR() ([]lumps.AmbientSample, error) {
	decoded, err := bsp.decodeLump(LumpLeafAmbientLightingHDR)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.LeafAmbientLightingHDR).GetData(), nil
}

// SetAmbientIndices commits LDR indices without implicitly migrating sample layout.
func (bsp *File) SetAmbientIndices(records []lumps.AmbientIndex) error {
	return bsp.setAmbientIndices(LumpLeafAmbientIndex, records)
}

// SetAmbientIndicesHDR commits the separate HDR indices.
func (bsp *File) SetAmbientIndicesHDR(records []lumps.AmbientIndex) error {
	return bsp.setAmbientIndices(LumpLeafAmbientIndexHDR, records)
}
func (bsp *File) setAmbientIndices(id LumpID, records []lumps.AmbientIndex) error {
	if err := bsp.checkTypedRevision(id); err != nil {
		return err
	}
	context := bsp.codecContext(id)
	if id == LumpLeafAmbientIndex {
		context.LegacyAmbient = len(records) == 0
	}
	raw, err := codec.EncodeAmbientIndices(records, context)
	if err != nil {
		return bsp.commitOrdinary(id, nil, err)
	}
	sampleID := LumpLeafAmbientLighting
	if id == LumpLeafAmbientIndexHDR {
		sampleID = LumpLeafAmbientLightingHDR
	}
	sampleContext := bsp.codecContext(sampleID)
	sampleContext.LegacyAmbient = context.LegacyAmbient
	if err := bsp.validateAmbientChannel(id, sampleID, &records, nil, sampleContext, -1); err != nil {
		return err
	}
	if id == LumpLeafAmbientIndex {
		otherContext := bsp.codecContext(LumpLeafAmbientLightingHDR)
		otherContext.LegacyAmbient = context.LegacyAmbient
		if err := bsp.validateAmbientChannel(LumpLeafAmbientIndexHDR, LumpLeafAmbientLightingHDR, nil, nil, otherContext, -1); err != nil {
			return err
		}
	}
	return bsp.commitOrdinary(id, raw, nil)
}

// SetAmbientSamples commits the current standalone LDR sample layout.
func (bsp *File) SetAmbientSamples(records []lumps.AmbientSample) error {
	return bsp.setAmbientSamples(LumpLeafAmbientLighting, records)
}

// SetAmbientSamplesHDR commits the separate standalone HDR sample layout.
func (bsp *File) SetAmbientSamplesHDR(records []lumps.AmbientSample) error {
	return bsp.setAmbientSamples(LumpLeafAmbientLightingHDR, records)
}
func (bsp *File) setAmbientSamples(id LumpID, records []lumps.AmbientSample) error {
	if err := bsp.checkTypedRevision(id); err != nil {
		return err
	}
	context := bsp.codecContext(id)
	raw, err := codec.EncodeAmbientSamples(records, context)
	if err != nil {
		return bsp.commitOrdinary(id, nil, err)
	}
	indexID := LumpLeafAmbientIndex
	if id == LumpLeafAmbientLightingHDR {
		indexID = LumpLeafAmbientIndexHDR
	}
	if err := bsp.validateAmbientChannel(indexID, id, nil, &records, context, -1); err != nil {
		return err
	}
	return bsp.commitOrdinary(id, raw, nil)
}

// Occlusion returns owned occluders, polygons and vertex indices with checked references.
func (bsp *File) Occlusion() (lumps.Occlusion, error) {
	decoded, err := bsp.decodeLump(LumpOcclusion)
	if err != nil {
		return lumps.Occlusion{}, err
	}
	return decoded.(*codec.Occlusion).GetData(), nil
}

// SetOcclusion validates counted spans and references before committing any bytes.
func (bsp *File) SetOcclusion(value lumps.Occlusion) error {
	if err := bsp.checkTypedRevision(LumpOcclusion); err != nil {
		return err
	}
	if err := bsp.validateOcclusionReferences(value); err != nil {
		return err
	}
	raw, err := codec.EncodeOcclusion(value, bsp.codecContext(LumpOcclusion))
	return bsp.commitOrdinary(LumpOcclusion, raw, err)
}

func (bsp *File) logicalRecordCount(id LumpID, width int) (int, error) {
	raw, err := bsp.decodedLumpView(id)
	if err != nil {
		return 0, err
	}
	if len(raw)%width != 0 {
		return 0, &FormatError{fmt.Sprintf("count lump %d", id), fmt.Errorf("length %d is not divisible by %d", len(raw), width)}
	}
	return len(raw) / width, nil
}
func (bsp *File) leafCount() (int, error) {
	if err := bsp.checkTypedRevision(LumpLeafs); err != nil {
		return 0, err
	}
	width, err := codec.LeafRecordSize(bsp.codecContext(LumpLeafs))
	if err != nil {
		return 0, err
	}
	count, err := bsp.logicalRecordCount(LumpLeafs, width)
	if err != nil {
		return 0, err
	}
	if count > 65536 {
		return 0, &FormatError{"ambient leaves", fmt.Errorf("leaf count %d exceeds65536", count)}
	}
	return count, nil
}
func (bsp *File) validateOcclusionReferences(value lumps.Occlusion) error {
	planes, vertices := 0, 0
	var err error
	if len(value.Polygons) > 0 {
		planes, err = bsp.logicalRecordCount(LumpPlanes, 20)
		if err != nil {
			return err
		}
	}
	if len(value.VertexIndices) > 0 {
		vertices, err = bsp.logicalRecordCount(LumpVertexes, 12)
		if err != nil {
			return err
		}
	}
	if err := codec.ValidateOcclusionReferences(value, planes, vertices); err != nil {
		return &FormatError{"occlusion references", err}
	}
	return nil
}
func (bsp *File) validateAmbientChannel(indexID, sampleID LumpID, indices *[]lumps.AmbientIndex, samples *[]lumps.AmbientSample, context format.Context, leafCount int) error {
	var indexValues []lumps.AmbientIndex
	var sampleValues []lumps.AmbientSample
	var err error
	if indices != nil {
		indexValues = *indices
	} else {
		raw, readErr := bsp.decodedLumpView(indexID)
		if readErr != nil {
			return readErr
		}
		indexContext := bsp.codecContext(indexID)
		indexContext.LegacyAmbient = context.LegacyAmbient
		indexValues, err = codec.DecodeAmbientIndices(raw, indexContext)
		if err != nil {
			return &FormatError{fmt.Sprintf("ambient indices %d", indexID), err}
		}
	}
	if samples != nil {
		sampleValues = *samples
	} else {
		raw, readErr := bsp.decodedLumpView(sampleID)
		if readErr != nil {
			return readErr
		}
		sampleValues, err = codec.DecodeAmbientSamples(raw, context)
		if err != nil {
			return &FormatError{fmt.Sprintf("ambient samples %d", sampleID), err}
		}
	}
	if len(indexValues) == 0 && len(sampleValues) == 0 {
		return nil
	}
	if leafCount < 0 {
		leafCount, err = bsp.leafCount()
		if err != nil {
			return err
		}
	}
	if err := codec.ValidateAmbient(indexValues, sampleValues, leafCount, context); err != nil {
		return &FormatError{fmt.Sprintf("ambient channel %d/%d", indexID, sampleID), err}
	}
	return nil
}

func (bsp *File) validateOrdinarySnapshot(id LumpID, decoded codec.Decoder) error {
	switch id {
	case LumpFaces, LumpOriginalFaces, LumpFacesHDR:
		return bsp.validateFacePrimitiveReferences(decoded.(*codec.Faces).GetData())
	case LumpDispInfo:
		return bsp.validateDisplacementReferences(decoded.(*codec.Displacements).GetData())
	case LumpWaterOverlays:
		return bsp.validateWaterOverlayReferences(decoded.(*codec.WaterOverlays).GetData())
	case LumpOverlaySystemLevels:
		return bsp.validateOverlaySystemLevelReferences(decoded.(*codec.OverlaySystemLevels).GetData())
	case LumpDispMultiBlend:
		return bsp.validateMultiblendReferences(decoded.(*codec.DisplacementMultiblends).GetData())
	case LumpOverlays:
		return bsp.validateOverlayReferences(decoded.(*codec.Overlays).GetData())
	case LumpPrimitives:
		return bsp.validatePrimitiveReferences(decoded.(*codec.Primitives).GetData())
	case LumpOcclusion:
		return bsp.validateOcclusionReferences(decoded.(*codec.Occlusion).GetData())
	case LumpLeafAmbientIndex:
		indices := decoded.(*codec.LeafAmbientIndex).GetData()
		return bsp.validateAmbientChannel(id, LumpLeafAmbientLighting, &indices, nil, bsp.codecContext(LumpLeafAmbientLighting), -1)
	case LumpLeafAmbientIndexHDR:
		indices := decoded.(*codec.LeafAmbientIndexHDR).GetData()
		return bsp.validateAmbientChannel(id, LumpLeafAmbientLightingHDR, &indices, nil, bsp.codecContext(LumpLeafAmbientLightingHDR), -1)
	case LumpLeafAmbientLighting:
		samples := decoded.(*codec.LeafAmbientLighting).GetData()
		return bsp.validateAmbientChannel(LumpLeafAmbientIndex, id, nil, &samples, bsp.codecContext(id), -1)
	case LumpLeafAmbientLightingHDR:
		samples := decoded.(*codec.LeafAmbientLightingHDR).GetData()
		return bsp.validateAmbientChannel(LumpLeafAmbientIndexHDR, id, nil, &samples, bsp.codecContext(id), -1)
	}
	return nil
}
