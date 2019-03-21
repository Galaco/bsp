package game

import "testing"

func TestHeader_SetLumpCount(t *testing.T) {
	sut := Header{}

	sut.SetLumpCount(11)

	if sut.LumpCount != 11 {
		t.Error("incorrect lump count for header")
	}

	if len(sut.GameLumps) != 11 {
		t.Error("failed to allocate correct number of lumps")
	}
}
