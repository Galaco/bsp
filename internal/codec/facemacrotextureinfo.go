package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

// DecodeFaceMacroTextureInfo decodes complete 2-byte wire records into owned values.
func DecodeFaceMacroTextureInfo(raw []byte) ([]lumps.FaceMacroTextureInfo, error) {
	count, err := format.FixedRecordCount(len(raw), 2)
	if err != nil {
		return nil, err
	}
	records := make([]lumps.FaceMacroTextureInfo, count)
	for i := range records {
		wire := raw[i*2 : (i+1)*2]
		records[i].MacroTextureNameID = binary.LittleEndian.Uint16(wire[0:])
	}
	return records, nil
}

// EncodeFaceMacroTextureInfo preserves every field, including named padding and float bits.
func EncodeFaceMacroTextureInfo(records []lumps.FaceMacroTextureInfo) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 2)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*2 : (i+1)*2]
		binary.LittleEndian.PutUint16(wire[0:], uint16(records[i].MacroTextureNameID))
	}
	return raw, nil
}

// FaceMacroTextureInfo is the private decoder adapter for fixed wire records.
type FaceMacroTextureInfo struct {
	Metadata
	data []lumps.FaceMacroTextureInfo
}

func (lump *FaceMacroTextureInfo) Decode(raw []byte) error {
	data, err := DecodeFaceMacroTextureInfo(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *FaceMacroTextureInfo) Encode() ([]byte, error) {
	return EncodeFaceMacroTextureInfo(lump.data)
}
func (lump *FaceMacroTextureInfo) GetData() []lumps.FaceMacroTextureInfo { return lump.data }
