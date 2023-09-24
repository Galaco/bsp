package lumps

import (
	"bytes"
	"encoding/binary"
)

// PrimIndice is Lump 39: PrimIndice
type PrimIndice struct {
	Metadata
	data []uint16
}

// Unmarshall Imports this lump from raw byte data
func (lump *PrimIndice) Unmarshall(raw []byte) (err error) {
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
func (lump *PrimIndice) GetData() []uint16 {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *PrimIndice) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}
