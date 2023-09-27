package lump

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"unsafe"

	game2 "github.com/galaco/bsp/lump/primitive/game"
)

// Game is Lump 35.
// @TODO NOTE: This really needs per-game implementations to be completely useful,
// otherwise we only get staticprop Data from it
type Game struct {
	Metadata
	Header    game2.Header            `json:"data"`
	GameLumps []game2.GenericGameLump `json:"gameLumps"`

	// offsetOfLumpIntoFile tracks the offset of the game lump into the whole BSP.
	// Game lump has special rules where it contains offsets into the file, not offsets into the lump,
	// so we need to know where the lump is in the file to actually use the offsets.
	offsetOfLumpIntoFile int32
}

// FromBytes imports this lump from raw byte Data
func (lump *Game) FromBytes(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)

	if len(raw) == 0 {
		return nil
	}

	// First reconstruct the header to be of the right size
	lumpCount := binary.LittleEndian.Uint32(raw[:4])
	lump.Header.SetLumpCount(int32(lumpCount))

	// Read header
	lump.Header.GameLumps = make([]game2.LumpDef, lumpCount)
	if err := binary.Read(
		bytes.NewBuffer(raw[4:4+(int(unsafe.Sizeof(game2.LumpDef{}))*int(lumpCount))]),
		binary.LittleEndian,
		&lump.Header.GameLumps,
	); err != nil {
		return err
	}

	// Correct file offsets.
	if lump.offsetOfLumpIntoFile == 0 {
		return fmt.Errorf("lump offset not set. Cannot correct offsets")
	}

	// This makes the absolute file offset into a relative offset of this lump.
	for index := range lump.Header.GameLumps {
		lump.Header.GameLumps[index].FileOffset -= lump.offsetOfLumpIntoFile
	}

	// Read gamelumps.
	lump.GameLumps = make([]game2.GenericGameLump, lumpCount)
	for i, lumpHeader := range lump.Header.GameLumps {
		lump.GameLumps[i].Data = raw[lumpHeader.FileOffset : lumpHeader.FileOffset+lumpHeader.FileLength]
	}

	return nil
}

// Contents returns internal format structure Data
func (lump *Game) Contents() *Game {
	return lump
}

// ToBytes converts this lump back to raw byte Data
func (lump *Game) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, lump.Header.LumpCount); err != nil {
		return nil, err
	}
	for _, lumpHeader := range lump.Header.GameLumps {
		h := lumpHeader

		// Set the file offset back to the actual file offset, not the relative offset of the game lump.
		h.FileOffset += lump.offsetOfLumpIntoFile
		if err := binary.Write(&buf, binary.LittleEndian, h); err != nil {
			return nil, err
		}
	}
	for _, l := range lump.GameLumps {
		if err := binary.Write(&buf, binary.LittleEndian, l.Data); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// UpdateInternalOffsets updates the lumps offsets to be relative to the lump, rather
// than the bsp start.
func (lump *Game) UpdateInternalOffsets(fileOffset int32) *Game {
	lump.offsetOfLumpIntoFile = fileOffset

	return lump
}

// GetStaticPropLump returns the staticprop lump.
func (lump *Game) GetStaticPropLump() *game2.StaticPropLump {
	for i, gameLump := range lump.Header.GameLumps {

		if gameLump.Id == game2.StaticPropLumpId {
			sprpLump := lump.GameLumps[i]

			offset := 0

			//dict
			numDicts := int32(0)
			err := binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+4]), binary.LittleEndian, &numDicts)
			if err != nil {
				return nil
			}
			offset += 4
			dicts := game2.StaticPropDictLump{
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
			leaf := game2.StaticPropLeafLump{
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
			props := make([]game2.IStaticPropDataLump, numProps)
			propLumpSize := 0
			switch gameLump.Version {
			case 4:
				propLumpSize = int(unsafe.Sizeof(game2.StaticPropV4{})) * int(numProps)
				vprops := make([]game2.StaticPropV4, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = game2.IStaticPropDataLump(&vprops[idx])
				}
			case 5:
				propLumpSize = int(unsafe.Sizeof(game2.StaticPropV5{})) * int(numProps)
				vprops := make([]game2.StaticPropV5, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = game2.IStaticPropDataLump(&vprops[idx])
				}
			case 6:
				propLumpSize = int(unsafe.Sizeof(game2.StaticPropV6{})) * int(numProps)
				vprops := make([]game2.StaticPropV6, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = game2.IStaticPropDataLump(&vprops[idx])
				}
			case 7:
				propLumpSize = int(unsafe.Sizeof(game2.StaticPropV7{})) * int(numProps)
				vprops := make([]game2.StaticPropV7, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = game2.IStaticPropDataLump(&vprops[idx])
				}
			case 8:
				propLumpSize = int(unsafe.Sizeof(game2.StaticPropV8{})) * int(numProps)
				vprops := make([]game2.StaticPropV8, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = game2.IStaticPropDataLump(&vprops[idx])
				}
			case 9:
				propLumpSize = int(unsafe.Sizeof(game2.StaticPropV9{})) * int(numProps)
				vprops := make([]game2.StaticPropV9, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = game2.IStaticPropDataLump(&vprops[idx])
				}
			case 10:
				// This switch is a major hackjob to avoid the need to know what game the bsp is for.
				// Because Valve in all their wisdom have multiple DIFFERENT v10 formats (a true v10,
				// and the MP2013 updated v6 which is REPORTED as v10 as well) we can attempt to infer
				// which format it actually is.
				switch {
				case offset+(int(unsafe.Sizeof(game2.StaticPropV10{}))*int(numProps)) <= len(sprpLump.Data):
					// Real v10 format
					propLumpSize = int(unsafe.Sizeof(game2.StaticPropV10{})) * int(numProps)
					vprops := make([]game2.StaticPropV10, numProps)
					err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
					if err != nil {
						return nil
					}
					for idx := range vprops {
						props[idx] = game2.IStaticPropDataLump(&vprops[idx])
					}
				case offset+(int(unsafe.Sizeof(game2.StaticPropV10MP2013{}))*int(numProps)) <= len(sprpLump.Data):
					// Fake v7* 2013MP format.
					propLumpSize = int(unsafe.Sizeof(game2.StaticPropV10MP2013{})) * int(numProps)
					vprops := make([]game2.StaticPropV10MP2013, numProps)
					err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
					if err != nil {
						return nil
					}
					for idx := range vprops {
						props[idx] = game2.IStaticPropDataLump(&vprops[idx])
					}
				default:
					panic("staticpropdata doesn't correspond to a known v10 format")
				}
			case 11:
				vprops := make([]game2.StaticPropV11, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = game2.IStaticPropDataLump(&vprops[idx])
				}
			}

			return &game2.StaticPropLump{
				DictLump:  dicts,
				LeafLump:  leaf,
				PropLumps: props,
			}
		}
	}

	return nil
}
