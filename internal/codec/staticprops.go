package codec

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

// StaticPropLayoutFor selects an explicit profile/revision, never a byte-size
// guess. SDK7's lightmap variant and other unproven branch layouts stay opaque.
// References: Valve SDK b8cfb12 gamebspfile.h (4/5/6 and distinct10);
// srctools cd7e841 bsp.py StaticPropVersion (common4–11); authored CS:GO v11
// fixtures and independent consumer observations are recorded in task08.
func StaticPropLayoutFor(context format.Context) (lumps.StaticPropLayout, error) {
	if err := validatePCContext(context); err != nil {
		return 0, err
	}
	version := context.LumpVersion
	switch context.Profile {
	case "sdk2013":
		if version >= 4 && version <= 6 {
			return lumps.StaticPropLayout(version), nil
		}
		if version == 10 {
			return lumps.StaticPropSDK2013V10, nil
		}
	case "l4d2":
		if version >= 4 && version <= 9 {
			return lumps.StaticPropLayout(version), nil
		}
	case "csgo":
		if version >= 4 && version <= 11 {
			return lumps.StaticPropLayout(version), nil
		}
	}
	return 0, fmt.Errorf("unsupported staticprop revision %d for %s BSP%d", version, context.Profile, context.BSPVersion)
}

func staticPropWidth(layout lumps.StaticPropLayout) int {
	switch layout {
	case lumps.StaticPropCommon4:
		return 56
	case lumps.StaticPropCommon5:
		return 60
	case lumps.StaticPropCommon6:
		return 64
	case lumps.StaticPropCommon7, lumps.StaticPropCommon8:
		return 68
	case lumps.StaticPropCommon9, lumps.StaticPropSDK2013V10:
		return 72
	case lumps.StaticPropCommon10:
		return 76
	case lumps.StaticPropCommon11:
		return 80
	}
	return 0
}

func StaticPropRecordSize(context format.Context) (int, error) {
	layout, err := StaticPropLayoutFor(context)
	if err != nil {
		return 0, err
	}
	return staticPropWidth(layout), nil
}

// staticPropSection validates the signed count and byte budget before any
// allocation. Division bounds multiplication even on 32-bit hosts.
func staticPropSection(raw []byte, offset, width, reserve int, section string) (int, int, error) {
	if len(raw)-offset < 4+reserve {
		return 0, 0, fmt.Errorf("staticprop %s: truncated count", section)
	}
	count := int32(binary.LittleEndian.Uint32(raw[offset:]))
	start := offset + 4
	if count < 0 || int64(count) > int64((len(raw)-start-reserve)/width) {
		return 0, 0, fmt.Errorf("staticprop %s: count %d exceeds byte budget", section, count)
	}
	return int(count), start, nil
}

