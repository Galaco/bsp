package lumps

import (
	"bytes"
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	"github.com/go-gl/mathgl/mgl32"
)

// ClipPortalVerts is Lump 41: ClipPortalVerts
type ClipPortalVerts struct {
	Generic
	data []mgl32.Vec3
}

// Unmarshall Imports this lump from raw byte data
func (lump *ClipPortalVerts) Unmarshall(raw []byte) error {
	count, err := format.FixedRecordCount(len(raw), 12)
	if err != nil {
		return err
	}
	data := make([]mgl32.Vec3, count)
	if err := binary.Read(bytes.NewReader(raw), binary.LittleEndian, data); err != nil {
		return err
	}
	lump.data = data
	lump.Metadata.SetLength(len(raw))
	return nil
}

// GetData gets internal format structure data
func (lump *ClipPortalVerts) GetData() []mgl32.Vec3 {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *ClipPortalVerts) Marshall() ([]byte, error) {
	length, err := format.FixedRecordBytes(len(lump.data), 12)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(make([]byte, 0, length))
	if err := binary.Write(buf, binary.LittleEndian, lump.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
