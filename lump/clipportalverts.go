package lump

import (
	"github.com/go-gl/mathgl/mgl32"
)

// ClipPortalVerts is Lump 41: ClipPortalVerts
type ClipPortalVerts struct {
	Metadata
	Data []mgl32.Vec3 `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *ClipPortalVerts) FromBytes(raw []byte) error {
	meta, data, err := unmarshallBasicLump[mgl32.Vec3](raw)
	if err != nil {
		return err
	}
	lump.Metadata = meta
	lump.Data = data
	return nil
}

// Contents returns internal format structure Data
func (lump *ClipPortalVerts) Contents() []mgl32.Vec3 {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *ClipPortalVerts) ToBytes() ([]byte, error) {
	return marshallBasicLump(lump.Data)
}
