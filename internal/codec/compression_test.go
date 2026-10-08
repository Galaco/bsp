package codec

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"testing"
)

// Independently generated using liblzma 5.8.4, not the Go decoder/encoder.
// lzma_raw_buffer_encode used LZMA_FILTER_LZMA1EXT, preset 6, dict_size=4096,
// lc=3, lp=0, pb=2, ext_flags=0 (no EOS) or ALLOW_EOPM=1 (EOS). The alternate
// fixture uses lc=2, lp=1, pb=1. A Source envelope was packed as ASCII LZMA,
// uint32LE output size, uint32LE raw stream size, property byte, uint32LE dict.
// Each stream was independently decoded with liblzma lzma_raw_buffer_decode,
// the exact ext_size, and EOS forbidden for the no-EOS fixtures.
// LZMA1EXT semantics: https://tukaani.org/xz/liblzma-api/lzma12_8h.html
const (
	sourceLZMANoEOSHex      = "4c5a4d414c0000002d0000005d0010000000299bcaa7464d459b75515be774a9dd46fd03fd3ec267f8dbe0401b97e01a618cf8676c183ed6e231866e0200"
	sourceLZMAEOSHex        = "4c5a4d414c000000320000005d0010000000299bcaa7464d459b75515be774a9dd46fd03fd3ec267f8dbe0401b97e01a618cf8676c183ed6e231ec4aa351ffd92e0000"
	sourceLZMAEmptyHex      = "4c5a4d41000000000a0000005d001000000083fffbffffc0000000"
	sourceLZMAEmptyNoEOSHex = "4c5a4d4100000000050000005d001000000000000000"
	sourceLZMAAlternateHex  = "4c5a4d414c0000002d000000380010000000299bcb1cccba7215cbd08225b362120bd019b2f88fce9df58cd81a501b20c75b8b040eb94e04fc80c5ed59d6"
	sourceLZMAPlain         = "Source LZMA fixture\nSource LZMA fixture\nSource LZMA fixture\n" +
		"\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f"
)

func sourceLZMAFixture(t testing.TB, encoded string) []byte {
	t.Helper()
	raw, err := hex.DecodeString(encoded)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func TestDecodeLZMAIndependentFixtures(t *testing.T) {
	for _, test := range []struct {
		name  string
		wire  string
		plain string
	}{
		{"known-size-no-EOS", sourceLZMANoEOSHex, sourceLZMAPlain},
		{"known-size-EOS", sourceLZMAEOSHex, sourceLZMAPlain},
		{"empty-EOS", sourceLZMAEmptyHex, ""},
		{"empty-no-EOS", sourceLZMAEmptyNoEOSHex, ""},
		{"alternate-properties-no-EOS", sourceLZMAAlternateHex, sourceLZMAPlain},
	} {
		t.Run(test.name, func(t *testing.T) {
			raw := sourceLZMAFixture(t, test.wire)
			before := append([]byte(nil), raw...)
			decoded, err := DecodeLZMA(raw, uint32(len(test.plain)), false, Limits{})
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(decoded, []byte(test.plain)) {
				t.Fatalf("decoded %x, want %x", decoded, test.plain)
			}
			if len(decoded) > 0 {
				decoded[0] ^= 0xff
			}
			if !bytes.Equal(raw, before) {
				t.Fatal("decode or returned-output mutation changed stored input")
			}
		})
	}
}

func TestDecodeLZMAEmptyNoEOSBoundary(t *testing.T) {
	// The canonical zero-symbol range state is independently valid in
	// liblzma's known-size LZMA1EXT decoder. Only that fully validated state
	// bypasses the pinned Go decoder's operation-before-size=0 gap.
	for _, test := range []struct {
		name     string
		change   func([]byte) []byte
		trailing bool
		limits   Limits
		success  bool
	}{
		{"canonical", func(raw []byte) []byte { return raw }, false, Limits{}, true},
		{"game-gap", func(raw []byte) []byte { return append(raw, 255, 254) }, true, Limits{}, true},
		{"ordinary-trailer", func(raw []byte) []byte { return append(raw, 255) }, false, Limits{}, false},
		{"bad-magic", func(raw []byte) []byte { raw[0] = 'X'; return raw }, true, Limits{}, false},
		{"bad-property", func(raw []byte) []byte { raw[12] = 225; return raw }, true, Limits{}, false},
		{"effective-dictionary-limit", func(raw []byte) []byte { return raw }, true, Limits{MaxLZMADictionaryBytes: 4095}, false},
		{"nonzero-range-prefix", func(raw []byte) []byte { raw[17] = 1; return raw }, true, Limits{}, false},
		{"unfinished-range-state", func(raw []byte) []byte { raw[21] = 1; return raw }, true, Limits{}, false},
		{"truncated-stored-range", func(raw []byte) []byte { return raw[:21] }, true, Limits{}, false},
		{"short-advertised-range", func(raw []byte) []byte { binary.LittleEndian.PutUint32(raw[8:12], 4); return raw[:21] }, true, Limits{}, false},
		{"extra-advertised-range-byte", func(raw []byte) []byte { binary.LittleEndian.PutUint32(raw[8:12], 6); return append(raw, 0) }, true, Limits{}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			raw := test.change(sourceLZMAFixture(t, sourceLZMAEmptyNoEOSHex))
			decoded, err := DecodeLZMA(raw, 0, test.trailing, test.limits)
			if test.success {
				if err != nil || decoded == nil || len(decoded) != 0 {
					t.Fatalf("got %x, %v; want verified empty output", decoded, err)
				}
			} else if err == nil || decoded != nil {
				t.Fatalf("got %x, %v; want no output and error", decoded, err)
			}
		})
	}
}

