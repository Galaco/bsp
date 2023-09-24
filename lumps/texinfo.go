package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/texinfo"
)

// TexInfo is Lump 6: TexInfo
type TexInfo struct {
	Metadata
	data []primitives.TexInfo
}

// Unmarshall Imports this lump from raw byte data
func (lump *TexInfo) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.data = make([]primitives.TexInfo, length/int(unsafe.Sizeof(primitives.TexInfo{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return err
}

// GetData gets internal format structure data
func (lump *TexInfo) GetData() []primitives.TexInfo {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *TexInfo) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}
