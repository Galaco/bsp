package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/dispinfo"
)

// DispInfo is Lump 26: DispInfo
type DispInfo struct {
	Metadata
	data []primitives.DispInfo
}

// Unmarshall Imports this lump from raw byte data
func (lump *DispInfo) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]primitives.DispInfo, length/int(unsafe.Sizeof(primitives.DispInfo{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return nil
}

// GetData gets internal format structure data
func (lump *DispInfo) GetData() []primitives.DispInfo {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *DispInfo) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}
