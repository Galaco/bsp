package lumps

import (
	"log"
)




/**
	Lump interface.
	Organise Lump data in a cleaner and more accessible manner
 */
type ILump interface {
	FromBytes([]byte, int32) ILump

	GetData() interface{}

	ToBytes() []byte
}

/**
	Helper info for a lump
 */
type LumpInfo struct {
	length int32
}
func (info LumpInfo) GetLength() int32 {
	return info.length
}
func (info LumpInfo) SetLength(length int32) {
	info.length = length
}





/**
	Return an instance of a Lump for a given offset.
 */
func GetLumpForIndex(index int) ILump {
	if index < 0 || index > 63 {
		log.Fatal("Invalid lump index provided.")
	}

	return lumpMap[index]
}


var lumpMap = [64]ILump{
	EntData{},
	Planes{},
	TexData{},
	Unimplemented{},
	Visibility{},
	Node{},
	TexInfo{},
	Face{},
	Unimplemented{},
	Unimplemented{},
	Leaf{},
	Unimplemented{},
	Edge{},
	Surfedge{},
	Model{},
	Unimplemented{},
	LeafFace{},
	LeafBrush{},
	Brush{},
	BrushSide{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Game{}, //
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	TexdataStringData{},
	TexDataStringTable{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
	Unimplemented{},
}

