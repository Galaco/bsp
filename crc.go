package bsp

import (
	"hash/crc32"
	"sort"
)

func (bsp *Bsp) CRC32() (uint32, error) {
	const lumpCount = 64

	crc := crc32.NewIEEE()

	lumpList := make([]LumpId, lumpCount)
	for i := 0; i < lumpCount; i++ {
		lumpList[i] = LumpId(i)
	}

	sort.Slice(lumpList, func(i, j int) bool {
		return bsp.Header.Lumps[lumpList[i]].Offset < bsp.Header.Lumps[lumpList[j]].Offset
	})

	for i := 0; i < lumpCount; i++ {
		l := lumpList[i]
		if l == LumpEntities {
			continue
		}

		raw, err := bsp.Lumps[l].ToBytes()
		if err != nil {
			return 0, err
		}
		if _, err := crc.Write(raw); err != nil {
			return 0, err
		}
	}

	// see CRC32_Final.
	res := crc.Sum32() ^ 0xFFFFFFFF

	return res, nil
}