func TestDecodeLZMATruncations(t *testing.T) {
	for _, wire := range []string{sourceLZMANoEOSHex, sourceLZMAEOSHex, sourceLZMAEmptyHex, sourceLZMAEmptyNoEOSHex, sourceLZMAAlternateHex} {
		raw := sourceLZMAFixture(t, wire)
		expected := binary.LittleEndian.Uint32(raw[4:8])
		for cut := 0; cut < len(raw); cut++ {
			for _, trailing := range []bool{false, true} {
				decoded, err := DecodeLZMA(raw[:cut], expected, trailing, Limits{})
				if !errors.Is(err, io.ErrUnexpectedEOF) || decoded != nil {
					t.Fatalf("cut=%d trailing=%v: got %x, %v; want no output and unexpected EOF", cut, trailing, decoded, err)
				}
			}
		}
	}
}

func TestDecodeLZMAMalformedEnvelopesAndStreams(t *testing.T) {
	for _, test := range []struct {
		name     string
		change   func([]byte) []byte
		expected uint32
		want     error
	}{
		{"magic", func(raw []byte) []byte { raw[0] = 'X'; return raw }, 76, nil},
		{"descriptor-size", func(raw []byte) []byte { return raw }, 75, nil},
		{"actual-size", func(raw []byte) []byte { binary.LittleEndian.PutUint32(raw[4:8], 75); return raw }, 76, nil},
		{"stream-too-long", func(raw []byte) []byte { binary.LittleEndian.PutUint32(raw[8:12], 46); return raw }, 76, io.ErrUnexpectedEOF},
		{"stream-size-overflow", func(raw []byte) []byte { binary.LittleEndian.PutUint32(raw[8:12], ^uint32(0)); return raw }, 76, io.ErrUnexpectedEOF},
		{"stream-too-short", func(raw []byte) []byte { binary.LittleEndian.PutUint32(raw[8:12], 44); return raw }, 76, nil},
		{"nonzero-range-prefix", func(raw []byte) []byte { raw[17] = 1; return raw }, 76, nil},
		{"invalid-range-state", func(raw []byte) []byte { copy(raw[18:22], []byte{255, 255, 255, 255}); return raw }, 76, nil},
		{"decoded-shorter-than-stream", func(raw []byte) []byte { binary.LittleEndian.PutUint32(raw[4:8], 75); return raw }, 75, nil},
		{"zero-size-with-nonempty-stream", func(raw []byte) []byte { binary.LittleEndian.PutUint32(raw[4:8], 0); return raw }, 0, nil},
		{"empty-range-stream", func(raw []byte) []byte { binary.LittleEndian.PutUint32(raw[8:12], 0); return raw[:17] }, 76, io.EOF},
		{"truncated-range-init", func(raw []byte) []byte { binary.LittleEndian.PutUint32(raw[8:12], 4); return raw[:21] }, 76, io.EOF},
		{"truncated-range-body", func(raw []byte) []byte { binary.LittleEndian.PutUint32(raw[8:12], 5); return raw[:22] }, 76, io.ErrUnexpectedEOF},
	} {
		t.Run(test.name, func(t *testing.T) {
			raw := test.change(sourceLZMAFixture(t, sourceLZMANoEOSHex))
			before := append([]byte(nil), raw...)
			decoded, err := DecodeLZMA(raw, test.expected, false, Limits{})
			if err == nil || decoded != nil || errors.Is(err, ErrLimitExceeded) {
				t.Fatalf("got %x, %v; want no output and a malformed-data error", decoded, err)
			}
			if test.want != nil && !errors.Is(err, test.want) {
				t.Fatalf("got %v; want error matching %v", err, test.want)
			}
			if !bytes.Equal(raw, before) {
				t.Fatal("failed decode changed stored input")
			}
		})
	}
	for property := 225; property <= 255; property++ {
		raw := sourceLZMAFixture(t, sourceLZMANoEOSHex)
		raw[12] = byte(property)
		decoded, err := DecodeLZMA(raw, 76, false, Limits{})
		if err == nil || decoded != nil || errors.Is(err, ErrLimitExceeded) {
			t.Fatalf("property=%d: got %x, %v; want invalid-property error", property, decoded, err)
		}
	}
}

