package lumps

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestClipPortalVerts_GetData(t *testing.T) {
	sut := ClipPortalVerts{}
	data := []mgl32.Vec3{
		{0, 1, 2},
		{5, 43, 156},
	}
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, data)
	if err != nil {
		t.Error(err)
	}
	err = sut.FromBytes(buf.Bytes())
	if err != nil {
		t.Error(err)
	}
	for idx, b := range sut.Contents() {
		if data[idx] != b {
			t.Error("mismatched between expected and actual unmarshalled bytes")
		}
	}
}

func TestClipPortalVerts_Marshall(t *testing.T) {
	sut := ClipPortalVerts{}
	data := []mgl32.Vec3{
		{0, 1, 2},
		{5, 43, 156},
	}
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, data)
	if err != nil {
		t.Error(err)
	}
	err = sut.FromBytes(buf.Bytes())
	if err != nil {
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

func TestClipPortalVerts_FromBytes(t *testing.T) {
	sut := ClipPortalVerts{}
	data := []mgl32.Vec3{
		{0, 1, 2},
		{5, 43, 156},
	}
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, data)
	if err != nil {
		t.Error(err)
	}
	err = sut.FromBytes(buf.Bytes())
	if err != nil {
		t.Error(err)
	}
	for idx, b := range sut.Data {
		if data[idx] != b {
			t.Error("mismatched between expected and actual unmarshalled bytes")
		}
	}
}
