package lumps

import "testing"

func TestDispLightmapSamplePosition_GetData(t *testing.T) {
	sut := DispLightmapSamplePosition{}
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

func TestDispLightmapSamplePosition_Marshall(t *testing.T) {
	sut := DispLightmapSamplePosition{}
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

func TestDispLightmapSamplePosition_FromBytes(t *testing.T) {
	sut := DispLightmapSamplePosition{}
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
