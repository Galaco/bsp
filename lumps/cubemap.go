package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/cubemap"
	"unsafe"
)

// Cubemap is Lump 42: Cubemaps
type Cubemap struct {
	Generic
	data []primitives.CubemapSample
}

// Unmarshall Imports this lump from raw byte data
func (lump *Cubemap) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.CubemapSample, length/int(unsafe.Sizeof(primitives.CubemapSample{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}

	return nil
}

// GetData gets internal format structure data
func (lump *Cubemap) GetData() []primitives.CubemapSample {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *Cubemap) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}
