package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/texdata"
	"unsafe"
)

// TexData is Lump 2: TexData
type TexData struct {
	Generic
	data []primitives.TexData
}

// Unmarshall Imports this lump from raw byte data
func (lump *TexData) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]primitives.TexData, length/int(unsafe.Sizeof(primitives.TexData{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return err
}

// GetData gets internal format structure data
func (lump *TexData) GetData() []primitives.TexData {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *TexData) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}
