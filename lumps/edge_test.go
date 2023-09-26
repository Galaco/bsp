package lumps

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestEdge_GetData(t *testing.T) {
	sut := Edge{}
	data := [][2]uint16{{0, 1}, {2, 5}, {43, 156}, {146, 3}, {3, 6}}
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, data)
	if err != nil {
		t.Error(err)
	}
	if err := sut.FromBytes(buf.Bytes()); err != nil {
		t.Error(err)
	}
	for idx, b := range sut.Contents() {
		if data[idx][0] != b[0] || data[idx][1] != b[1] {
			t.Error("mismatched between expected and actual unmarshalled bytes")
		}
	}
}

func TestEdge_Marshall(t *testing.T) {
	sut := Edge{}
	data := [][2]uint16{{0, 1}, {2, 5}, {43, 156}, {146, 3}, {3, 6}}
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, data)
	if err != nil {
		t.Error(err)
	}
	if err := sut.FromBytes(buf.Bytes()); err != nil {
		t.Error(err)
	}
	for idx, b := range sut.Data {
		if data[idx] != b {
			t.Error("mismatched between expected and actual unmarshalled bytes")
		}
	}
	res, err := sut.ToBytes()
	if err != nil {
		t.Errorf("unexpected error during marshalling")
	}
	for idx, b := range res {
		if buf.Bytes()[idx] != b {
			t.Error("mismatched between expected and actual marshalled bytes")
		}
	}
}

func TestEdge_FromBytes(t *testing.T) {
	sut := Edge{}
	data := [][2]uint16{{0, 1}, {2, 5}, {43, 156}, {146, 3}, {3, 6}}
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, data)
	if err != nil {
		t.Error(err)
	}
	if err := sut.FromBytes(buf.Bytes()); err != nil {
		t.Error(err)
	}
	for idx, b := range sut.Data {
		if data[idx] != b {
			t.Error("mismatched between expected and actual unmarshalled bytes")
		}
	}
}
