package lumps

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"testing"
)

// This inventory follows the retained fixed arrays in internal/versions. Wire
// fixtures use independent offsets and bytes, never production record encoding.
func TestRetainedFixedArrayWireLayouts(t *testing.T) {
	cases := []struct {
		id         int
		name       string
		width      int
		offset     int
		fieldWidth int
		newCodec   func() ILump
		samples    func(ILump) []uint32
	}{
		{1, "planes", 20, 12, 4, func() ILump { return new(Planes) }, func(c ILump) []uint32 {
			data := c.(*Planes).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = math.Float32bits(record.Distance)
			}
			return values
		}},
		{2, "texdata", 32, 28, 4, func() ILump { return new(TexData) }, func(c ILump) []uint32 {
			data := c.(*TexData).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.ViewHeight)
			}
			return values
		}},
		{3, "vertex", 12, 8, 4, func() ILump { return new(Vertex) }, func(c ILump) []uint32 {
			data := c.(*Vertex).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = math.Float32bits(record[2])
			}
			return values
		}},
		{5, "node", 32, 8, 4, func() ILump { return new(Node) }, func(c ILump) []uint32 {
			data := c.(*Node).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.Children[1])
			}
			return values
		}},
		{6, "texinfo", 72, 68, 4, func() ILump { return new(TexInfo) }, func(c ILump) []uint32 {
			data := c.(*TexInfo).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.TexData)
			}
			return values
		}},
		{7, "face", 56, 52, 4, func() ILump { return new(Face) }, func(c ILump) []uint32 {
			data := c.(*Face).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.SmoothingGroups)
			}
			return values
		}},
		{8, "lighting", 4, 3, 1, func() ILump { return new(Lighting) }, func(c ILump) []uint32 {
			data := c.(*Lighting).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(uint8(record.Exponent))
			}
			return values
		}},
		{11, "faceid", 2, 0, 2, func() ILump { return new(FaceId) }, func(c ILump) []uint32 {
			data := c.(*FaceId).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.HammerFaceId)
			}
			return values
		}},
		{12, "edge", 4, 2, 2, func() ILump { return new(Edge) }, func(c ILump) []uint32 {
			data := c.(*Edge).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record[1])
			}
			return values
		}},
		{13, "surfedge", 4, 0, 4, func() ILump { return new(Surfedge) }, func(c ILump) []uint32 {
			data := c.(*Surfedge).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record)
			}
			return values
		}},
		{14, "model", 48, 44, 4, func() ILump { return new(Model) }, func(c ILump) []uint32 {
			data := c.(*Model).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.NumFaces)
			}
			return values
		}},
		{16, "leafface", 2, 0, 2, func() ILump { return new(LeafFace) }, func(c ILump) []uint32 {
			data := c.(*LeafFace).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record)
			}
			return values
		}},
		{17, "leafbrush", 2, 0, 2, func() ILump { return new(LeafBrush) }, func(c ILump) []uint32 {
			data := c.(*LeafBrush).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record)
			}
			return values
		}},
		{18, "brush", 12, 8, 4, func() ILump { return new(Brush) }, func(c ILump) []uint32 {
			data := c.(*Brush).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.Contents)
			}
			return values
		}},
		{20, "area", 8, 4, 4, func() ILump { return new(Area) }, func(c ILump) []uint32 {
			data := c.(*Area).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.FirstAreaPortal)
			}
			return values
		}},
		{21, "areaportal", 12, 8, 4, func() ILump { return new(AreaPortal) }, func(c ILump) []uint32 {
			data := c.(*AreaPortal).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.PlaneNum)
			}
			return values
		}},
		{30, "vertnormal", 12, 8, 4, func() ILump { return new(VertNormal) }, func(c ILump) []uint32 {
			data := c.(*VertNormal).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = math.Float32bits(record.Pos[2])
			}
			return values
		}},
		{31, "vertnormalindice", 2, 0, 2, func() ILump { return new(VertNormalIndice) }, func(c ILump) []uint32 {
			data := c.(*VertNormalIndice).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record)
			}
			return values
		}},
		{33, "dispvert", 20, 16, 4, func() ILump { return new(DispVert) }, func(c ILump) []uint32 {
			data := c.(*DispVert).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = math.Float32bits(record.Alpha)
			}
			return values
		}},
		{36, "leafwaterdata", 12, 8, 2, func() ILump { return new(LeafWaterData) }, func(c ILump) []uint32 {
			data := c.(*LeafWaterData).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(uint16(record.SurfaceTexInfoID))
			}
			return values
		}},
		{38, "primvert", 12, 8, 4, func() ILump { return new(PrimVert) }, func(c ILump) []uint32 {
			data := c.(*PrimVert).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = math.Float32bits(record.Pos[2])
			}
			return values
		}},
		{39, "primindice", 2, 0, 2, func() ILump { return new(PrimIndice) }, func(c ILump) []uint32 {
			data := c.(*PrimIndice).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record)
			}
			return values
		}},
		{41, "clipportalverts", 12, 8, 4, func() ILump { return new(ClipPortalVerts) }, func(c ILump) []uint32 {
			data := c.(*ClipPortalVerts).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = math.Float32bits(record[2])
			}
			return values
		}},
		{44, "texdatastringtable", 4, 0, 4, func() ILump { return new(TexDataStringTable) }, func(c ILump) []uint32 {
			data := c.(*TexDataStringTable).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record)
			}
			return values
		}},
		{46, "leafmindisttowater", 2, 0, 2, func() ILump { return new(LeafMinDistToWater) }, func(c ILump) []uint32 {
			data := c.(*LeafMinDistToWater).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record)
			}
			return values
		}},
		{47, "facemacrotextureinfo", 2, 0, 2, func() ILump { return new(FaceMacroTextureInfo) }, func(c ILump) []uint32 {
			data := c.(*FaceMacroTextureInfo).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.MacroTextureNameId)
			}
			return values
		}},
		{48, "disptris", 2, 0, 2, func() ILump { return new(DispTris) }, func(c ILump) []uint32 {
			data := c.(*DispTris).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.Tags)
			}
			return values
		}},
		{58, "facehdr", 56, 52, 4, func() ILump { return new(FaceHDR) }, func(c ILump) []uint32 {
			data := c.(*FaceHDR).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.SmoothingGroups)
			}
			return values
		}},
		{60, "overlayfade", 8, 4, 4, func() ILump { return new(OverlayFade) }, func(c ILump) []uint32 {
			data := c.(*OverlayFade).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = math.Float32bits(record.FadeDistMaxSq)
			}
			return values
		}},
		{27, "face", 56, 52, 4, func() ILump { return new(Face) }, func(c ILump) []uint32 {
			data := c.(*Face).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.SmoothingGroups)
			}
			return values
		}},
		{53, "lighting", 4, 3, 1, func() ILump { return new(Lighting) }, func(c ILump) []uint32 {
			data := c.(*Lighting).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(uint8(record.Exponent))
			}
			return values
		}},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("ID%d/%s", tc.id, tc.name), func(t *testing.T) {
			c := tc.newCodec()
			// Two records precede empty so a valid empty decode must clear data.
			for _, count := range []int{1, 2, 0} {
				raw := make([]byte, count*tc.width)
				want := make([]uint32, count)
				for i := 0; i < count; i++ {
					record := raw[i*tc.width : (i+1)*tc.width]
					for j := range record {
						record[j] = byte(17 + i*53 + j*29)
					}
					// Distinct finite float bits, negative signed values and high
					// unsigned values exercise representation rather than zeroes.
					value := uint32(0xc1234567 + i*0x01010101)
					switch tc.fieldWidth {
					case 1:
						record[tc.offset] = byte(value)
						want[i] = uint32(byte(value))
					case 2:
						binary.LittleEndian.PutUint16(record[tc.offset:], uint16(value))
						want[i] = uint32(uint16(value))
					case 4:
						binary.LittleEndian.PutUint32(record[tc.offset:], value)
						want[i] = value
					}
				}
				if err := c.Unmarshall(raw); err != nil {
					t.Fatalf("decode %d records: %v", count, err)
				}
				got := tc.samples(c)
				if len(got) != count {
					t.Fatalf("decoded %d records, want %d", len(got), count)
				}
				for i, value := range got {
					if value != want[i] {
						t.Errorf("record %d field bits %#x, want %#x", i, value, want[i])
					}
				}
				encoded, err := c.Marshall()
				if err != nil || !bytes.Equal(encoded, raw) {
					t.Fatalf("encode %d records: error %v, payload differs: %x != %x", count, err, encoded, raw)
				}
			}
			for _, length := range []int{tc.width - 1, tc.width + 1, 2*tc.width - 1, 2*tc.width + 1} {
				if err := c.Unmarshall(make([]byte, length)); err == nil {
					t.Errorf("accepted %d bytes at width %d", length, tc.width)
				}
			}
		})
	}
}

