package lump

import (
	"bytes"
	"encoding/binary"

	primitives "github.com/galaco/bsp/lump/primitive/physcollide"
)

// PhysCollide is Lump 20: PhysCollide
type PhysCollide struct {
	Metadata
	Data []primitives.PhysCollideEntry `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *PhysCollide) FromBytes(raw []byte) error {
	meta, data, err := unmarshallBasicLump[primitives.PhysCollideEntry](raw)
	if err != nil {
		return err
	}
	lump.Metadata = meta
	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *PhysCollide) Contents() []primitives.PhysCollideEntry {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *PhysCollide) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	for _, entry := range lump.Data {
		if err := binary.Write(&buf, binary.LittleEndian, entry.ModelHeader); err != nil {
			return nil, err
		}
		for _, solid := range entry.Solids {
			if err := binary.Write(&buf, binary.LittleEndian, solid.Size); err != nil {
				return nil, err
			}
			if err := binary.Write(&buf, binary.LittleEndian, solid.CollisionBinary); err != nil {
				return nil, err
			}
		}
		if err := binary.Write(&buf, binary.LittleEndian, []byte(entry.TextBuffer)); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
