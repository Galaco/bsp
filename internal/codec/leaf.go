package codec

import (
	"encoding/binary"
	"fmt"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
	"github.com/galaco/bsp/primitives/common"
)

// LeafRecordSize selects the explicit wire width from the leaf descriptor,
// independently of the BSP version and the ambient sibling descriptors.
func LeafRecordSize(context format.Context) (int, error) {
	if err := validatePCContext(context); err != nil {
		return 0, err
	}
	switch context.LumpVersion {
	case 0:
		return 56, nil
	case 1:
		return 32, nil
	default:
		return 0, fmt.Errorf("unsupported leaf lump version %d", context.LumpVersion)
	}
}

// DecodeLeaves preserves every field, including the selected layout and padding.
func DecodeLeaves(raw []byte, context format.Context) ([]lumps.Leaf, error) {
	stride, err := LeafRecordSize(context)
	if err != nil {
		return nil, err
	}
	if len(raw)%stride != 0 {
		return nil, fmt.Errorf("leaf length %d is not a multiple of %d", len(raw), stride)
	}
	count := len(raw) / stride
	if count > lumps.MaxMapLeafs {
		return nil, fmt.Errorf("leaf count %d exceeds %d", count, lumps.MaxMapLeafs)
	}
	data := make([]lumps.Leaf, count)
	for i := range data {
		record := raw[i*stride : (i+1)*stride]
		leaf := &data[i]
		leaf.Contents = int32(binary.LittleEndian.Uint32(record[0:4]))
		leaf.Cluster = int16(binary.LittleEndian.Uint16(record[4:6]))
		leaf.BitField = binary.LittleEndian.Uint16(record[6:8])
		for axis := range leaf.Mins {
			leaf.Mins[axis] = int16(binary.LittleEndian.Uint16(record[8+axis*2 : 10+axis*2]))
			leaf.Maxs[axis] = int16(binary.LittleEndian.Uint16(record[14+axis*2 : 16+axis*2]))
		}
		leaf.FirstLeafFace = binary.LittleEndian.Uint16(record[20:22])
		leaf.NumLeafFaces = binary.LittleEndian.Uint16(record[22:24])
		leaf.FirstLeafBrush = binary.LittleEndian.Uint16(record[24:26])
		leaf.NumLeafBrushes = binary.LittleEndian.Uint16(record[26:28])
		leaf.LeafWaterDataID = int16(binary.LittleEndian.Uint16(record[28:30]))
		leaf.HasInlineLighting = context.LumpVersion == 0
		if leaf.HasInlineLighting {
			leaf.LightSample = decodeAmbientCube(record[30:54])
		}
		copy(leaf.Padding[:], record[stride-2:])
	}
	return data, nil
}

// EncodeLeaves rejects records whose explicit inline-lighting layout or contents
// cannot be represented by the selected descriptor, instead of converting them.
func EncodeLeaves(data []lumps.Leaf, context format.Context) ([]byte, error) {
	stride, err := LeafRecordSize(context)
	if err != nil {
		return nil, err
	}
	if len(data) > lumps.MaxMapLeafs {
		return nil, fmt.Errorf("leaf count %d exceeds %d", len(data), lumps.MaxMapLeafs)
	}
	inline := context.LumpVersion == 0
	for i := range data {
		if data[i].HasInlineLighting != inline {
			return nil, fmt.Errorf("leaf %d inline-lighting layout does not match leaf version %d", i, context.LumpVersion)
		}
		if !inline && data[i].LightSample != (common.CompressedLightCube{}) {
			return nil, fmt.Errorf("leaf %d has nonrepresentable inline lighting in version 1", i)
		}
	}
	raw := make([]byte, len(data)*stride)
	for i := range data {
		record := raw[i*stride : (i+1)*stride]
		leaf := &data[i]
		binary.LittleEndian.PutUint32(record[0:4], uint32(leaf.Contents))
		binary.LittleEndian.PutUint16(record[4:6], uint16(leaf.Cluster))
		binary.LittleEndian.PutUint16(record[6:8], leaf.BitField)
		for axis := range leaf.Mins {
			binary.LittleEndian.PutUint16(record[8+axis*2:10+axis*2], uint16(leaf.Mins[axis]))
			binary.LittleEndian.PutUint16(record[14+axis*2:16+axis*2], uint16(leaf.Maxs[axis]))
		}
		binary.LittleEndian.PutUint16(record[20:22], leaf.FirstLeafFace)
		binary.LittleEndian.PutUint16(record[22:24], leaf.NumLeafFaces)
		binary.LittleEndian.PutUint16(record[24:26], leaf.FirstLeafBrush)
		binary.LittleEndian.PutUint16(record[26:28], leaf.NumLeafBrushes)
		binary.LittleEndian.PutUint16(record[28:30], uint16(leaf.LeafWaterDataID))
		if inline {
			encodeAmbientCube(record[30:54], leaf.LightSample)
		}
		copy(record[stride-2:], leaf.Padding[:])
	}
	return raw, nil
}

// Leaf wraps the stateless leaf codec for the internal lump registry.
type Leaf struct {
	lumps.Metadata
	data []lumps.Leaf
}

func (lump *Leaf) Unmarshall(raw []byte) error {
	data, err := DecodeLeaves(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *Leaf) Marshall() ([]byte, error) {
	return EncodeLeaves(lump.data, lump.Context())
}

func (lump *Leaf) GetData() []lumps.Leaf { return lump.data }
