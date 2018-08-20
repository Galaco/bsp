package lumps

import (
	"encoding/binary"
	"bytes"
	"log"
	primitives "github.com/galaco/bsp/primitives/visibility"
)

// Lump 4: Visibility
type Visibility struct {
	LumpGeneric
	data primitives.Vis
}

// FromBytes
// Populate receiver lump from byte slice
func (lump *Visibility) FromBytes(raw []byte, length int32) ILump {
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

// GetData
// Get internal lump data structure
// Returns interface{} to fulfill interface
// Should be typecasted to expected type
func (lump *Visibility) GetData() interface{} {
	return &lump.data
}

// ToBytes
// Convert internal data structure into a byte slice
func (lump *Visibility) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, lump.data.NumClusters)
	for _, offset := range lump.data.ByteOffset {
		binary.Write(&buf, binary.LittleEndian, offset)
	}
	binary.Write(&buf, binary.LittleEndian, lump.data.BitVectors)
	return buf.Bytes()
}

// GetVisCache
// Determines the Potential Visible Set for a cluster?
func (lump *Visibility) GetVisCache(lastOffset int, cluster int, pvs []byte) int {
	// get the PVS for the pos to limit the number of checks
	if lump.data.NumClusters == 0 {
		for i := range pvs {
			if i < int((lump.data.NumClusters + 7) / 8) {
				pvs[i] = 255
			} else {
				break
			}
		}
		lastOffset = -1
	} else {
		if cluster < 0 {
			// Error, point embedded in wall
			// sampled[0][1] = 255;
			for i := range pvs {
				if i < int((lump.data.NumClusters + 7) / 8) {
					pvs[i] = 255
				} else {
					break
				}
			}
			lastOffset = -1
		} else {
			thisOffset := int(lump.data.ByteOffset[cluster][primitives.DVIS_PVS])
			if thisOffset != lastOffset {
				if thisOffset == -1 {
					log.Fatalf("visofs == -1\n")
				}

				visRunlength := lump.ToBytes()[thisOffset:]
				pvs = lump.DecompressVis(visRunlength, len(pvs))
			}
			lastOffset = thisOffset
		}
	}
	return lastOffset
}

// DecompressVis
// Decompress Visibility BitVectors
// Note: Often we want to decompress only a subset of the compressed data the lump contains. As such,
// target compressed data is passed in rather than derived from the receiver.
func (lump *Visibility) DecompressVis(in []byte, length int) []byte {
	var c int
	var out = make([]byte, length)
	var row int
	var inOffset = 0
	var outOffset = 0

	row = int(lump.data.NumClusters + 7) >> 3

	hasSimulatedDoWhile := false
	for (outOffset < len(out) && int(out[outOffset]) < row) || hasSimulatedDoWhile == false {
		hasSimulatedDoWhile = true

		// @NOTE: The ++ operations may need to shift to the stop
		// In this case, that will cause an out-of-bounds unless we compare to len()-1
		if inOffset < len(in) {
			out[outOffset] = in[inOffset]
			inOffset++
			outOffset++
			continue
		}

		c = int(in[1])
		if c == 0 {
			log.Fatalf("DecompressVis: 0 repeat")
		}

		inOffset += 2
		if (int(out[outOffset])) + c > row {
			c = row - int(out[outOffset])
			log.Printf("warning: Vis decompression overrun\n")
		}

		for c > 0 {
			outOffset++
			out[outOffset] = 0
			c--
		}
	}

	return out
}
