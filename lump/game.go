package lump

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"unsafe"

	primitive "github.com/galaco/bsp/lump/primitive/game"
)

// Game is Lump 35.
// @TODO NOTE: This really needs per-game implementations to be completely useful,
// otherwise we only get staticprop Data from it
type Game struct {
	Metadata
	Header    primitive.Header            `json:"data"`
	GameLumps []primitive.GenericGameLump `json:"gameLumps"`

	// absoluteFileOffset tracks the offset of the game lump into the whole BSP.
	// Game lump has special rules where it contains offsets into the file, not offsets into the lump,
	// so we need to know where the lump is in the file to actually use the offsets.
	// This is required for read game lumps which are offset based.
	absoluteFileOffset int
}

// FromBytes imports this lump from raw byte Data
func (lump *Game) FromBytes(raw []byte) (err error) {
	if len(raw) == 0 {
		return nil
	}

	// First reconstruct the header to be of the right size
	lumpCount := binary.LittleEndian.Uint32(raw[:4])
	lump.Header.SetLumpCount(int32(lumpCount))

	// Read header
	lump.Header.GameLumps = make([]primitive.LumpDef, lumpCount)
	if err := binary.Read(
		bytes.NewBuffer(raw[4:4+(int(unsafe.Sizeof(primitive.LumpDef{}))*int(lumpCount))]),
		binary.LittleEndian,
		&lump.Header.GameLumps,
	); err != nil {
		return err
	}

	// Correct file offsets.
	if lump.absoluteFileOffset == 0 {
		return fmt.Errorf("lump offset not set. Cannot correct offsets")
	}

	// Read gamelumps.
	lump.GameLumps = make([]primitive.GenericGameLump, lumpCount)
	for i, lumpHeader := range lump.Header.GameLumps {
		offset := lump.absoluteToRelativeOffset(int(lumpHeader.FileOffset))
		lump.GameLumps[i].Data = raw[offset : offset+int(lumpHeader.FileLength)]
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
		if err := binary.Write(&buf, binary.LittleEndian, lumpHeader); err != nil {
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

// SetAbsoluteFileOffset updates the lumps offsets to be relative to the lump, rather
// than the bsp start.
func (lump *Game) SetAbsoluteFileOffset(fileOffset int) {
	lump.absoluteFileOffset = fileOffset
}

// GetStaticPropLump returns the staticprop lump.
func (lump *Game) GetStaticPropLump() *primitive.StaticPropLump {
	for i, gameLump := range lump.Header.GameLumps {

		if gameLump.Id == primitive.StaticPropLumpId {
			sprpLump := lump.GameLumps[i]

			offset := 0

			//dict
			numDicts := int32(0)
			err := binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+4]), binary.LittleEndian, &numDicts)
			if err != nil {
				return nil
			}
			offset += 4
			dicts := primitive.StaticPropDictLump{
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
			leaf := primitive.StaticPropLeafLump{
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
			props := make([]primitive.IStaticPropDataLump, numProps)
			propLumpSize := 0
			switch gameLump.Version {
			case 4:
				propLumpSize = int(unsafe.Sizeof(primitive.StaticPropV4{})) * int(numProps)
				vprops := make([]primitive.StaticPropV4, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = primitive.IStaticPropDataLump(&vprops[idx])
				}
			case 5:
				propLumpSize = int(unsafe.Sizeof(primitive.StaticPropV5{})) * int(numProps)
				vprops := make([]primitive.StaticPropV5, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = primitive.IStaticPropDataLump(&vprops[idx])
				}
			case 6:
				propLumpSize = int(unsafe.Sizeof(primitive.StaticPropV6{})) * int(numProps)
				vprops := make([]primitive.StaticPropV6, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = primitive.IStaticPropDataLump(&vprops[idx])
				}
			case 7:
				propLumpSize = int(unsafe.Sizeof(primitive.StaticPropV7{})) * int(numProps)
				vprops := make([]primitive.StaticPropV7, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = primitive.IStaticPropDataLump(&vprops[idx])
				}
			case 8:
				propLumpSize = int(unsafe.Sizeof(primitive.StaticPropV8{})) * int(numProps)
				vprops := make([]primitive.StaticPropV8, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = primitive.IStaticPropDataLump(&vprops[idx])
				}
			case 9:
				propLumpSize = int(unsafe.Sizeof(primitive.StaticPropV9{})) * int(numProps)
				vprops := make([]primitive.StaticPropV9, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = primitive.IStaticPropDataLump(&vprops[idx])
				}
			case 10:
				// This switch is a major hackjob to avoid the need to know what game the bsp is for.
				// Because Valve in all their wisdom have multiple DIFFERENT v10 formats (a true v10,
				// and the MP2013 updated v6 which is REPORTED as v10 as well) we can attempt to infer
				// which format it actually is.
				switch {
				case offset+(int(unsafe.Sizeof(primitive.StaticPropV10{}))*int(numProps)) <= len(sprpLump.Data):
					// Real v10 format
					propLumpSize = int(unsafe.Sizeof(primitive.StaticPropV10{})) * int(numProps)
					vprops := make([]primitive.StaticPropV10, numProps)
					err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
					if err != nil {
						return nil
					}
					for idx := range vprops {
						props[idx] = primitive.IStaticPropDataLump(&vprops[idx])
					}
				case offset+(int(unsafe.Sizeof(primitive.StaticPropV10MP2013{}))*int(numProps)) <= len(sprpLump.Data):
					// Fake v7* 2013MP format.
					propLumpSize = int(unsafe.Sizeof(primitive.StaticPropV10MP2013{})) * int(numProps)
					vprops := make([]primitive.StaticPropV10MP2013, numProps)
					err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
					if err != nil {
						return nil
					}
					for idx := range vprops {
						props[idx] = primitive.IStaticPropDataLump(&vprops[idx])
					}
				default:
					panic("staticpropdata doesn't correspond to a known v10 format")
				}
			case 11:
				vprops := make([]primitive.StaticPropV11, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:]), binary.LittleEndian, &vprops)
				if err != nil {
					return nil
				}
				for idx := range vprops {
					props[idx] = primitive.IStaticPropDataLump(&vprops[idx])
				}
			}

			return &primitive.StaticPropLump{
				DictLump:  dicts,
				LeafLump:  leaf,
				PropLumps: props,
			}
		}
	}

	return nil
}

// absoluteToRelativeOffset converts an absolute offset into the file into a relative offset into the lump.
func (lump *Game) absoluteToRelativeOffset(absoluteOffset int) int {
	return absoluteOffset - lump.absoluteFileOffset
}
