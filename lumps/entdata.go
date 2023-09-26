package lumps

// EntData is Lump 0: Entdata
type EntData struct {
	Metadata
	Data string `json:"data"`
}

// FromBytes imports this lump from raw byte Data
func (lump *EntData) FromBytes(raw []byte) (err error) {
	lump.Data = string(raw)
	length := len(raw)
	lump.Metadata.SetLength(length)

	return err
}

// Contents returns internal format structure Data
func (lump *EntData) Contents() string {
	return lump.Data
}

// ToBytes converts this lump back to raw byte Data
func (lump *EntData) ToBytes() ([]byte, error) {
	return []byte(lump.Data), nil
}