func TestDecodeLZMAGameTrailerIsNotDecoded(t *testing.T) {
	raw := sourceLZMAFixture(t, sourceLZMANoEOSHex)
	stored := append(append([]byte(nil), raw...), []byte{0xff, 0xde, 0xad, 0xbe, 0xef}...)
	before := append([]byte(nil), stored...)
	decoded, err := DecodeLZMA(stored, 76, true, Limits{})
	if err != nil || string(decoded) != sourceLZMAPlain {
		t.Fatalf("game trailer: got %x, %v", decoded, err)
	}
	if !bytes.Equal(stored, before) {
		t.Fatal("game trailer changed")
	}
	if decoded, err := DecodeLZMA(stored, 76, false, Limits{}); err == nil || decoded != nil {
		t.Fatalf("ordinary trailer: got %x, %v; want failure", decoded, err)
	}
	// The remaining genuine compressed bytes now reside in a game gap. A
	// decoder must not borrow them to repair the advertised truncated stream.
	binary.LittleEndian.PutUint32(stored[8:12], 5)
	if decoded, err := DecodeLZMA(stored, 76, true, Limits{}); !errors.Is(err, io.ErrUnexpectedEOF) || decoded != nil {
		t.Fatalf("stream escaped advertised boundary: got %x, %v", decoded, err)
	}
}

