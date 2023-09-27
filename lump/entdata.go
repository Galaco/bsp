package lump

// EntData is Lump 0: Entdata
type EntData struct {
	Metadata
	Data string `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *EntData) FromBytes(raw []byte) error {
	lump.Data = string(raw)

	return nil
}

// Contents returns internal format structure Data
func (lump *EntData) Contents() string {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *EntData) ToBytes() ([]byte, error) {
	return []byte(lump.Data), nil
}
