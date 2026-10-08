package lumps

import (
	"bytes"
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	primitives "github.com/galaco/bsp/primitives/node"
)

// Node is Lump 5: Node
type Node struct {
	Generic
	data []primitives.Node // MAP_MAX_NODES = 65536
}

// Unmarshall Imports this lump from raw byte data
func (lump *Node) Unmarshall(raw []byte) error {
	count, err := format.FixedRecordCount(len(raw), 32)
	if err != nil {
		return err
	}
	data := make([]primitives.Node, count)
	if err := binary.Read(bytes.NewReader(raw), binary.LittleEndian, data); err != nil {
		return err
	}
	lump.data = data
	lump.Metadata.SetLength(len(raw))
	return nil
}

// GetData gets internal format structure data
func (lump *Node) GetData() []primitives.Node {
	return lump.data
}

// Marshall dumps this lump back to raw byte data
func (lump *Node) Marshall() ([]byte, error) {
	length, err := format.FixedRecordBytes(len(lump.data), 32)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(make([]byte, 0, length))
	if err := binary.Write(buf, binary.LittleEndian, lump.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
