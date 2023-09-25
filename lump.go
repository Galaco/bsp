package bsp

import (
	"fmt"

	"github.com/galaco/bsp/internal/versions"
	"github.com/galaco/bsp/lumps"
)

// getReferenceLumpByIndex Return an instance of a Lump for a given offset.
func getReferenceLumpByIndex(index int, version int32) (lumps.ILump, error) {
	if index < 0 || index > 63 {
		return nil, fmt.Errorf("invalid lump id: %d provided", index)
	}

	l, err := versions.GetLumpForVersion(int(version), index)
	if err != nil {
		return nil, err
	}

	l.SetVersion(version)

	return l, nil
}
