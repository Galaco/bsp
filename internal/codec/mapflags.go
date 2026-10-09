package codec

import (
	"encoding/binary"
	"fmt"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

func validateMapFlagsContext(context format.Context) error {
	if err := validatePCContext(context); err != nil {
		return fmt.Errorf("mapflags: %w", err)
	}
	if context.LumpVersion != 0 {
		return fmt.Errorf("mapflags: unsupported revision %d", context.LumpVersion)
	}
	return nil
}

// DecodeMapFlags distinguishes absent bytes from an explicitly present zero word.
func DecodeMapFlags(raw []byte, context format.Context) (lumps.MapFlags, bool, error) {
	if err := validateMapFlagsContext(context); err != nil {
		return lumps.MapFlags{}, false, err
	}
	switch len(raw) {
	case 0:
		return lumps.MapFlags{}, false, nil
	case 4:
		return lumps.MapFlags{LevelFlags: binary.LittleEndian.Uint32(raw)}, true, nil
	default:
		return lumps.MapFlags{}, false, fmt.Errorf("mapflags: %d bytes; expected absent payload or exactly 4 bytes", len(raw))
	}
}

// EncodeMapFlags retains explicit presence and rejects loss of an absent value.
func EncodeMapFlags(value lumps.MapFlags, present bool, context format.Context) ([]byte, error) {
	if err := validateMapFlagsContext(context); err != nil {
		return nil, err
	}
	if !present {
		if value.LevelFlags != 0 {
			return nil, fmt.Errorf("mapflags: nonzero flags cannot be encoded as absent")
		}
		return nil, nil
	}
	raw := make([]byte, 4)
	binary.LittleEndian.PutUint32(raw, value.LevelFlags)
	return raw, nil
}

// MapFlags adapts the optional map-wide flag word for the internal registry.
// Failed decoding preserves the previously decoded value, presence and length.
type MapFlags struct {
	Metadata
	data    lumps.MapFlags
	present bool
}

func (lump *MapFlags) Decode(raw []byte) error {
	data, present, err := DecodeMapFlags(raw, lump.Context())
	if err != nil {
		return err
	}
	lump.data = data
	lump.present = present
	lump.SetLength(len(raw))
	return nil
}

func (lump *MapFlags) Encode() ([]byte, error) {
	return EncodeMapFlags(lump.data, lump.present, lump.Context())
}

func (lump *MapFlags) GetData() lumps.MapFlags { return lump.data }
func (lump *MapFlags) Present() bool           { return lump.present }
