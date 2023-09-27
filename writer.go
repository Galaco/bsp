package bsp

import (
	"bytes"
	"encoding/binary"
	"sort"
)

// Writer is a Bsp export writer.
type Writer struct{}

// NewWriter Returns a new bsp writer instance.
func NewWriter() *Writer {
	return &Writer{}
}

// toBytes bsp to []byte.
func (w *Writer) toBytes(data *Bsp) ([]byte, error) {
	// First we need to update the header to reflect any lump changes
	// At the same time we can dump our lumps as bytes to write later
	marshalledLumps := make([][]byte, 64)
	var err error

	currentOffset := 1036 // Header always 1036bytes, so we start immediately afterward.
	order := resolveLumpExportOrder(&data.Header)

	for _, index := range order {
		// Export lump to bytes.
		if marshalledLumps[index], err = data.Lumps[index].ToBytes(); err != nil {
			return nil, err
		}

		// If the lump is empty, we can skip it.
		if len(marshalledLumps[index]) == 0 {
			// 0 bytes is a valid lump, but all fields are 0 in this case.
			// @NOTE. Apparently this isn't actually true; sometimes the offset and version is non-zero.
			//data.Header.Lumps[index].Offset = 0
			data.Header.Lumps[index].Length = 0
			data.Header.Lumps[index].ID = [4]byte{0, 0, 0, 0}
			//data.Header.Lumps[index].Version = 0
			continue
		}

		// Update header with new lump offset and length.
		data.Header.Lumps[index].Offset = int32(currentOffset)
		data.Header.Lumps[index].Length = int32(len(marshalledLumps[index]))

		currentOffset += int(data.Header.Lumps[index].Length)

		// @TODO WTF is this.
		weirdPad := extraWeirdPad(index)
		currentOffset += weirdPad

		// Finally 4byte align the data and current offset.
		// Note that we don't adjust the lump length in the header to reflect this;
		// it's for padding reasons.
		marshalledLumps[index] = append(marshalledLumps[index], make([]byte, weirdPad+(currentOffset%4))...)
		currentOffset += currentOffset % 4
	}

	// Now we can export our bsp.
	var buf bytes.Buffer

	// toBytes Header.
	if err := binary.Write(&buf, binary.LittleEndian, data.Header); err != nil {
		return nil, err
	}

	// toBytes lumps.
	for _, index := range order {
		if _, err := buf.Write(marshalledLumps[index]); err != nil {
			return nil, err
		}
	}
	//for _, lumpData := range marshalledLumps {
	//	if _, err := buf.Write(lumpData); err != nil {
	//		return nil, err
	//	}
	//}

	return buf.Bytes(), nil
}

// resolveLumpExportOrder returns a lump export order that matches the order of the original file
// based on lump offsets.
func resolveLumpExportOrder(header *Header) [64]LumpId {
	vals := make([]struct {
		Offset int32
		Id     int
	}, 64)
	for idx, l := range header.Lumps {
		vals[idx].Offset = l.Offset
		vals[idx].Id = idx
	}

	sort.Slice(vals, func(i, j int) bool {
		return vals[i].Offset < vals[j].Offset
	})

	res := [64]LumpId{}

	for idx := range vals {
		res[idx] = LumpId(vals[idx].Id)
	}
	return res
}

// addWeirdPad is a helper function for adding weird padding to the end of lumps.
// @TODO figure out why this is necessary and remove.
func extraWeirdPad(index LumpId) int {
	switch index {
	case LumpEntities:
		return 1
	case LumpPhysCollide:
		return 1
	case LumpTexDataStringData:
		return -1
	case LumpTexDataStringTable:
		return 0
	default:
		return 0
	}
}
