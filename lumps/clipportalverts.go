package lumps

import (
	"github.com/go-gl/mathgl/mgl32"
)

// ClipPortalVerts is Lump 41: ClipPortalVerts
type ClipPortalVerts struct {
	Metadata
	data []mgl32.Vec3
}

// FromBytes imports this lump from raw byte data
func (lump *ClipPortalVerts) FromBytes(raw []byte) (err error) {
	meta, data, err := unmarshallBasicLump[mgl32.Vec3](raw)
	lump.Metadata = meta
	if err != nil {
		return err
	}

	lump.data = data
	return nil
}

// Contents returns internal format structure data
func (lump *ClipPortalVerts) Contents() []mgl32.Vec3 {
	return lump.data
}

// ToBytes converts this lump back to raw byte data
func (lump *ClipPortalVerts) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.data)
}
