package bsp

import "hash/crc32"

// CRC32 computes the SDK map checksum of the current planned export: numeric
// IDs1..63, final stored bytes including relocated game offsets, no header or
// alignment padding. Untouched compressed/opaque payloads are not decoded.
//
// Valve SDK b8cfb12c0e083a2ef5b2f9f9b50f3902fa034474 bsplib.cpp initializes
// CRC_MapFile's state with CRC32_Init and consumes it without CRC32_Final.
// Go's accumulator is already finalized, so the final complement below restores
// that SDK map-state convention; this is not the standard finalized ZIP CRC.
func (file *Bsp) CRC32() (uint32, error) {
	plan, err := file.prepareExport()
	if err != nil {
		return 0, err
	}
	var crc uint32
	for id := LumpId(1); id < LumpCount; id++ {
		if id == LumpGame {
			crc = plan.game.UpdateCRC32(crc)
		} else {
			crc = crc32.Update(crc, crc32.IEEETable, file.lumps[id].raw)
		}
	}
	return crc ^ 0xffffffff, nil
}
