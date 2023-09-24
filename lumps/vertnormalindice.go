package lumps

import (
	"bytes"
	"encoding/binary"
)

// VertNormalIndice is Lump 31: VertNormalIndice
type VertNormalIndice struct {
	Metadata
	data []uint16
}

// Unmarshall Imports this lump from raw byte data
func (lump *VertNormalIndice) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]uint16, length/2)
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	return err
}

// GetData gets internal format structure data
func (lump *VertNormalIndice) GetData() []uint16 {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *VertNormalIndice) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}
