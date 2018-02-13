package lumps

import (
	"bytes"
	"log"
	"encoding/binary"
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
	Generic Lump
	Used for non-existant lumps in the bsp structure (and lumps without an implementation
 */
type Lump struct {
	RawData []byte
}

func (lump Lump) FromRaw(raw []byte, length int32) ILump {
	lump.RawData = raw

	return lump
}

func (lump Lump) GetData() interface{} {
	return lump.RawData
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
type plane_t struct {
	normal [3]float32 // normal vector
	distance float32  // distance from origin
	axisType int32	  // plane axis identifier
}
type Planes struct {
	data []plane_t // MAP_MAX_PLANES = 65536
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
type texData_t struct {
	reflectivity [3]float32
	nameStringTableID int32
	width int32
	height int32
	viewWidth int32
	viewHeight int32
}
type TexData struct {
	data []texData_t
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
	Lump 5: Node
 */
type node_t struct{
	planeNum int32
	children [2]int32
	mins [3]int16
	maxs [3]int16
	firstFace uint16
	numFaces uint16
	area int16
	padding int16
}
type Node struct {
	data []node_t // MAP_MAX_NODES = 65536
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
 type texInfo_t struct {
 	textureVecs [2][4]float32
 	lightmapVecs [2][4]float32
 	flags int32
 	texData int32
 }
type TexInfo struct {
	data []texInfo_t
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
type face_t struct {
 	planenum uint16
 	side byte
 	onNode byte
 	firstEdge int32
 	numEdges int16
 	texInfo int16
 	dispInfo int16
 	surfaceFogVolumeID int16
 	styles [4]byte
 	lightofs int32
 	area float32
 	lightmapTextureMinsInLuxels [2]int32
 	lightmapTextureSizeInLuxels [2]int32
 	origFace int32
 	numPrims uint16
 	firstPrimID uint16
 	smoothingGroups uint32
}

type Face struct {
	data []face_t
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
type leaf_t struct {
	contents int32
	cluster int16
	area int8 // NOTE: Actually first 9 bits of a short, but not implemented
	flags int8 // NOTE: Actually second 7 bits of a short, but not implemented
	mins [3]int16
	maxs [3]int16
	firstLeafFace uint16
	numLeafFaces uint16
	firstLeafBrush uint16
	numLeafBrushes uint16
	leafWaterDataID int16
}
type Leaf struct {
	data []leaf_t
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
 type brush_t struct {
 	firstSize int32
 	numSides  int32
 	contents int32
 }
type Brush struct {
	data []brush_t
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
type brushSide_t struct {
	planeNum uint16
	texInfo int16
	dispInfo int16
	bevel int16
}
type BrushSide struct {
	data []brushSide_t // MAX_MAP_BRUSHSIDES = 65536
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
	Lump :
 */
type Generic struct {
	data []byte
}

func (lump Generic) FromRaw(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data)
	if err != nil {
		log.Fatal(err)
	}

	return lump
}

func (lump Generic) GetData() interface{} {
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
	Lump{},
	Node{},
	TexInfo{},
	Face{},
	Lump{},
	Lump{},
	Leaf{}, //10
	Lump{},
	Edge{},
	Surfedge{},
	Lump{},
	Lump{},
	LeafFace{},
	LeafBrush{},
	Brush{},
	BrushSide{},
	Lump{}, //20
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{}, //30
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{},
	Lump{}, //40
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

