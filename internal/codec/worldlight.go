package codec

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
	"github.com/go-gl/mathgl/mgl32"
)

// WorldLightRecordSize selects the descriptor revision without payload guessing.
// SDK2013 defines v0; CS:GO and Reactive Drop retain v0 and add v1.
func WorldLightRecordSize(context format.Context) (int, error) {
	if err := validatePCContext(context); err != nil {
		return 0, err
	}
	switch context.LumpVersion {
	case 0:
		return 88, nil
	case 1:
		if context.Profile == "csgo" || context.Profile == "reactivedrop" {
			return 100, nil
		}
	}
	return 0, fmt.Errorf("unsupported worldlight revision %d for %s", context.LumpVersion, context.Profile)
}

func lightSidePayloadSize(count, width int) (int, error) {
	if count > int(^uint(0)>>1)/width || int64(count) > (1<<31-1)/int64(width) {
		return 0, fmt.Errorf("ordinary payload exceeds native or signed 32-bit size")
	}
	return count * width, nil
}

func worldLightVector(raw []byte) mgl32.Vec3 {
	return mgl32.Vec3{
		math.Float32frombits(binary.LittleEndian.Uint32(raw)),
		math.Float32frombits(binary.LittleEndian.Uint32(raw[4:])),
		math.Float32frombits(binary.LittleEndian.Uint32(raw[8:])),
	}
}

func putWorldLightVector(raw []byte, value mgl32.Vec3) {
	binary.LittleEndian.PutUint32(raw, math.Float32bits(value[0]))
	binary.LittleEndian.PutUint32(raw[4:], math.Float32bits(value[1]))
	binary.LittleEndian.PutUint32(raw[8:], math.Float32bits(value[2]))
}

// DecodeWorldLights decodes either channel, preserving every 32-bit emission
// value (including values unknown to the SDK) and all floating-point bits.
func DecodeWorldLights(raw []byte, context format.Context) ([]lumps.WorldLight, error) {
	width, err := WorldLightRecordSize(context)
	if err != nil {
		return nil, err
	}
	if len(raw)%width != 0 {
		return nil, fmt.Errorf("worldlight revision %d: %d bytes is not a multiple of %d", context.LumpVersion, len(raw), width)
	}
	records := make([]lumps.WorldLight, len(raw)/width)
	for i := range records {
		wire := raw[i*width : (i+1)*width]
		record := &records[i]
		record.Origin = worldLightVector(wire)
		record.Intensity = worldLightVector(wire[12:])
		record.Normal = worldLightVector(wire[24:])
		tail := wire[36:]
		if context.LumpVersion == 1 {
			record.ShadowCastOffset = worldLightVector(tail)
			tail = tail[12:]
		}
		record.Cluster = int32(binary.LittleEndian.Uint32(tail))
		record.Type = lumps.EmitType(int32(binary.LittleEndian.Uint32(tail[4:])))
		record.Style = int32(binary.LittleEndian.Uint32(tail[8:]))
		record.Stopdot = math.Float32frombits(binary.LittleEndian.Uint32(tail[12:]))
		record.Stopdot2 = math.Float32frombits(binary.LittleEndian.Uint32(tail[16:]))
		record.Exponent = math.Float32frombits(binary.LittleEndian.Uint32(tail[20:]))
		record.Radius = math.Float32frombits(binary.LittleEndian.Uint32(tail[24:]))
		record.ConstantAttenuation = math.Float32frombits(binary.LittleEndian.Uint32(tail[28:]))
		record.LinearAttenuation = math.Float32frombits(binary.LittleEndian.Uint32(tail[32:]))
		record.QuadraticAttenuation = math.Float32frombits(binary.LittleEndian.Uint32(tail[36:]))
		record.Flags = int32(binary.LittleEndian.Uint32(tail[40:]))
		record.TexInfo = int32(binary.LittleEndian.Uint32(tail[44:]))
		record.Owner = int32(binary.LittleEndian.Uint32(tail[48:]))
	}
	return records, nil
}

// EncodeWorldLights validates representability before allocating output. Even
// negative zero in the absent v0 shadow vector is rejected rather than lost.
func EncodeWorldLights(records []lumps.WorldLight, context format.Context) ([]byte, error) {
	width, err := WorldLightRecordSize(context)
	if err != nil {
		return nil, err
	}
	size, err := lightSidePayloadSize(len(records), width)
	if err != nil {
		return nil, err
	}
	if context.LumpVersion == 0 {
		for i := range records {
			for _, value := range records[i].ShadowCastOffset {
				if math.Float32bits(value) != 0 {
					return nil, fmt.Errorf("worldlight %d: shadow cast offset is not representable in revision 0", i)
				}
			}
		}
	}
	raw := make([]byte, size)
	for i := range records {
		wire := raw[i*width : (i+1)*width]
		record := &records[i]
		putWorldLightVector(wire, record.Origin)
		putWorldLightVector(wire[12:], record.Intensity)
		putWorldLightVector(wire[24:], record.Normal)
		tail := wire[36:]
		if context.LumpVersion == 1 {
			putWorldLightVector(tail, record.ShadowCastOffset)
			tail = tail[12:]
		}
		binary.LittleEndian.PutUint32(tail, uint32(record.Cluster))
		binary.LittleEndian.PutUint32(tail[4:], uint32(record.Type))
		binary.LittleEndian.PutUint32(tail[8:], uint32(record.Style))
		binary.LittleEndian.PutUint32(tail[12:], math.Float32bits(record.Stopdot))
		binary.LittleEndian.PutUint32(tail[16:], math.Float32bits(record.Stopdot2))
		binary.LittleEndian.PutUint32(tail[20:], math.Float32bits(record.Exponent))
		binary.LittleEndian.PutUint32(tail[24:], math.Float32bits(record.Radius))
		binary.LittleEndian.PutUint32(tail[28:], math.Float32bits(record.ConstantAttenuation))
		binary.LittleEndian.PutUint32(tail[32:], math.Float32bits(record.LinearAttenuation))
		binary.LittleEndian.PutUint32(tail[36:], math.Float32bits(record.QuadraticAttenuation))
		binary.LittleEndian.PutUint32(tail[40:], uint32(record.Flags))
		binary.LittleEndian.PutUint32(tail[44:], uint32(record.TexInfo))
		binary.LittleEndian.PutUint32(tail[48:], uint32(record.Owner))
	}
	return raw, nil
}

// WorldLight wraps the LDR worldlight lump for the internal registry.
type WorldLight struct {
	lumps.Metadata
	data []lumps.WorldLight
}

func (lump *WorldLight) Unmarshall(raw []byte) error {
	data, err := DecodeWorldLights(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *WorldLight) Marshall() ([]byte, error) {
	return EncodeWorldLights(lump.data, lump.Context())
}

func (lump *WorldLight) GetData() []lumps.WorldLight { return lump.data }

// WorldLightHDR is a distinct HDR channel with the same versioned wire layout.
type WorldLightHDR struct {
	lumps.Metadata
	data []lumps.WorldLight
}

func (lump *WorldLightHDR) Unmarshall(raw []byte) error {
	data, err := DecodeWorldLights(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *WorldLightHDR) Marshall() ([]byte, error) {
	return EncodeWorldLights(lump.data, lump.Context())
}

func (lump *WorldLightHDR) GetData() []lumps.WorldLight { return lump.data }
