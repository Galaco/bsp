package codec

import (
	"bytes"
	"encoding/binary"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

// This fixture is authored directly at the pinned C++ wire offsets. No
// production values or encoders participate in creating these two records.
func displacementWireFixture() []byte {
	raw := make([]byte, 352)
	for i := 0; i < 2; i++ {
		wire := raw[i*176 : (i+1)*176]
		for slot, bits := range []uint32{0x80000000, 0x7fc12345, 0xff800000} {
			binary.LittleEndian.PutUint32(wire[slot*4:], bits+uint32(i))
		}
		binary.LittleEndian.PutUint32(wire[12:], uint32(25*i))
		binary.LittleEndian.PutUint32(wire[16:], uint32(32*i))
		binary.LittleEndian.PutUint32(wire[20:], uint32(2+2*i))
		binary.LittleEndian.PutUint32(wire[24:], []uint32{3, 0xc1234567}[i])
		binary.LittleEndian.PutUint32(wire[28:], 0x7fcabcde+uint32(i))
		binary.LittleEndian.PutUint32(wire[32:], 0xfedcba98+uint32(i))
		binary.LittleEndian.PutUint16(wire[36:], uint16(501+i))
		wire[38], wire[39] = byte(0xa1+i), byte(0xf2+i)
		binary.LittleEndian.PutUint32(wire[40:], 0xffffffff-uint32(i))
		binary.LittleEndian.PutUint32(wire[44:], uint32(14*i))
		for edge := 0; edge < 4; edge++ {
			for sub := 0; sub < 2; sub++ {
				offset := 48 + edge*12 + sub*6
				binary.LittleEndian.PutUint16(wire[offset:], 0xffff)
				wire[offset+2], wire[offset+3], wire[offset+4], wire[offset+5] = 0xe1, 0xf2, 0xd3, byte(0x81+i*8+edge*2+sub)
				if sub == i {
					binary.LittleEndian.PutUint16(wire[offset:], uint16(1-i))
					wire[offset+2], wire[offset+3], wire[offset+4] = byte(edge), byte(edge%3), byte((edge+1)%3)
				}
			}
		}
		for corner := 0; corner < 4; corner++ {
			offset := 96 + corner*10
			for slot := 0; slot < 4; slot++ {
				value := uint16(0xbe00 + i*16 + corner*4 + slot)
				if slot <= corner {
					value = uint16(1 - i)
				}
				// Even an active sentinel must remain exactly a sentinel.
				if corner == 3 && slot == 3 {
					value = 0xffff
				}
				binary.LittleEndian.PutUint16(wire[offset+slot*2:], value)
			}
			wire[offset+8], wire[offset+9] = byte(corner+1), byte(0xd1+i*4+corner)
		}
		for word := 0; word < 10; word++ {
			binary.LittleEndian.PutUint32(wire[136+word*4:], 0x81000001+uint32(i*16+word*2))
		}
	}
	return raw
}

func TestDisplacementsIndependentWire(t *testing.T) {
	for _, context := range spatialContexts() {
		raw := displacementWireFixture()
		records, err := DecodeDisplacements(raw, context)
		if err != nil || len(records) != 2 {
			t.Fatalf("%+v: decode len %d: %v", context, len(records), err)
		}
		for i, record := range records {
			for slot, bits := range []uint32{0x80000000, 0x7fc12345, 0xff800000} {
				if math.Float32bits(record.StartPosition[slot]) != bits+uint32(i) {
					t.Fatalf("record %d start slot %d bits lost", i, slot)
				}
			}
			if record.DispVertStart != int32(25*i) || record.DispTriStart != int32(32*i) || record.Power != int32(2+2*i) || record.Contents != -19088744+int32(i) || record.MapFace != uint16(501+i) || record.MapFacePadding != [2]byte{byte(0xa1 + i), byte(0xf2 + i)} || record.LightmapAlphaStart != int32(-1-i) || record.LightmapSampleStart != int32(14*i) || math.Float32bits(record.SmoothingAngle) != 0x7fcabcde+uint32(i) {
				t.Fatalf("record %d header fields lost: %+v", i, record)
			}
			if i == 0 && (record.MinTess != 3 || record.Flags != 0) || i == 1 && (record.MinTess != 0 || record.Flags != 0xc1234567) {
				t.Fatalf("record %d byte-24 interpretation: MinTess %d Flags %#x", i, record.MinTess, record.Flags)
			}
			for edge, neighbors := range record.EdgeNeighbors {
				for sub, neighbor := range neighbors.SubNeighbors {
					want := lumps.DisplacementSubNeighbor{Index: 0xffff, Orientation: 0xe1, Span: 0xf2, NeighborSpan: 0xd3, Padding: byte(0x81 + i*8 + edge*2 + sub)}
					if sub == i {
						want.Index, want.Orientation, want.Span, want.NeighborSpan = uint16(1-i), byte(edge), byte(edge%3), byte((edge+1)%3)
					}
					if neighbor != want {
						t.Fatalf("record %d edge %d sub %d: %+v want %+v", i, edge, sub, neighbor, want)
					}
				}
			}
			for corner, neighbor := range record.CornerNeighbors {
				if neighbor.Count != byte(corner+1) || neighbor.Padding != byte(0xd1+i*4+corner) {
					t.Fatalf("record %d corner %d count/padding: %+v", i, corner, neighbor)
				}
				for slot, value := range neighbor.Neighbors {
					want := uint16(0xbe00 + i*16 + corner*4 + slot)
					if slot <= corner {
						want = uint16(1 - i)
					}
					if corner == 3 && slot == 3 {
						want = 0xffff
					}
					if value != want {
						t.Fatalf("record %d corner %d slot %d: %#x want %#x", i, corner, slot, value, want)
					}
				}
			}
			for word, value := range record.AllowedVerts {
				if want := uint32(0x81000001 + i*16 + word*2); value != want {
					t.Fatalf("record %d allowed word %d: %#x want %#x", i, word, value, want)
				}
			}
		}
		encoded, err := EncodeDisplacements(records, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("%+v: bytewise reencode: %v", context, err)
		}
		records[1].AllowedVerts[9] = 0xfedcba01
		records[1].EdgeNeighbors[3].SubNeighbors[1].Padding = 0xee
		records[1].CornerNeighbors[0].Neighbors[3] = 0xabcd
		records[0].MapFace = 17
		records[1].Flags = 0x80000000
		binary.LittleEndian.PutUint32(raw[176+172:], 0xfedcba01)
		raw[176+48+3*12+6+5] = 0xee
		binary.LittleEndian.PutUint16(raw[176+96+6:], 0xabcd)
		binary.LittleEndian.PutUint16(raw[36:], 17)
		binary.LittleEndian.PutUint32(raw[176+24:], 0x80000000)
		encoded, err = EncodeDisplacements(records, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("%+v: field edits changed unrelated bytes: %v", context, err)
		}
	}
}

func TestDisplacementRejectsInvalidWireAndValues(t *testing.T) {
	context := format.Context{Profile: "sdk2013", BSPVersion: 20}
	for _, test := range []struct {
		name    string
		offset  int
		width   int
		value   uint32
		message string
	}{
		{"negative power", 20, 4, 0xffffffff, "power"},
		{"small power", 20, 4, 1, "power"},
		{"large power", 20, 4, 5, "power"},
		{"shift overflow power", 20, 4, 0x7fffffff, "power"},
		{"negative vertex start", 12, 4, 0xffffffff, "vertex span"},
		{"vertex span overflow", 12, 4, 0x7fffffff, "vertex span"},
		{"negative tri start", 16, 4, 0x80000000, "triangle span"},
		{"triangle span overflow", 16, 4, 0x7fffffff, "triangle span"},
		{"neighbor index", 48, 2, 2, "index"},
		{"orientation", 50, 1, 4, "orientation"},
		{"span", 51, 1, 3, "span"},
		{"neighbor span", 52, 1, 255, "span"},
		{"corner count", 104, 1, 5, "count"},
		{"corner index", 96, 2, 2, "index"},
	} {
		t.Run(test.name, func(t *testing.T) {
			raw := displacementWireFixture()
			switch test.width {
			case 1:
				raw[test.offset] = byte(test.value)
			case 2:
				binary.LittleEndian.PutUint16(raw[test.offset:], uint16(test.value))
			case 4:
				binary.LittleEndian.PutUint32(raw[test.offset:], test.value)
			}
			if _, err := DecodeDisplacements(raw, context); err == nil || !strings.Contains(err.Error(), test.message) || !strings.Contains(err.Error(), "displacement 0") {
				t.Fatalf("invalid wire error = %v", err)
			}
		})
	}
	for _, length := range []int{1, 175, 177, 351, 353} {
		raw := make([]byte, length)
		if _, err := DecodeDisplacements(raw, context); err == nil {
			t.Fatalf("partial record length %d accepted", length)
		}
	}
	for _, mutate := range []func(*lumps.Displacement){
		func(r *lumps.Displacement) { r.Power = -1 },
		func(r *lumps.Displacement) { r.Power = 31 },
		func(r *lumps.Displacement) { r.DispVertStart = -1 },
		func(r *lumps.Displacement) { r.DispTriStart = math.MaxInt32 },
		func(r *lumps.Displacement) { r.EdgeNeighbors[0].SubNeighbors[0].Index = 2 },
		func(r *lumps.Displacement) { r.EdgeNeighbors[0].SubNeighbors[0].Orientation = 4 },
		func(r *lumps.Displacement) { r.EdgeNeighbors[0].SubNeighbors[0].Span = 3 },
		func(r *lumps.Displacement) { r.EdgeNeighbors[0].SubNeighbors[0].NeighborSpan = 3 },
		func(r *lumps.Displacement) { r.CornerNeighbors[0].Count = 5 },
		func(r *lumps.Displacement) { r.CornerNeighbors[0].Neighbors[0] = 2 },
		func(r *lumps.Displacement) { r.Flags = 1 },
		func(r *lumps.Displacement) { r.Flags = 0x80000000 },
		func(r *lumps.Displacement) { r.MinTess = -1 },
	} {
		records, err := DecodeDisplacements(displacementWireFixture(), context)
		if err != nil {
			t.Fatal(err)
		}
		mutate(&records[0])
		if raw, err := EncodeDisplacements(records, context); err == nil || raw != nil {
			t.Fatalf("invalid value encoded: %x, %v", raw, err)
		}
	}
}

func TestDisplacementContextAndAdapter(t *testing.T) {
	for _, context := range []format.Context{
		{Profile: "sdk2013", BSPVersion: 21},
		{Profile: "csgo", BSPVersion: 20},
		{Profile: "unknown", BSPVersion: 20},
		{Profile: "sdk2013", BSPVersion: 20, LumpVersion: 1},
		{Profile: "sdk2013", BSPVersion: 20, ByteOrder: binary.BigEndian},
	} {
		if _, err := DecodeDisplacements(nil, context); err == nil {
			t.Fatalf("decode accepted context %+v", context)
		}
		if _, err := EncodeDisplacements(nil, context); err == nil {
			t.Fatalf("encode accepted context %+v", context)
		}
	}
	adapter := new(Displacements)
	adapter.SetContext(format.Context{Profile: "sdk2013", BSPVersion: 20})
	raw := displacementWireFixture()
	if err := adapter.Decode(raw); err != nil {
		t.Fatal(err)
	}
	if adapter.Length() != 352 || len(adapter.GetData()) != 2 {
		t.Fatal("adapter metadata/data incorrect")
	}
	if err := adapter.Decode(raw[:351]); err == nil {
		t.Fatal("adapter accepted partial record")
	}
	encoded, err := adapter.Encode()
	if err != nil || !bytes.Equal(encoded, raw) || adapter.Length() != 352 {
		t.Fatalf("failed adapter decode changed previous state: %v", err)
	}
}

func TestDisplacementSamplesTransitionsAndConsumption(t *testing.T) {
	raw := []byte{254, 11, 22, 33, 255, 0, 44, 55, 66, 255, 255, 77, 88, 99, 0xde, 0xad}
	want := []lumps.DisplacementSample{
		{Index: 254, Barycentric: [3]byte{11, 22, 33}},
		{Index: 255, Barycentric: [3]byte{44, 55, 66}},
		{Index: 510, Barycentric: [3]byte{77, 88, 99}},
	}
	records, used, err := DecodeDisplacementSamples(raw, 3, 512)
	if err != nil || used != 14 || !reflect.DeepEqual(records, want) {
		t.Fatalf("samples = %+v consumed %d err %v", records, used, err)
	}
	encoded, err := EncodeDisplacementSamples(records, 512)
	if err != nil || !bytes.Equal(encoded, raw[:14]) {
		t.Fatalf("sample reencode = %x, %v", encoded, err)
	}
	records[0].Index, records[1].Index = 255, 254
	records[2].Barycentric = [3]byte{0, 255, 7}
	encoded, err = EncodeDisplacementSamples(records, 512)
	wantBytes := []byte{255, 0, 11, 22, 33, 254, 44, 55, 66, 255, 255, 0, 255, 7}
	if err != nil || !bytes.Equal(encoded, wantBytes) {
		t.Fatalf("sample edits = %x, %v", encoded, err)
	}
	records, used, err = DecodeDisplacementSamples(raw, 1, 512)
	if err != nil || used != 4 || len(records) != 1 || records[0] != want[0] {
		t.Fatalf("exact-count sample decode = %+v consumed %d err %v", records, used, err)
	}
	for _, triangleCount := range []int{0, 512} {
		records, used, err = DecodeDisplacementSamples(raw, 0, triangleCount)
		if err != nil || used != 0 || len(records) != 0 {
			t.Fatalf("empty samples = %+v consumed %d err %v", records, used, err)
		}
		if encoded, err = EncodeDisplacementSamples(nil, triangleCount); err != nil || len(encoded) != 0 {
			t.Fatalf("empty sample encode = %x, %v", encoded, err)
		}
	}
}

func TestDisplacementSampleErrorsAndBudgets(t *testing.T) {
	for _, test := range []struct {
		name             string
		raw              []byte
		count, triangles int
		message          string
	}{
		{"negative count", nil, -1, 32, "negative count"},
		{"negative triangles", nil, 0, -1, "negative triangle"},
		{"huge count", []byte{0, 0, 0, 0}, math.MaxInt, 32, "capacity"},
		{"short index", nil, 1, 32, "capacity"},
		{"short escape", []byte{255}, 1, 512, "capacity"},
		{"short coordinates", []byte{0, 1, 2}, 1, 32, "capacity"},
		{"escaped coordinates", []byte{255, 0, 1, 2}, 1, 512, "barycentric"},
		{"late escape", []byte{255, 0, 1, 2, 3, 255, 0, 1, 2, 3, 255, 0, 1, 2, 3, 255}, 4, 512, "escape"},
		{"late index", []byte{255, 0, 1, 2, 3, 255, 0, 1, 2, 3, 255, 0, 1, 2, 3, 255, 0, 1, 2, 3}, 5, 512, "missing triangle"},
		{"zero triangles", []byte{0, 1, 2, 3}, 1, 0, "triangle"},
		{"index at count", []byte{32, 1, 2, 3}, 1, 32, "triangle"},
		{"escaped index at count", []byte{255, 0, 1, 2, 3}, 1, 255, "triangle"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if records, _, err := DecodeDisplacementSamples(test.raw, test.count, test.triangles); err == nil || records != nil || !strings.Contains(err.Error(), test.message) {
				t.Fatalf("sample error = %v records %+v", err, records)
			}
			if _, err := MeasureDisplacementSamples(test.raw, test.count, test.triangles); err == nil || !strings.Contains(err.Error(), test.message) {
				t.Fatalf("sample measurement error = %v", err)
			}
		})
	}
	for _, test := range []struct {
		index     uint16
		triangles int
	}{
		{511, 512}, {65535, 65536}, {32, 32}, {0, 0}, {0, -1},
	} {
		if raw, err := EncodeDisplacementSamples([]lumps.DisplacementSample{{Index: test.index}}, test.triangles); err == nil || raw != nil {
			t.Fatalf("invalid sample %d/%d encoded %x err %v", test.index, test.triangles, raw, err)
		}
	}
	for _, test := range []struct{ count, escaped int }{
		{-1, 0}, {math.MaxInt, 0}, {math.MaxInt32 / 4, 4}, {1, -1}, {1, 2},
	} {
		if _, err := displacementSampleBytes(test.count, test.escaped); err == nil {
			t.Fatalf("sample output budget accepted %+v", test)
		}
	}
	if size, err := displacementSampleBytes(3, 2); err != nil || size != 14 {
		t.Fatalf("sample output byte size = %d, %v", size, err)
	}
}

