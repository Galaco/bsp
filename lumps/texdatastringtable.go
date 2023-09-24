package lumps

import (
	"bytes"
	"encoding/binary"
)

// TexDataStringTable is Lump 44: TexDataStringTable
type TexDataStringTable struct {
	Metadata
	data []int32 // MAX_MAP_TEXINFO = 2048
}

// Unmarshall Imports this lump from raw byte data
func (lump *TexDataStringTable) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]int32, length/4)
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return nil
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
