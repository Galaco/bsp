package bsp

import (
	"fmt"

	"github.com/galaco/bsp/internal/codec"
	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

func (bsp *File) staticPropEntry() (codec.GameDirectory, int, format.Context, error) {
	directory, err := bsp.gameDirectory()
	if err != nil {
		return codec.GameDirectory{}, 0, format.Context{}, err
	}
	index, err := gameEntryIndex(directory.Entries, lumps.StaticPropID)
	if err != nil {
		return codec.GameDirectory{}, 0, format.Context{}, err
	}
	context := bsp.codecContext(LumpGame)
	context.LumpVersion = int32(directory.Entries[index].Version)
	if _, err := codec.StaticPropLayoutFor(context); err != nil {
		return codec.GameDirectory{}, 0, format.Context{}, fmt.Errorf("%w: %s", ErrUnsupportedLump, err)
	}
	return directory, index, context, nil
}

// StaticProps returns an owned, profile-selected snapshot of the unique sprp entry.
// Models retain all original slot bytes; Layout describes available record fields.
// Missing or duplicate sprp IDs return errors rather than empty placeholders.
func (bsp *File) StaticProps() (lumps.StaticProps, error) {
	directory, index, context, err := bsp.staticPropEntry()
	if err != nil {
		return lumps.StaticProps{}, err
	}
	raw, err := bsp.decodedGameEntry(directory.Entries[index])
	if err != nil {
		return lumps.StaticProps{}, err
	}
	value, err := codec.DecodeStaticProps(raw, context)
	if err != nil {
		return lumps.StaticProps{}, &FormatError{"decode static props", err}
	}
	if err := bsp.validateStaticProps(value); err != nil {
		return lumps.StaticProps{}, err
	}
	return value, nil
}

// SetStaticProps validates and atomically commits the existing sprp dialect.
// Counts derive from slices. Nonrepresentable fields, invalid model/leaf references
// or a mismatched Layout fail without changing stored bytes or directory metadata.
// Modified compressed props emit uncompressed; siblings, other flags and gaps remain.
func (bsp *File) SetStaticProps(value lumps.StaticProps) error {
	directory, index, context, err := bsp.staticPropEntry()
	if err != nil {
		return err
	}
	if err := bsp.validateStaticProps(value); err != nil {
		return err
	}
	raw, err := codec.EncodeStaticProps(value, context)
	if err != nil {
		return &FormatError{"encode static props", err}
	}
	return bsp.commitDecodedGameEntry(directory, index, raw)
}

func (bsp *File) validateStaticProps(value lumps.StaticProps) error {
	count := 0
	if len(value.LeafIndices) != 0 {
		var err error
		count, err = bsp.leafCount()
		if err != nil {
			return err
		}
	}
	if err := codec.ValidateStaticPropReferences(value, count); err != nil {
		return &FormatError{"static prop references", err}
	}
	return nil
}
