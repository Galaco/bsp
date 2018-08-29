package lumps

import (
	"encoding/binary"
	"bytes"
	"log"
	"unsafe"
	primitives "github.com/galaco/bsp/primitives/game"
)



/**
	Lump 35.
	@TODO NOTE: This really needs per-game implementations to be useful, otherwise we might as well skip reading this lump entirely
 */
type Game struct {
	LumpGeneric
	Header primitives.Header
	GameLumps []byte
}

func (lump *Game) FromBytes(raw []byte, length int32) {
	lump.LumpInfo.SetLength(length)

	if len(raw) == 0 {
		return
	}

	// First reconstruct the header to be of the right size
	lumpCount := binary.LittleEndian.Uint32(raw[:4])
	lump.Header = lump.Header.SetLumpCount(int32(lumpCount))

	// Read header
	lump.Header.GameLumps = make([]primitives.LumpDef, lumpCount)
	headerSize := 4 + (int32(unsafe.Sizeof(primitives.LumpDef{}))*int32(lumpCount))
	err := binary.Read(bytes.NewBuffer(raw[4:headerSize]), binary.LittleEndian, &lump.Header.GameLumps)
	if err != nil {
		log.Fatal(err)
	}

	// Read gamelumps
	// @TODO We dont care about the contents right now
	// In future we should create a set of lumps based on file version
	// For now implementation just copies all lumps to preserve the file, similar to Unimplemented lumps
	payloadLength := length-headerSize
	lump.GameLumps = make([]byte, payloadLength)

	err = binary.Read(
		bytes.NewBuffer(raw[headerSize:length]),
		binary.LittleEndian, &lump.GameLumps)
}

func (lump *Game) GetData() *Game {
	return lump
}

func (lump *Game) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.Header.LumpCount)
	for _,lumpHeader := range lump.Header.GameLumps {
		binary.Write(&buf, binary.LittleEndian, lumpHeader)
	}
	payload := append(buf.Bytes(), lump.GameLumps...)
	return payload
}

func (lump *Game) UpdateInternalOffsets(offsetAdjustment int32) *Game {
	if offsetAdjustment == 0 {
		return lump
	}

	for index,header := range lump.Header.GameLumps {
		header.FileOffset += offsetAdjustment
		lump.Header.GameLumps[index] = header
	}

	return lump
}
