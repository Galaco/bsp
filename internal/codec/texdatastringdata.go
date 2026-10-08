package codec

// DecodeTextureStringData retains the exact string bytes, including NUL terminators.
func DecodeTextureStringData(raw []byte) (string, error) { return string(raw), nil }

func EncodeTextureStringData(data string) ([]byte, error) { return []byte(data), nil }

// TextureStringData is the private string-byte decoder adapter.
type TextureStringData struct {
	Metadata
	data string
}

func (lump *TextureStringData) Decode(raw []byte) error {
	data, err := DecodeTextureStringData(raw)
	if err != nil {
		return err
	}
	lump.data = data
	lump.SetLength(len(raw))
	return nil
}

func (lump *TextureStringData) Encode() ([]byte, error) { return EncodeTextureStringData(lump.data) }
func (lump *TextureStringData) GetData() string         { return lump.data }
