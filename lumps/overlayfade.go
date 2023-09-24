package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/overlayfade"
)

// OverlayFade is Lump 60: Overlayfades
type OverlayFade struct {
	Metadata
	data []primitives.OverlayFade
}

// Unmarshall Imports this lump from raw byte data
func (lump *OverlayFade) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.OverlayFade, length/int(unsafe.Sizeof(primitives.OverlayFade{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// GetData gets internal format structure data
func (lump *OverlayFade) GetData() []primitives.OverlayFade {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *OverlayFade) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}
