package lumps

import (
	"encoding/binary"
	"fmt"

	primitives "github.com/galaco/bsp/primitives/leaf"
)

const (
	// MaxMapLeafs is the maximum number of leafs the engine can handle.
	MaxMapLeafs  = 65536
	leafV0Stride = 56
	leafV1Stride = 32
)

// Leaf is Lump 10: Leaf.
type Leaf struct {
	Generic
	data []primitives.Leaf
}

func (lump *Leaf) stride() (int, error) {
	switch lump.Context().LumpVersion {
	case 0:
		return leafV0Stride, nil
	case 1:
		return leafV1Stride, nil
	default:
		return 0, fmt.Errorf("unsupported leaf lump version %d", lump.Context().LumpVersion)
	}
}

// Unmarshall imports leaf records using the actual lump descriptor version.
func (lump *Leaf) Unmarshall(raw []byte) error {
	stride, err := lump.stride()
	if err != nil {
		return err
	}
	if len(raw)%stride != 0 {
		return fmt.Errorf("leaf lump length %d is not a multiple of record size %d", len(raw), stride)
	}
	count := len(raw) / stride
	if count > MaxMapLeafs {
		return fmt.Errorf("leaf count overflows maximum allowed size of %d", MaxMapLeafs)
	}
	data := make([]primitives.Leaf, count)
	for i := range data {
		record := raw[i*stride : (i+1)*stride]
		leaf := &data[i]
		leaf.Contents = int32(binary.LittleEndian.Uint32(record[0:4]))
		leaf.Cluster = int16(binary.LittleEndian.Uint16(record[4:6]))
		leaf.BitField = int16(binary.LittleEndian.Uint16(record[6:8]))
		for axis := range leaf.Mins {
			leaf.Mins[axis] = int16(binary.LittleEndian.Uint16(record[8+axis*2 : 10+axis*2]))
			leaf.Maxs[axis] = int16(binary.LittleEndian.Uint16(record[14+axis*2 : 16+axis*2]))
		}
		leaf.FirstLeafFace = binary.LittleEndian.Uint16(record[20:22])
		leaf.NumLeafFaces = binary.LittleEndian.Uint16(record[22:24])
		leaf.FirstLeafBrush = binary.LittleEndian.Uint16(record[24:26])
		leaf.NumLeafBrushes = binary.LittleEndian.Uint16(record[26:28])
		leaf.LeafWaterDataID = int16(binary.LittleEndian.Uint16(record[28:30]))
		if stride == leafV0Stride {
			for side := range leaf.LightSample.Color {
				color := &leaf.LightSample.Color[side]
				offset := 30 + side*4
				color.R = record[offset]
				color.G = record[offset+1]
				color.B = record[offset+2]
				color.Exponent = int8(record[offset+3])
			}
		}
		// The last two bytes are wire padding, not a light sample.
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

// GetData gets internal format structure data.
func (lump *Leaf) GetData() []primitives.Leaf {
	return lump.data
}

// Marshall exports leaf records with two zero-filled tail padding bytes.
func (lump *Leaf) Marshall() ([]byte, error) {
	stride, err := lump.stride()
	if err != nil {
		return nil, err
	}
	if len(lump.data) > MaxMapLeafs {
		return nil, fmt.Errorf("leaf count overflows maximum allowed size of %d", MaxMapLeafs)
	}
	raw := make([]byte, len(lump.data)*stride)
	for i := range lump.data {
		record := raw[i*stride : (i+1)*stride]
		leaf := &lump.data[i]
		binary.LittleEndian.PutUint32(record[0:4], uint32(leaf.Contents))
		binary.LittleEndian.PutUint16(record[4:6], uint16(leaf.Cluster))
		binary.LittleEndian.PutUint16(record[6:8], uint16(leaf.BitField))
		for axis := range leaf.Mins {
			binary.LittleEndian.PutUint16(record[8+axis*2:10+axis*2], uint16(leaf.Mins[axis]))
			binary.LittleEndian.PutUint16(record[14+axis*2:16+axis*2], uint16(leaf.Maxs[axis]))
		}
		binary.LittleEndian.PutUint16(record[20:22], leaf.FirstLeafFace)
		binary.LittleEndian.PutUint16(record[22:24], leaf.NumLeafFaces)
		binary.LittleEndian.PutUint16(record[24:26], leaf.FirstLeafBrush)
		binary.LittleEndian.PutUint16(record[26:28], leaf.NumLeafBrushes)
		binary.LittleEndian.PutUint16(record[28:30], uint16(leaf.LeafWaterDataID))
		if stride == leafV0Stride {
			for side := range leaf.LightSample.Color {
				color := &leaf.LightSample.Color[side]
				offset := 30 + side*4
				record[offset] = color.R
				record[offset+1] = color.G
				record[offset+2] = color.B
				record[offset+3] = byte(color.Exponent)
			}
		}
	}
	return raw, nil
}
