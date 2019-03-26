package lumps

import (
	"bytes"
	"encoding/binary"
	"github.com/galaco/bsp/primitives/mapflags"
	"testing"
)

func TestMapFlags_GetData(t *testing.T) {
	sut := MapFlags{}
	data := mapflags.MapFlags{
		LevelFlags: 523423,
	}
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, data)
	if err != nil {
		t.Error(err)
	}

	if err := sut.Unmarshall(buf.Bytes()); err != nil {
		t.Error(err)
	}
	if sut.GetData().LevelFlags != data.LevelFlags {
		t.Error("mismatched between expected and actual unmarshalled bytes")
	}
}

func TestMapFlags_Marshall(t *testing.T) {
	sut := MapFlags{}
	data := mapflags.MapFlags{
		LevelFlags: 523423,
	}
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, data)
	if err != nil {
		t.Error(err)
	}

	if err := sut.Unmarshall(buf.Bytes()); err != nil {
		t.Error(err)
	}
	if sut.data.LevelFlags != data.LevelFlags {
		t.Error("mismatched between expected and actual unmarshalled bytes")
	}
	res, err := sut.Marshall()
	if err != nil {
		t.Error(err)
	}
	for idx, b := range res {
		if buf.Bytes()[idx] != b {
			t.Error("mismatched between expected and actual marshalled bytes")
		}
	}

}

func TestMapFlags_Unmarshall(t *testing.T) {
	sut := MapFlags{}
	data := mapflags.MapFlags{
		LevelFlags: 523423,
	}
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, data)
	if err != nil {
		t.Error(err)
	}

	if err := sut.Unmarshall(buf.Bytes()); err != nil {
		t.Error(err)
	}
	if sut.data.LevelFlags != data.LevelFlags {
		t.Error("mismatched between expected and actual unmarshalled bytes")
	}
}
