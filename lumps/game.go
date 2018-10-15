package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/game"
	"log"
	"strings"
	"unsafe"
)

/**
Lump 35.
@TODO NOTE: This really needs per-game implementations to be useful, otherwise we might as well skip reading this lump entirely
*/
type Game struct {
	LumpGeneric
	Header              primitives.Header
	GameLumps           []primitives.GenericGameLump
	LumpOffset          int32
	areOffsetsCorrected bool
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
	headerSize := 4 + (int32(unsafe.Sizeof(primitives.LumpDef{})) * int32(lumpCount))
	err := binary.Read(bytes.NewBuffer(raw[4:headerSize]), binary.LittleEndian, &lump.Header.GameLumps)
	if err != nil {
		log.Fatal(err)
	}

	// Correct file offsets
	if lump.areOffsetsCorrected == false {
		for index := range lump.Header.GameLumps {
			lump.Header.GameLumps[index].FileOffset -= lump.LumpOffset
		}
		lump.areOffsetsCorrected = true
	}

	// Read gamelumps
	lump.GameLumps = make([]primitives.GenericGameLump, lumpCount)
	for i, lumpHeader := range lump.Header.GameLumps {
		lump.GameLumps[i].Length = lumpHeader.FileLength
		lump.GameLumps[i].Data = raw[lumpHeader.FileOffset : lumpHeader.FileOffset+lumpHeader.FileLength]
	}
}

func (lump *Game) GetData() *Game {
	return lump
}

func (lump *Game) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.Header.LumpCount)
	for _, lumpHeader := range lump.Header.GameLumps {
		binary.Write(&buf, binary.LittleEndian, lumpHeader)
	}
	for _, l := range lump.GameLumps {
		binary.Write(&buf, binary.LittleEndian, l)
	}
	return buf.Bytes()
}

// This update the lumps offsets to be relative to the lump, rather
// than the bsp start
func (lump *Game) UpdateInternalOffsets(fileOffset int32) *Game {
	lump.LumpOffset = fileOffset

	return lump
}

func (lump *Game) GetStaticPropLump() *primitives.StaticPropLump {
	for i, gameLump := range lump.Header.GameLumps {
		if gameLump.Id == primitives.StaticPropLumpId {
			sprpLump := lump.GameLumps[i]

			offset := 0

			//dict
			numDicts := int32(0)
			err := binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+4]), binary.LittleEndian, &numDicts)
			if err != nil {
				return nil
			}
			offset += 4
			dicts := primitives.StaticPropDictLump{
				DictEntries: numDicts,
			}
			dictNames := make([]string, numDicts)
			for i := int32(0); i < numDicts; i++ {
				t := make([]byte, 128)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+128]), binary.LittleEndian, &t)
				if err != nil {
					return nil
				}
				dictNames[i] = strings.TrimRight(string(t), "\x00")
				offset += 128
			}
			dicts.Name = dictNames

			//leaf
			numLeafs := int32(0)
			err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+4]), binary.LittleEndian, &numLeafs)
			if err != nil {
				return nil
			}
			offset += 4
			leaf := primitives.StaticPropLeafLump{
				LeafEntries: numLeafs,
			}
			leafs := make([]uint16, numLeafs)
			err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+int(2*numLeafs)]), binary.LittleEndian, &leafs)
			if err != nil {
				return nil
			}
			leaf.Leaf = leafs
			offset += int(2 * numLeafs)

			//props
			numProps := int32(0)
			err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+4]), binary.LittleEndian, &numProps)
			if err != nil {
				return nil
			}
			offset += 4
			props := make([]primitives.IStaticPropDataLump, numProps)
			propLumpSize := 0
			switch gameLump.Version {
			case 4:
				propLumpSize = int(unsafe.Sizeof(primitives.StaticPropV4{})) * int(numProps)
				vprops := make([]primitives.StaticPropV4, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = primitives.IStaticPropDataLump(&vprops[idx])
				}
			case 5:
				propLumpSize = int(unsafe.Sizeof(primitives.StaticPropV5{})) * int(numProps)
				vprops := make([]primitives.StaticPropV5, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = primitives.IStaticPropDataLump(&vprops[idx])
				}
			case 6:
				propLumpSize = int(unsafe.Sizeof(primitives.StaticPropV6{})) * int(numProps)
				vprops := make([]primitives.StaticPropV6, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = primitives.IStaticPropDataLump(&vprops[idx])
				}
			case 7:
				propLumpSize = int(unsafe.Sizeof(primitives.StaticPropV7{})) * int(numProps)
				vprops := make([]primitives.StaticPropV7, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = primitives.IStaticPropDataLump(&vprops[idx])
				}
			case 8:
				propLumpSize = int(unsafe.Sizeof(primitives.StaticPropV8{})) * int(numProps)
				vprops := make([]primitives.StaticPropV8, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = primitives.IStaticPropDataLump(&vprops[idx])
				}
			case 9:
				propLumpSize = int(unsafe.Sizeof(primitives.StaticPropV9{})) * int(numProps)
				vprops := make([]primitives.StaticPropV9, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = primitives.IStaticPropDataLump(&vprops[idx])
				}
			case 10:
				propLumpSize = int(unsafe.Sizeof(primitives.StaticPropV10{})) * int(numProps)
				vprops := make([]primitives.StaticPropV10, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = primitives.IStaticPropDataLump(&vprops[idx])
				}
			case 11:
				propLumpSize = int(unsafe.Sizeof(primitives.StaticPropV11{})) * int(numProps)
				vprops := make([]primitives.StaticPropV11, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = primitives.IStaticPropDataLump(&vprops[idx])
				}
			}

			return &primitives.StaticPropLump{
				DictLump:  dicts,
				LeafLump:  leaf,
				PropLumps: props,
			}
		}
	}

	return nil
}
