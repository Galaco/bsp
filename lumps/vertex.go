package lumps

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Vertex is Lump 3: Vertex
type Vertex struct {
	Metadata
	data []mgl32.Vec3
}

// FromBytes imports this lump from raw byte data
func (lump *Vertex) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[mgl32.Vec3](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *Vertex) Contents() []mgl32.Vec3 {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *Vertex) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}
