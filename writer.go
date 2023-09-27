package bsp

import (
	"bytes"
	"encoding/binary"
	"sort"
	"unsafe"
)

// Writer is a Bsp export writer.
type Writer struct{}

// NewWriter Returns a new bsp writer instance.
func NewWriter() *Writer {
	return &Writer{}
}

// toBytes bsp to []byte.
func (w *Writer) toBytes(data *Bsp) ([]byte, error) {
	var err error

	currentOffset := int(unsafe.Sizeof(Header{})) // Header always 1036bytes, so we start immediately afterward.

	// Now we can export our lumps.
	var lumpBuffer bytes.Buffer

	var temp []byte
	for _, index := range resolveLumpExportOrder(&data.Header) {
		if temp, err = data.Lumps[index].ToBytes(); err != nil {
			return nil, err
		}
		lumpSize, err := lumpBuffer.Write(temp)
		if err != nil {
			return nil, err
		}

		// If the lump is empty, we can skip it.
		if lumpSize == 0 {
			// 0 bytes is a valid lump, but some fields are actually set to 0, other not.
			// @TODO why are some offsets and versions non-0 when length is 0. (e.g. ar_baggage)
			data.Header.Lumps[index].Length = 0
			data.Header.Lumps[index].ID = [4]byte{0, 0, 0, 0}
			continue
		}
		// Update header with new lump offset and lumpSize.
		data.Header.Lumps[index].Offset = int32(currentOffset)
		data.Header.Lumps[index].Length = int32(lumpSize)

		// Update current offset for next iteration.
		currentOffset += lumpSize

		// Finally 4byte align the data and current offset.
		// Note that we don't adjust the lump lumpSize in the header to reflect this;
		// it's for padding reasons.
		if pad := currentOffset % 4; pad != 0 {
			lumpBuffer.Write(make([]byte, 4-pad))
			currentOffset += 4 - pad
		}
	}

	var buffer bytes.Buffer
	if err := binary.Write(&buffer, binary.LittleEndian, data.Header); err != nil {
		return nil, err
	}

	if _, err := buffer.Write(lumpBuffer.Bytes()); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
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
