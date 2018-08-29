package lumps

import (
	primitives "github.com/galaco/bsp/primitives/texinfo"
	"encoding/binary"
	"bytes"
	"log"
	"unsafe"
)

/**
	Lump 6: TexInfo
 */

type TexInfo struct {
	LumpGeneric
	data []primitives.TexInfo
}

func (lump *TexInfo) FromBytes(raw []byte, length int32) {
	lump.data = make([]primitives.TexInfo, length/int32(unsafe.Sizeof(primitives.TexInfo{})))
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

func (lump *TexInfo) GetData() []primitives.TexInfo {
	return lump.data
}

func (lump *TexInfo) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}
