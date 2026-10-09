package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/ulikunitz/xz/lzma"
)

const (
	// DefaultMaxDecodedLumpBytes is a safety policy, not a Source format limit.
	DefaultMaxDecodedLumpBytes uint64 = 256 << 20
	// DefaultMaxLZMADictionaryBytes bounds the decoder's effective dictionary.
	DefaultMaxLZMADictionaryBytes uint64 = 64 << 20

	sourceLZMAHeaderBytes = 17
)

// ErrLimitExceeded identifies decoded-size, dictionary, and address-space limits.
var ErrLimitExceeded = errors.New("compression limit exceeded")

// Limits bounds allocations made when decoding a native Source LZMA envelope.
// Zero fields select the corresponding safety defaults.
type Limits struct {
	MaxDecodedLumpBytes    uint64
	MaxLZMADictionaryBytes uint64
}

// DecodeLZMA decodes a packed Source lzma_header_t and its advertised stream.
// expected is the descriptor's logical length, including zero for game entries.
// Ordinary lumps require an exact stored length; allowTrailing permits a game's
// retained gap/alignment bytes, which are never passed to the decoder. Neither
// the header nor the compressed payload is modified or copied.
// Known-size streams need not contain EOS, including the canonical zero-symbol
// stream whose five-byte range initialization is entirely zero.
//
// The envelope is defined by ValveSoftware/source-sdk-2013 at
// b2705ba55b3b802b86ef2b2dbf97939c9d4fb685, src/public/tier1/lzmaDecoder.h.
func DecodeLZMA(raw []byte, expected uint32, allowTrailing bool, limits Limits) ([]byte, error) {
	streamEnd, err := LZMAStoredSize(raw, expected, allowTrailing)
	if err != nil {
		return nil, err
	}
	actual := expected
	streamBytes := binary.LittleEndian.Uint32(raw[8:12])
	if _, err := lzma.PropertiesForCode(raw[12]); err != nil {
		return nil, err
	}

	if limits.MaxDecodedLumpBytes == 0 {
		limits.MaxDecodedLumpBytes = DefaultMaxDecodedLumpBytes
	}
	if limits.MaxLZMADictionaryBytes == 0 {
		limits.MaxLZMADictionaryBytes = DefaultMaxLZMADictionaryBytes
	}
	if uint64(actual) > limits.MaxDecodedLumpBytes {
		return nil, fmt.Errorf("decoded size %d exceeds limit %d: %w", actual, limits.MaxDecodedLumpBytes, ErrLimitExceeded)
	}
	dictionaryBytes := uint64(binary.LittleEndian.Uint32(raw[13:17]))
	if dictionaryBytes < uint64(lzma.MinDictCap) {
		dictionaryBytes = uint64(lzma.MinDictCap)
	}
	if dictionaryBytes > limits.MaxLZMADictionaryBytes {
		return nil, fmt.Errorf("effective dictionary size %d exceeds limit %d: %w", dictionaryBytes, limits.MaxLZMADictionaryBytes, ErrLimitExceeded)
	}
	maxInt := uint64(^uint(0) >> 1)
	// The dependency's dictionary ring allocates DictCap+1 bytes.
	if uint64(actual) > maxInt || dictionaryBytes >= maxInt {
		return nil, fmt.Errorf("decoded size or dictionary exceeds address space: %w", ErrLimitExceeded)
	}
	if actual == 0 && streamBytes == 5 && string(raw[sourceLZMAHeaderBytes:int(streamEnd)]) == "\x00\x00\x00\x00\x00" {
		// A known zero size and finished zero range state encode no symbols.
		// v0.5.11 instead reads an operation before checking size=0. Validate
		// this wire boundary directly; do not catch or suppress decoder errors.
		return []byte{}, nil
	}

	var header [lzma.HeaderLen]byte
	copy(header[:5], raw[12:17])
	binary.LittleEndian.PutUint64(header[5:], uint64(actual))
	stream := io.MultiReader(bytes.NewReader(header[:]), bytes.NewReader(raw[sourceLZMAHeaderBytes:int(streamEnd)]))
	// DictCap is a minimum, not a ceiling. The effective dictionary was
	// checked above; using the library minimum avoids its 8 MiB default.
	reader, err := (lzma.ReaderConfig{DictCap: lzma.MinDictCap}).NewReader(stream)
	if err != nil {
		return nil, err
	}
	decoded := make([]byte, int(actual))
	if _, err := io.ReadFull(reader, decoded); err != nil {
		return nil, err
	}
	// A known-size stream can finish without EOS. Do not read unboundedly
	// or require EOSMarker; do check errors after the requested output.
	var extra [1]byte
	n, err := reader.Read(extra[:])
	if err != nil && err != io.EOF {
		return nil, err
	}
	if n != 0 || err != io.EOF {
		return nil, errors.New("Source LZMA stream did not finish at its logical size")
	}
	return decoded, nil
}

// LZMAStoredSize validates framing and identifies the end of the advertised
// stream without decoding it. Game edits retain the bytes beyond this boundary.
func LZMAStoredSize(raw []byte, expected uint32, allowTrailing bool) (int, error) {
	if len(raw) < sourceLZMAHeaderBytes {
		return 0, io.ErrUnexpectedEOF
	}
	if string(raw[:4]) != "LZMA" {
		return 0, errors.New("invalid Source LZMA magic")
	}
	actual := binary.LittleEndian.Uint32(raw[4:8])
	if actual != expected {
		return 0, fmt.Errorf("Source LZMA size %d differs from logical size %d", actual, expected)
	}
	end := uint64(sourceLZMAHeaderBytes) + uint64(binary.LittleEndian.Uint32(raw[8:12]))
	if end > uint64(len(raw)) {
		return 0, io.ErrUnexpectedEOF
	}
	if !allowTrailing && end != uint64(len(raw)) {
		return 0, errors.New("Source LZMA stored length differs from advertised stream length")
	}
	return int(end), nil
}
