package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

func visibilityTestContext() format.Context {
	return format.Context{Profile: "sdk2013", BSPVersion: 20, ByteOrder: binary.LittleEndian}
}

// Build the wire independently from the production encoder.
func visibilityWire(count int, pairs [][2]int32, data []byte) []byte {
	raw := make([]byte, 4+8*count+len(data))
	binary.LittleEndian.PutUint32(raw[:4], uint32(count))
	for cluster, pair := range pairs {
		for kind, offset := range pair {
			binary.LittleEndian.PutUint32(raw[4+8*cluster+4*kind:], uint32(offset))
		}
	}
	copy(raw[4+8*count:], data)
	return raw
}

func TestVisibilityIndependentWireCounts(t *testing.T) {
	for _, count := range []int{0, 1, 8, 9, 65536} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			pairs := make([][2]int32, count)
			var data []byte
			var want []int
			switch count {
			case 1:
				data, want = []byte{0xff}, []int{0}
			case 8:
				data, want = []byte{0x81}, []int{0, 7}
			case 9:
				data, want = []byte{1, 0x81}, []int{0, 8}
			case 65536:
				for remaining := 8191; remaining > 0; {
					run := remaining
					if run > 255 {
						run = 255
					}
					data = append(data, 0, byte(run))
					remaining -= run
				}
				data = append(data, 0x80)
				want = []int{65535}
			}
			for cluster := range pairs {
				pairs[cluster] = [2]int32{int32(4 + 8*count), int32(4 + 8*count)}
			}
			raw := visibilityWire(count, pairs, data)
			value, err := DecodeVisibility(raw, visibilityTestContext())
			if err != nil || len(value.Offsets) != count || !bytes.Equal(value.Data, data) {
				t.Fatalf("decode count %d: %+v %v", count, value, err)
			}
			encoded, err := EncodeVisibility(value, visibilityTestContext())
			if err != nil || !bytes.Equal(encoded, raw) {
				t.Fatalf("wire differs for count %d: %v", count, err)
			}
			if count == 0 {
				if len(encoded) != 4 {
					t.Fatal("zero count must be a present four-byte record")
				}
				return
			}
			for _, kind := range []lumps.VisibilityKind{lumps.PVS, lumps.PAS} {
				row, present, err := value.Row(count-1, kind)
				if err != nil || !present || len(row) != (count+7)/8 {
					t.Fatalf("row: bytes=%d present=%v error=%v", len(row), present, err)
				}
				visible, present, err := value.VisibleClusters(count-1, kind)
				if err != nil || !present || !reflect.DeepEqual(visible, want) {
					t.Fatalf("visible=%v present=%v error=%v want=%v", visible, present, err, want)
				}
			}
		})
	}
}

func TestVisibilityDistinctPASMissingAndExactRuns(t *testing.T) {
	pairs := make([][2]int32, 9)
	for i := range pairs {
		pairs[i] = [2]int32{-1, -1}
	}
	pairs[0] = [2]int32{76, 78}
	pairs[1] = pairs[0]
	pairs[8] = [2]int32{80, -1}
	raw := visibilityWire(9, pairs, []byte{1, 0x81, 0x80, 0x80, 0, 2})
	value, err := DecodeVisibility(raw, visibilityTestContext())
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		cluster int
		kind    lumps.VisibilityKind
		want    []int
		present bool
	}{
		{0, lumps.PVS, []int{0, 8}, true},
		{1, lumps.PAS, []int{7}, true},
		{8, lumps.PVS, nil, true},
		{8, lumps.PAS, nil, false},
		{2, lumps.PVS, nil, false},
	} {
		visible, present, err := value.VisibleClusters(tc.cluster, tc.kind)
		if err != nil || present != tc.present || len(visible) != len(tc.want) || len(visible) != 0 && !reflect.DeepEqual(visible, tc.want) {
			t.Fatalf("cluster %d kind %d: %v %v %v", tc.cluster, tc.kind, visible, present, err)
		}
	}
	encoded, err := EncodeVisibility(value, visibilityTestContext())
	if err != nil || !bytes.Equal(encoded, raw) {
		t.Fatalf("shared/missing rows changed: %v", err)
	}
}

