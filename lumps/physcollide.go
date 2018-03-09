package lumps

import (
	datatypes "github.com/galaco/bsp/lumps/datatypes/physcollide"
	"unsafe"
	"encoding/binary"
	"bytes"
	"log"
)
/**
	Lump 20: PhysCollide
 */
type PhysCollide struct {
	LumpInfo
	data []datatypes.PhysCollideEntry
}

func (lump PhysCollide) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}

	lump.data = make([]datatypes.PhysCollideEntry, length/int32(unsafe.Sizeof(datatypes.PhysCollideEntry{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump PhysCollide) GetData() interface{} {
	return lump.data
}

func (lump PhysCollide) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}
