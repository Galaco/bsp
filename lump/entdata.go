package lump

import (
	"github.com/galaco/bsp/lump/primitive/entities"
	"strings"
)

// EntData is Lump 0: Entdata
type EntData struct {
	Metadata
	Data []entities.Entity `json:"data"`
}

// FromBytes imports this lump from raw byte Data.
// @TODO this is a very naive implementation.
func (lump *EntData) FromBytes(raw []byte) error {
	var ents []entities.Entity

	// split on quotes after trimming escape char.
	components := strings.Split(string(raw), "\"")

	currentEnt := entities.Entity{}
	// Start at 1 to skip the opening "{".
	for i := 1; i < len(components); i += 4 {

		// format is "key" + " " + "value" + "\n"
		currentEnt = append(currentEnt, entities.KeyValue{Key: components[i], Value: components[i+2]})

		// New entity)
		components[i+3] = strings.Trim(components[i+3], " \n\r\t")
		if components[i+3] == "}\n{" || components[i+3] == "}\n\x00" {
			ents = append(ents, currentEnt)
			currentEnt = entities.Entity{}
			continue
		}
	}

	lump.Data = ents

	return nil
}

// Contents returns internal format structure Data
func (lump *EntData) Contents() []entities.Entity {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data.
// @TODO this is a very naive implementation.
func (lump *EntData) ToBytes() ([]byte, error) {
	if len(lump.Data) == 0 {
		return []byte{}, nil
	}

	var raw string
	for _, ent := range lump.Data {
		raw += "{\n"
		for _, v := range ent {
			raw += "\"" + v.Key + "\" \"" + v.Value + "\"\n"
		}
		raw += "}\n"
	}
	raw += "\x00"

	return []byte(raw), nil
}
