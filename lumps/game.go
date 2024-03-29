package lumps

import (
	"bytes"
	"encoding/binary"
	primitives "github.com/galaco/bsp/primitives/game"
	"strings"
	"unsafe"
)

// Game is Lump 35.
// @TODO NOTE: This really needs per-game implementations to be completely useful,
// otherwise we only get staticprop data from it
type Game struct {
	Generic
	Header              primitives.Header
	GameLumps           []primitives.GenericGameLump
	LumpOffset          int32
	areOffsetsCorrected bool
}

// Unmarshall Imports this lump from raw byte data
func (lump *Game) Unmarshall(raw []byte) (err error) {
	length := len(raw)
	lump.Metadata.SetLength(length)

	if len(raw) == 0 {
		return
	}

	// First reconstruct the header to be of the right size
	lumpCount := binary.LittleEndian.Uint32(raw[:4])
	lump.Header.SetLumpCount(int32(lumpCount))

	// Read header
	lump.Header.GameLumps = make([]primitives.LumpDef, lumpCount)
	headerSize := 4 + (int32(unsafe.Sizeof(primitives.LumpDef{})) * int32(lumpCount))
	err = binary.Read(bytes.NewBuffer(raw[4:headerSize]), binary.LittleEndian, &lump.Header.GameLumps)
	if err != nil {
		return err
	}

	// Correct file offsets
	if !lump.areOffsetsCorrected {
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

	return err
}

// GetData gets internal format structure data
func (lump *Game) GetData() *Game {
	return lump
}

// Marshall dumps this lump back to raw byte data
func (lump *Game) Marshall() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, lump.Header.LumpCount)
	if err != nil {
		return nil, err
	}
	for _, lumpHeader := range lump.Header.GameLumps {
		if err = binary.Write(&buf, binary.LittleEndian, lumpHeader); err != nil {
			return nil, err
		}
	}
	for _, l := range lump.GameLumps {
		if err = binary.Write(&buf, binary.LittleEndian, l.Length); err != nil {
			return nil, err
		}
		if err = binary.Write(&buf, binary.LittleEndian, l.Data); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), err
}

// UpdateInternalOffsets updates the lumps offsets to be relative to the lump, rather
// than the bsp start
func (lump *Game) UpdateInternalOffsets(fileOffset int32) *Game {
	lump.LumpOffset = fileOffset

	return lump
}

// GetStaticPropLump returns the staticprop lump
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
				// This switch is a major hackjob to avoid the need to know what game the bsp is for.
				// Because Valve in all their wisdom have multiple DIFFERENT v10 formats (a true v10,
				// and the MP2013 updated v6 which is REPORTED as v10 as well) we can attempt to infer
				// which format it actually is.
				switch {
				case offset+(int(unsafe.Sizeof(primitives.StaticPropV10{}))*int(numProps)) <= len(sprpLump.Data):
					// Real v10 format
					propLumpSize = int(unsafe.Sizeof(primitives.StaticPropV10{})) * int(numProps)
					vprops := make([]primitives.StaticPropV10, numProps)
					err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
					if err != nil {
						return nil
					}
					for idx := range vprops {
						props[idx] = primitives.IStaticPropDataLump(&vprops[idx])
					}
				case offset+(int(unsafe.Sizeof(primitives.StaticPropV10MP2013{}))*int(numProps)) <= len(sprpLump.Data):
					// Fake v7* 2013MP format.
					propLumpSize = int(unsafe.Sizeof(primitives.StaticPropV10MP2013{})) * int(numProps)
					vprops := make([]primitives.StaticPropV10MP2013, numProps)
					err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:offset+propLumpSize]), binary.LittleEndian, &vprops)
					if err != nil {
						return nil
					}
					for idx := range vprops {
						props[idx] = primitives.IStaticPropDataLump(&vprops[idx])
					}
				default:
					panic("staticpropdata doesn't correspond to a known v10 format")
				}
			case 11:
				vprops := make([]primitives.StaticPropV11, numProps)
				err = binary.Read(bytes.NewBuffer(sprpLump.Data[offset:]), binary.LittleEndian, &vprops)
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
