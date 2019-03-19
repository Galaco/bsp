package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/texinfo"
	"log"
	"unsafe"
)

// TexInfo is Lump 6: TexInfo
type TexInfo struct {
	LumpGeneric
	data []primitives.TexInfo
}

// Unmarshall Imports this lump from raw byte data
func (lump *TexInfo) Unmarshall(raw []byte, length int32) {
	lump.data = make([]primitives.TexInfo, length/int32(unsafe.Sizeof(primitives.TexInfo{})))
	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
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
