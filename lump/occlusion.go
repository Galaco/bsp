package lump

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/lump/primitive/occlusion"
)

// Occlusion is Lump 9: Occlusion
type Occlusion struct {
	Metadata
	Count            int32                          `json:"count"`
	Data             []primitives.OcclusionData     `json:"data"` // len(slice) = Count
	PolyDataCount    int32                          `json:"polyDataCount"`
	PolyData         []primitives.OcclusionPolyData `json:"polyData"` //len(slice) = PolyDataCount
	VertexIndexCount int32                          `json:"VertexIndexCount"`
	VertexIndices    []int32                        `json:"VertexIndices"` //len(slice) = VertexIndexCount
}

// FromBytes imports this lump from raw byte Data
func (lump *Occlusion) FromBytes(raw []byte) error {
	if len(raw) == 0 {
		return nil
	}
	offset := 0
	// Data
	if err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.Count); err != nil {
		return err
	}
	offset += 4
	lump.Data = make([]primitives.OcclusionData, lump.Count)
	if err := binary.Read(bytes.NewBuffer(raw[offset:]), binary.LittleEndian, &lump.Data); err != nil {
		return err
	}
	offset += int(unsafe.Sizeof(primitives.OcclusionData{})) * int(lump.Count)

	// polydata
	if err := binary.Read(bytes.NewBuffer(raw[offset:]), binary.LittleEndian, &lump.PolyDataCount); err != nil {
		return err
	}
	offset += 4
	lump.PolyData = make([]primitives.OcclusionPolyData, lump.PolyDataCount)
	if err := binary.Read(bytes.NewBuffer(raw[offset:]), binary.LittleEndian, &lump.PolyData); err != nil {
		return err
	}
	offset += int(unsafe.Sizeof(primitives.OcclusionPolyData{})) * int(lump.PolyDataCount)

	// vertexdata
	if err := binary.Read(bytes.NewBuffer(raw[offset:]), binary.LittleEndian, &lump.VertexIndexCount); err != nil {
		return err
	}
	offset += 4
	lump.VertexIndices = make([]int32, lump.VertexIndexCount)
	if err := binary.Read(bytes.NewBuffer(raw[offset:]), binary.LittleEndian, &lump.VertexIndices); err != nil {
		return err
	}

	return nil
}

// Contents returns internal format structure Data
func (lump *Occlusion) Contents() *Occlusion {
	return lump
}

// ToBytes converts this lump back to raw byte Data
func (lump *Occlusion) ToBytes() ([]byte, error) {
	var buf bytes.Buffer

	// write Data
	if err := binary.Write(&buf, binary.LittleEndian, lump.Count); err != nil {
		return nil, err
	}
	for _, data := range lump.Data {
		if err := binary.Write(&buf, binary.LittleEndian, data); err != nil {
			return nil, err
		}
	}

	// write polydata
	if err := binary.Write(&buf, binary.LittleEndian, lump.PolyDataCount); err != nil {
		return nil, err
	}
	for _, data := range lump.PolyData {
		if err := binary.Write(&buf, binary.LittleEndian, data); err != nil {
			return nil, err
		}
	}

	// write indices
	if err := binary.Write(&buf, binary.LittleEndian, lump.VertexIndexCount); err != nil {
		return nil, err
	}
	for _, data := range lump.VertexIndices {
		if err := binary.Write(&buf, binary.LittleEndian, data); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
