package codec

import (
	"encoding/binary"
	"fmt"

	"github.com/galaco/bsp/internal/format"
)

func validatePCContext(context format.Context) error {
	if context.ByteOrder != nil && context.ByteOrder != binary.LittleEndian {
		return fmt.Errorf("unsupported ordinary byte order")
	}
	if context.Profile == "reactivedrop" && context.HeaderLayout != format.StandardHeader {
		return fmt.Errorf("unsupported Reactive Drop header layout")
	}
	switch context.Profile {
	case "sdk2013":
		if context.BSPVersion == 19 || context.BSPVersion == 20 {
			return nil
		}
	case "csgo", "l4d2", "reactivedrop":
		if context.BSPVersion == 21 {
			return nil
		}
	}
	return fmt.Errorf("unsupported ordinary profile %q BSP version %d", context.Profile, context.BSPVersion)
}