// DecodeStaticProps preflights every counted section and exact consumption,
// then returns new owned arrays. Original slot bytes, padding and float bits
// remain untouched. External leaf IDs are checked by the root against its leaves.
func DecodeStaticProps(raw []byte, context format.Context) (lumps.StaticProps, error) {
	layout, err := StaticPropLayoutFor(context)
	if err != nil {
		return lumps.StaticProps{}, err
	}
	width := staticPropWidth(layout)
	models, modelStart, err := staticPropSection(raw, 0, 128, 8, "models")
	if err != nil {
		return lumps.StaticProps{}, err
	}
	leaves, leafStart, err := staticPropSection(raw, modelStart+models*128, 2, 4, "leaves")
	if err != nil {
		return lumps.StaticProps{}, err
	}
	props, propStart, err := staticPropSection(raw, leafStart+leaves*2, width, 0, "props")
	if err != nil {
		return lumps.StaticProps{}, err
	}
	if props*width != len(raw)-propStart {
		return lumps.StaticProps{}, fmt.Errorf("staticprop: trailing bytes or mismatched record width")
	}
	for i := 0; i < models; i++ {
		if bytes.IndexByte(raw[modelStart+i*128:modelStart+(i+1)*128], 0) < 0 {
			return lumps.StaticProps{}, fmt.Errorf("staticprop model %d: unterminated slot", i)
		}
	}
	value := lumps.StaticProps{
		Layout:      layout,
		Models:      make([]lumps.StaticPropModel, models),
		LeafIndices: make([]uint16, leaves),
		Props:       make([]lumps.StaticProp, props),
	}
	for i := range value.Models {
		copy(value.Models[i][:], raw[modelStart+i*128:modelStart+(i+1)*128])
	}
	for i := range value.LeafIndices {
		value.LeafIndices[i] = binary.LittleEndian.Uint16(raw[leafStart+i*2:])
	}
	for i := range value.Props {
		wire := raw[propStart+i*width : propStart+(i+1)*width]
		record := &value.Props[i]
		record.Origin = worldLightVector(wire)
		record.Angles = worldLightVector(wire[12:])
		record.ModelIndex = binary.LittleEndian.Uint16(wire[24:])
		record.FirstLeaf = binary.LittleEndian.Uint16(wire[26:])
		record.LeafCount = binary.LittleEndian.Uint16(wire[28:])
		record.Solid = wire[30]
		record.Skin = int32(binary.LittleEndian.Uint32(wire[32:]))
		record.FadeMinDist = math.Float32frombits(binary.LittleEndian.Uint32(wire[36:]))
		record.FadeMaxDist = math.Float32frombits(binary.LittleEndian.Uint32(wire[40:]))
		record.LightingOrigin = worldLightVector(wire[44:])
		if layout != lumps.StaticPropCommon4 {
			record.ForcedFadeScale = math.Float32frombits(binary.LittleEndian.Uint32(wire[56:]))
		}
		if layout == lumps.StaticPropCommon6 || layout == lumps.StaticPropCommon7 || layout == lumps.StaticPropSDK2013V10 {
			record.MinDXLevel = binary.LittleEndian.Uint16(wire[60:])
			record.MaxDXLevel = binary.LittleEndian.Uint16(wire[62:])
		}
		if layout == lumps.StaticPropSDK2013V10 {
			record.BasePadding = wire[31]
			record.Flags = binary.LittleEndian.Uint32(wire[64:])
			record.LightmapResX = binary.LittleEndian.Uint16(wire[68:])
			record.LightmapResY = binary.LittleEndian.Uint16(wire[70:])
		} else {
			record.Flags = uint32(wire[31])
			if layout >= lumps.StaticPropCommon8 {
				record.MinCPULevel, record.MaxCPULevel = wire[60], wire[61]
				record.MinGPULevel, record.MaxGPULevel = wire[62], wire[63]
			}
			if layout >= lumps.StaticPropCommon7 {
				copy(record.DiffuseModulation[:], wire[64:68])
			}
			if layout >= lumps.StaticPropCommon9 {
				record.DisableXBox360 = wire[68]
				copy(record.Padding[:], wire[69:72])
			}
			if layout >= lumps.StaticPropCommon10 {
				record.ExtendedFlags = binary.LittleEndian.Uint32(wire[72:])
			}
			if layout == lumps.StaticPropCommon11 {
				record.UniformScale = math.Float32frombits(binary.LittleEndian.Uint32(wire[76:]))
			}
		}
	}
	if err := validateStaticPropSelfReferences(value); err != nil {
		return lumps.StaticProps{}, err
	}
	return value, nil
}

func validateStaticPropSelfReferences(value lumps.StaticProps) error {
	for i := range value.Props {
		record := &value.Props[i]
		if int(record.ModelIndex) >= len(value.Models) {
			return fmt.Errorf("staticprop %d: model index %d outside %d slots", i, record.ModelIndex, len(value.Models))
		}
		first, count := int(record.FirstLeaf), int(record.LeafCount)
		if first > len(value.LeafIndices) || count > len(value.LeafIndices)-first {
			return fmt.Errorf("staticprop %d: leaf span %d+%d outside %d indices", i, first, count, len(value.LeafIndices))
		}
	}
	return nil
}

