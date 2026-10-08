[![GoDoc](https://godoc.org/github.com/Galaco/bsp?status.svg)](https://godoc.org/github.com/Galaco/bsp)
[![Go report card](https://goreportcard.com/badge/github.com/galaco/bsp)](https://goreportcard.com/badge/github.com/galaco/bsp)
[![GolangCI](https://golangci.com/badges/github.com/galaco/bsp.svg)](https://golangci.com)
[![Build Status](https://travis-ci.com/Galaco/bsp.svg?branch=master)](https://travis-ci.com/Galaco/bsp)
[![CircleCI](https://circleci.com/gh/Galaco/bsp/tree/master.svg?style=svg)](https://circleci.com/gh/Galaco/bsp/tree/master)
[![codecov](https://codecov.io/gh/Galaco/bsp/branch/master/graph/badge.svg)](https://codecov.io/gh/Galaco/bsp)

# Bsp
Go library for handling Source Engine .bsp map files.

### Features

* Explicit little-endian PC profiles for SDK2013 BSP19/20, CS:GO BSP21 and L4D2 BSP21.
* Owned raw replacement and independent typed snapshots; edits require setters.
* Canonical full-file export covering all 64 IDs, preserving opaque bytes and metadata.
* Bounded game directories and pure absolute-offset relocation, including unknown children.
* Native Source LZMA decoding with bounded output/dictionary allocations; unchanged streams stay stored.

Container/raw support and typed record support are separate. See the profile
matrix below. Branch-specific ordinary records and static-prop codecs are tracked
by the ordered implementation tasks; directory acceptance does not establish a
correct typed layout for every payload.

Canonical export intentionally uses numeric ID order, four-byte alignment and
zero locations for empty descriptors. It does not preserve original physical
ordering, inter-lump gaps or unrelated trailing bytes. Game absolute-offset
relocation changes directory bytes; checksum convention/order is a separate task.

# Usage

Minimal example of obtaining entdata from a BSP. The following will print the entdata
lump (entdata is a single json-like string) of a specified .bsp to console.

```go
package main

import (
	"github.com/galaco/bsp"
	"log"
)

func main() {
	file, err := bsp.ReadFromFile("de_dust2.bsp", bsp.ReadOptions{Profile: bsp.ProfileCSGO})
	if err != nil {
		log.Fatal(err)
	}
	text, err := file.Entities()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(text)
}
```

Raw and decoded values are independent snapshots. Use `SetRawLump`, `SetEntities`,
or `SetPlanes` to commit an edit; modifying a returned slice does not edit the file.
`Header`/`Descriptor` return copied metadata, and stored lengths are derived from
committed bytes. `NewWriter(file).Write()` plans checked sizes before allocating,
then emits committed payloads without typed decoding or source mutation.
`WriteLump(id)` returns stored bytes; `MarshalGameLump(base)` explicitly relocates
game bytes. Input descriptor offsets remain historical, not planned output locations.

## Format selection

Pass `ReadOptions{Profile: ...}` explicitly. File version alone cannot identify an
engine dialect; zero options do not guess. Header identity/profile checks happen
before payload buffering. `HeaderLump.FourCC` preserves the fourth wire word
(ordinary decoded-size metadata), not a lump ID. Revision is map metadata,
not a decoder selector.

| Profile | BSP versions | Descriptor order | Exercised container/record baseline |
| --- | --- | --- | --- |
| `ProfileSDK2013` | 19–20 | offset, length, version, fourth word | Valve `dm_lockdown.bsp`: BSP19, leaf0; entities, planes and leaves decoded |
| `ProfileCSGO` | 21 | standard | `ar_baggage.bsp.gz`: leaf1; entities, planes and leaves decoded; game entries `sprp`10/`dprp`4 preserved |
| `ProfileL4D2` | 21 | version, offset, length, fourth word | BSPSource's game-compiled `l4d2_tooltextures.bsp`: leaf1; entities, planes and leaves decoded |
| Reactive Drop | not admitted yet | source-defined standard | Public declarations exist; genuine compiled-fixture gate remains open |

Codec selection uses the actual descriptor revision and retained engine context
for both decode and encode. Leaf revisions 0/1 have explicit 56/32-byte layouts.
Unknown revisions and opaque IDs return `ErrUnsupportedLump` on typed access;
`RawLump`/`SetRawLump` still preserve their bytes and metadata. Container acceptance
does not assert every codec's correctness. Big-endian/console files are rejected
deliberately; unproven branch layouts stay opaque.

Pinned evidence: [Valve SDK header](https://github.com/ValveSoftware/source-sdk-2013/blob/b8cfb12c0e083a2ef5b2f9f9b50f3902fa034474/src/public/bspfile.h),
[Valve game records](https://github.com/ValveSoftware/source-sdk-2013/blob/b8cfb12c0e083a2ef5b2f9f9b50f3902fa034474/src/public/gamebspfile.h),
[L4D2 directory definitions](https://github.com/snake-biscuits/bsp_tool/blob/655627ea8704353ac14bd753df8ca05c15838eda/bsp_tool/branches/valve/left4dead2.py),
[BSPSource compiler-fixture provenance](https://github.com/ata4/bspsrc/blob/eaac417f16d3a4d8ef55bb637ecd7bff28740209/bspsrc-decompiler/src/test/java/info/ata4/bspsrc/decompiler/modules/texture/ToolTextureMatcherTestHelper.java),
and [Reactive Drop source definitions](https://github.com/ReactiveDrop/reactivedrop_public_src/blob/af3be02c1d7db2cd07c41105e7519cf8edd66ce3/src/public/bspfile.h).

## Game directory ownership

`GameLumps()` returns ordered owned entry snapshots, including unknown IDs and a
terminal boundary when present. Use directory indices for duplicates;
`GameLumpByID` returns `ErrAmbiguousGameLumpID` rather than dropping entries.
`StoredData` is opaque physical bytes; compressed entries' `LogicalLength` can
differ. Decompression/static-prop record decoding is separate from the directory.

`SetGameLump(index, entry)` or `SetGameLumps(entries)` commits atomically, derives
uncompressed lengths and keeps sibling metadata/payloads. Indexed edits preserve
unclaimed game gaps; rebuilding aligns payload starts to four bytes. Whole-directory
replacement uses one count, one 16-byte-entry directory and payloads—no extra
length prefixes. `MarshalGameLump(base)` relocates an unchanged container by
patching checked nonzero absolute offsets; gaps, trailing bytes, shared views,
unknown flags and zero-offset terminal conventions remain unchanged.

## Native lump compression

`RawLump` and game `StoredData` return owned exact stored bytes. `DecodedLump`,
`DecodedGameLump` and typed getters decode only explicitly marked Source LZMA
envelopes; unmarked `LZMA` prefixes are ordinary data. Game compression belongs to
individual children, not the outer lump. Known-size streams with or without EOS
are supported; game alignment/trailer bytes are retained but never decoded.

`ReadOptions.MaxDecodedLumpBytes` defaults to 256 MiB and
`MaxLZMADictionaryBytes` to 64 MiB. Zero selects defaults; negative options return
`ErrInvalidOptions`. Limits are overridable safety policy, not format limits.
Logical access above a limit returns `ErrLimitExceeded` before decoder allocation.
Malformed frames return contextual errors; raw reads/export retain opaque bytes.

Typed setters emit uncompressed data atomically and clear only applicable
compression metadata. `SetDecodedGameLump(index, bytes)` retains ID/version, other
flags and stream trailer/gaps. Unchanged compressed siblings remain stored; no
recompression or external executable/cgo is required. The Go decoder is pinned
to `github.com/ulikunitz/xz` v0.5.11.


## Versioned ordinary records

Public records live in `lumps`; private stateless wire codecs live in
`internal/codec`. Use attached `Leaves`, `WorldLights`/`WorldLightsHDR`,
`BrushSides`, `AmbientIndices`/`AmbientSamples` (and HDR counterparts), and
`Occlusion` getters/setters rather than decoder objects.

| Records | Supported wire layouts |
| --- | --- |
| Leaves | descriptor0:56 bytes with inline LDR cube; descriptor1:32 bytes |
| Worldlights | revision0:88 bytes; CS:GO revision1:100 bytes with shadow offset; SDK/L4D2 revision1 remains opaque |
| Brush sides | eight bytes; SDK/L4D2 signed16 bevel, CS:GO byte bevel/thin |
| Ambient | index4 bytes; legacy cube24, indexed positioned sample28 |
| Occlusion | revision1 object36/revision2 object40; counted polygon12 and vertex-index4 arrays |

Leaf/sample records expose layout presence and padding. Setters reject fields
absent from the selected layout; they never silently upgrade versions or discard
data. Worldlight emission enums retain the complete signed32-bit value.

Ambient selection follows sibling context, not divisibility: leaf0 uses inline
LDR lighting and rejects nonempty ignored standalone channels on typed access;
compact leaves with absent LDR index52 use legacy cubes for both channels,
regardless unreliable sample revisions. Available indices/cubes must match
leaves; positive spans address samples, zero-count indices reference leaves.
Absent channels stay absent, with no fallback or implicit layout migration.
Occlusion getters/setters validate polygon spans and actual plane/vertex references.

## Real World examples
* Proof of concept BSP viewer: [https://github.com/Galaco/kero](https://github.com/Galaco/kero)
* Insert game_text newline placeholder characters (avoids Hammer crash) as a compile step: [https://github.com/Galaco/CS-GO-game_text-newline-inserter/tree/golang](https://github.com/Galaco/CS-GO-game_text-newline-inserter/tree/golang)
* Bspzip filelist generator from a mountable resource directory: [https://github.com/Galaco/bspzip-traverser](https://github.com/Galaco/bspzip-traverser)


# Contributing
All contributions welcome. Known unsupported games/maps are especially useful.
