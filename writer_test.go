package bsp_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"github.com/galaco/bsp"
	"github.com/galaco/bsp/primitives/plane"
)

// Fixtures and inspections use the PC wire fields directly, not the writer's
// layout planner or header encoder. Game ID 35 contains a real directory with
// an opaque child and two gap bytes which relocation must also preserve.
func writerAllLumps(t *testing.T, profile bsp.Profile) (*bsp.Bsp, [bsp.LumpCount][]byte) {
	t.Helper()
	input := containerBytes(-1, nil)
	if profile != bsp.ProfileSDK2013 {
		binary.LittleEndian.PutUint32(input[4:], 21)
	}
	binary.LittleEndian.PutUint32(input[1032:], 4321)
	file, err := bsp.ReadFromStream(bytes.NewReader(input), bsp.ReadOptions{Profile: profile})
	if err != nil {
		t.Fatal(err)
	}
	var payloads [bsp.LumpCount][]byte
	for id := 0; id < bsp.LumpCount; id++ {
		raw := make([]byte, 1+(id+1)%8)
		for j := range raw {
			raw[j] = byte(1 + id*3 + j)
		}
		desc := bsp.HeaderLump{
			Offset: int32(4096 + id*64), Length: -123, Version: int32(1000 + id),
			FourCC: [4]byte{byte(id + 1), 0xa5, 0x5a, byte(255 - id)},
		}
		if id == 35 {
			desc.Offset = 8192
			desc.FourCC = [4]byte{}
			raw = make([]byte, 27)
			binary.LittleEndian.PutUint32(raw, 1)
			putGameEntry(raw, 0, 0xdecafbad, 0x4000, 0x1234, 8192+22, 5)
			copy(raw[20:], []byte{0xec, 0xed, 0x35, 0xaa, 0xbb, 0xcc, 0xdd})
		}
		payloads[id] = append([]byte(nil), raw...)
		if err := file.SetRawLump(bsp.LumpId(id), desc, raw); err != nil {
			t.Fatalf("set lump %d: %v", id, err)
		}
		// The recognizable bytes above must remain owned by the file.
		raw[0] ^= 0xff
	}
	return file, payloads
}

func writerWireDescriptor(raw []byte, id int, profile bsp.Profile) bsp.HeaderLump {
	entry := raw[8+16*id : 8+16*(id+1)]
	var desc bsp.HeaderLump
	if profile == bsp.ProfileL4D2 {
		desc.Version = int32(binary.LittleEndian.Uint32(entry))
		desc.Offset = int32(binary.LittleEndian.Uint32(entry[4:]))
		desc.Length = int32(binary.LittleEndian.Uint32(entry[8:]))
	} else {
		desc.Offset = int32(binary.LittleEndian.Uint32(entry))
		desc.Length = int32(binary.LittleEndian.Uint32(entry[4:]))
		desc.Version = int32(binary.LittleEndian.Uint32(entry[8:]))
	}
	copy(desc.FourCC[:], entry[12:])
	return desc
}

func requireWriterSource(t *testing.T, file *bsp.Bsp, header bsp.Header, payloads [bsp.LumpCount][]byte) {
	t.Helper()
	if file.Header() != header {
		t.Fatal("export changed the source header")
	}
	for id := 0; id < bsp.LumpCount; id++ {
		desc, err := file.Descriptor(bsp.LumpId(id))
		if err != nil || desc != header.Lumps[id] {
			t.Fatalf("source descriptor %d=%+v, want %+v; error=%v", id, desc, header.Lumps[id], err)
		}
		stored, err := file.RawLump(bsp.LumpId(id))
		if err != nil || !bytes.Equal(stored, payloads[id]) {
			t.Fatalf("source lump %d=%x, want %x; error=%v", id, stored, payloads[id], err)
		}
	}
}

