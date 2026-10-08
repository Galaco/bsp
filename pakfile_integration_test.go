package bsp_test

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"errors"
	"io/fs"
	"testing"

	"github.com/galaco/bsp"
)

type pakTestEntry struct {
	name, text string
	method     uint16
}

func pakTestZIP(t *testing.T, entries []pakTestEntry) []byte {
	t.Helper()
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	for _, entry := range entries {
		file, err := writer.CreateHeader(&zip.FileHeader{Name: entry.name, Method: entry.method})
		if err != nil {
			t.Fatal(err)
		}
		if _, err = file.Write([]byte(entry.text)); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}
func requirePakContents(t *testing.T, pak *bsp.Pakfile, name, want string) {
	t.Helper()
	data, err := pak.ReadFile(name)
	if err != nil || string(data) != want {
		t.Fatalf("resource %q=%q error=%v want=%q", name, data, err, want)
	}
}

func TestPakfileReadOnlyNamesAndMissingVersusEmpty(t *testing.T) {
	raw := pakTestZIP(t, []pakTestEntry{{"Dir\\Mixed////file.TXT", "first", zip.Deflate}, {"dir/mixed/file.txt", "collision", zip.Store}, {"Empty", "", zip.Store}, {"dir/../literal", "literal", zip.Store}, {"literal", "resolved", zip.Store}})
	for _, profile := range []bsp.Profile{bsp.ProfileSDK2013, bsp.ProfileCSGO, bsp.ProfileL4D2, bsp.ProfileReactiveDrop} {
		file := spatialFixture(t, profile, map[bsp.LumpId][]byte{bsp.LumpPakfile: raw})
		pak, err := file.Pakfile()
		if err != nil {
			t.Fatal(err)
		}
		requirePakContents(t, pak, "DIR////MIXED\\\\FILE.txt", "first")
		requirePakContents(t, pak, "empty", "")
		requirePakContents(t, pak, "dir/../literal", "literal")
		requirePakContents(t, pak, "literal", "resolved")
		if _, err := pak.ReadFile("missing"); !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("missing error=%v", err)
		}
		metadata := pak.Entries()
		if len(metadata) != 5 || metadata[0].Name != "Dir\\Mixed////file.TXT" || metadata[0].UncompressedSize != 5 || metadata[2].UncompressedSize != 0 {
			t.Fatalf("metadata=%+v", metadata)
		}
		metadata[0].Name = "changed"
		metadata[0].UncompressedSize = 999
		requirePakContents(t, pak, "dir/mixed/file.txt", "first")
		next := rereadOrdinary(t, file)
		stored, err := next.RawLump(bsp.LumpPakfile)
		if err != nil || !bytes.Equal(stored, raw) {
			t.Fatalf("resource reads changed ZIP=%v", err)
		}
		nextPak, err := next.Pakfile()
		if err != nil {
			t.Fatal(err)
		}
		requirePakContents(t, nextPak, "DIR/MIXED/FILE.TXT", "first")
	}
}

func TestPakfileAbsentCorruptAndReplacedSnapshots(t *testing.T) {
	file := spatialFixture(t, bsp.ProfileCSGO, nil)
	absent, err := file.Pakfile()
	if err != nil {
		t.Fatal(err)
	}
	if len(absent.Entries()) != 0 {
		t.Fatal("absent archive fabricated entries")
	}
	if _, err := absent.ReadFile("x"); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("absent lookup=%v", err)
	}
	emptyZIP := pakTestZIP(t, nil)
	if err := file.SetRawLump(bsp.LumpPakfile, bsp.HeaderLump{}, emptyZIP); err != nil {
		t.Fatal(err)
	}
	empty, err := file.Pakfile()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := empty.ReadFile("x"); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("empty ZIP lookup=%v", err)
	}
	original := pakTestZIP(t, []pakTestEntry{{"resource", "old", zip.Deflate}})
	if err := file.SetRawLump(bsp.LumpPakfile, bsp.HeaderLump{}, original); err != nil {
		t.Fatal(err)
	}
	snapshot, err := file.Pakfile()
	if err != nil {
		t.Fatal(err)
	}
	replacement := pakTestZIP(t, []pakTestEntry{{"resource", "new", zip.Store}})
	want := append([]byte(nil), replacement...)
	if err := file.SetRawLump(bsp.LumpPakfile, bsp.HeaderLump{}, replacement); err != nil {
		t.Fatal(err)
	}
	replacement[0] ^= 0xff
	current, err := file.Pakfile()
	if err != nil {
		t.Fatal(err)
	}
	requirePakContents(t, current, "resource", "new")
	requirePakContents(t, snapshot, "resource", "old")
	next := rereadOrdinary(t, file)
	raw, err := next.RawLump(bsp.LumpPakfile)
	if err != nil || !bytes.Equal(raw, want) {
		t.Fatalf("replacement export=%v", err)
	}
	malformed := []byte("not ZIP")
	if err := file.SetRawLump(bsp.LumpPakfile, bsp.HeaderLump{}, malformed); err != nil {
		t.Fatal(err)
	}
	if bad, err := file.Pakfile(); bad != nil || !errors.Is(err, bsp.ErrMalformed) || !errors.Is(err, zip.ErrFormat) {
		t.Fatalf("corrupt archive=%+v error=%v", bad, err)
	}
	requirePakContents(t, snapshot, "resource", "old")
	next = rereadOrdinary(t, file)
	raw, err = next.RawLump(bsp.LumpPakfile)
	if err != nil || !bytes.Equal(raw, malformed) {
		t.Fatalf("opaque invalid ZIP altered=%v", err)
	}
}

func TestPakfileCorruptEntriesPreserveUnderlyingErrors(t *testing.T) {
	for _, method := range []uint16{zip.Store, zip.Deflate} {
		raw := pakTestZIP(t, []pakTestEntry{{"broken", "contents with a valid original checksum", method}})
		reader, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
		if err != nil {
			t.Fatal(err)
		}
		offset, err := reader.File[0].DataOffset()
		if err != nil {
			t.Fatal(err)
		}
		if method == zip.Store {
			raw[offset] ^= 0xff
		} else {
			raw[offset] = 0xff
		}
		file := spatialFixture(t, bsp.ProfileCSGO, map[bsp.LumpId][]byte{bsp.LumpPakfile: raw})
		pak, err := file.Pakfile()
		if err != nil {
			t.Fatal(err)
		}
		_, err = pak.ReadFile("BROKEN")
		if method == zip.Store {
			if !errors.Is(err, zip.ErrChecksum) {
				t.Fatalf("stored corruption error=%v", err)
			}
		} else {
			var corrupt flate.CorruptInputError
			if !errors.As(err, &corrupt) {
				t.Fatalf("deflate corruption error=%v", err)
			}
		}
		next := rereadOrdinary(t, file)
		stored, err := next.RawLump(bsp.LumpPakfile)
		if err != nil || !bytes.Equal(stored, raw) {
			t.Fatalf("corrupt entry read altered ZIP=%v", err)
		}
	}
}
