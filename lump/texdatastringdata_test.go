package lump

import "testing"

func TestTexDataStringData_GetData(t *testing.T) {
	sut := TexDataStringData{}
	data := []byte{0, 1, 2, 5, 43, 156, 146, 3, 3, 6}
	if err := sut.FromBytes(data); err != nil {
		t.Error(err)
	}
	if sut.Contents() != string(data) {
		t.Error("mismatched between expected and actual unmarshalled bytes")
	}
}

func TestTexDataStringData_Marshall(t *testing.T) {
	sut := TexDataStringData{}
	data := []byte{0, 1, 2, 5, 43, 156, 146, 3, 3, 6}
	if err := sut.FromBytes(data); err != nil {
		t.Error(err)
	}
	if sut.Contents() != string(data) {
		t.Error("mismatched between expected and actual unmarshalled bytes")
	}
	res, err := sut.ToBytes()
	if err != nil {
		t.Errorf("unexpected error during marshalling")
	}
	if string(res) != string(data) {
		t.Error("mismatched between expected and actual marshalled bytes")
	}
}

func TestTexDataStringData_FromBytes(t *testing.T) {
	sut := TexDataStringData{}
	data := []byte{0, 1, 2, 5, 43, 156, 146, 3, 3, 6}
	if err := sut.FromBytes(data); err != nil {
		t.Error(err)
	}
	if sut.Data != string(data) {
		t.Error("mismatched between expected and actual unmarshalled bytes")
	}
}
