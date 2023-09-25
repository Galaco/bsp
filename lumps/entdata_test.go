package lumps

import "testing"

func TestEntData_GetData(t *testing.T) {
	sut := EntData{}
	data := []byte{0, 1, 2, 5, 43, 156, 146, 3, 3, 6}
	if err := sut.FromBytes(data); err != nil {
		t.Error(err)
	}
	if string(sut.Contents()) != string(data) {
		t.Error("mismatched between expected and actual unmarshalled bytes")
	}
}

func TestEntData_Marshall(t *testing.T) {
	sut := EntData{}
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

func TestEntData_FromBytes(t *testing.T) {
	sut := EntData{}
	data := []byte{0, 1, 2, 5, 43, 156, 146, 3, 3, 6}
	if err := sut.FromBytes(data); err != nil {
		t.Error(err)
	}
	if string(sut.data) != string(data) {
		t.Error("mismatched between expected and actual unmarshalled bytes")
	}
}
