package lumps

import (
	"encoding/binary"
	"bytes"
	"log"
	primitives "github.com/galaco/bsp/primitives/visibility"
)

/**
	Lump 4: Visibility
 */
type Visibility struct {
	LumpInfo
	data primitives.Vis
}

func (lump Visibility) FromBytes(raw []byte, length int32) ILump {
	err := binary.Read(bytes.NewBuffer(raw[:]), binary.LittleEndian, &lump.data.NumClusters)
	if err != nil {
		log.Fatal(err)
	}
	lump.data.ByteOffset = make([][2]int32, lump.data.NumClusters)
	err = binary.Read(bytes.NewBuffer(raw[4:]), binary.LittleEndian, &lump.data.ByteOffset)
	if err != nil {
		log.Fatal(err)
	}
	bitVectorLength := (length - 4) - (8 * lump.data.NumClusters)
	lump.data.BitVectors = make([]byte, bitVectorLength)
	err = binary.Read(bytes.NewBuffer(raw[length-bitVectorLength:]), binary.LittleEndian, &lump.data.BitVectors)
	if err != nil {
		log.Fatal(err)
	}

	lump.LumpInfo.SetLength(length)

	return lump
}

func (lump Visibility) GetData() interface{} {
	return &lump.data
}

func (lump Visibility) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data.NumClusters)
	for _, offset := range lump.data.ByteOffset {
		binary.Write(&buf, binary.LittleEndian, offset)
	}
	binary.Write(&buf, binary.LittleEndian, lump.data.BitVectors)
	return buf.Bytes()
}

/*

/*
===================
DecompressVis
===================
*/
func (l Visibility) DecompressVis(in *[]byte, decompressed *[]byte) *[]byte {
	var c int
	var out []byte
	var row int
	var inOffset = 0
	var outOffset = 0

	row = int(l.data.NumClusters + 7) >> 3
	out = *decompressed

	hasSimulatedDoWhile := false
	for (int(out[outOffset])) < row || hasSimulatedDoWhile == false {
		hasSimulatedDoWhile = true

		if int((*in)[inOffset]) > 0 {
			inOffset++
			outOffset++
			out[outOffset] = (*in)[inOffset]
			continue
		}

		c = int((*in)[1])
		if c == 0 {
			log.Fatalf("DecompressVis: 0 repeat")
		}

		inOffset += 2
		if (int(out[outOffset])) + c > row {
			c = row - int(out[outOffset])
			log.Printf("warning: Vis decompression overrun\n")
		}

		for c != 0 {
			outOffset++
			out[outOffset] = 0
			c--
		}
	}

	return &out
}
