package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/mapflags"
	"log"
)

// Lump 59: MapFlags
type MapFlags struct {
	LumpGeneric
	data primitives.MapFlags
}

// Import this lump from raw byte data
func (lump *MapFlags) FromBytes(raw []byte, length int32) {
	if length == 0 {
		return
	}

	err := binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}
	lump.LumpInfo.SetLength(length)
}

// Get internal format structure data
func (lump *MapFlags) GetData() *primitives.MapFlags {
	return &lump.data
}

// Dump this lump back to raw byte data
func (lump *MapFlags) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes()
}
