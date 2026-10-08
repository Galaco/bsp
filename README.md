[![GoDoc](https://godoc.org/github.com/Galaco/bsp?status.svg)](https://godoc.org/github.com/Galaco/bsp)
[![Go report card](https://goreportcard.com/badge/github.com/galaco/bsp)](https://goreportcard.com/badge/github.com/galaco/bsp)
[![GolangCI](https://golangci.com/badges/github.com/galaco/bsp.svg)](https://golangci.com)
[![Build Status](https://travis-ci.com/Galaco/bsp.svg?branch=master)](https://travis-ci.com/Galaco/bsp)
[![CircleCI](https://circleci.com/gh/Galaco/bsp/tree/master.svg?style=svg)](https://circleci.com/gh/Galaco/bsp/tree/master)
[![codecov](https://codecov.io/gh/Galaco/bsp/branch/master/graph/badge.svg)](https://codecov.io/gh/Galaco/bsp)

# Bsp
Go library for handling Source Engine .bsp map files.

### Features:
* Read support for (most) non-xbox360 BSPs (v20,21). v19 support limited, may work
* Freely modify and resize any Lump data
* Limited write support, mostly untested

##### Not all lumps are currently supported, but can be freely read and modified, as they are treated as `[]byte`

The following lumps currently have a full implementation for v20 & v21 BSPs (tested against CS:S & CS:GO):

```
0: Entdata
1: Planes
2: Texdata
3: Vertexes
4: Visibility
5: Nodes
6: Texinfo
7: Faces
8: Lighting
9: Occlusion
10: Leafs
11: FaceId
12: Edges
13: Surfedges
14: Models
15: WorldLight
16: Leaffaces
17: LeafBrushes
18: Brushes
19: Brushsides
20: Areas
21: AreaPortals
26: DispInfo
27: OriginalFaces
28: PhysDisp
30: VertNormals
31: VertNormalIndices
33: DispVerts
34: DispLightmapSamplePosition
35: Game lump (partial: sprp only)
36: LeafWaterData
37: Primitives
38: PrimVerts
39: PrimIndices
40: Pakfile
41: ClipPortalVerts
42: Cubemaps
43: Texdatastringdata
44: Texdatastringtable
45: Overlays
46: LeafMinDistToWater
47: FaceMacroTextureInfo
48: DispTris
51: LeafAmbientIndexHDR
52: LeafAmbientIndex
53: WorldLightHDR
54: WorldLightHDR
55: LeafAmbientLightingHDR
56: LeafAmbientLighting
58: FacesHDR
59: MapFlags
60: OverlayFades
```

##### This library may reorganise lump order during the first export. This is intentional to handle lump resizing, but will change your checksum if you export without changes.

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
committed bytes. `NewWriter(file).WriteLump(id)` preserves stored bytes without
requiring typed decoding. Whole-file layout and game relocation remain subject to
the write-support limitations above.

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
ordinary records, game/static-prop decoding and compression are completed by the
corresponding ordered tasks; container acceptance does not assert every codec's
correctness. Big-endian/console files are rejected deliberately.

Pinned evidence: [Valve SDK header](https://github.com/ValveSoftware/source-sdk-2013/blob/b8cfb12c0e083a2ef5b2f9f9b50f3902fa034474/src/public/bspfile.h),
[Valve game records](https://github.com/ValveSoftware/source-sdk-2013/blob/b8cfb12c0e083a2ef5b2f9f9b50f3902fa034474/src/public/gamebspfile.h),
[L4D2 directory definitions](https://github.com/snake-biscuits/bsp_tool/blob/655627ea8704353ac14bd753df8ca05c15838eda/bsp_tool/branches/valve/left4dead2.py),
[BSPSource compiler-fixture provenance](https://github.com/ata4/bspsrc/blob/eaac417f16d3a4d8ef55bb637ecd7bff28740209/bspsrc-decompiler/src/test/java/info/ata4/bspsrc/decompiler/modules/texture/ToolTextureMatcherTestHelper.java),
and [Reactive Drop source definitions](https://github.com/ReactiveDrop/reactivedrop_public_src/blob/af3be02c1d7db2cd07c41105e7519cf8edd66ce3/src/public/bspfile.h).

## Real World examples
* Proof of concept BSP viewer: [https://github.com/Galaco/kero](https://github.com/Galaco/kero)
* Insert game_text newline placeholder characters (avoids Hammer crash) as a compile step: [https://github.com/Galaco/CS-GO-game_text-newline-inserter/tree/golang](https://github.com/Galaco/CS-GO-game_text-newline-inserter/tree/golang)
* Bspzip filelist generator from a mountable resource directory: [https://github.com/Galaco/bspzip-traverser](https://github.com/Galaco/bspzip-traverser)


# Contributing
All contributions welcome. Known unsupported games/maps are especially useful.
