package codec

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"testing"
)

// This inventory follows the retained fixed arrays in internal/codec. Wire
// fixtures use independent offsets and bytes, never production record encoding.
func TestRetainedFixedArrayWireLayouts(t *testing.T) {
	cases := []struct {
		id         int
		name       string
		width      int
		offset     int
		fieldWidth int
		newCodec   func() Decoder
		samples    func(Decoder) []uint32
	}{
		{1, "planes", 20, 12, 4, func() Decoder { return new(Planes) }, func(c Decoder) []uint32 {
			data := c.(*Planes).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = math.Float32bits(record.Distance)
			}
			return values
		}},
		{2, "texdata", 32, 28, 4, func() Decoder { return new(TextureData) }, func(c Decoder) []uint32 {
			data := c.(*TextureData).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.ViewHeight)
			}
			return values
		}},
		{3, "vertex", 12, 8, 4, func() Decoder { return new(Vertices) }, func(c Decoder) []uint32 {
			data := c.(*Vertices).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = math.Float32bits(record[2])
			}
			return values
		}},
		{5, "node", 32, 8, 4, func() Decoder { return new(Nodes) }, func(c Decoder) []uint32 {
			data := c.(*Nodes).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.Children[1])
			}
			return values
		}},
		{6, "texinfo", 72, 68, 4, func() Decoder { return new(TextureInfo) }, func(c Decoder) []uint32 {
			data := c.(*TextureInfo).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.TexData)
			}
			return values
		}},
		{8, "lighting", 4, 3, 1, func() Decoder { return new(Lighting) }, func(c Decoder) []uint32 {
			data := c.(*Lighting).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(uint8(record.Exponent))
			}
			return values
		}},
		{11, "faceid", 2, 0, 2, func() Decoder { return new(FaceIDs) }, func(c Decoder) []uint32 {
			data := c.(*FaceIDs).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.HammerFaceID)
			}
			return values
		}},
		{12, "edge", 4, 2, 2, func() Decoder { return new(Edges) }, func(c Decoder) []uint32 {
			data := c.(*Edges).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record[1])
			}
			return values
		}},
		{13, "surfedge", 4, 0, 4, func() Decoder { return new(SurfaceEdges) }, func(c Decoder) []uint32 {
			data := c.(*SurfaceEdges).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record)
			}
			return values
		}},
		{14, "model", 48, 44, 4, func() Decoder { return new(Models) }, func(c Decoder) []uint32 {
			data := c.(*Models).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.NumFaces)
			}
			return values
		}},
		{16, "leafface", 2, 0, 2, func() Decoder { return new(LeafFaces) }, func(c Decoder) []uint32 {
			data := c.(*LeafFaces).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record)
			}
			return values
		}},
		{17, "leafbrush", 2, 0, 2, func() Decoder { return new(LeafBrushes) }, func(c Decoder) []uint32 {
			data := c.(*LeafBrushes).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record)
			}
			return values
		}},
		{18, "brush", 12, 8, 4, func() Decoder { return new(Brushes) }, func(c Decoder) []uint32 {
			data := c.(*Brushes).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.Contents)
			}
			return values
		}},
		{20, "area", 8, 4, 4, func() Decoder { return new(Areas) }, func(c Decoder) []uint32 {
			data := c.(*Areas).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.FirstAreaPortal)
			}
			return values
		}},
		{21, "areaportal", 12, 8, 4, func() Decoder { return new(AreaPortals) }, func(c Decoder) []uint32 {
			data := c.(*AreaPortals).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.PlaneNum)
			}
			return values
		}},
		{30, "vertnormal", 12, 8, 4, func() Decoder { return new(VertexNormals) }, func(c Decoder) []uint32 {
			data := c.(*VertexNormals).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = math.Float32bits(record.Pos[2])
			}
			return values
		}},
		{31, "vertnormalindice", 2, 0, 2, func() Decoder { return new(VertexNormalIndices) }, func(c Decoder) []uint32 {
			data := c.(*VertexNormalIndices).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record)
			}
			return values
		}},
		{33, "dispvert", 20, 16, 4, func() Decoder { return new(DisplacementVertices) }, func(c Decoder) []uint32 {
			data := c.(*DisplacementVertices).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = math.Float32bits(record.Alpha)
			}
			return values
		}},
		{36, "leafwaterdata", 12, 8, 2, func() Decoder { return new(LeafWaterData) }, func(c Decoder) []uint32 {
			data := c.(*LeafWaterData).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(uint16(record.SurfaceTexInfoID))
			}
			return values
		}},
		{38, "primvert", 12, 8, 4, func() Decoder { return new(PrimitiveVertices) }, func(c Decoder) []uint32 {
			data := c.(*PrimitiveVertices).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = math.Float32bits(record.Pos[2])
			}
			return values
		}},
		{39, "primindice", 2, 0, 2, func() Decoder { return new(PrimitiveIndices) }, func(c Decoder) []uint32 {
			data := c.(*PrimitiveIndices).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record)
			}
			return values
		}},
		{41, "clipportalverts", 12, 8, 4, func() Decoder { return new(ClipPortalVertices) }, func(c Decoder) []uint32 {
			data := c.(*ClipPortalVertices).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = math.Float32bits(record[2])
			}
			return values
		}},
		{44, "texdatastringtable", 4, 0, 4, func() Decoder { return new(TextureStringTable) }, func(c Decoder) []uint32 {
			data := c.(*TextureStringTable).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record)
			}
			return values
		}},
		{46, "leafmindisttowater", 2, 0, 2, func() Decoder { return new(LeafMinDistancesToWater) }, func(c Decoder) []uint32 {
			data := c.(*LeafMinDistancesToWater).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record)
			}
			return values
		}},
		{47, "facemacrotextureinfo", 2, 0, 2, func() Decoder { return new(FaceMacroTextureInfo) }, func(c Decoder) []uint32 {
			data := c.(*FaceMacroTextureInfo).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.MacroTextureNameID)
			}
			return values
		}},
		{48, "disptris", 2, 0, 2, func() Decoder { return new(DisplacementTriangles) }, func(c Decoder) []uint32 {
			data := c.(*DisplacementTriangles).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = uint32(record.Tags)
			}
			return values
		}},
		{60, "overlayfade", 8, 4, 4, func() Decoder { return new(OverlayFades) }, func(c Decoder) []uint32 {
			data := c.(*OverlayFades).GetData()
			values := make([]uint32, len(data))
			for i, record := range data {
				values[i] = math.Float32bits(record.FadeDistMaxSq)
			}
			return values
		}},
		{53, "lighting", 4, 3, 1, func() Decoder { return new(Lighting) }, func(c Decoder) []uint32 {
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
				if err := c.Decode(raw); err != nil {
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
				encoded, err := c.Encode()
				if err != nil || !bytes.Equal(encoded, raw) {
					t.Fatalf("encode %d records: error %v, payload differs: %x != %x", count, err, encoded, raw)
				}
			}
			for _, length := range []int{tc.width - 1, tc.width + 1, 2*tc.width - 1, 2*tc.width + 1} {
				if err := c.Decode(make([]byte, length)); err == nil {
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
		c := new(Nodes)
		if err := c.Decode(raw); err != nil {
			t.Fatal(err)
		}
		if uint16(c.GetData()[0].Padding) != 0xe792 {
			t.Fatalf("node padding = %#x", c.GetData()[0].Padding)
		}
		c.GetData()[0].PlaneNum = -123
		binary.LittleEndian.PutUint32(raw, 0xffffff85)
		got, err := c.Encode()
		if err != nil || !bytes.Equal(got, raw) {
			t.Fatalf("edited node discarded padding: %x, %v", got, err)
		}
	})
	t.Run("LeafWaterData12", func(t *testing.T) {
		raw := make([]byte, 12)
		raw[10], raw[11] = 0x9b, 0xf2
		c := new(LeafWaterData)
		if err := c.Decode(raw); err != nil {
			t.Fatal(err)
		}
		if c.GetData()[0].Padding != [2]byte{0x9b, 0xf2} {
			t.Fatalf("water padding = %x", c.GetData()[0].Padding)
		}
		c.GetData()[0].SurfaceTexInfoID = -2
		binary.LittleEndian.PutUint16(raw[8:], 0xfffe)
		got, err := c.Encode()
		if err != nil || !bytes.Equal(got, raw) {
			t.Fatalf("edited water data discarded padding: %x, %v", got, err)
		}
	})
}
