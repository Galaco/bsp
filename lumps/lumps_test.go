package lumps

import (
	"log"
	"testing"
)

func TestUnimplemented_GetData(t *testing.T) {
	sut := RawBytes{}
	data := []byte{0, 1, 2, 5, 43, 156, 146, 3, 3, 6}
	if err := sut.FromBytes(data); err != nil {
		t.Error(err)
	}
	for idx, b := range sut.Contents() {
		if data[idx] != b {
			t.Error("mismatched between expected and actual unmarshalled bytes")
		}
	}
}

func TestUnimplemented_Marshall(t *testing.T) {
	sut := RawBytes{}
	data := []byte{0, 1, 2, 5, 43, 156, 146, 3, 3, 6}
	if err := sut.FromBytes(data); err != nil {
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
		if data[idx] != b {
			t.Error("mismatched between expected and actual marshalled bytes")
		}
	}
}

func TestUnimplemented_FromBytes(t *testing.T) {
	sut := RawBytes{}
	data := []byte{0, 1, 2, 5, 43, 156, 146, 3, 3, 6}
	if err := sut.FromBytes(data); err != nil {
		t.Error(err)
	}
	for idx, b := range sut.Data {
		if data[idx] != b {
			t.Error("mismatched between expected and actual unmarshalled bytes")
		}
	}
}

func TestMetadata_Length(t *testing.T) {
	sut := Metadata{}
	sut.SetLength(32)
	if sut.Length() != 32 {
		t.Errorf("wrong length returned fron Length")
	}
}

func TestMetadata_SetLength(t *testing.T) {
	sut := Metadata{}
	sut.SetLength(32)
	if sut.Length() != 32 {
		log.Println(sut.length)
		t.Errorf("wrong length returned fron Length")
	}
}

func TestMetadata_Version(t *testing.T) {
	sut := Metadata{}
	sut.SetVersion(11)
	if sut.Version() != 11 {
		t.Errorf("wrong version set")
	}
}

func TestMetadata_SetVersion(t *testing.T) {
	sut := Metadata{}
	sut.SetVersion(11)
	if sut.version != 11 {
		t.Errorf("wrong version set")
	}
}