func TestWriterAll64LumpsCanonicalAndOpaque(t *testing.T) {
	for _, profile := range []bsp.Profile{bsp.ProfileSDK2013, bsp.ProfileCSGO, bsp.ProfileL4D2} {
		t.Run(string(profile), func(t *testing.T) {
			file, payloads := writerAllLumps(t, profile)
			header := file.Header()
			writer := bsp.NewWriter(file)
			output, err := writer.Write()
			if err != nil {
				t.Fatal(err)
			}
			// Eight cycles of residues, replacing ID 35 with the 27-byte game
			// container, yield 404 payload/padding bytes after the PC header.
			if len(output) != 1440 || binary.LittleEndian.Uint32(output) != 0x50534256 || int32(binary.LittleEndian.Uint32(output[4:])) != header.Version || binary.LittleEndian.Uint32(output[1032:]) != 4321 {
				t.Fatalf("output length=%d, want 1440, or file metadata changed", len(output))
			}
			cursor := 1036
			var emitted [bsp.LumpCount][]byte
			var wireHeader bsp.Header
			wireHeader.Id, wireHeader.Version, wireHeader.Revision = header.Id, header.Version, 4321
			for id := 0; id < bsp.LumpCount; id++ {
				desc := writerWireDescriptor(output, id, profile)
				wireHeader.Lumps[id] = desc
				wantFourCC := [4]byte{byte(id + 1), 0xa5, 0x5a, byte(255 - id)}
				if id == 35 {
					wantFourCC = [4]byte{}
				}
				if desc.Offset != int32(cursor) || desc.Offset%4 != 0 || desc.Length != int32(len(payloads[id])) || desc.Version != int32(1000+id) || desc.FourCC != wantFourCC {
					t.Fatalf("wire descriptor %d=%+v, want offset=%d length=%d version=%d FourCC=%x", id, desc, cursor, len(payloads[id]), 1000+id, wantFourCC)
				}
				want := append([]byte(nil), payloads[id]...)
				if id == 35 {
					// Only the absolute directory offset may move. Count, flags,
					// child version/length, gaps and opaque bytes stay exact.
					if cursor != 1240 {
						t.Fatalf("game base=%d, want 1240", cursor)
					}
					binary.LittleEndian.PutUint32(want[12:], 1262)
				}
				end := cursor + len(want)
				if end > len(output) || !bytes.Equal(output[cursor:end], want) {
					t.Fatalf("lump %d payload at offset %d does not match its declared bytes", id, cursor)
				}
				emitted[id] = want
				cursor = end
				for cursor%4 != 0 {
					if cursor >= len(output) || output[cursor] != 0 {
						t.Fatalf("lump %d padding byte at %d is missing or nonzero", id, cursor)
					}
					cursor++
				}
			}
			if cursor != len(output) || !bytes.Equal(output[len(output)-3:], []byte{0, 0, 0}) {
				t.Fatalf("EOF=%d, want aligned final cursor=%d with three trailing padding bytes", len(output), cursor)
			}
			reloaded, err := bsp.ReadFromStream(bytes.NewReader(output), bsp.ReadOptions{Profile: profile})
			if err != nil {
				t.Fatal(err)
			}
			if reloaded.Profile() != profile || reloaded.Header() != wireHeader {
				t.Fatalf("reimport did not normalize the %s wire header", profile)
			}
			for id := 0; id < bsp.LumpCount; id++ {
				got, err := reloaded.RawLump(bsp.LumpId(id))
				if err != nil || !bytes.Equal(got, emitted[id]) {
					t.Fatalf("reimport lump %d=%x, want %x; error=%v", id, got, emitted[id], err)
				}
				desc, err := reloaded.Descriptor(bsp.LumpId(id))
				if err != nil || desc != wireHeader.Lumps[id] {
					t.Fatalf("reimport descriptor %d=%+v, want %+v; error=%v", id, desc, wireHeader.Lumps[id], err)
				}
			}
			// Explicitly guard the six IDs the former ordering table omitted.
			for _, id := range []int{32, 49, 57, 61, 62, 63} {
				desc := wireHeader.Lumps[id]
				if desc.Length == 0 || !bytes.Equal(output[int(desc.Offset):int(desc.Offset+desc.Length)], payloads[id]) {
					t.Fatalf("formerly omitted lump %d did not survive", id)
				}
			}
			entries := requireGameEntries(t, reloaded, 1)
			if entries[0].ID != 0xdecafbad || entries[0].Flags != 0x4000 || entries[0].Version != 0x1234 || entries[0].Offset != 1262 || entries[0].LogicalLength != 5 || !bytes.Equal(entries[0].StoredData, []byte{0x35, 0xaa, 0xbb, 0xcc, 0xdd}) {
				t.Fatalf("reimport game entry=%+v", entries[0])
			}

			// Getter access, even unsupported typed access and mutations to
			// returned snapshots, must not affect subsequent exports.
			if _, err := file.Lump(bsp.LumpPlanes); !errors.Is(err, bsp.ErrUnsupportedLump) {
				t.Fatalf("unknown plane layout error=%v", err)
			}
			for id := 0; id < bsp.LumpCount; id++ {
				stored, err := writer.WriteLump(bsp.LumpId(id))
				if err != nil || !bytes.Equal(stored, payloads[id]) {
					t.Fatalf("stored export %d=%x, want %x; error=%v", id, stored, payloads[id], err)
				}
				stored[0] ^= 0xff
				raw, err := file.RawLump(bsp.LumpId(id))
				if err != nil {
					t.Fatal(err)
				}
				raw[0] ^= 0xff
			}
			snapshot := file.Header()
			snapshot.Lumps[63].Length = 0
			desc, err := file.Descriptor(63)
			if err != nil {
				t.Fatal(err)
			}
			desc.Offset = 0
			sourceEntries := requireGameEntries(t, file, 1)
			if sourceEntries[0].Offset != 8214 {
				t.Fatalf("export mutated original game position: %+v", sourceEntries[0])
			}
			sourceEntries[0].StoredData[0] ^= 0xff
			again, err := writer.Write()
			if err != nil || !bytes.Equal(again, output) {
				t.Fatalf("repeated export changed after getter access; error=%v", err)
			}
			// The full-file result is owned too, including relocated game data.
			output[0] = 0
			output[1262] ^= 0xff
			third, err := writer.Write()
			if err != nil || !bytes.Equal(third, again) {
				t.Fatalf("mutating an exported buffer affected another export; error=%v", err)
			}
			requireWriterSource(t, file, header, payloads)
		})
	}
}

