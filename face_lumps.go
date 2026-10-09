package bsp

import (
	"encoding/binary"
	"fmt"

	"github.com/galaco/bsp/internal/codec"
	"github.com/galaco/bsp/lumps"
)

// Faces returns owned ordinary face records; HDR never substitutes for LDR.
func (bsp *File) Faces() ([]lumps.Face, error) { return bsp.faces(LumpFaces) }

// OriginalFaces returns the independent pre-split face channel.
func (bsp *File) OriginalFaces() ([]lumps.Face, error) { return bsp.faces(LumpOriginalFaces) }

// FacesHDR returns only the independent HDR channel.
func (bsp *File) FacesHDR() ([]lumps.Face, error) { return bsp.faces(LumpFacesHDR) }
func (bsp *File) faces(id LumpID) ([]lumps.Face, error) {
	decoded, err := bsp.decodeLump(id)
	if err != nil {
		return nil, err
	}
	return decoded.(*codec.Faces).GetData(), nil
}

// SetFaces commits ordinary records after masked primitive-range validation.
func (bsp *File) SetFaces(values []lumps.Face) error { return bsp.setFaces(LumpFaces, values) }

// SetOriginalFaces commits the independent original-face channel.
func (bsp *File) SetOriginalFaces(values []lumps.Face) error {
	return bsp.setFaces(LumpOriginalFaces, values)
}

// SetFacesHDR commits HDR records without changing ordinary faces.
func (bsp *File) SetFacesHDR(values []lumps.Face) error { return bsp.setFaces(LumpFacesHDR, values) }
func (bsp *File) setFaces(id LumpID, values []lumps.Face) error {
	if err := bsp.checkTypedRevision(id); err != nil {
		return err
	}
	if err := bsp.validateFacePrimitiveReferences(values); err != nil {
		return err
	}
	raw, err := codec.EncodeFaces(values, bsp.codecContext(id))
	return bsp.commitOrdinary(id, raw, err)
}

func (bsp *File) validateFacePrimitiveReferences(values []lumps.Face) error { // Absent primitive references must not force interpretation of opaque peers.
	used := false
	for i := range values {
		if values[i].PrimitiveCount() != 0 {
			used = true
			break
		}
	}
	if !used {
		return nil
	}
	if err := bsp.checkTypedRevision(LumpPrimitives); err != nil {
		return err
	}
	raw, err := bsp.decodedLumpView(LumpPrimitives)
	if err != nil {
		return err
	}
	if len(raw)%10 != 0 {
		return &FormatError{"face primitive references", fmt.Errorf("primitive length %d is not divisible by 10", len(raw))}
	}
	total := len(raw) / 10
	var indices []byte
	indicesLoaded := false
	for i := range values {
		record := &values[i]
		count, first := int(record.PrimitiveCount()), int(record.FirstPrimID)
		if count == 0 {
			continue
		}
		if first > total || count > total-first {
			return &FormatError{"face primitive references", fmt.Errorf("face %d span %d+%d outside %d primitives", i, first, count, total)}
		}
		for j := first; j < first+count; j++ {
			primitive := raw[j*10 : (j+1)*10]
			if primitive[0] > 1 {
				return fmt.Errorf("%w: referenced primitive %d type %d", ErrUnsupportedLump, j, primitive[0])
			}
			// Positive vertex count uses absolute ID38 indices; zero count uses this
			// owner's local winding, not an absent vertex table or another channel.
			if binary.LittleEndian.Uint16(primitive[8:]) != 0 {
				continue
			}
			if !indicesLoaded {
				if err := bsp.checkTypedRevision(LumpPrimIndices); err != nil {
					return err
				}
				indices, err = bsp.decodedLumpView(LumpPrimIndices)
				if err != nil {
					return err
				}
				if len(indices)%2 != 0 {
					return &FormatError{"face primitive indices", fmt.Errorf("length %d is not divisible by 2", len(indices))}
				}
				indicesLoaded = true
			}
			start, n := int(binary.LittleEndian.Uint16(primitive[2:])), int(binary.LittleEndian.Uint16(primitive[4:]))
			if start > len(indices)/2 || n > len(indices)/2-start {
				return &FormatError{"face primitive indices", fmt.Errorf("face %d primitive %d span %d+%d outside %d indices", i, j, start, n, len(indices)/2)}
			}
			for k := start; k < start+n; k++ {
				index := int(binary.LittleEndian.Uint16(indices[k*2:]))
				if index >= int(record.NumEdges) {
					return &FormatError{"face-local primitive indices", fmt.Errorf("face %d primitive %d index %d outside %d winding edges", i, j, index, record.NumEdges)}
				}
			}
		}
	}
	return nil
}
