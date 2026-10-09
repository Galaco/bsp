package codec

// DecodeEntData retains the exact string bytes, including NUL terminators.
func DecodeEntData(raw []byte) (string, error) { return string(raw), nil }

func EncodeEntData(data string) ([]byte, error) { return []byte(data), nil }

// EntData is the private string-byte decoder adapter.
type EntData struct {
	Metadata
	data string
}

func (lump *EntData) Decode(raw []byte) error {
	data, err := DecodeEntData(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *EntData) Encode() ([]byte, error) { return EncodeEntData(lump.data) }
func (lump *EntData) GetData() string         { return lump.data }
