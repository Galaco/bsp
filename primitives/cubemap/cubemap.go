package cubemap

import "github.com/galaco/bsp/primitives/common"

type CubemapSample struct {
	Origin common.VectorInt32
	Size byte
	AlignmentPadding [3]byte
}