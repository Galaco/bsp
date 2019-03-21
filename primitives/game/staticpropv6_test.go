package game

import (
	"github.com/go-gl/mathgl/mgl32"
	"testing"
)

func TestStaticPropV6_GetAngles(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetAngles().X() != 2 &&
		sut.GetAngles().Y() != 5 &&
		sut.GetAngles().Z() != 8 {
		t.Error("unexpected value for angles property")
	}
}

func TestStaticPropV6_GetDiffuseModulation(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetDiffuseModulation() != 0 {
		t.Error("unexpected value for diffuseModulation property")
	}

}

func TestStaticPropV6_GetDisableXBox360(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetDisableXBox360() != false {
		t.Error("unexpected value for unknown property")
	}

}

func TestStaticPropV6_GetFadeMaxDist(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetFadeMaxDist() != 256 {
		t.Error("unexpected value for fadeMaxDist property")
	}

}

func TestStaticPropV6_GetFadeMinDist(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetFadeMinDist() != 84.5 {
		t.Error("unexpected value for fadeMinDist property")
	}

}

func TestStaticPropV6_GetFirstLeaf(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetFirstLeaf() != 21 {
		t.Error("unexpected value for firstLeaf property")
	}

}

func TestStaticPropV6_GetFlags(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetFlags() != 85 {
		t.Error("unexpected value for flags property")
	}

}

func TestStaticPropV6_GetForcedFadeScale(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetForcedFadeScale() != 65 {
		t.Error("unexpected value for forcedFadeScale property")
	}

}

func TestStaticPropV6_GetLeafCount(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetLeafCount() != 65 {
		t.Error("unexpected value for leafCount property")
	}

}

func TestStaticPropV6_GetLightingOrigin(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetLightingOrigin().X() != 32 &&
		sut.GetLightingOrigin().Y() != 64 &&
		sut.GetLightingOrigin().Z() != 128 {
		t.Error("unexpected value for angles property")
	}
}

func TestStaticPropV6_GetMaxCPULevel(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetMaxCPULevel() != 0 {
		t.Error("unexpected value for maxCPULevel property")
	}

}

func TestStaticPropV6_GetMaxDXLevel(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetMaxDXLevel() != 0 {
		t.Error("unexpected value for maxDXLevel property")
	}

}

func TestStaticPropV6_GetMaxGPULevel(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetMaxGPULevel() != 0 {
		t.Error("unexpected value for maxGPULevel property")
	}

}

func TestStaticPropV6_GetMinCPULevel(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetMinCPULevel() != 0 {
		t.Error("unexpected value for minCPULevel property")
	}

}

func TestStaticPropV6_GetMinDXLevel(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetMinDXLevel() != 0 {
		t.Error("unexpected value for minDXLevel property")
	}

}

func TestStaticPropV6_GetMinGPULevel(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetMinGPULevel() != 0 {
		t.Error("unexpected value for minGPULevel property")
	}

}

func TestStaticPropV6_GetOrigin(t *testing.T) {
	sut := getStaticPropV6()

	if sut.GetOrigin().X() != 1 &&
		sut.GetOrigin().Y() != 3 &&
		sut.GetOrigin().Z() != 6 {
		t.Error("unexpected value for origin property")
	}
}

func TestStaticPropV6_GetPropType(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetPropType() != 11 {
		t.Error("unexpected value for propType property")
	}

}

func TestStaticPropV6_GetSkin(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetSkin() != 2 {
		t.Error("unexpected value for skin property")
	}

}

func TestStaticPropV6_GetSolid(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetSolid() != 1 {
		t.Error("unexpected value for solid property")
	}

}

func TestStaticPropV6_GetUniformScale(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetUniformScale() != 1 {
		t.Error("unexpected value for uniformScale property")
	}

}

func TestStaticPropV6_GetUnknown(t *testing.T) {
	sut := getStaticPropV6()
	if sut.GetUnknown() != 0 {
		t.Error("unexpected value for unknown property")
	}

}

func getStaticPropV6() *StaticPropV6 {
	return &StaticPropV6{
		Origin:          mgl32.Vec3{1, 3, 6},
		Angles:          mgl32.Vec3{2, 5, 8},
		PropType:        11,
		FirstLeaf:       21,
		LeafCount:       65,
		Solid:           1,
		Flags:           85,
		Skin:            2,
		FadeMinDist:     84.5,
		FadeMaxDist:     256,
		LightingOrigin:  mgl32.Vec3{32, 64, 128},
		ForcedFadeScale: 65,
	}
}
