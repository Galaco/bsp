package lumps

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
)

// StaticPropLayout describes field availability, not a request to convert a file.
type StaticPropLayout uint8

const (
	StaticPropCommon4    StaticPropLayout = 4
	StaticPropCommon5    StaticPropLayout = 5
	StaticPropCommon6    StaticPropLayout = 6
	StaticPropCommon7    StaticPropLayout = 7
	StaticPropCommon8    StaticPropLayout = 8
	StaticPropCommon9    StaticPropLayout = 9
	StaticPropCommon10   StaticPropLayout = 10
	StaticPropCommon11   StaticPropLayout = 11
	StaticPropSDK2013V10 StaticPropLayout = 110
)

// StaticPropModel retains the entire original model slot, including bytes after NUL.
type StaticPropModel [128]byte

// Name interprets the slot through its first NUL. Unterminated slots are invalid.
func (model *StaticPropModel) Name() (string, error) {
	end := bytes.IndexByte(model[:], 0)
	if end < 0 {
		return "", fmt.Errorf("staticprop model slot has no NUL terminator")
	}
	return string(model[:end]), nil
}

// SetName replaces the slot atomically, NUL-terminating and zeroing its tail.
// The limit is encoded bytes, not Unicode code points; embedded NUL is rejected.
func (model *StaticPropModel) SetName(name string) error {
	if len(name) > 127 || strings.IndexByte(name, 0) >= 0 {
		return fmt.Errorf("staticprop model name must be at most 127 bytes without NUL")
	}
	var replacement StaticPropModel
	copy(replacement[:], name)
	*model = replacement
	return nil
}

// StaticProps is an owned snapshot. Counts are derived exclusively from slices.
type StaticProps struct {
	Layout      StaticPropLayout
	Models      []StaticPropModel
	LeafIndices []uint16
	Props       []StaticProp
}

// StaticProp retains the selected dialect's fields without flag, color or float
// normalization. Fields absent from Layout must have zero bits when encoding.
// UniformScale is absent (zero), not implicitly 1, before Common11.
type StaticProp struct {
	Origin, Angles                                     mgl32.Vec3
	ModelIndex, FirstLeaf, LeafCount                   uint16
	Solid                                              uint8
	Flags                                              uint32
	BasePadding                                        uint8
	Skin                                               int32
	FadeMinDist, FadeMaxDist                           float32
	LightingOrigin                                     mgl32.Vec3
	ForcedFadeScale                                    float32
	MinDXLevel, MaxDXLevel                             uint16
	MinCPULevel, MaxCPULevel, MinGPULevel, MaxGPULevel uint8
	DiffuseModulation                                  [4]byte
	DisableXBox360                                     uint8
	Padding                                            [3]byte
	ExtendedFlags                                      uint32
	UniformScale                                       float32
	LightmapResX, LightmapResY                         uint16
}
