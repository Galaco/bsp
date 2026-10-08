package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

// DecodeNodes decodes complete 32-byte wire records into owned values.
func DecodeNodes(raw []byte) ([]lumps.Node, error) {
	count, err := format.FixedRecordCount(len(raw), 32)
	if err != nil {
		return nil, err
	}
	records := make([]lumps.Node, count)
	for i := range records {
		wire := raw[i*32 : (i+1)*32]
		records[i].PlaneNum = int32(binary.LittleEndian.Uint32(wire[0:]))
		records[i].Children[0] = int32(binary.LittleEndian.Uint32(wire[4:]))
		records[i].Children[1] = int32(binary.LittleEndian.Uint32(wire[8:]))
		records[i].Mins[0] = int16(binary.LittleEndian.Uint16(wire[12:]))
		records[i].Mins[1] = int16(binary.LittleEndian.Uint16(wire[14:]))
		records[i].Mins[2] = int16(binary.LittleEndian.Uint16(wire[16:]))
		records[i].Maxs[0] = int16(binary.LittleEndian.Uint16(wire[18:]))
		records[i].Maxs[1] = int16(binary.LittleEndian.Uint16(wire[20:]))
		records[i].Maxs[2] = int16(binary.LittleEndian.Uint16(wire[22:]))
		records[i].FirstFace = binary.LittleEndian.Uint16(wire[24:])
		records[i].NumFaces = binary.LittleEndian.Uint16(wire[26:])
		records[i].Area = int16(binary.LittleEndian.Uint16(wire[28:]))
		records[i].Padding = int16(binary.LittleEndian.Uint16(wire[30:]))
	}
	return records, nil
}

// EncodeNodes preserves every field, including named padding and float bits.
func EncodeNodes(records []lumps.Node) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 32)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*32 : (i+1)*32]
		binary.LittleEndian.PutUint32(wire[0:], uint32(records[i].PlaneNum))
		binary.LittleEndian.PutUint32(wire[4:], uint32(records[i].Children[0]))
		binary.LittleEndian.PutUint32(wire[8:], uint32(records[i].Children[1]))
		binary.LittleEndian.PutUint16(wire[12:], uint16(records[i].Mins[0]))
		binary.LittleEndian.PutUint16(wire[14:], uint16(records[i].Mins[1]))
		binary.LittleEndian.PutUint16(wire[16:], uint16(records[i].Mins[2]))
		binary.LittleEndian.PutUint16(wire[18:], uint16(records[i].Maxs[0]))
		binary.LittleEndian.PutUint16(wire[20:], uint16(records[i].Maxs[1]))
		binary.LittleEndian.PutUint16(wire[22:], uint16(records[i].Maxs[2]))
		binary.LittleEndian.PutUint16(wire[24:], uint16(records[i].FirstFace))
		binary.LittleEndian.PutUint16(wire[26:], uint16(records[i].NumFaces))
		binary.LittleEndian.PutUint16(wire[28:], uint16(records[i].Area))
		binary.LittleEndian.PutUint16(wire[30:], uint16(records[i].Padding))
	}
	return raw, nil
}

// Nodes is the private decoder adapter for fixed wire records.
type Nodes struct {
	Metadata
	data []lumps.Node
}

func (lump *Nodes) Decode(raw []byte) error {
	data, err := DecodeNodes(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *Nodes) Encode() ([]byte, error) { return EncodeNodes(lump.data) }
func (lump *Nodes) GetData() []lumps.Node   { return lump.data }
