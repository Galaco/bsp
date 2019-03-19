package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/occlusion"
	"log"
	"unsafe"
)

// Occlusion is Lump 9: Occlusion
type Occlusion struct {
	LumpGeneric
	Count            int32
	Data             []primitives.OcclusionData // len(slice) = Count
	PolyDataCount    int32
	PolyData         []primitives.OcclusionPolyData //len(slice) = PolyDataCount
	VertexIndexCount int32
	VertexIndices    []int32 //len(slice) = VertexIndexCount
}

// Unmarshall Imports this lump from raw byte data
func (lump *Occlusion) Unmarshall(raw []byte, length int32) {
	if length == 0 {
		return
	}
	offset := 0
	// data
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.Count)
	if err != nil {
		log.Fatal(err)
	}
	offset += 4
	lump.Data = make([]primitives.OcclusionData, lump.Count)
	err = binary.Read(bytes.NewBuffer(raw[offset:]), binary.LittleEndian, &lump.Data)
	if err != nil {
		log.Fatal(err)
	}
	offset += int(unsafe.Sizeof(primitives.OcclusionData{})) * int(lump.Count)

	// polydata
	err = binary.Read(bytes.NewBuffer(raw[offset:]), binary.LittleEndian, &lump.PolyDataCount)
	if err != nil {
		log.Fatal(err)
	}
	offset += 4
	lump.PolyData = make([]primitives.OcclusionPolyData, lump.PolyDataCount)
	err = binary.Read(bytes.NewBuffer(raw[offset:]), binary.LittleEndian, &lump.PolyData)
	if err != nil {
		log.Fatal(err)
	}
	offset += int(unsafe.Sizeof(primitives.OcclusionPolyData{})) * int(lump.PolyDataCount)

	// vertexdata
	err = binary.Read(bytes.NewBuffer(raw[offset:]), binary.LittleEndian, &lump.VertexIndexCount)
	if err != nil {
		log.Fatal(err)
	}
	offset += 4
	lump.VertexIndices = make([]int32, lump.VertexIndexCount)
	err = binary.Read(bytes.NewBuffer(raw[offset:]), binary.LittleEndian, &lump.VertexIndices)
	if err != nil {
		log.Fatal(err)
	}

	lump.LumpInfo.SetLength(length)
}

// GetData gets internal format structure data
func (lump *Occlusion) GetData() *Occlusion {
	return lump
}

// Marshall dumps this lump back to raw byte data
func (lump *Occlusion) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	// write data
	err := binary.Write(&buf, binary.LittleEndian, lump.Count)
	if err != nil {
		return nil, err
	}
	for _, data := range lump.Data {
		if err = binary.Write(&buf, binary.LittleEndian, data); err != nil {
			return nil, err
		}
	}

	// write polydata
	if err = binary.Write(&buf, binary.LittleEndian, lump.PolyDataCount); err != nil {
		return nil, err
	}
	for _, data := range lump.PolyData {
		if err = binary.Write(&buf, binary.LittleEndian, data); err != nil {
			return nil, err
		}
	}

	// write indices
	err = binary.Write(&buf, binary.LittleEndian, lump.VertexIndexCount)
	for _, data := range lump.VertexIndices {
		if err = binary.Write(&buf, binary.LittleEndian, data); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), err
}
