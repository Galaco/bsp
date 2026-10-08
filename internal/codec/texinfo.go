package codec

import (
	"encoding/binary"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
	"math"
)

// DecodeTextureInfo decodes complete 72-byte wire records into owned values.
func DecodeTextureInfo(raw []byte) ([]lumps.TextureInfo, error) {
	count, err := format.FixedRecordCount(len(raw), 72)
	if err != nil {
		return nil, err
	}
	records := make([]lumps.TextureInfo, count)
	for i := range records {
		wire := raw[i*72 : (i+1)*72]
		records[i].TextureVecsTexelsPerWorldUnits[0][0] = math.Float32frombits(binary.LittleEndian.Uint32(wire[0:]))
		records[i].TextureVecsTexelsPerWorldUnits[0][1] = math.Float32frombits(binary.LittleEndian.Uint32(wire[4:]))
		records[i].TextureVecsTexelsPerWorldUnits[0][2] = math.Float32frombits(binary.LittleEndian.Uint32(wire[8:]))
		records[i].TextureVecsTexelsPerWorldUnits[0][3] = math.Float32frombits(binary.LittleEndian.Uint32(wire[12:]))
		records[i].TextureVecsTexelsPerWorldUnits[1][0] = math.Float32frombits(binary.LittleEndian.Uint32(wire[16:]))
		records[i].TextureVecsTexelsPerWorldUnits[1][1] = math.Float32frombits(binary.LittleEndian.Uint32(wire[20:]))
		records[i].TextureVecsTexelsPerWorldUnits[1][2] = math.Float32frombits(binary.LittleEndian.Uint32(wire[24:]))
		records[i].TextureVecsTexelsPerWorldUnits[1][3] = math.Float32frombits(binary.LittleEndian.Uint32(wire[28:]))
		records[i].LightmapVecsLuxelsPerWorldUnits[0][0] = math.Float32frombits(binary.LittleEndian.Uint32(wire[32:]))
		records[i].LightmapVecsLuxelsPerWorldUnits[0][1] = math.Float32frombits(binary.LittleEndian.Uint32(wire[36:]))
		records[i].LightmapVecsLuxelsPerWorldUnits[0][2] = math.Float32frombits(binary.LittleEndian.Uint32(wire[40:]))
		records[i].LightmapVecsLuxelsPerWorldUnits[0][3] = math.Float32frombits(binary.LittleEndian.Uint32(wire[44:]))
		records[i].LightmapVecsLuxelsPerWorldUnits[1][0] = math.Float32frombits(binary.LittleEndian.Uint32(wire[48:]))
		records[i].LightmapVecsLuxelsPerWorldUnits[1][1] = math.Float32frombits(binary.LittleEndian.Uint32(wire[52:]))
		records[i].LightmapVecsLuxelsPerWorldUnits[1][2] = math.Float32frombits(binary.LittleEndian.Uint32(wire[56:]))
		records[i].LightmapVecsLuxelsPerWorldUnits[1][3] = math.Float32frombits(binary.LittleEndian.Uint32(wire[60:]))
		records[i].Flags = int32(binary.LittleEndian.Uint32(wire[64:]))
		records[i].TexData = int32(binary.LittleEndian.Uint32(wire[68:]))
	}
	return records, nil
}

// EncodeTextureInfo preserves every field, including named padding and float bits.
func EncodeTextureInfo(records []lumps.TextureInfo) ([]byte, error) {
	size, err := format.FixedRecordBytes(len(records), 72)
	if err != nil {
		return nil, err
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*72 : (i+1)*72]
		binary.LittleEndian.PutUint32(wire[0:], math.Float32bits(records[i].TextureVecsTexelsPerWorldUnits[0][0]))
		binary.LittleEndian.PutUint32(wire[4:], math.Float32bits(records[i].TextureVecsTexelsPerWorldUnits[0][1]))
		binary.LittleEndian.PutUint32(wire[8:], math.Float32bits(records[i].TextureVecsTexelsPerWorldUnits[0][2]))
		binary.LittleEndian.PutUint32(wire[12:], math.Float32bits(records[i].TextureVecsTexelsPerWorldUnits[0][3]))
		binary.LittleEndian.PutUint32(wire[16:], math.Float32bits(records[i].TextureVecsTexelsPerWorldUnits[1][0]))
		binary.LittleEndian.PutUint32(wire[20:], math.Float32bits(records[i].TextureVecsTexelsPerWorldUnits[1][1]))
		binary.LittleEndian.PutUint32(wire[24:], math.Float32bits(records[i].TextureVecsTexelsPerWorldUnits[1][2]))
		binary.LittleEndian.PutUint32(wire[28:], math.Float32bits(records[i].TextureVecsTexelsPerWorldUnits[1][3]))
		binary.LittleEndian.PutUint32(wire[32:], math.Float32bits(records[i].LightmapVecsLuxelsPerWorldUnits[0][0]))
		binary.LittleEndian.PutUint32(wire[36:], math.Float32bits(records[i].LightmapVecsLuxelsPerWorldUnits[0][1]))
		binary.LittleEndian.PutUint32(wire[40:], math.Float32bits(records[i].LightmapVecsLuxelsPerWorldUnits[0][2]))
		binary.LittleEndian.PutUint32(wire[44:], math.Float32bits(records[i].LightmapVecsLuxelsPerWorldUnits[0][3]))
		binary.LittleEndian.PutUint32(wire[48:], math.Float32bits(records[i].LightmapVecsLuxelsPerWorldUnits[1][0]))
		binary.LittleEndian.PutUint32(wire[52:], math.Float32bits(records[i].LightmapVecsLuxelsPerWorldUnits[1][1]))
		binary.LittleEndian.PutUint32(wire[56:], math.Float32bits(records[i].LightmapVecsLuxelsPerWorldUnits[1][2]))
		binary.LittleEndian.PutUint32(wire[60:], math.Float32bits(records[i].LightmapVecsLuxelsPerWorldUnits[1][3]))
		binary.LittleEndian.PutUint32(wire[64:], uint32(records[i].Flags))
		binary.LittleEndian.PutUint32(wire[68:], uint32(records[i].TexData))
	}
	return raw, nil
}

// TextureInfo is the private decoder adapter for fixed wire records.
type TextureInfo struct {
	Metadata
	data []lumps.TextureInfo
}

func (lump *TextureInfo) Decode(raw []byte) error {
	data, err := DecodeTextureInfo(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *TextureInfo) Encode() ([]byte, error)      { return EncodeTextureInfo(lump.data) }
func (lump *TextureInfo) GetData() []lumps.TextureInfo { return lump.data }
