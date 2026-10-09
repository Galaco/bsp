package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

// DecodeAreaPortals decodes complete 12-byte wire records into owned values.
func DecodeAreaPortals(raw []byte) ([]lumps.AreaPortal, error) {
	count, err := format.FixedRecordCount(len(raw), 12)
	if err != nil {
		return nil, err
	}
	records := make([]lumps.AreaPortal, count)
	for i := range records {
		wire := raw[i*12 : (i+1)*12]
		records[i].PortalKey = binary.LittleEndian.Uint16(wire[0:])
		records[i].OtherArea = binary.LittleEndian.Uint16(wire[2:])
		records[i].FirstClipPortalVert = binary.LittleEndian.Uint16(wire[4:])
		records[i].NumClipPortalVerts = binary.LittleEndian.Uint16(wire[6:])
		records[i].PlaneNum = int32(binary.LittleEndian.Uint32(wire[8:]))
	}
	return records, nil
}

// EncodeAreaPortals preserves every field, including named padding and float bits.
func EncodeAreaPortals(records []lumps.AreaPortal) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 12)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*12 : (i+1)*12]
		binary.LittleEndian.PutUint16(wire[0:], uint16(records[i].PortalKey))
		binary.LittleEndian.PutUint16(wire[2:], uint16(records[i].OtherArea))
		binary.LittleEndian.PutUint16(wire[4:], uint16(records[i].FirstClipPortalVert))
		binary.LittleEndian.PutUint16(wire[6:], uint16(records[i].NumClipPortalVerts))
		binary.LittleEndian.PutUint32(wire[8:], uint32(records[i].PlaneNum))
	}
	return raw, nil
}

// AreaPortals is the private decoder adapter for fixed wire records.
type AreaPortals struct {
	Metadata
	data []lumps.AreaPortal
}

func (lump *AreaPortals) Decode(raw []byte) error {
	data, err := DecodeAreaPortals(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *AreaPortals) Encode() ([]byte, error)     { return EncodeAreaPortals(lump.data) }
func (lump *AreaPortals) GetData() []lumps.AreaPortal { return lump.data }
