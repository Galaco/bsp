package lumps

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	primitives "github.com/galaco/bsp/primitives/physcollide"
)

// PhysCollide is Lump 20: PhysCollide
type PhysCollide struct {
	Metadata
	data []primitives.PhysCollideEntry
}

// FromBytes imports this lump from raw byte data
func (lump *PhysCollide) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)
	if length == 0 {
		return
	}

	lump.data = make([]primitives.PhysCollideEntry, length/int(unsafe.Sizeof(primitives.PhysCollideEntry{})))
	err = binary.Read(bytes.NewBuffer(raw), binary.LittleEndian, &lump.data)

	return err
}

// Contents returns internal format structure data
func (lump *PhysCollide) Contents() []primitives.PhysCollideEntry {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *PhysCollide) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	for _, entry := range lump.data {
		err := binary.Write(&buf, binary.LittleEndian, entry.ModelHeader)
		if err != nil {
			return nil, err
		}
		for _, solid := range entry.Solids {
			if err = binary.Write(&buf, binary.LittleEndian, solid.Size); err != nil {
				return nil, err
			}
			if err = binary.Write(&buf, binary.LittleEndian, solid.CollisionBinary); err != nil {
				return nil, err
			}
		}
		if err = binary.Write(&buf, binary.LittleEndian, []byte(entry.TextBuffer)); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