func TestVisibilityMalformedCountsAndTables(t *testing.T) {
	for _, count := range []uint32{0xffffffff, 65537, 0x7fffffff, 1, 9, 65536} {
		raw := make([]byte, 4)
		binary.LittleEndian.PutUint32(raw, count)
		if _, err := DecodeVisibility(raw, visibilityTestContext()); err == nil {
			t.Fatalf("accepted count %d with no table", count)
		}
	}
	raw := visibilityWire(1, [][2]int32{{12, 12}}, []byte{1})
	for length := 0; length < len(raw); length++ {
		if _, err := DecodeVisibility(raw[:length], visibilityTestContext()); err == nil {
			t.Fatalf("accepted truncated prefix %d", length)
		}
	}
	for length := 0; length < 4; length++ {
		if _, err := DecodeVisibility(raw[:length], visibilityTestContext()); !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("missing count error=%v", err)
		}
	}
}

func TestVisibilityInvalidOffsetsAndRunsBothColumns(t *testing.T) {
	for _, tc := range []struct {
		name   string
		count  int
		offset int32
		data   []byte
	}{
		{"negative", 1, -2, []byte{1}},
		{"header", 1, 0, []byte{1}},
		{"table", 1, 11, []byte{1}},
		{"end", 1, 13, []byte{1}},
		{"far", 1, 0x7fffffff, []byte{1}},
		{"run-count-missing", 1, 12, []byte{0}},
		{"run-zero", 1, 12, []byte{0, 0}},
		{"run-overflow", 1, 12, []byte{0, 2}},
		{"literal-truncated", 9, 76, []byte{1}},
		{"run-after-literal-overflow", 9, 76, []byte{1, 0, 2}},
	} {
		for kind := 0; kind < 2; kind++ {
			t.Run(fmt.Sprintf("%s/%d", tc.name, kind), func(t *testing.T) {
				pairs := make([][2]int32, tc.count)
				for i := range pairs {
					pairs[i] = [2]int32{-1, -1}
				}
				pairs[0][kind] = tc.offset
				raw := visibilityWire(tc.count, pairs, tc.data)
				if _, err := DecodeVisibility(raw, visibilityTestContext()); err == nil {
					t.Fatal("decode accepted malformed row")
				}
				value := lumps.Visibility{Offsets: pairs, Data: tc.data}
				if _, err := EncodeVisibility(value, visibilityTestContext()); err == nil {
					t.Fatal("encode accepted malformed value")
				}
				if _, _, err := value.Row(0, lumps.VisibilityKind(kind)); err == nil {
					t.Fatal("public Row accepted malformed value")
				}
			})
		}
	}
}

func TestVisibilityAdapterDecodeAtomicAndOwned(t *testing.T) {
	var adapter Visibility
	adapter.SetContext(visibilityTestContext())
	raw := visibilityWire(1, [][2]int32{{12, -1}}, []byte{0x80})
	if err := adapter.Decode(raw); err != nil {
		t.Fatal(err)
	}
	before := lumps.Visibility{Offsets: [][2]int32{{12, -1}}, Data: []byte{0x80}}
	length := adapter.Length()
	raw[4], raw[12] = 0, 0
	if adapter.GetData().Offsets[0][0] != 12 || adapter.GetData().Data[0] != 0x80 {
		t.Fatal("decode retained input aliases")
	}
	for _, malformed := range [][]byte{{}, {255, 255, 255, 255}, visibilityWire(1, [][2]int32{{12, 12}}, []byte{0, 0})} {
		if err := adapter.Decode(malformed); err == nil {
			t.Fatal("accepted malformed replacement")
		}
		if adapter.Length() != length || !reflect.DeepEqual(adapter.GetData(), before) {
			t.Fatal("failed decode changed adapter data or length")
		}
	}
}

func TestVisibilityContextsAndSizeBudget(t *testing.T) {
	for _, mutate := range []func(*format.Context){
		func(c *format.Context) { c.LumpVersion = 1 },
		func(c *format.Context) { c.Profile = "unknown" },
		func(c *format.Context) { c.BSPVersion = 21 },
		func(c *format.Context) { c.ByteOrder = binary.BigEndian },
	} {
		context := visibilityTestContext()
		mutate(&context)
		if _, err := DecodeVisibility([]byte{0, 0, 0, 0}, context); err == nil {
			t.Fatal("accepted invalid context")
		}
		if _, err := EncodeVisibility(lumps.Visibility{}, context); err == nil {
			t.Fatal("encode accepted invalid context")
		}
	}
	if _, err := format.VisibilitySize(65537, 0); err == nil {
		t.Fatal("accepted over-limit count")
	}
	if _, err := format.VisibilitySize(1, 1<<31-1); err == nil {
		t.Fatal("accepted overflow payload")
	}
	if size, err := format.VisibilitySize(0, 1<<31-5); err != nil || size != 1<<31-1 {
		t.Fatalf("exact signed budget: %d %v", size, err)
	}
}
