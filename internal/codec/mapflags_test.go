package codec

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/galaco/bsp/internal/format"
	"github.com/galaco/bsp/lumps"
)

func mapFlagsTestContexts() []format.Context {
	return []format.Context{
		{Profile: "sdk2013", BSPVersion: 19, ByteOrder: binary.LittleEndian},
		{Profile: "sdk2013", BSPVersion: 20, ByteOrder: binary.LittleEndian},
		{Profile: "csgo", BSPVersion: 21, ByteOrder: binary.LittleEndian},
		{Profile: "l4d2", BSPVersion: 21, ByteOrder: binary.LittleEndian, HeaderLayout: format.L4D2Header},
		{Profile: "reactivedrop", BSPVersion: 21, ByteOrder: binary.LittleEndian},
	}
}

func TestMapFlagsWirePresenceAndAllBits(t *testing.T) {
	cases := []struct {
		raw     []byte
		value   lumps.MapFlags
		present bool
	}{
		{nil, lumps.MapFlags{}, false},
		{[]byte{}, lumps.MapFlags{}, false},
		{[]byte{0, 0, 0, 0}, lumps.MapFlags{}, true},
		{[]byte{0x67, 0x45, 0x23, 0x81}, lumps.MapFlags{LevelFlags: 0x81234567}, true},
		{[]byte{0xff, 0xff, 0xff, 0xff}, lumps.MapFlags{LevelFlags: 0xffffffff}, true},
	}
	for _, context := range mapFlagsTestContexts() {
		for _, tc := range cases {
			value, present, err := DecodeMapFlags(tc.raw, context)
			if err != nil || value != tc.value || present != tc.present {
				t.Fatalf("%+v raw=%x: value=%+v present=%v error=%v", context, tc.raw, value, present, err)
			}
			raw, err := EncodeMapFlags(tc.value, tc.present, context)
			if err != nil || !bytes.Equal(raw, tc.raw) {
				t.Fatalf("%+v value=%+v present=%v: raw=%x error=%v", context, tc.value, tc.present, raw, err)
			}
		}
		if raw, err := EncodeMapFlags(lumps.MapFlags{LevelFlags: 0x80000000}, false, context); err == nil || raw != nil {
			t.Fatalf("%+v: absent nonzero value was discarded: raw=%x error=%v", context, raw, err)
		}
	}
}

func TestMapFlagsRejectsMalformedLengthAndUnsupportedContext(t *testing.T) {
	for _, context := range mapFlagsTestContexts() {
		for _, length := range []int{1, 2, 3, 5, 6, 7, 8, 16} {
			value, present, err := DecodeMapFlags(make([]byte, length), context)
			if err == nil || present || value != (lumps.MapFlags{}) {
				t.Fatalf("%+v accepted length %d: value=%+v present=%v error=%v", context, length, value, present, err)
			}
		}
	}
	invalid := []format.Context{
		{},
		{Profile: "unknown", BSPVersion: 21},
		{Profile: "sdk2013", BSPVersion: 21},
		{Profile: "csgo", BSPVersion: 20},
		{Profile: "l4d2", BSPVersion: 20},
		{Profile: "reactivedrop", BSPVersion: 19},
		{Profile: "reactivedrop", BSPVersion: 20},
		{Profile: "reactivedrop", BSPVersion: 21, HeaderLayout: format.L4D2Header},
	}
	for _, context := range mapFlagsTestContexts() {
		bigEndian := context
		bigEndian.ByteOrder = binary.BigEndian
		invalid = append(invalid, bigEndian)
		for _, revision := range []int32{-1, 1, 2} {
			unsupported := context
			unsupported.LumpVersion = revision
			invalid = append(invalid, unsupported)
		}
	}
	for _, context := range invalid {
		for _, raw := range [][]byte{nil, {0, 0, 0, 0}} {
			value, present, err := DecodeMapFlags(raw, context)
			if err == nil || present || value != (lumps.MapFlags{}) {
				t.Fatalf("decode accepted invalid context %+v: value=%+v present=%v error=%v", context, value, present, err)
			}
		}
		for _, present := range []bool{false, true} {
			if raw, err := EncodeMapFlags(lumps.MapFlags{}, present, context); err == nil || raw != nil {
				t.Fatalf("encode accepted invalid context %+v: present=%v raw=%x error=%v", context, present, raw, err)
			}
		}
	}
}

func TestMapFlagsAdapterTransactionalDecodeAndSnapshots(t *testing.T) {
	for _, context := range mapFlagsTestContexts() {
		var adapter MapFlags
		adapter.SetContext(context)
		for _, initial := range [][]byte{nil, {0, 0, 0, 0}, {0x67, 0x45, 0x23, 0x81}} {
			if err := adapter.Unmarshall(initial); err != nil {
				t.Fatal(err)
			}
			before, present, length := adapter.GetData(), adapter.Present(), adapter.Length()
			for _, malformed := range [][]byte{{1}, {1, 2}, {1, 2, 3}, {1, 2, 3, 4, 5}, make([]byte, 8)} {
				if err := adapter.Unmarshall(malformed); err == nil {
					t.Fatalf("accepted malformed raw=%x", malformed)
				}
				if adapter.GetData() != before || adapter.Present() != present || adapter.Length() != length {
					t.Fatal("failed decode partially published state")
				}
				if raw, err := adapter.Marshall(); err != nil || !bytes.Equal(raw, initial) {
					t.Fatalf("failed decode changed encoding: raw=%x error=%v", raw, err)
				}
			}
		}
		if err := adapter.Unmarshall(nil); err != nil || adapter.Present() || adapter.Length() != 0 || adapter.GetData() != (lumps.MapFlags{}) {
			t.Fatalf("decode removal: value=%+v present=%v length=%d error=%v", adapter.GetData(), adapter.Present(), adapter.Length(), err)
		}
		if raw, err := adapter.Marshall(); err != nil || len(raw) != 0 {
			t.Fatalf("removed encoding: raw=%x error=%v", raw, err)
		}
	}
}
