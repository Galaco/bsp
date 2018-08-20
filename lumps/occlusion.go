package lumps

import (
	primitives "github.com/galaco/bsp/primitives/occlusion"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)

/**
	Lump 9: Occlusion
 */
type Occlusion struct {
	LumpGeneric
	Count int32
	Data []primitives.OcclusionData // len(slice) = Count
	PolyDataCount int32
	PolyData []primitives.OcclusionPolyData //len(slice) = PolyDataCount
	VertexIndexCount int32
	VertexIndices []int32 //len(slice) = VertexIndexCount
}

func (lump *Occlusion) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}
	offset := 0
	// data
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.Count)
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

	return lump
}

func (lump *Occlusion) GetData() interface{} {
	return &lump
}

func (lump *Occlusion) ToBytes() []byte {
	var buf bytes.Buffer
	// write data
	binary.Write(&buf, binary.LittleEndian, lump.Count)
	for _,data := range lump.Data {
		binary.Write(&buf, binary.LittleEndian, data)
	}

	// write polydata
	binary.Write(&buf, binary.LittleEndian, lump.PolyDataCount)
	for _,data := range lump.PolyData {
		binary.Write(&buf, binary.LittleEndian, data)
	}

	// write indices
	binary.Write(&buf, binary.LittleEndian, lump.VertexIndexCount)
	for _,data := range lump.VertexIndices {
		binary.Write(&buf, binary.LittleEndian, data)
	}
	return buf.Bytes()
}
