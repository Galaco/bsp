package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
	"math"
)

// DecodeTextureData decodes complete 32-byte wire records into owned values.
func DecodeTextureData(raw []byte) ([]lumps.TextureData, error) {
	count, err := format.FixedRecordCount(len(raw), 32)
	if err != nil {
		return nil, err
	}
	records := make([]lumps.TextureData, count)
	for i := range records {
		wire := raw[i*32 : (i+1)*32]
		records[i].Reflectivity[0] = math.Float32frombits(binary.LittleEndian.Uint32(wire[0:]))
		records[i].Reflectivity[1] = math.Float32frombits(binary.LittleEndian.Uint32(wire[4:]))
		records[i].Reflectivity[2] = math.Float32frombits(binary.LittleEndian.Uint32(wire[8:]))
		records[i].NameStringTableID = int32(binary.LittleEndian.Uint32(wire[12:]))
		records[i].Width = int32(binary.LittleEndian.Uint32(wire[16:]))
		records[i].Height = int32(binary.LittleEndian.Uint32(wire[20:]))
		records[i].ViewWidth = int32(binary.LittleEndian.Uint32(wire[24:]))
		records[i].ViewHeight = int32(binary.LittleEndian.Uint32(wire[28:]))
	}
	return records, nil
}

// EncodeTextureData preserves every field, including named padding and float bits.
func EncodeTextureData(records []lumps.TextureData) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 32)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*32 : (i+1)*32]
		binary.LittleEndian.PutUint32(wire[0:], math.Float32bits(records[i].Reflectivity[0]))
		binary.LittleEndian.PutUint32(wire[4:], math.Float32bits(records[i].Reflectivity[1]))
		binary.LittleEndian.PutUint32(wire[8:], math.Float32bits(records[i].Reflectivity[2]))
		binary.LittleEndian.PutUint32(wire[12:], uint32(records[i].NameStringTableID))
		binary.LittleEndian.PutUint32(wire[16:], uint32(records[i].Width))
		binary.LittleEndian.PutUint32(wire[20:], uint32(records[i].Height))
		binary.LittleEndian.PutUint32(wire[24:], uint32(records[i].ViewWidth))
		binary.LittleEndian.PutUint32(wire[28:], uint32(records[i].ViewHeight))
	}
	return raw, nil
}

// TextureData is the private decoder adapter for fixed wire records.
type TextureData struct {
	Metadata
	data []lumps.TextureData
}

func (lump *TextureData) Decode(raw []byte) error {
	data, err := DecodeTextureData(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *TextureData) Encode() ([]byte, error)      { return EncodeTextureData(lump.data) }
func (lump *TextureData) GetData() []lumps.TextureData { return lump.data }
