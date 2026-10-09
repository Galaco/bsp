package codec

import (
	"testing"

	"github.com/galaco/bsp/internal/format"
)

func TestPlanLayoutHeaderOnly(t *testing.T) {
	got, err := PlanLayout([format.LumpCount]int64{})
	if err != nil {
		t.Fatal(err)
	}
	want := Layout{Size: 1036}
	if got != want {
		t.Fatalf("layout = %+v, want %+v", got, want)
	}
}

func TestPlanLayoutNumericOrderAndPadding(t *testing.T) {
	var lengths [format.LumpCount]int64
	lengths[0], lengths[1], lengths[2], lengths[3], lengths[63] = 1, 2, 3, 4, 5
	want := Layout{Size: 1060}
	want.Offsets[0], want.Lengths[0] = 1036, 1
	want.Offsets[1], want.Lengths[1] = 1040, 2
	want.Offsets[2], want.Lengths[2] = 1044, 3
	want.Offsets[3], want.Lengths[3] = 1048, 4
	want.Offsets[63], want.Lengths[63] = 1052, 5
	for attempt := 0; attempt < 2; attempt++ {
		got, err := PlanLayout(lengths)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("attempt %d: layout = %+v, want %+v", attempt, got, want)
		}
	}
}

func TestPlanLayoutAllIDs(t *testing.T) {
	var lengths [format.LumpCount]int64
	want := Layout{Size: 1292}
	for id := range lengths {
		lengths[id] = 4
		want.Offsets[id] = int32(1036 + id*4)
		want.Lengths[id] = 4
	}
	got, err := PlanLayout(lengths)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("layout = %+v, want %+v", got, want)
	}
}

func TestPlanLayoutEmptySlotsAndFinalPadding(t *testing.T) {
	cases := []struct {
		name   string
		length int64
		size   int32
	}{
		{name: "one", length: 1, size: 1040},
		{name: "two", length: 2, size: 1040},
		{name: "three", length: 3, size: 1040},
		{name: "four", length: 4, size: 1040},
		{name: "five", length: 5, size: 1044},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var lengths [format.LumpCount]int64
			lengths[31] = tc.length
			want := Layout{Size: tc.size}
			want.Offsets[31], want.Lengths[31] = 1036, int32(tc.length)
			got, err := PlanLayout(lengths)
			if err != nil {
				t.Fatal(err)
			}
			if got != want {
				t.Fatalf("layout = %+v, want %+v", got, want)
			}
		})
	}
}

func TestPlanLayoutHighestAlignedSize(t *testing.T) {
	cases := []struct {
		name    string
		lengths [format.LumpCount]int64
		want    Layout
	}{
		{
			name:    "exact aligned end",
			lengths: [format.LumpCount]int64{63: 2147482608},
			want: Layout{
				Offsets: [format.LumpCount]int32{63: 1036},
				Lengths: [format.LumpCount]int32{63: 2147482608},
				Size:    2147483644,
			},
		},
		{
			name:    "padding to highest end",
			lengths: [format.LumpCount]int64{63: 2147482605},
			want: Layout{
				Offsets: [format.LumpCount]int32{63: 1036},
				Lengths: [format.LumpCount]int32{63: 2147482605},
				Size:    2147483644,
			},
		},
		{
			name:    "last ID at highest offset",
			lengths: [format.LumpCount]int64{0: 2147482604, 63: 4},
			want: Layout{
				Offsets: [format.LumpCount]int32{0: 1036, 63: 2147483640},
				Lengths: [format.LumpCount]int32{0: 2147482604, 63: 4},
				Size:    2147483644,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := PlanLayout(tc.lengths)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("layout = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestPlanLayoutInvalidSizesReturnNoPartialPlan(t *testing.T) {
	cases := []struct {
		name    string
		lengths [format.LumpCount]int64
	}{
		{name: "negative", lengths: [format.LumpCount]int64{0: -1}},
		{name: "minimum int64", lengths: [format.LumpCount]int64{0: -1 << 63}},
		{name: "length above int32", lengths: [format.LumpCount]int64{0: 1 << 31}},
		{name: "maximum int64", lengths: [format.LumpCount]int64{0: 1<<63 - 1}},
		{name: "header plus maximum length", lengths: [format.LumpCount]int64{0: 1<<31 - 1}},
		{name: "final alignment overflow", lengths: [format.LumpCount]int64{0: 2147482609}},
		{name: "late negative", lengths: [format.LumpCount]int64{0: 4, 63: -1}},
		{name: "late oversized length", lengths: [format.LumpCount]int64{0: 4, 63: 1<<63 - 1}},
		{name: "late sum overflow", lengths: [format.LumpCount]int64{0: 2147482604, 63: 8}},
		{name: "late alignment overflow", lengths: [format.LumpCount]int64{0: 2147482604, 63: 5}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := PlanLayout(tc.lengths)
			if err == nil {
				t.Fatal("accepted an unrepresentable layout")
			}
			if got != (Layout{}) {
				t.Fatalf("returned partial layout on error: %+v", got)
			}
		})
	}
}
