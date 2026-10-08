package lumps

import (
	"bytes"
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	primitives "github.com/galaco/bsp/primitives/face"
)

// FaceHDR is Lump 58: FaceHDR
type FaceHDR struct {
	Generic
	data []primitives.Face
}

// Unmarshall Imports this lump from raw byte data
func (lump *FaceHDR) Unmarshall(raw []byte) error {
	count, err := format.FixedRecordCount(len(raw), 56)
	if err != nil {
		return err
	}
	data := make([]primitives.Face, count)
	if err := binary.Read(bytes.NewReader(raw), binary.LittleEndian, data); err != nil {
		return err
	}
	lump.data = data
	lump.Metadata.SetLength(len(raw))
	return nil
}

// GetData gets internal format structure data
func (lump *FaceHDR) GetData() []primitives.Face {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *FaceHDR) Marshall() ([]byte, error) {
	length, err := format.FixedRecordBytes(len(lump.data), 56)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(make([]byte, 0, length))
	if err := binary.Write(buf, binary.LittleEndian, lump.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