// ValidateStaticPropReferences validates every leaf ID, including unreferenced
// entries, against the root's selected ordinary leaf count (exclusive bound).
func ValidateStaticPropReferences(value lumps.StaticProps, leafCount int) error {
	if leafCount < 0 {
		return fmt.Errorf("staticprop: negative external leaf count")
	}
	if err := validateStaticPropSelfReferences(value); err != nil {
		return err
	}
	for i, leaf := range value.LeafIndices {
		if int(leaf) >= leafCount {
			return fmt.Errorf("staticprop leaf index %d: leaf %d outside %d leaves", i, leaf, leafCount)
		}
	}
	return nil
}

func staticPropRepresentable(record *lumps.StaticProp, layout lumps.StaticPropLayout) error {
	sdk := layout == lumps.StaticPropSDK2013V10
	if !sdk && (record.Flags > 255 || record.BasePadding != 0 || record.LightmapResX != 0 || record.LightmapResY != 0) {
		return fmt.Errorf("SDK full flags, base padding or lightmap fields unavailable in common layout")
	}
	if layout == lumps.StaticPropCommon4 && math.Float32bits(record.ForcedFadeScale) != 0 {
		return fmt.Errorf("forced fade scale unavailable")
	}
	if !sdk && layout != lumps.StaticPropCommon6 && layout != lumps.StaticPropCommon7 && (record.MinDXLevel != 0 || record.MaxDXLevel != 0) {
		return fmt.Errorf("DX levels unavailable")
	}
	if (sdk || layout < lumps.StaticPropCommon8) && (record.MinCPULevel != 0 || record.MaxCPULevel != 0 || record.MinGPULevel != 0 || record.MaxGPULevel != 0) {
		return fmt.Errorf("CPU/GPU levels unavailable")
	}
	if (sdk || layout < lumps.StaticPropCommon7) && record.DiffuseModulation != [4]byte{} {
		return fmt.Errorf("diffuse modulation unavailable")
	}
	if (sdk || layout < lumps.StaticPropCommon9) && (record.DisableXBox360 != 0 || record.Padding != [3]byte{}) {
		return fmt.Errorf("Xbox disable byte or padding unavailable")
	}
	if (sdk || layout < lumps.StaticPropCommon10) && record.ExtendedFlags != 0 {
		return fmt.Errorf("extended flags unavailable")
	}
	if layout != lumps.StaticPropCommon11 && math.Float32bits(record.UniformScale) != 0 {
		return fmt.Errorf("uniform scale unavailable")
	}
	return nil
}

func staticPropPayloadSize(models, leaves, props, width int) (int, error) {
	if models < 0 || leaves < 0 || props < 0 || width <= 0 {
		return 0, fmt.Errorf("staticprop: invalid count or width")
	}
	// Counts and encoded payloads must fit signed 32-bit Source fields and the
	// host's int. Add only individually bounded products to avoid uint64 overflow.
	size := uint64(12)
	for _, section := range [3][2]int{{models, 128}, {leaves, 2}, {props, width}} {
		if uint64(section[0]) > math.MaxInt32 || uint64(section[0]) > (math.MaxInt32-size)/uint64(section[1]) {
			return 0, fmt.Errorf("staticprop: count or payload exceeds signed 32-bit size")
		}
		size += uint64(section[0]) * uint64(section[1])
	}
	if size > uint64(^uint(0)>>1) {
		return 0, fmt.Errorf("staticprop: payload exceeds native size")
	}
	return int(size), nil
}

