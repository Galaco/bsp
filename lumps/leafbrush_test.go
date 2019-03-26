package lumps

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestLeafBrush_GetData(t *testing.T) {
	sut := LeafBrush{}
	data := []uint16{0, 1, 2, 5, 43, 156, 146, 3, 3, 6}
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, data)
	if err != nil {
		t.Error(err)
	}
	if err := sut.Unmarshall(buf.Bytes()); err != nil {
		t.Error(err)
	}
	for idx, b := range sut.GetData() {
		if data[idx] != b {
			t.Error("mismatched between expected and actual unmarshalled bytes")
		}
	}
}

func TestLeafBrush_Marshall(t *testing.T) {
	sut := LeafBrush{}
	data := []uint16{0, 1, 2, 5, 43, 156, 146, 3, 3, 6}
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, data)
	if err != nil {
		t.Error(err)
	}
	if err := sut.Unmarshall(buf.Bytes()); err != nil {
		t.Error(err)
	}
	for idx, b := range sut.data {
		if data[idx] != b {
			t.Error("mismatched between expected and actual unmarshalled bytes")
		}
	}
	res, err := sut.Marshall()
	if err != nil {
		t.Errorf("unexpected error during marshalling")
	}
	for idx, b := range res {
		if buf.Bytes()[idx] != b {
			t.Error("mismatched between expected and actual marshalled bytes")
		}
	}
}

func TestLeafBrush_Unmarshall(t *testing.T) {
	sut := LeafBrush{}
	data := []uint16{0, 1, 2, 5, 43, 156, 146, 3, 3, 6}
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, data)
	if err != nil {
		t.Error(err)
	}
	if err := sut.Unmarshall(buf.Bytes()); err != nil {
		t.Error(err)
	}
	for idx, b := range sut.data {
		if data[idx] != b {
			t.Error("mismatched between expected and actual unmarshalled bytes")
		}
	}
}
