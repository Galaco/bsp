package lumps

import (
	"bytes"
	"reflect"
	"testing"
)

func TestVisibilityValueRowsSelectorsAndBounds(t *testing.T) {
	value := Visibility{Offsets: [][2]int32{{12, -1}}, Data: []byte{0xff}}
	row, present, err := value.Row(0, PVS)
	if err != nil || !present || !bytes.Equal(row, []byte{0xff}) {
		t.Fatalf("row=%v present=%v error=%v", row, present, err)
	}
	row[0] = 0
	if value.Data[0] != 0xff {
		t.Fatal("row aliases compressed data")
	}
	visible, present, err := value.VisibleClusters(0, PVS)
	if err != nil || !present || !reflect.DeepEqual(visible, []int{0}) {
		t.Fatalf("unused bits interpreted: %v %v %v", visible, present, err)
	}
	row, present, err = value.Row(0, PAS)
	if err != nil || present || row != nil {
		t.Fatalf("missing PAS fabricated: %v %v %v", row, present, err)
	}
	for _, tc := range []struct {
		cluster int
		kind    VisibilityKind
	}{{-1, PVS}, {1, PVS}, {0, VisibilityKind(2)}, {0, VisibilityKind(255)}} {
		if _, _, err := value.Row(tc.cluster, tc.kind); err == nil {
			t.Fatalf("accepted cluster %d kind %d", tc.cluster, tc.kind)
		}
		if _, _, err := value.VisibleClusters(tc.cluster, tc.kind); err == nil {
			t.Fatal("VisibleClusters swallowed bounds error")
		}
	}
	if _, _, err := (Visibility{}).Row(0, PVS); err == nil {
		t.Fatal("zero-count value fabricated a row")
	}
	tooMany := Visibility{Offsets: make([][2]int32, 65537)}
	if _, _, err := tooMany.Row(0, PVS); err == nil {
		t.Fatal("public row accepted over-limit count")
	}
}

func TestVisibilityTableEditsRemainLumpRelative(t *testing.T) {
	value := Visibility{Offsets: [][2]int32{{12, -1}}, Data: []byte{1}}
	value.Offsets = append(value.Offsets, [2]int32{-1, -1})
	if _, _, err := value.Row(0, PVS); err == nil {
		t.Fatal("table growth silently rebased old offset")
	}
	value.Offsets[0][0] = 20
	row, present, err := value.Row(0, PVS)
	if err != nil || !present || !bytes.Equal(row, []byte{1}) {
		t.Fatalf("explicit rebasing failed: %v %v %v", row, present, err)
	}
}
