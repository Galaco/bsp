package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/area"
)

// Area is Lump 20: Areas
type Area struct {
	Metadata
	data []primitives.Area
}

// Unmarshall Imports this lump from raw byte data
func (lump *Area) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return nil
	}

	lump.data = make([]primitives.Area, length/int(unsafe.Sizeof(primitives.Area{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return nil
}

// GetData gets internal format structure data
func (lump *Area) GetData() []primitives.Area {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *Area) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, lump.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
