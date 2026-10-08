package versions

import (
	"fmt"
	"github.com/galaco/bsp/internal/codec"
	"github.com/galaco/bsp/internal/format"
)

// Branch-invariant records share codecs; only actual revisions differ here.
// New codecs extend this policy together with their decode/encode implementation.
func supportedRevision(index int, context format.Context) bool {
	switch index {
	case 22, 23, 24, 25, 49:
		return context.Profile == "reactivedrop" && context.LumpVersion == 0
	case 61:
		return (context.Profile == "l4d2" || context.Profile == "csgo" || context.Profile == "reactivedrop") && context.LumpVersion == 0
	case 63:
		return (context.Profile == "csgo" || context.Profile == "reactivedrop") && context.LumpVersion == 0
	case 10, 7, 8, 27, 53, 58:
		return context.LumpVersion == 0 || context.LumpVersion == 1
	case 9:
		if context.Profile == "reactivedrop" {
			return context.LumpVersion == 2
		}
		return context.LumpVersion == 1 || context.LumpVersion == 2
	case 55, 56:
		return context.LeafVersion == 0 || context.LegacyAmbient || context.LumpVersion == 0 || context.LumpVersion == 1
	case 51, 52:
		return context.LeafVersion == 0 || context.LegacyAmbient || context.LumpVersion == 0
	case 15, 54:
		_, err := codec.WorldLightRecordSize(context)
		return err == nil
	case 19:
		return context.LumpVersion == 0
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
