package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

// DecodeFaceIDs decodes complete 2-byte wire records into owned values.
func DecodeFaceIDs(raw []byte) ([]lumps.FaceID, error) {
	count, err := format.FixedRecordCount(len(raw), 2)
	if err != nil {
		return nil, err
	}
	records := make([]lumps.FaceID, count)
	for i := range records {
		wire := raw[i*2 : (i+1)*2]
		records[i].HammerFaceID = binary.LittleEndian.Uint16(wire[0:])
	}
	return records, nil
}

// EncodeFaceIDs preserves every field, including named padding and float bits.
func EncodeFaceIDs(records []lumps.FaceID) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 2)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*2 : (i+1)*2]
		binary.LittleEndian.PutUint16(wire[0:], uint16(records[i].HammerFaceID))
	}
	return raw, nil
}

// FaceIDs is the private decoder adapter for fixed wire records.
type FaceIDs struct {
	Metadata
	data []lumps.FaceID
}

func (lump *FaceIDs) Decode(raw []byte) error {
	data, err := DecodeFaceIDs(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *FaceIDs) Encode() ([]byte, error) { return EncodeFaceIDs(lump.data) }
func (lump *FaceIDs) GetData() []lumps.FaceID { return lump.data }
