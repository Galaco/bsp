package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

// DecodeAreas decodes complete 8-byte wire records into owned values.
func DecodeAreas(raw []byte) ([]lumps.Area, error) {
	count, err := format.FixedRecordCount(len(raw), 8)
	if err != nil {
		return nil, err
	}
	records := make([]lumps.Area, count)
	for i := range records {
		wire := raw[i*8 : (i+1)*8]
		records[i].NumAreaPortals = int32(binary.LittleEndian.Uint32(wire[0:]))
		records[i].FirstAreaPortal = int32(binary.LittleEndian.Uint32(wire[4:]))
	}
	return records, nil
}

// EncodeAreas preserves every field, including named padding and float bits.
func EncodeAreas(records []lumps.Area) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 8)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*8 : (i+1)*8]
		binary.LittleEndian.PutUint32(wire[0:], uint32(records[i].NumAreaPortals))
		binary.LittleEndian.PutUint32(wire[4:], uint32(records[i].FirstAreaPortal))
	}
	return raw, nil
}

// Areas is the private decoder adapter for fixed wire records.
type Areas struct {
	Metadata
	data []lumps.Area
}

func (lump *Areas) Decode(raw []byte) error {
	data, err := DecodeAreas(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *Areas) Encode() ([]byte, error) { return EncodeAreas(lump.data) }
func (lump *Areas) GetData() []lumps.Area   { return lump.data }
