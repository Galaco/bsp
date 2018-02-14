package lumps

import (
	"bytes"
	"log"
	"encoding/binary"
	"github.com/galaco/bsp/lumps/lumpdata"
)

/**
	Lump interface.
	Organise Lump data in a cleaner and more accessible manner
 */
type ILump interface {
	FromRaw([]byte, int32) ILump

	GetData() interface{}
}


/**
	Lump 0: Entdata
 */
type EntData struct {
	data string
}

func (lump EntData) FromRaw(raw []byte, length int32) ILump {
	lump.data = string(raw[:length])

	return lump
}

func (lump EntData) GetData() interface{} {
	return lump.data
}


/**
	Lump 1: Planes
 */
type Planes struct {
	data []lumpdata.Plane // MAP_MAX_PLANES = 65536
}

func (lump Planes) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump Planes) GetData() interface{} {
	return lump.data
}


/**
	Lump 2: TexData
 */
type TexData struct {
	data []lumpdata.TexData
}
func (lump TexData) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump TexData) GetData() interface{} {
	return lump.data
}


/**
	Lump 3: Vertex
 */

type Vertex struct {
	data [][3]float32
}
func (lump Vertex) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump Vertex) GetData() interface{} {
	return lump.data
}


/**
	Lump 4: Visibility
 */
type Visibility struct {
	data []byte
}

func (lump Visibility) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump Visibility) GetData() interface{} {
	return lump.data
}


/**
	Lump 5: Node
 */
type Node struct {
	data []lumpdata.Node // MAP_MAX_NODES = 65536
}

func (lump Node) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump Node) GetData() interface{} {
	return lump.data
}


/**
	Lump 6: TexInfo
 */

type TexInfo struct {
	data []lumpdata.TexInfo
}

func (lump TexInfo) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump TexInfo) GetData() interface{} {
	return lump.data
}



/**
	Lump 7: Face
 */

type Face struct {
	data []lumpdata.Face
}

func (lump Face) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump Face) GetData() interface{} {
	return lump.data
}


/**
	Lump 10: Leaf
 */
type Leaf struct {
	data []lumpdata.Leaf
}

func (lump Leaf) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump Leaf) GetData() interface{} {
	return lump.data
}


/**
	Lump 12: Edge
 */
type Edge struct {
	data [][2]uint16 // MAX_MAP_EDGES = 256000
}

func (lump Edge) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump Edge) GetData() interface{} {
	return lump.data
}


/**
	Lump 13: Surfedge
 */
type Surfedge struct {
	data []int32 // MAX_MAP_SURFEDGES = 512000
}

func (lump Surfedge) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump Surfedge) GetData() interface{} {
	return lump.data
}


/**
	Lump 14: Model
 */
type Model struct {
	data []lumpdata.Model
}

func (lump Model) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump Model) GetData() interface{} {
	return lump.data
}


/**
	Lump 16: LeafFace
 */
type LeafFace struct {
	data []uint16 // MAX_MAP_LEAFFACES = 65536
}

func (lump LeafFace) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump LeafFace) GetData() interface{} {
	return lump.data
}


/**
	Lump 17: LeafBrush
 */
type LeafBrush struct {
	data []uint16 // MAX_MAP_LEAFBRUSHES = 65536
}

func (lump LeafBrush) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump LeafBrush) GetData() interface{} {
	return lump.data
}


/**
	Lump 18: Brush
 */
type Brush struct {
	data []lumpdata.Brush
}

func (lump Brush) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump Brush) GetData() interface{} {
	return lump.data
}


/**
	Lump 19: BrushSide
 */
type BrushSide struct {
	data []lumpdata.BrushSide // MAX_MAP_BRUSHSIDES = 65536
}

func (lump BrushSide) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump BrushSide) GetData() interface{} {
	return lump.data
}


/**
	Lump 40: Pakfile
 */
type Pakfile struct {
	data []byte
}

func (lump Pakfile) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump Pakfile) GetData() interface{} {
	return lump.data
}

/**
	Lump 35: GameDataLump
	@TODO Isn't fully implemented
 */
type GameData struct {
	data []byte
}

func (lump GameData) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump GameData) GetData() interface{} {
	return lump.data
}



/**
	Lump 43: TexdataStringData
 */
type TexdataStringData struct {
	data string // MAX_MAP_TEXDATA_STRING_DATA = 256000, TEXTURE_NAME_LENGTH = 128
}

func (lump TexdataStringData) FromRaw(raw []byte, length int32) ILump {
	lump.data = string(raw[:length])

	return lump
}

func (lump TexdataStringData) GetData() interface{} {
	return lump.data
}


/**
	Lump 44: TexDataStringTable
 */
type TexDataStringTable struct {
	data []int32 // MAX_MAP_TEXINFO = 2048
}

func (lump TexDataStringTable) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump TexDataStringTable) GetData() interface{} {
	return lump.data
}









/**
	Lump n:
 */
type Lump struct {
	data []byte
}

func (lump Lump) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump Lump) GetData() interface{} {
	return lump.data
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
	Vertex{},
	Visibility{},
	Node{},
	TexInfo{},
	Face{},
	Lump{},
	Lump{},
	Leaf{},
	Lump{},
	Edge{},
	Surfedge{},
	Model{},
	Lump{},
	LeafFace{},
	LeafBrush{},
	Brush{},
	BrushSide{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	GameData{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	TexdataStringData{},
	TexDataStringTable{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
}

