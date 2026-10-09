package format

import "testing"

func TestFixedRecordCount(t *testing.T) {
	cases := []struct {
		name      string
		length    int
		width     int
		want      int
		wantError bool
	}{
		{"empty", 0, 20, 0, false},
		{"one", 20, 20, 1, false},
		{"two", 40, 20, 2, false},
		{"short", 19, 20, 0, true},
		{"plane21", 21, 20, 0, true},
		{"negative length", -1, 20, 0, true},
		{"zero width", 0, 0, 0, true},
		{"negative width", 0, -1, 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := FixedRecordCount(tc.length, tc.width)
			if (err != nil) != tc.wantError || got != tc.want {
				t.Fatalf("count = %d, error = %v; want %d, error %t", got, err, tc.want, tc.wantError)
			}
		})
	}
}

func TestFixedRecordBytes(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	cases := []struct {
		name      string
		count     int
		width     int
		want      int
		wantError bool
	}{
		{"empty", 0, 176, 0, false},
		{"one", 1, 176, 176, false},
		{"two", 2, 176, 352, false},
		{"signed32 maximum", 1<<31 - 1, 1, 1<<31 - 1, false},
		{"signed32 overflow", 1 << 30, 2, 0, true},
		{"record rounding boundary", (1<<31 - 1) / 176, 176, ((1<<31 - 1) / 176) * 176, false},
		{"record rounding overflow", (1<<31-1)/176 + 1, 176, 0, true},
		{"address multiplication overflow", maxInt, maxInt, 0, true},
		{"negative count", -1, 20, 0, true},
		{"zero width", 0, 0, 0, true},
		{"negative width", 0, -1, 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := FixedRecordBytes(tc.count, tc.width)
			if (err != nil) != tc.wantError || got != tc.want {
				t.Fatalf("bytes = %d, error = %v; want %d, error %t", got, err, tc.want, tc.wantError)
			}
		})
	}
}
