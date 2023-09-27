package lump

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/galaco/bsp/lump/primitive/mapflags"
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

	if err := sut.FromBytes(buf.Bytes()); err != nil {
		t.Error(err)
	}
	if sut.Contents().LevelFlags != data.LevelFlags {
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

	if err := sut.FromBytes(buf.Bytes()); err != nil {
		t.Error(err)
	}
	if sut.Data.LevelFlags != data.LevelFlags {
		t.Error("mismatched between expected and actual unmarshalled bytes")
	}
	res, err := sut.ToBytes()
	if err != nil {
		t.Error(err)
	}
	for idx, b := range res {
		if buf.Bytes()[idx] != b {
			t.Error("mismatched between expected and actual marshalled bytes")
		}
	}

}

func TestMapFlags_FromBytes(t *testing.T) {
	sut := MapFlags{}
	data := mapflags.MapFlags{
		LevelFlags: 523423,
	}
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, data)
	if err != nil {
		t.Error(err)
	}

	if err := sut.FromBytes(buf.Bytes()); err != nil {
		t.Error(err)
	}
	if sut.Data.LevelFlags != data.LevelFlags {
		t.Error("mismatched between expected and actual unmarshalled bytes")
	}
}