func TestDisplacementInactiveSlotsAndSpanBoundaries(t *testing.T) {
	context := format.Context{Profile: "sdk2013", BSPVersion: 20}
	for _, power := range []uint32{2, 3, 4} {
		raw := displacementWireFixture()
		binary.LittleEndian.PutUint32(raw[20:], power)
		side := uint32(1) << power
		binary.LittleEndian.PutUint32(raw[12:], math.MaxInt32-(side+1)*(side+1))
		binary.LittleEndian.PutUint32(raw[16:], math.MaxInt32-2*side*side)
		binary.LittleEndian.PutUint32(raw[24:], math.MaxInt32)
		// All four corner slots are inactive when Count is zero, even when
		// their indices cannot reference any displacement in this payload.
		for corner := 0; corner < 4; corner++ {
			offset := 96 + corner*10
			for slot := 0; slot < 4; slot++ {
				binary.LittleEndian.PutUint16(raw[offset+slot*2:], uint16(0xff00+slot))
			}
			raw[offset+8] = 0
		}
		records, err := DecodeDisplacements(raw, context)
		if err != nil {
			t.Fatalf("power %d valid signed span boundary: %v", power, err)
		}
		encoded, err := EncodeDisplacements(records, context)
		if err != nil || !bytes.Equal(encoded, raw) {
			t.Fatalf("power %d inactive slots/boundary changed: %v", power, err)
		}
	}
	// Marker-only flags are distinct from unmarked MinTess zero.
	raw := displacementWireFixture()
	binary.LittleEndian.PutUint32(raw[24:], 0x80000000)
	records, err := DecodeDisplacements(raw, context)
	if err != nil || records[0].Flags != 0x80000000 || records[0].MinTess != 0 {
		t.Fatalf("marker-only flags: %+v, %v", records, err)
	}
	records[0].Flags = 0
	encoded, err := EncodeDisplacements(records, context)
	binary.LittleEndian.PutUint32(raw[24:], 0)
	if err != nil || !bytes.Equal(encoded, raw) {
		t.Fatalf("flags-to-MinTess cutover: %v", err)
	}
}