// EncodeStaticProps validates the entire snapshot before allocating output.
// Layout must match the existing descriptor's context; absent fields cannot be
// silently discarded, including negative-zero float bits.
func EncodeStaticProps(value lumps.StaticProps, context format.Context) ([]byte, error) {
	layout, err := StaticPropLayoutFor(context)
	if err != nil {
		return nil, err
	}
	if value.Layout != layout {
		return nil, fmt.Errorf("staticprop: snapshot layout %d does not match selected layout %d", value.Layout, layout)
	}
	width := staticPropWidth(layout)
	size, err := staticPropPayloadSize(len(value.Models), len(value.LeafIndices), len(value.Props), width)
	if err != nil {
		return nil, err
	}
	for i := range value.Models {
		if bytes.IndexByte(value.Models[i][:], 0) < 0 {
			return nil, fmt.Errorf("staticprop model %d: unterminated slot", i)
		}
	}
	if err := validateStaticPropSelfReferences(value); err != nil {
		return nil, err
	}
	for i := range value.Props {
		if err := staticPropRepresentable(&value.Props[i], layout); err != nil {
			return nil, fmt.Errorf("staticprop %d: %w", i, err)
		}
	}
	raw := make([]byte, size)
	binary.LittleEndian.PutUint32(raw, uint32(len(value.Models)))
	offset := 4
	for i := range value.Models {
		copy(raw[offset:offset+128], value.Models[i][:])
		offset += 128
	}
	binary.LittleEndian.PutUint32(raw[offset:], uint32(len(value.LeafIndices)))
	offset += 4
	for _, leaf := range value.LeafIndices {
		binary.LittleEndian.PutUint16(raw[offset:], leaf)
		offset += 2
	}
	binary.LittleEndian.PutUint32(raw[offset:], uint32(len(value.Props)))
	offset += 4
	for i := range value.Props {
		record := &value.Props[i]
		wire := raw[offset : offset+width]
		offset += width
		putWorldLightVector(wire, record.Origin)
		putWorldLightVector(wire[12:], record.Angles)
		binary.LittleEndian.PutUint16(wire[24:], record.ModelIndex)
		binary.LittleEndian.PutUint16(wire[26:], record.FirstLeaf)
		binary.LittleEndian.PutUint16(wire[28:], record.LeafCount)
		wire[30] = record.Solid
		binary.LittleEndian.PutUint32(wire[32:], uint32(record.Skin))
		binary.LittleEndian.PutUint32(wire[36:], math.Float32bits(record.FadeMinDist))
		binary.LittleEndian.PutUint32(wire[40:], math.Float32bits(record.FadeMaxDist))
		putWorldLightVector(wire[44:], record.LightingOrigin)
		if layout != lumps.StaticPropCommon4 {
			binary.LittleEndian.PutUint32(wire[56:], math.Float32bits(record.ForcedFadeScale))
		}
		if layout == lumps.StaticPropCommon6 || layout == lumps.StaticPropCommon7 || layout == lumps.StaticPropSDK2013V10 {
			binary.LittleEndian.PutUint16(wire[60:], record.MinDXLevel)
			binary.LittleEndian.PutUint16(wire[62:], record.MaxDXLevel)
		}
		if layout == lumps.StaticPropSDK2013V10 {
			wire[31] = record.BasePadding
			binary.LittleEndian.PutUint32(wire[64:], record.Flags)
			binary.LittleEndian.PutUint16(wire[68:], record.LightmapResX)
			binary.LittleEndian.PutUint16(wire[70:], record.LightmapResY)
		} else {
			wire[31] = byte(record.Flags)
			if layout >= lumps.StaticPropCommon8 {
				wire[60], wire[61] = record.MinCPULevel, record.MaxCPULevel
				wire[62], wire[63] = record.MinGPULevel, record.MaxGPULevel
			}
			if layout >= lumps.StaticPropCommon7 {
				copy(wire[64:68], record.DiffuseModulation[:])
			}
			if layout >= lumps.StaticPropCommon9 {
				wire[68] = record.DisableXBox360
				copy(wire[69:72], record.Padding[:])
			}
			if layout >= lumps.StaticPropCommon10 {
				binary.LittleEndian.PutUint32(wire[72:], record.ExtendedFlags)
			}
			if layout == lumps.StaticPropCommon11 {
				binary.LittleEndian.PutUint32(wire[76:], math.Float32bits(record.UniformScale))
			}
		}
	}
	return raw, nil
}
