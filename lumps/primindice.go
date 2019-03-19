package lumps

import (
	"bytes"
	"encoding/binary"
	"log"
)

// PrimIndice is Lump 39: PrimIndice
type PrimIndice struct {
	LumpGeneric
	data []uint16
}

// Unmarshall Imports this lump from raw byte data
func (lump *PrimIndice) Unmarshall(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]uint16, length/int32(2))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
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
