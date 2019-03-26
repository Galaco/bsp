package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/worldlight"
	"unsafe"
)

// WorldLightHDR is Lump 15: Worldlight
type WorldLightHDR struct {
	Generic
	data []primitives.WorldLight
}

// Unmarshall Imports this lump from raw byte data
func (lump *WorldLightHDR) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.WorldLight, length/int(unsafe.Sizeof(primitives.WorldLight{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	return err
}

// GetData gets internal format structure data
func (lump *WorldLightHDR) GetData() []primitives.WorldLight {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *WorldLightHDR) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}
