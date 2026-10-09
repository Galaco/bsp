package bsp_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"

	"github.com/galaco/bsp"
)

func streamingFixture(t *testing.T) *bsp.File {
	t.Helper()
	wire := make([]byte, 2024)
	copy(wire, "VBSP")
	binary.LittleEndian.PutUint32(wire[4:], 20)
	for _, entry := range []struct {
		id     bsp.LumpID
		offset int
		data   string
	}{
		{bsp.LumpEntities, 1500, "{\x00}"},
		{bsp.LumpPhysCollide, 1600, "OPAQUE"},
		{bsp.LumpDispMultiBlend, 1700, "x"},
	} {
		position := 8 + int(entry.id)*16
		binary.LittleEndian.PutUint32(wire[position:], uint32(entry.offset))
		binary.LittleEndian.PutUint32(wire[position+4:], uint32(len(entry.data)))
		copy(wire[entry.offset:], entry.data)
	}
	binary.LittleEndian.PutUint32(wire[8+int(bsp.LumpGame)*16:], 2000)
	binary.LittleEndian.PutUint32(wire[12+int(bsp.LumpGame)*16:], 24)
	game := wire[2000:]
	binary.LittleEndian.PutUint32(game, 1)
	binary.LittleEndian.PutUint32(game[4:], 0x12345678)
	binary.LittleEndian.PutUint16(game[10:], 7)
	binary.LittleEndian.PutUint32(game[12:], 2020)
	binary.LittleEndian.PutUint32(game[16:], 4)
	copy(game[20:], "GAME")
	file, err := bsp.Read(bytes.NewReader(wire), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	return file
}

type streamSnapshot struct {
	header bsp.Header
	raw    [bsp.LumpCount][]byte
}

func snapshotStreamSource(t *testing.T, file *bsp.File) streamSnapshot {
	t.Helper()
	snapshot := streamSnapshot{header: file.Header()}
	for id := bsp.LumpID(0); id < bsp.LumpCount; id++ {
		var err error
		snapshot.raw[id], err = file.RawLump(id)
		if err != nil {
			t.Fatal(err)
		}
	}
	return snapshot
}

func requireStreamSourceUnchanged(t *testing.T, file *bsp.File, before streamSnapshot) {
	t.Helper()
	if file.Header() != before.header {
		t.Fatal("export changed source metadata")
	}
	for id := bsp.LumpID(0); id < bsp.LumpCount; id++ {
		raw, err := file.RawLump(id)
		if err != nil || !bytes.Equal(raw, before.raw[id]) {
			t.Fatalf("export changed source lump %d: %v", id, err)
		}
	}
}

func TestWriteToMatchesMarshalBinaryAndPreservesSource(t *testing.T) {
	file := streamingFixture(t)
	before := snapshotStreamSource(t, file)
	want, err := file.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	for attempt := 0; attempt < 2; attempt++ {
		var output bytes.Buffer
		n, err := file.WriteTo(&output)
		if err != nil || n != int64(len(want)) || !bytes.Equal(output.Bytes(), want) {
			t.Fatalf("WriteTo count=%d, length=%d, error=%v; want %d identical bytes", n, output.Len(), err, len(want))
		}
		requireStreamSourceUnchanged(t, file, before)
	}
	copyFile, err := bsp.Read(bytes.NewReader(want), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	entities, err := copyFile.Entities()
	if err != nil || entities != "{\x00}" {
		t.Fatalf("exact entity bytes=%q, error=%v", entities, err)
	}
	gameDescriptor, err := copyFile.Descriptor(bsp.LumpGame)
	if err != nil || gameDescriptor.Offset == before.header.Lumps[bsp.LumpGame].Offset {
		t.Fatalf("game was not relocated: descriptor=%+v error=%v", gameDescriptor, err)
	}
	game, err := copyFile.RawLump(bsp.LumpGame)
	if err != nil || len(game) != 24 {
		t.Fatalf("game bytes length=%d error=%v", len(game), err)
	}
	if offset := int32(binary.LittleEndian.Uint32(game[12:])); offset != gameDescriptor.Offset+20 {
		t.Fatalf("relocated entry offset=%d, want=%d", offset, gameDescriptor.Offset+20)
	}
	if !bytes.Equal(game[20:], []byte("GAME")) {
		t.Fatal("relocation changed game payload")
	}
}

type limitedStreamWriter struct {
	output    bytes.Buffer
	remaining int
	failure   error
}

func (writer *limitedStreamWriter) Write(raw []byte) (int, error) {
	if len(raw) <= writer.remaining {
		writer.remaining -= len(raw)
		return writer.output.Write(raw)
	}
	n, _ := writer.output.Write(raw[:writer.remaining])
	writer.remaining = 0
	return n, writer.failure
}

func TestWriteToReportsActualPartialCountsAndShortWrites(t *testing.T) {
	file := streamingFixture(t)
	before := snapshotStreamSource(t, file)
	want, err := file.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	gameBase := int(binary.LittleEndian.Uint32(want[8+int(bsp.LumpGame)*16:]))
	failure := errors.New("destination failed")
	for _, limit := range []int{0, 19, bsp.HeaderSize, bsp.HeaderSize + 2, bsp.HeaderSize + 3, gameBase + 2, gameBase + 13, gameBase + 22, len(want) - 1} {
		for _, returnedError := range []error{failure, nil} {
			writer := limitedStreamWriter{remaining: limit, failure: returnedError}
			n, err := file.WriteTo(&writer)
			wantError := returnedError
			if wantError == nil {
				wantError = io.ErrShortWrite
			}
			if !errors.Is(err, wantError) || n != int64(limit) || writer.output.Len() != limit {
				t.Fatalf("limit=%d returned error=%v: count=%d actual=%d error=%v; want count=%d error=%v", limit, returnedError, n, writer.output.Len(), err, limit, wantError)
			}
			if !bytes.Equal(writer.output.Bytes(), want[:limit]) {
				t.Fatalf("limit=%d: partial output differs from planned prefix", limit)
			}
			requireStreamSourceUnchanged(t, file, before)
		}
	}
}

func TestWriteToValidationFailureProducesNoOutputAndPreservesSource(t *testing.T) {
	file := streamingFixture(t)
	if err := file.SetRawLump(bsp.LumpGame, bsp.HeaderLump{Offset: 2000}, []byte{255, 255, 255, 255}); err != nil {
		t.Fatal(err)
	}
	before := snapshotStreamSource(t, file)
	var output bytes.Buffer
	if n, err := file.WriteTo(&output); n != 0 || output.Len() != 0 || !errors.Is(err, bsp.ErrMalformed) {
		t.Fatalf("invalid game export count=%d output=%d error=%v", n, output.Len(), err)
	}
	if raw, err := file.MarshalBinary(); raw != nil || !errors.Is(err, bsp.ErrMalformed) {
		t.Fatalf("invalid game MarshalBinary bytes=%d error=%v", len(raw), err)
	}
	requireStreamSourceUnchanged(t, file, before)
}
