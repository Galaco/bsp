package codec

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/galaco/bsp/lumps"
)

// Pakfile binds ZIP decoding to an owned, immutable payload. Entry readers are
// opened and closed inside ReadFile; no mutable zip.Reader escapes this view.
type Pakfile struct {
	Metadata
	raw      []byte
	reader   *zip.Reader
	names    []string
	parseErr error
}

func (pak *Pakfile) Decode(raw []byte) error {
	// A failed import must not make an earlier archive appear to be the new one.
	pak.reader, pak.names, pak.parseErr = nil, nil, nil
	pak.raw = nil
	pak.SetLength(0)
	if err := validatePCContext(pak.Context()); err != nil {
		pak.parseErr = err
		return err
	}
	if pak.Context().LumpVersion != 0 {
		pak.parseErr = fmt.Errorf("pakfile: unsupported revision %d", pak.Context().LumpVersion)
		return pak.parseErr
	}
	pak.raw = make([]byte, len(raw))
	copy(pak.raw, raw)
	pak.SetLength(len(raw))
	if len(raw) == 0 {
		return nil
	}
	reader, err := zip.NewReader(bytes.NewReader(pak.raw), int64(len(pak.raw)))
	if err != nil {
		pak.parseErr = err
		return err
	}
	names := make([]string, len(reader.File))
	for i, entry := range reader.File {
		names[i] = normalizePakName(entry.Name)
	}
	pak.reader, pak.names = reader, names
	return nil
}

func (pak *Pakfile) Encode() ([]byte, error) {
	raw := make([]byte, len(pak.raw))
	copy(raw, pak.raw)
	return raw, nil
}

// Entries returns metadata only, never externally mutable ZIP decoder state.
func (pak *Pakfile) Entries() []lumps.PakEntry {
	if pak == nil || pak.reader == nil {
		return nil
	}
	entries := make([]lumps.PakEntry, len(pak.reader.File))
	for i, entry := range pak.reader.File {
		entries[i] = lumps.PakEntry{Name: entry.Name, Method: entry.Method, Flags: entry.Flags, CRC32: entry.CRC32, CompressedSize: entry.CompressedSize64, UncompressedSize: entry.UncompressedSize64, Directory: entry.FileInfo().IsDir()}
	}
	return entries
}

// ReadFile uses case-insensitive slash-normalized lookup. The first matching
// directory entry wins; '..' and leading/trailing separators are not resolved.
// Missing paths wrap fs.ErrNotExist; empty files are successful empty reads.
func (pak *Pakfile) ReadFile(name string) ([]byte, error) {
	if pak != nil && pak.parseErr != nil {
		return nil, fmt.Errorf("pakfile lookup %q: %w", name, pak.parseErr)
	}
	if pak == nil || pak.reader == nil {
		return nil, &fs.PathError{Op: "read pakfile", Path: name, Err: fs.ErrNotExist}
	}
	normalized := normalizePakName(name)
	for i, entry := range pak.reader.File {
		if !strings.EqualFold(normalized, pak.names[i]) {
			continue
		}
		reader, err := entry.Open()
		if err != nil {
			return nil, fmt.Errorf("pakfile entry %q: %w", entry.Name, err)
		}
		data, readErr := io.ReadAll(reader)
		closeErr := reader.Close()
		if readErr != nil {
			return nil, fmt.Errorf("pakfile entry %q: %w", entry.Name, readErr)
		}
		if closeErr != nil {
			return nil, fmt.Errorf("pakfile entry %q: %w", entry.Name, closeErr)
		}
		return data, nil
	}
	return nil, &fs.PathError{Op: "read pakfile", Path: name, Err: fs.ErrNotExist}
}

func normalizePakName(name string) string {
	if !strings.Contains(name, "\\") && !strings.Contains(name, "//") {
		return name
	}
	var normalized strings.Builder
	normalized.Grow(len(name))
	slash := false
	for i := 0; i < len(name); i++ {
		value := name[i]
		if value == '\\' || value == '/' {
			if slash {
				continue
			}
			value = '/'
			slash = true
		} else {
			slash = false
		}
		normalized.WriteByte(value)
	}
	return normalized.String()
}