func TestDisplacementSelectedRecordAndSampleMeasurement(t *testing.T) {
	context := format.Context{Profile: "sdk2013", BSPVersion: 20}
	raw := displacementWireFixture()
	record, err := DecodeDisplacementRecord(raw, 1, context)
	if err != nil || record.Power != 4 || record.Flags != 0xc1234567 || record.AllowedVerts[9] != 0x81000023 || record.EdgeNeighbors[3].SubNeighbors[1].Index != 0 || record.CornerNeighbors[3].Neighbors[3] != 0xffff {
		t.Fatalf("selected record = %+v, %v", record, err)
	}
	for _, index := range []int{-1, 2, math.MaxInt} {
		if _, err := DecodeDisplacementRecord(raw, index, context); err == nil {
			t.Fatalf("selected index %d accepted", index)
		}
	}
	if _, err := DecodeDisplacementRecord(raw[:351], 0, context); err == nil {
		t.Fatal("selected decoder accepted a partial final record")
	}
	if _, err := DecodeDisplacementRecord(raw, 0, format.Context{Profile: "sdk2013", BSPVersion: 20, LumpVersion: 1}); err == nil {
		t.Fatal("selected decoder accepted unsupported version")
	}
	binary.LittleEndian.PutUint32(raw[20:], 0xffffffff)
	if _, err := DecodeDisplacementRecord(raw, 1, context); err != nil {
		t.Fatalf("selected decoder interpreted unrelated record: %v", err)
	}
	if _, err := DecodeDisplacementRecord(raw, 0, context); err == nil {
		t.Fatal("selected decoder accepted invalid selected power")
	}
	binary.LittleEndian.PutUint16(raw[176+48+6:], 2)
	if _, err := DecodeDisplacementRecord(raw, 1, context); err == nil {
		t.Fatal("selected decoder accepted neighbor outside full record count")
	}
	stream := []byte{254, 1, 2, 3, 255, 0, 4, 5, 6, 255, 255, 7, 8, 9, 0xab}
	for count, want := range []int{0, 4, 9, 14} {
		if used, err := MeasureDisplacementSamples(stream, count, 512); err != nil || used != want {
			t.Fatalf("measure count %d = %d, %v; want %d", count, used, err, want)
		}
	}
	displacementMeasureRaw := displacementWireFixture()
	if allocations := testing.AllocsPerRun(10, func() {
		_, _ = DecodeDisplacementRecord(displacementMeasureRaw, 1, context)
		_, _ = MeasureDisplacementSamples(stream, 3, 512)
	}); allocations != 0 {
		t.Fatalf("selected decode/sample measurement allocated %g times", allocations)
	}
}