func TestFixedArrayNamedPadding(t *testing.T) {
	t.Run("Node32", func(t *testing.T) {
		raw := make([]byte, 32)
		raw[30], raw[31] = 0x92, 0xe7
		c := new(Node)
		if err := c.Unmarshall(raw); err != nil {
			t.Fatal(err)
		}
		if uint16(c.GetData()[0].Padding) != 0xe792 {
			t.Fatalf("node padding = %#x", c.GetData()[0].Padding)
		}
		c.GetData()[0].PlaneNum = -123
		binary.LittleEndian.PutUint32(raw, 0xffffff85)
		got, err := c.Marshall()
		if err != nil || !bytes.Equal(got, raw) {
			t.Fatalf("edited node discarded padding: %x, %v", got, err)
		}
	})
	t.Run("LeafWaterData12", func(t *testing.T) {
		raw := make([]byte, 12)
		raw[10], raw[11] = 0x9b, 0xf2
		c := new(LeafWaterData)
		if err := c.Unmarshall(raw); err != nil {
			t.Fatal(err)
		}
		if c.GetData()[0].Padding != [2]byte{0x9b, 0xf2} {
			t.Fatalf("water padding = %x", c.GetData()[0].Padding)
		}
		c.GetData()[0].SurfaceTexInfoID = -2
		binary.LittleEndian.PutUint16(raw[8:], 0xfffe)
		got, err := c.Marshall()
		if err != nil || !bytes.Equal(got, raw) {
			t.Fatalf("edited water data discarded padding: %x, %v", got, err)
		}
	})
}
