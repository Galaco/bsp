package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/cubemap"
	"log"
	"unsafe"
)

// Cubemap is Lump 42: Cubemaps
type Cubemap struct {
	LumpGeneric
	data []primitives.CubemapSample
}

// Unmarshall Imports this lump from raw byte data
func (lump *Cubemap) Unmarshall(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)
	if length == 0 {
		return
	}
	lump.data = make([]primitives.CubemapSample, length/int32(unsafe.Sizeof(primitives.CubemapSample{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
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
