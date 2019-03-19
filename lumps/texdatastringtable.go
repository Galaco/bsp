package lumps

import (
	"bytes"
	"encoding/binary"
	"log"
)

// TexDataStringTable is Lump 44: TexDataStringTable
type TexDataStringTable struct {
	LumpGeneric
	data []int32 // MAX_MAP_TEXINFO = 2048
}

// Unmarshall Imports this lump from raw byte data
func (lump *TexDataStringTable) Unmarshall(raw []byte, length int32) {
	lump.data = make([]int32, length/4)
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

// GetData gets internal format structure data
func (lump *TexDataStringTable) GetData() []int32 {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *TexDataStringTable) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}