func TestWriterCommittedEntityAndPlaneEdits(t *testing.T) {
	file, payloads := writerAllLumps(t, bsp.ProfileSDK2013)
	for _, id := range []bsp.LumpId{bsp.LumpEntities, bsp.LumpPlanes} {
		desc, err := file.Descriptor(id)
		if err != nil {
			t.Fatal(err)
		}
		desc.Version, desc.FourCC = 0, [4]byte{}
		if err := file.SetRawLump(id, desc, nil); err != nil {
			t.Fatal(err)
		}
	}
	if err := file.SetEntities("old entity"); err != nil {
		t.Fatal(err)
	}
	if err := file.SetPlanes([]plane.Plane{{Normal: [3]float32{1, 0, 0}, Distance: 7, AxisType: 0}}); err != nil {
		t.Fatal(err)
	}
	// An ordinary future version with no compression metadata cannot be
	// decoded, but must not prevent the full-file write below.
	if err := file.SetRawLump(19, bsp.HeaderLump{Version: 987}, []byte{0xf1, 0x82, 0x73}); err != nil {
		t.Fatal(err)
	}
	payloads[19] = []byte{0xf1, 0x82, 0x73}
	if _, err := file.Lump(19); !errors.Is(err, bsp.ErrUnsupportedLump) {
		t.Fatalf("future ordinary layout error=%v", err)
	}
	writer := bsp.NewWriter(file)
	before, err := writer.Write()
	if err != nil {
		t.Fatal(err)
	}
	initial, err := bsp.ReadFromStream(bytes.NewReader(before), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	if text, err := initial.Entities(); err != nil || text != "old entity" {
		t.Fatalf("initial entities=%q error=%v", text, err)
	}
	initialPlanes, err := initial.Planes()
	if err != nil || len(initialPlanes) != 1 || initialPlanes[0].Distance != 7 {
		t.Fatalf("initial planes=%+v error=%v", initialPlanes, err)
	}
	planes, err := file.Planes()
	if err != nil {
		t.Fatal(err)
	}
	planes[0].Distance = 23.5
	planes = append(planes, plane.Plane{Normal: [3]float32{0, -1, 0}, Distance: -4, AxisType: 1})
	if err := file.SetPlanes(planes); err != nil {
		t.Fatal(err)
	}
	planes[0].Distance = 999
	const editedEntities = "{\n\"classname\" \"worldspawn\"\n\"message\" \"committed edit\"\n}\n"
	if err := file.SetEntities(editedEntities); err != nil {
		t.Fatal(err)
	}
	header := file.Header()
	for _, id := range []bsp.LumpId{bsp.LumpEntities, bsp.LumpPlanes} {
		payloads[id], err = file.RawLump(id)
		if err != nil {
			t.Fatal(err)
		}
	}
	output, err := writer.Write()
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(output, before) {
		t.Fatal("committed edits did not change the full-file export")
	}
	reloaded, err := bsp.ReadFromStream(bytes.NewReader(output), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	if text, err := reloaded.Entities(); err != nil || text != editedEntities {
		t.Fatalf("edited entities=%q error=%v", text, err)
	}
	got, err := reloaded.Planes()
	if err != nil || len(got) != 2 {
		t.Fatalf("edited planes=%+v error=%v", got, err)
	}
	if got[0].Normal != [3]float32{1, 0, 0} || got[0].Distance != 23.5 || got[0].AxisType != 0 || got[1].Normal != [3]float32{0, -1, 0} || got[1].Distance != -4 || got[1].AxisType != 1 {
		t.Fatalf("committed plane values=%+v", got)
	}
	if reloaded.Header().Lumps[0].Length != int32(len(editedEntities)) || reloaded.Header().Lumps[1].Length != 40 {
		t.Fatalf("edited lengths=%d/%d", reloaded.Header().Lumps[0].Length, reloaded.Header().Lumps[1].Length)
	}
	for id := 2; id < bsp.LumpCount; id++ {
		current, err := reloaded.RawLump(bsp.LumpId(id))
		if err != nil {
			t.Fatal(err)
		}
		want := append([]byte(nil), payloads[id]...)
		if id == 35 {
			base := reloaded.Header().Lumps[id].Offset
			binary.LittleEndian.PutUint32(want[12:], uint32(base+22))
		}
		if !bytes.Equal(current, want) {
			t.Fatalf("editing entities/planes changed unrelated lump %d: %x, want %x", id, current, want)
		}
		oldDesc, newDesc := initial.Header().Lumps[id], reloaded.Header().Lumps[id]
		if newDesc.Length != oldDesc.Length || newDesc.Version != oldDesc.Version || newDesc.FourCC != oldDesc.FourCC {
			t.Fatalf("editing changed unrelated metadata %d: %+v, was %+v", id, newDesc, oldDesc)
		}
	}
	got[0].Distance = -99
	storedPlanes, err := file.Planes()
	if err != nil || len(storedPlanes) != 2 || storedPlanes[0].Distance != 23.5 {
		t.Fatalf("source planes after export=%+v error=%v", storedPlanes, err)
	}
	requireWriterSource(t, file, header, payloads)
}

func TestWriterMalformedGameFailureIsAtomic(t *testing.T) {
	for _, bad := range []string{"truncated directory", "child beyond container"} {
		t.Run(bad, func(t *testing.T) {
			file, payloads := writerAllLumps(t, bsp.ProfileSDK2013)
			writer := bsp.NewWriter(file)
			valid, err := writer.Write()
			if err != nil {
				t.Fatal(err)
			}
			original := append([]byte(nil), payloads[35]...)
			desc, err := file.Descriptor(35)
			if err != nil {
				t.Fatal(err)
			}
			broken := append([]byte(nil), original...)
			if bad == "truncated directory" {
				binary.LittleEndian.PutUint32(broken, 2) // Needs 36 bytes; only 27 exist.
			} else {
				binary.LittleEndian.PutUint32(broken[12:], uint32(desc.Offset+28))
			}
			if err := file.SetRawLump(35, desc, broken); err != nil {
				t.Fatal(err)
			}
			payloads[35] = broken
			header := file.Header()
			for attempt := 0; attempt < 2; attempt++ {
				output, err := writer.Write()
				if output != nil || !errors.Is(err, bsp.ErrMalformed) {
					t.Fatalf("failed export output=%x error=%v, want nil and ErrMalformed", output, err)
				}
				requireWriterSource(t, file, header, payloads)
			}
			// A failed full export must not poison the stored-byte API or the
			// writer; committing a valid directory restores identical output.
			stored, err := writer.WriteLump(35)
			if err != nil || !bytes.Equal(stored, broken) {
				t.Fatalf("stored malformed game=%x error=%v", stored, err)
			}
			if err := file.SetRawLump(35, desc, original); err != nil {
				t.Fatal(err)
			}
			restored, err := writer.Write()
			if err != nil || !bytes.Equal(restored, valid) {
				t.Fatalf("export after repair changed; error=%v", err)
			}
		})
	}
}

func TestWriterEmptyLumpsCanonicalizeLocationsOnly(t *testing.T) {
	file, err := bsp.ReadFromStream(bytes.NewReader(containerBytes(-1, nil)), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	for id := 0; id < bsp.LumpCount; id++ {
		desc := bsp.HeaderLump{Offset: int32(7000 + id), Length: 123, Version: int32(1000 + id), FourCC: [4]byte{byte(id + 1), 9, 8, 7}}
		if id == 35 {
			desc.FourCC = [4]byte{}
		}
		if err := file.SetRawLump(bsp.LumpId(id), desc, nil); err != nil {
			t.Fatal(err)
		}
	}
	header := file.Header()
	output, err := bsp.NewWriter(file).Write()
	if err != nil {
		t.Fatal(err)
	}
	if len(output) != 1036 {
		t.Fatalf("all-empty export length=%d, want 1036", len(output))
	}
	var normalized = header
	for id := 0; id < bsp.LumpCount; id++ {
		want := header.Lumps[id]
		want.Offset, want.Length = 0, 0
		normalized.Lumps[id] = want
		if got := writerWireDescriptor(output, id, bsp.ProfileSDK2013); got != want {
			t.Fatalf("empty wire descriptor %d=%+v, want %+v", id, got, want)
		}
	}
	reloaded, err := bsp.ReadFromStream(bytes.NewReader(output), bsp.ReadOptions{Profile: bsp.ProfileSDK2013})
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.Header() != normalized {
		t.Fatalf("empty normalized reimport=%+v, want %+v", reloaded.Header(), normalized)
	}
	requireWriterSource(t, file, header, [bsp.LumpCount][]byte{})
}
