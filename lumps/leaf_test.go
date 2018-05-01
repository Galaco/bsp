package lumps

import (
	"testing"
	"unsafe"
	primitives "github.com/galaco/bsp/primitives/leaf"
	"log"
)

const C_STRUCT_SIZE = 32

// Assert leaf data when read from bytes is valid
func TestLeafDataFromBytes(t *testing.T) {
	l := primitives.Leaf{}

	if unsafe.Sizeof(l) != C_STRUCT_SIZE {
		t.Errorf("Leaf struct is of incorrect size, expected: %d, actual: %d", C_STRUCT_SIZE, unsafe.Sizeof(l))
	}

	lump := Leaf{}
	data := lump.FromBytes(GetTestDataBytes(), int32(len(GetTestDataBytes())))
	expected := GetTestLeafData()
	actual := (*(data.GetData().(*([]primitives.Leaf))))[0]


	if actual != expected {
		log.Println("Expected: ")
		log.Println(expected)
		log.Println("Actual: ")
		log.Println(actual)
		t.Errorf("Imported Leaf data mismatch.")
	}
}

func GetTestLeafData() primitives.Leaf {
	l := primitives.Leaf{}
	l.Contents = 1
	l.Cluster = 2
	//l.Name = 61 // ("a")[0]
	l.SetArea(61)
	l.SetFlags(46)
	l.Mins[0] = 3
	l.Mins[1] = 7
	l.Mins[2] = 14
	l.Maxs[0] = 13
	l.Maxs[1] = 36
	l.Maxs[2] = 51
	l.FirstLeafFace = 32
	l.NumLeafFaces = 11
	l.FirstLeafBrush = 21
	l.NumLeafBrushes = 76
	l.LeafWaterDataID = 621

	return l
}


func GetTestDataBytes() []byte{
	return []byte{
		1,	// contents
		0,
		0,
		0,
		2, // cluster
		0,
		61, // BitField
		92,
		3, // Mins
		0,
		7,
		0,
		14,
		0,
		13, // Maxs
		0,
		36,
		0,
		51,
		0,
		32, // FirstLeafFace
		0,
		11, // NumLeafFaces
		0,
		21, // FirstLeafBrush
		0,
		76, // NumLeafBrushes
		0,
		109, // leafwaterdataid
		2,
		0,
		0,
	}
}