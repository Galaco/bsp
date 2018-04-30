package primitive


type Primitive struct {
	Type byte
	_ byte
	FirstIndex uint16
	IndexCount uint16
	FirstVert uint16
	VertCount uint16
	_ [2]byte
}