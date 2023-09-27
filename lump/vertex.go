package lump

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Vertex is Lump 3: Vertex
type Vertex struct {
	Metadata
	Data []mgl32.Vec3 `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *Vertex) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[mgl32.Vec3](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *Vertex) Contents() []mgl32.Vec3 {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *Vertex) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}
