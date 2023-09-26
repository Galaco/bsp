package game

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestStaticPropV7_GetAngles(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetAngles().X() != 2 &&
		sut.GetAngles().Y() != 5 &&
		sut.GetAngles().Z() != 8 {
		t.Error("unexpected value for angles property")
	}
}

func TestStaticPropV7_GetDiffuseModulation(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetDiffuseModulation() != 986 {
		t.Error("unexpected value for diffuseModulation property")
	}

}

func TestStaticPropV7_GetDisableXBox360(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetDisableXBox360() != false {
		t.Error("unexpected value for unknown property")
	}

}

func TestStaticPropV7_GetFadeMaxDist(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetFadeMaxDist() != 256 {
		t.Error("unexpected value for fadeMaxDist property")
	}

}

func TestStaticPropV7_GetFadeMinDist(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetFadeMinDist() != 84.5 {
		t.Error("unexpected value for fadeMinDist property")
	}

}

func TestStaticPropV7_GetFirstLeaf(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetFirstLeaf() != 21 {
		t.Error("unexpected value for firstLeaf property")
	}

}

func TestStaticPropV7_GetFlags(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetFlags() != 85 {
		t.Error("unexpected value for flags property")
	}

}

func TestStaticPropV7_GetForcedFadeScale(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetForcedFadeScale() != 65 {
		t.Error("unexpected value for forcedFadeScale property")
	}

}

func TestStaticPropV7_GetLeafCount(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetLeafCount() != 65 {
		t.Error("unexpected value for leafCount property")
	}

}

func TestStaticPropV7_GetLightingOrigin(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetLightingOrigin().X() != 32 &&
		sut.GetLightingOrigin().Y() != 64 &&
		sut.GetLightingOrigin().Z() != 128 {
		t.Error("unexpected value for angles property")
	}
}

func TestStaticPropV7_GetMaxCPULevel(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetMaxCPULevel() != 0 {
		t.Error("unexpected value for maxCPULevel property")
	}

}

func TestStaticPropV7_GetMaxDXLevel(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetMaxDXLevel() != 0 {
		t.Error("unexpected value for maxDXLevel property")
	}

}

func TestStaticPropV7_GetMaxGPULevel(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetMaxGPULevel() != 0 {
		t.Error("unexpected value for maxGPULevel property")
	}

}

func TestStaticPropV7_GetMinCPULevel(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetMinCPULevel() != 0 {
		t.Error("unexpected value for minCPULevel property")
	}

}

func TestStaticPropV7_GetMinDXLevel(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetMinDXLevel() != 0 {
		t.Error("unexpected value for minDXLevel property")
	}

}

func TestStaticPropV7_GetMinGPULevel(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetMinGPULevel() != 0 {
		t.Error("unexpected value for minGPULevel property")
	}

}

func TestStaticPropV7_GetOrigin(t *testing.T) {
	sut := getStaticPropV7()

	if sut.GetOrigin().X() != 1 &&
		sut.GetOrigin().Y() != 3 &&
		sut.GetOrigin().Z() != 6 {
		t.Error("unexpected value for origin property")
	}
}

func TestStaticPropV7_GetPropType(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetPropType() != 11 {
		t.Error("unexpected value for propType property")
	}

}

func TestStaticPropV7_GetSkin(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetSkin() != 2 {
		t.Error("unexpected value for skin property")
	}

}

func TestStaticPropV7_GetSolid(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetSolid() != 1 {
		t.Error("unexpected value for solid property")
	}

}

func TestStaticPropV7_GetUniformScale(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetUniformScale() != 1 {
		t.Error("unexpected value for uniformScale property")
	}

}

func TestStaticPropV7_GetUnknown(t *testing.T) {
	sut := getStaticPropV7()
	if sut.GetUnknown() != 0 {
		t.Error("unexpected value for unknown property")
	}

}

func getStaticPropV7() *StaticPropV7 {
	return &StaticPropV7{
		Origin:            mgl32.Vec3{1, 3, 6},
		Angles:            mgl32.Vec3{2, 5, 8},
		PropType:          11,
		FirstLeaf:         21,
		LeafCount:         65,
		Solid:             1,
		Flags:             85,
		Skin:              2,
		FadeMinDist:       84.5,
		FadeMaxDist:       256,
		LightingOrigin:    mgl32.Vec3{32, 64, 128},
		ForcedFadeScale:   65,
		DiffuseModulation: 986,
	}
}
