package lumps

import (
	"bytes"
	"encoding/binary"
	"log"
)

// Lump 44: TexDataStringTable
type TexDataStringTable struct {
	LumpGeneric
	data []int32 // MAX_MAP_TEXINFO = 2048
}

// Import this lump from raw byte data
func (lump *TexDataStringTable) FromBytes(raw []byte, length int32) {
	lump.data = make([]int32, length/4)
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

// Get internal format structure data
func (lump *TexDataStringTable) GetData() []int32 {
	return lump.data
}

// Dump this lump back to raw byte data
func (lump *TexDataStringTable) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}
