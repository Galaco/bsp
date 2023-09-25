package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/cubemap"
)

// Cubemap is Lump 42: Cubemaps
type Cubemap struct {
	Metadata
	data []primitives.CubemapSample
}

// FromBytes imports this lump from raw byte data
func (lump *Cubemap) FromBytes(raw []byte) (err error) {
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

// Contents returns internal format structure data
func (lump *Cubemap) Contents() []primitives.CubemapSample {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Cubemap) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}
