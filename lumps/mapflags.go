package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/mapflags"
)

// MapFlags is Lump 59: MapFlags
type MapFlags struct {
	Generic
	data primitives.MapFlags
}

// Unmarshall Imports this lump from raw byte data
func (lump *MapFlags) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	if length == 0 {
		return
	}

	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)
	if err != nil {
		return err
	}
	lump.Metadata.SetLength(length)

	return err
}

// GetData gets internal format structure data
func (lump *MapFlags) GetData() *primitives.MapFlags {
	return &lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *MapFlags) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.data)
	return buf.Bytes(), err
}
