package lumps

import (
	"bytes"
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	primitives "github.com/galaco/bsp/primitives/faceid"
)

// FaceId is Lump 11: FaceIds
type FaceId struct {
	Generic
	data []primitives.FaceId
}

// Unmarshall Imports this lump from raw byte data
func (lump *FaceId) Unmarshall(raw []byte) error {
	count, err := format.FixedRecordCount(len(raw), 2)
	if err != nil {
		return err
	}
	data := make([]primitives.FaceId, count)
	if err := binary.Read(bytes.NewReader(raw), binary.LittleEndian, data); err != nil {
		return err
	}
	lump.data = data
	lump.Metadata.SetLength(len(raw))
	return nil
}

// GetData gets internal format structure data
func (lump *FaceId) GetData() []primitives.FaceId {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *FaceId) Marshall() ([]byte, error) {
	length, err := format.FixedRecordBytes(len(lump.data), 2)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(make([]byte, 0, length))
	if err := binary.Write(buf, binary.LittleEndian, lump.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
