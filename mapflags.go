package bsp

import (
	"github.com/galaco/bsp/internal/codec"
	"github.com/galaco/bsp/lumps"
)

// MapFlags returns an independent map-wide flag value and its explicit presence.
// An absent payload has present=false; a present zero word has present=true.
func (bsp *File) MapFlags() (lumps.MapFlags, bool, error) {
	decoded, err := bsp.decodeLump(LumpMapFlags)
	if err != nil {
		return lumps.MapFlags{}, false, err
	}
	flags := decoded.(*codec.MapFlags)
	return flags.GetData(), flags.Present(), nil
}

// SetMapFlags atomically stores an explicitly present word, including zero.
// It preserves the selected descriptor revision and clears compression metadata.
func (bsp *File) SetMapFlags(value lumps.MapFlags) error {
	if err := bsp.checkTypedRevision(LumpMapFlags); err != nil {
		return err
	}
	raw, err := codec.EncodeMapFlags(value, true, bsp.codecContext(LumpMapFlags))
	return bsp.commitOrdinary(LumpMapFlags, raw, err)
}

// RemoveMapFlags atomically removes the optional word rather than setting zero.
// It preserves the selected descriptor revision and clears compression metadata.
func (bsp *File) RemoveMapFlags() error {
	if err := bsp.checkTypedRevision(LumpMapFlags); err != nil {
		return err
	}
	raw, err := codec.EncodeMapFlags(lumps.MapFlags{}, false, bsp.codecContext(LumpMapFlags))
	return bsp.commitOrdinary(LumpMapFlags, raw, err)
}
