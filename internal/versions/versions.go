package versions

import (
	"fmt"
	"github.com/galaco/bsp/internal/format"
)

// Branch-invariant records share codecs; only actual revisions differ here.
// New codecs extend this policy together with their decode/encode implementation.
func supportedRevision(index int, context format.Context) bool {
	switch index {
	case 10, 7, 8, 27, 53, 58:
		return context.LumpVersion == 0 || context.LumpVersion == 1
	case 9:
		return context.LumpVersion == 2
	case 55, 56:
		return (context.LumpVersion == 0 || context.LumpVersion == 1) && !context.LegacyAmbient
	case 19:
		return context.Profile == "sdk2013" && context.LumpVersion == 0
	default:
		return context.LumpVersion == 0
	}
}

func ValidateRevision(index int, context format.Context) error {
	if !supportedRevision(index, context) {
		return fmt.Errorf("unsupported lump %d revision %d for %s", index, context.LumpVersion, context.Profile)
	}
	return nil
}
