package lumps

import (
datatypes "github.com/galaco/bsp/lumps/datatypes/leafambientlighting"
"unsafe"
"encoding/binary"
"bytes"
"log"
)
/**
	Lump n:
 */
type LeafAmbientLighting struct {
	LumpInfo
	data []datatypes.LeafAmbientLighting
}

func (lump LeafAmbientLighting) FromBytes(raw []byte, length int32) ILump {
	if length == 0 {
		return lump
	}
	lump.data = make([]datatypes.LeafAmbientLighting, length/int32(unsafe.Sizeof(datatypes.LeafAmbientLighting{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump LeafAmbientLighting) GetData() interface{} {
	return lump.data
}

func (lump LeafAmbientLighting) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}