func TestDecodeLZMALimits(t *testing.T) {
	for _, test := range []struct {
		name       string
		actual     uint32
		dictionary uint32
		limits     Limits
	}{
		{"decoded-one-over", 76, 4096, Limits{MaxDecodedLumpBytes: 75}},
		{"dictionary-one-over", 76, 4096, Limits{MaxLZMADictionaryBytes: 4095}},
		{"zero-dictionary-effective-minimum", 76, 0, Limits{MaxLZMADictionaryBytes: 4095}},
		{"small-dictionary-effective-minimum", 76, 1, Limits{MaxLZMADictionaryBytes: 4095}},
		{"decoded-default", 256<<20 + 1, 4096, Limits{}},
		{"decoded-default-with-custom-dictionary", 256<<20 + 1, 4096, Limits{MaxLZMADictionaryBytes: 4096}},
		{"dictionary-default", 76, 64<<20 + 1, Limits{}},
		{"dictionary-default-with-custom-output", 76, 64<<20 + 1, Limits{MaxDecodedLumpBytes: 76}},
		{"maximum-advertised-output", ^uint32(0), 4096, Limits{}},
		{"maximum-advertised-dictionary", 76, ^uint32(0), Limits{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			raw := sourceLZMAFixture(t, sourceLZMANoEOSHex)
			binary.LittleEndian.PutUint32(raw[4:8], test.actual)
			binary.LittleEndian.PutUint32(raw[13:17], test.dictionary)
			decoded, err := DecodeLZMA(raw, test.actual, false, test.limits)
			if !errors.Is(err, ErrLimitExceeded) || decoded != nil {
				t.Fatalf("got %x, %v; want no output and limit error", decoded, err)
			}
		})
	}
	for _, dictionary := range []uint32{0, 1, 4095, 4096, 4097} {
		t.Run(fmt.Sprintf("accepted-dictionary-%d", dictionary), func(t *testing.T) {
			raw := sourceLZMAFixture(t, sourceLZMANoEOSHex)
			binary.LittleEndian.PutUint32(raw[13:17], dictionary)
			limit := uint64(4096)
			if uint64(dictionary) > limit {
				limit = uint64(dictionary)
			}
			decoded, err := DecodeLZMA(raw, 76, false, Limits{MaxDecodedLumpBytes: 76, MaxLZMADictionaryBytes: limit})
			if err != nil || string(decoded) != sourceLZMAPlain {
				t.Fatalf("got %x, %v; want exact-boundary success", decoded, err)
			}
		})
	}
	// Raising the output default admits the size check, but the invalid range
	// prefix fails decoder construction before a large output can be allocated.
	raw := sourceLZMAFixture(t, sourceLZMANoEOSHex)
	binary.LittleEndian.PutUint32(raw[4:8], 256<<20+1)
	raw[17] = 1
	if decoded, err := DecodeLZMA(raw, 256<<20+1, false, Limits{MaxDecodedLumpBytes: 256<<20 + 1}); err == nil || errors.Is(err, ErrLimitExceeded) || decoded != nil {
		t.Fatalf("raised decoded limit: got %x, %v; want decoder error, not limit error", decoded, err)
	}
	// On 32-bit systems the dependency's dictionary allocation is DictCap+1,
	// so maxInt itself must be rejected even with caller policy limits lifted.
	maxInt := uint64(^uint(0) >> 1)
	if maxInt <= uint64(^uint32(0)) {
		raw := sourceLZMAFixture(t, sourceLZMANoEOSHex)
		binary.LittleEndian.PutUint32(raw[13:17], uint32(maxInt))
		if decoded, err := DecodeLZMA(raw, 76, false, Limits{MaxDecodedLumpBytes: ^uint64(0), MaxLZMADictionaryBytes: ^uint64(0)}); !errors.Is(err, ErrLimitExceeded) || decoded != nil {
			t.Fatalf("native dictionary bound: got %x, %v", decoded, err)
		}
		binary.LittleEndian.PutUint32(raw[4:8], ^uint32(0))
		binary.LittleEndian.PutUint32(raw[13:17], 4096)
		if decoded, err := DecodeLZMA(raw, ^uint32(0), false, Limits{MaxDecodedLumpBytes: ^uint64(0), MaxLZMADictionaryBytes: ^uint64(0)}); !errors.Is(err, ErrLimitExceeded) || decoded != nil {
			t.Fatalf("native output bound: got %x, %v", decoded, err)
		}
	}
}

func FuzzDecodeLZMA(f *testing.F) {
	f.Add([]byte{}, uint32(0), false)
	for _, wire := range []string{sourceLZMANoEOSHex, sourceLZMAEOSHex, sourceLZMAEmptyHex, sourceLZMAEmptyNoEOSHex, sourceLZMAAlternateHex} {
		raw := sourceLZMAFixture(f, wire)
		f.Add(raw, binary.LittleEndian.Uint32(raw[4:8]), false)
		f.Add(append(raw, 0xff, 0xff), binary.LittleEndian.Uint32(raw[4:8]), true)
	}
	f.Add([]byte("LZMA"), uint32(0), false)
	malformed := sourceLZMAFixture(f, sourceLZMANoEOSHex)
	binary.LittleEndian.PutUint32(malformed[4:8], 4097)
	f.Add(malformed, uint32(4097), false)
	f.Fuzz(func(t *testing.T, raw []byte, expected uint32, trailing bool) {
		if len(raw) > 4096 {
			return
		}
		before := append([]byte(nil), raw...)
		// Explicit small policy defaults keep fuzz mutations of advertised
		// lengths/dictionaries from requesting the normal production limits.
		decoded, err := DecodeLZMA(raw, expected, trailing, Limits{MaxDecodedLumpBytes: 4096, MaxLZMADictionaryBytes: 4096})
		if !bytes.Equal(raw, before) {
			t.Fatal("input mutated")
		}
		if err != nil {
			if decoded != nil {
				t.Fatal("partial output returned on failure")
			}
			return
		}
		if len(decoded) != int(expected) || len(decoded) > 4096 {
			t.Fatalf("decoded length %d differs from bounded expected size %d", len(decoded), expected)
		}
		repeated, err := DecodeLZMA(raw, expected, trailing, Limits{MaxDecodedLumpBytes: 4096, MaxLZMADictionaryBytes: 4096})
		if err != nil || !bytes.Equal(repeated, decoded) || !bytes.Equal(raw, before) {
			t.Fatalf("repeat decode changed output or source: %v", err)
		}
	})
}
