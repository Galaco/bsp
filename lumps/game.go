package lumps

import (
	"encoding/binary"
	"bytes"
	"log"
	"unsafe"
	"github.com/galaco/bsp/lumps/lumpdata"
)



/**
	Lump 35.
	@TODO NOTE: This really needs per-game implementations to be useful, otherwise we might as well skip reading this lump entirely
 */
type Game struct {
	LumpInfo
	Header lumpdata.GameHeader
	GameLumps []byte
}

func (lump Game) FromBytes(raw []byte, length int32) ILump {
	if len(raw) == 0 {
		lump.LumpInfo.SetLength(0)
		return lump
	}

	// First reconstruct the header to be of the right size
	lumpCount := binary.LittleEndian.Uint32(raw[:4])
	lump.Header = lump.Header.SetLumpCount(int32(lumpCount))

	// Read header
	lump.Header.GameLumps = make([]lumpdata.GameLump, lumpCount)
	headerSize := 4 + (int32(unsafe.Sizeof(lumpdata.GameHeader{}))*int32(lumpCount))
	err := binary.Read(bytes.NewBuffer(raw[4:headerSize]), binary.LittleEndian, &lump.Header.GameLumps)
	if err != nil {
		log.Fatal(err)
	}

	// Read gamelumps
	// @TODO We dont care about the contents right now
	// In future we should create a set of lumps based on file version
	// For now implementation just copies all lumps to preserve the file, similar to Unimplemented lumps
	payloadLength := len(raw)-int(headerSize)
	lump.GameLumps = make([]byte, payloadLength)

	err = binary.Read(
		bytes.NewBuffer(raw[headerSize:payloadLength]),
		binary.LittleEndian, &lump.GameLumps)

	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump Game) GetData() interface{} {
	return lump
}

func (lump Game) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.Header)
	binary.Write(&buf, binary.LittleEndian, lump.GameLumps)
	return buf.Bytes()
}

func (lump Game) UpdateInternalOffsets(offsetAdjustment int32) {
	if offsetAdjustment == 0 {
		return
	}

	for index,header := range lump.Header.GameLumps {
		header.FileOffset += offsetAdjustment
		lump.Header.GameLumps[index] = header
	}
}
