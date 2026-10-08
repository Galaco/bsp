package codec

import (
	"archive/zip"
	"bytes"
	"errors"
	"io/fs"
	"testing"
)

func TestPakfileFailedImportNeverExposesPreviousArchive(t *testing.T) {
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	entry, err := writer.Create("old")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = entry.Write([]byte("previous bytes")); err != nil {
		t.Fatal(err)
	}
	if err = writer.Close(); err != nil {
		t.Fatal(err)
	}
	var pak Pakfile
	pak.SetContext(visibilityTestContext())
	if err = pak.Unmarshall(buffer.Bytes()); err != nil {
		t.Fatal(err)
	}
	malformed := []byte("new corrupt archive")
	if err = pak.Unmarshall(malformed); !errors.Is(err, zip.ErrFormat) {
		t.Fatalf("decode error=%v", err)
	}
	if len(pak.Entries()) != 0 {
		t.Fatal("failed import exposed old directory")
	}
	if _, err = pak.ReadFile("old"); !errors.Is(err, zip.ErrFormat) {
		t.Fatalf("lookup hid corrupt replacement=%v", err)
	}
	raw, err := pak.Marshall()
	if err != nil || !bytes.Equal(raw, malformed) {
		t.Fatalf("failed import serialized old bytes=%v", err)
	}
	if err = pak.Unmarshall(nil); err != nil {
		t.Fatal(err)
	}
	if _, err = pak.ReadFile("old"); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("absent replacement retained corrupt/previous state=%v", err)
	}
}
