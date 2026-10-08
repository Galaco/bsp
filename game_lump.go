package bsp

import (
	"encoding/binary"
	"fmt"
	"github.com/galaco/bsp/internal/codec"
	"github.com/galaco/bsp/lumps"
)

func (bsp *File) gameDirectoryView() (codec.GameDirectory, error) {
	desc := bsp.header.Lumps[LumpGame]
	if binary.LittleEndian.Uint32(desc.FourCC[:]) != 0 {
		return codec.GameDirectory{}, fmt.Errorf("%w: compressed outer game lump", ErrUnsupportedLump)
	}
	return codec.GameDirectory{Raw: bsp.lumps[LumpGame].raw, Base: desc.Offset}, nil
}
func (bsp *File) gameDirectory() (codec.GameDirectory, error) {
	view, err := bsp.gameDirectoryView()
	if err != nil {
		return codec.GameDirectory{}, err
	}
	directory, err := codec.ParseGameDirectory(view.Raw, view.Base)
	if err != nil {
		return codec.GameDirectory{}, &FormatError{"game directory", err}
	}
	return directory, nil
}
func gameSnapshot(entry lumps.GameLump) lumps.GameLump {
	entry.StoredData = cloneBytes(entry.StoredData)
	return entry
}
func terminalGameEntry(entry lumps.GameLump) bool {
	return entry.ID == 0 && entry.Flags == 0 && entry.Version == 0 && entry.LogicalLength == 0
}

// GameLumps returns ordered, independently owned stored-entry snapshots.
// The terminal boundary is preserved; compressed payloads are not decoded here.
func (bsp *File) GameLumps() ([]lumps.GameLump, error) {
	directory, err := bsp.gameDirectory()
	if err != nil {
		return nil, err
	}
	entries := directory.Entries
	for i := range entries {
		entries[i].StoredData = cloneBytes(entries[i].StoredData)
	}
	return entries, nil
}

// GameLump addresses a directory index, preserving duplicate IDs.
func (bsp *File) GameLump(index int) (lumps.GameLump, error) {
	directory, err := bsp.gameDirectory()
	if err != nil {
		return lumps.GameLump{}, err
	}
	if index < 0 || index >= len(directory.Entries) {
		return lumps.GameLump{}, fmt.Errorf("%w: %d", ErrInvalidGameLumpIndex, index)
	}
	return gameSnapshot(directory.Entries[index]), nil
}

// GameLumpByID excludes the terminal boundary and rejects ambiguous duplicate IDs.
func (bsp *File) GameLumpByID(id uint32) (lumps.GameLump, error) {
	directory, err := bsp.gameDirectory()
	if err != nil {
		return lumps.GameLump{}, err
	}
	index, err := gameEntryIndex(directory.Entries, id)
	if err != nil {
		return lumps.GameLump{}, err
	}
	return gameSnapshot(directory.Entries[index]), nil
}

func gameEntryIndex(entries []lumps.GameLump, id uint32) (int, error) {
	found := -1
	for i, entry := range entries {
		if i == len(entries)-1 && terminalGameEntry(entry) {
			continue
		}
		if entry.ID != id {
			continue
		}
		if found >= 0 {
			return -1, fmt.Errorf("%w: %#x", ErrAmbiguousGameLumpID, id)
		}
		found = i
	}
	if found < 0 {
		return -1, fmt.Errorf("%w: game ID %#x", ErrLumpNotFound, id)
	}
	return found, nil
}

// SetGameLumps atomically encodes an ordered directory at an internal base of zero.
// Uncompressed lengths and all nonterminal positions are derived from stored bytes.
func (bsp *File) SetGameLumps(entries []lumps.GameLump) error {
	return bsp.commitGameEntries(entries)
}

// SetGameLump replaces one indexed entry without losing siblings or duplicate IDs.
func (bsp *File) SetGameLump(index int, entry lumps.GameLump) error {
	directory, err := bsp.gameDirectory()
	if err != nil {
		return err
	}
	if index < 0 || index >= len(directory.Entries) {
		return fmt.Errorf("%w: %d", ErrInvalidGameLumpIndex, index)
	}
	directory.Entries[index] = entry
	return bsp.commitGameDirectoryEdit(directory, index, nil)
}
func (bsp *File) commitGameEntries(entries []lumps.GameLump) error {
	raw, err := codec.EncodeGameDirectory(entries, 0)
	if err != nil {
		return &FormatError{"encode game directory", err}
	}
	bsp.storeGameDirectory(raw)
	return nil
}

func (bsp *File) commitGameDirectoryEdit(directory codec.GameDirectory, index int, tail []byte) error {
	raw, err := codec.EncodeGameDirectoryEdit(directory, index, tail, 0)
	if err != nil {
		return &FormatError{"encode game directory edit", err}
	}
	bsp.storeGameDirectory(raw)
	return nil
}

func (bsp *File) storeGameDirectory(raw []byte) {
	desc := bsp.header.Lumps[LumpGame]
	desc.Offset = 0
	desc.FourCC = [4]byte{}
	bsp.storeOwned(LumpGame, desc, raw)
}

// MarshalGameLump relocates an unchanged stored directory without mutating the file.
// It preserves shared views, unknown bytes, gaps and terminal offset conventions.
func (bsp *File) MarshalGameLump(base int32) ([]byte, error) {
	directory, err := bsp.gameDirectoryView()
	if err != nil {
		return nil, err
	}
	raw, err := codec.RelocateGameDirectory(directory, base)
	if err != nil {
		return nil, &FormatError{"relocate game directory", err}
	}
	return raw, nil
}
