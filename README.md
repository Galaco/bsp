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

Container/raw support and typed record support are separate. See the profile
matrix below. Branch-specific ordinary records, compression and static-prop
codecs are tracked by the ordered implementation tasks; directory acceptance
does not establish a correct typed layout for every payload.

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
(ordinary compressed size metadata), not a lump ID. Revision is map metadata,
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
`RawLump`/`SetRawLump` still preserve their bytes and metadata. Branch-specific
ordinary records, static-prop decoding and compression are completed by the
corresponding ordered tasks; container acceptance does not assert every codec's
correctness. Big-endian/console files are rejected deliberately.

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
uncompressed lengths and keeps sibling metadata/payloads. Canonical replacement
uses exactly one count, one 16-byte-entry directory and payload bytes—no extra
length prefixes. `MarshalGameLump(base)` relocates an unchanged container by
patching checked nonzero absolute offsets; gaps, trailing bytes, shared views,
unknown flags and zero-offset terminal conventions remain unchanged.

## Real World examples
* Proof of concept BSP viewer: [https://github.com/Galaco/kero](https://github.com/Galaco/kero)
* Insert game_text newline placeholder characters (avoids Hammer crash) as a compile step: [https://github.com/Galaco/CS-GO-game_text-newline-inserter/tree/golang](https://github.com/Galaco/CS-GO-game_text-newline-inserter/tree/golang)
* Bspzip filelist generator from a mountable resource directory: [https://github.com/Galaco/bspzip-traverser](https://github.com/Galaco/bspzip-traverser)


# Contributing
All contributions welcome. Known unsupported games/maps are especially useful.
