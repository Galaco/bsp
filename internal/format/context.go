// Package format describes wire choices independently of public file metadata.
package format

import "encoding/binary"

const (
	HeaderSize = 1036
	LumpCount  = 64
)

type HeaderLayout uint8

const (
	StandardHeader HeaderLayout = iota
	L4D2Header
)

type GameOffsets uint8

const AbsoluteGameOffsets GameOffsets = 0

// Context is fixed when a codec is selected. It is not a conversion request.
type Context struct {
	Profile         string
	BSPVersion      int32
	LumpVersion     int32
	LegacyAmbient   bool
	ByteOrder       binary.ByteOrder
	HeaderLayout    HeaderLayout
	GameOffsets     GameOffsets
	CompressionSize uint32
}
