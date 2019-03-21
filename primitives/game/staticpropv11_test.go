package game

import (
	"github.com/go-gl/mathgl/mgl32"
	"testing"
)

func TestStaticPropV11_GetAngles(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetAngles().X() != 2 &&
		sut.GetAngles().Y() != 5 &&
		sut.GetAngles().Z() != 8 {
		t.Error("unexpected value for angles property")
	}
}

func TestStaticPropV11_GetDiffuseModulation(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetDiffuseModulation() != 986 {
		t.Error("unexpected value for diffuseModulation property")
	}

}

func TestStaticPropV11_GetDisableXBox360(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetDisableXBox360() != true {
		t.Error("unexpected value for unknown property")
	}

}

func TestStaticPropV11_GetFadeMaxDist(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetFadeMaxDist() != 256 {
		t.Error("unexpected value for fadeMaxDist property")
	}

}

func TestStaticPropV11_GetFadeMinDist(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetFadeMinDist() != 84.5 {
		t.Error("unexpected value for fadeMinDist property")
	}

}

func TestStaticPropV11_GetFirstLeaf(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetFirstLeaf() != 21 {
		t.Error("unexpected value for firstLeaf property")
	}

}

func TestStaticPropV11_GetFlags(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetFlags() != 85 {
		t.Error("unexpected value for flags property")
	}

}

func TestStaticPropV11_GetForcedFadeScale(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetForcedFadeScale() != 65 {
		t.Error("unexpected value for forcedFadeScale property")
	}

}

func TestStaticPropV11_GetLeafCount(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetLeafCount() != 65 {
		t.Error("unexpected value for leafCount property")
	}

}

func TestStaticPropV11_GetLightingOrigin(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetLightingOrigin().X() != 32 &&
		sut.GetLightingOrigin().Y() != 64 &&
		sut.GetLightingOrigin().Z() != 128 {
		t.Error("unexpected value for angles property")
	}
}

func TestStaticPropV11_GetMaxCPULevel(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetMaxCPULevel() != 31 {
		t.Error("unexpected value for maxCPULevel property")
	}

}

func TestStaticPropV11_GetMaxDXLevel(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetMaxDXLevel() != 0 {
		t.Error("unexpected value for maxDXLevel property")
	}

}

func TestStaticPropV11_GetMaxGPULevel(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetMaxGPULevel() != 14 {
		t.Error("unexpected value for maxGPULevel property")
	}

}

func TestStaticPropV11_GetMinCPULevel(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetMinCPULevel() != 1 {
		t.Error("unexpected value for minCPULevel property")
	}

}

func TestStaticPropV11_GetMinDXLevel(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetMinDXLevel() != 0 {
		t.Error("unexpected value for minDXLevel property")
	}

}

func TestStaticPropV11_GetMinGPULevel(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetMinGPULevel() != 6 {
		t.Error("unexpected value for minGPULevel property")
	}

}

func TestStaticPropV11_GetOrigin(t *testing.T) {
	sut := getStaticPropV11()

	if sut.GetOrigin().X() != 1 &&
		sut.GetOrigin().Y() != 3 &&
		sut.GetOrigin().Z() != 6 {
		t.Error("unexpected value for origin property")
	}
}

func TestStaticPropV11_GetPropType(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetPropType() != 11 {
		t.Error("unexpected value for propType property")
	}

}

func TestStaticPropV11_GetSkin(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetSkin() != 2 {
		t.Error("unexpected value for skin property")
	}

}

func TestStaticPropV11_GetSolid(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetSolid() != 1 {
		t.Error("unexpected value for solid property")
	}

}

func TestStaticPropV11_GetUniformScale(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetUniformScale() != 82 {
		t.Error("unexpected value for uniformScale property")
	}

}

func TestStaticPropV11_GetUnknown(t *testing.T) {
	sut := getStaticPropV11()
	if sut.GetUnknown() != 0 {
		t.Error("unexpected value for unknown property")
	}

}

func getStaticPropV11() *StaticPropV11 {
	return &StaticPropV11{
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
		MinCPULevel:       1,
		MaxCPULevel:       31,
		MinGPULevel:       6,
		MaxGPULevel:       14,
		DiffuseModulation: 986,
		DisableXBox360:    true,
		UniformScale:      82,
	}
}
