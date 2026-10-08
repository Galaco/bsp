[![GoDoc](https://godoc.org/github.com/Galaco/bsp?status.svg)](https://godoc.org/github.com/Galaco/bsp)
[![Go report card](https://goreportcard.com/badge/github.com/galaco/bsp)](https://goreportcard.com/badge/github.com/galaco/bsp)
[![GolangCI](https://golangci.com/badges/github.com/galaco/bsp.svg)](https://golangci.com)
[![Build Status](https://travis-ci.com/Galaco/bsp.svg?branch=master)](https://travis-ci.com/Galaco/bsp)
[![CircleCI](https://circleci.com/gh/Galaco/bsp/tree/master.svg?style=svg)](https://circleci.com/gh/Galaco/bsp/tree/master)
[![codecov](https://codecov.io/gh/Galaco/bsp/branch/master/graph/badge.svg)](https://codecov.io/gh/Galaco/bsp)

# Bsp
Go library for handling Source Engine .bsp map files.

### Features

* Explicit little-endian PC profiles for SDK2013 BSP19/20 and CS:GO/L4D2/Reactive Drop BSP21.
* Owned raw replacement and independent typed snapshots; edits require setters.
* Canonical full-file export covering all 64 IDs, preserving opaque bytes and metadata.
* Bounded game directories and pure absolute-offset relocation, including unknown children.
* Native Source LZMA decoding with bounded output/dictionary allocations; unchanged streams stay stored.

Container/raw support and typed record support are separate. See the profile
matrix and record sections below. Remaining layout fixes follow the ordered
implementation tasks; directory acceptance does not establish a typed layout
for every payload.

Canonical export intentionally uses numeric ID order, four-byte alignment and
zero locations for empty descriptors. It does not preserve original physical
ordering, inter-lump gaps or unrelated trailing bytes. Game absolute-offset
relocation changes directory bytes; checksum convention/order is a separate task.

# Usage

Minimal example printing the exact entity text of a specified `.bsp`:

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
| `ProfileCSGO` | 21 | standard | `ar_baggage.bsp.gz`: leaf1 and correctly aligned sprp10; authored v11 maps independently corroborated |
| `ProfileL4D2` | 21 | version, offset, length, fourth word | BSPSource's game-compiled `l4d2_tooltextures.bsp`: leaf1; entities, planes and leaves decoded |
| `ProfileReactiveDrop` | 21 | standard | Steam-distributed `lobby`, `dm_testlab`, `dm_desert`: sprp9, worldlight1, thin brush sides, displacement streams and system levels |

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
| Worldlights | revision0:88 bytes; CS:GO/Reactive Drop revision1:100 bytes with shadow offset; SDK/L4D2 revision1 remains opaque |
| Brush sides | eight bytes; SDK/L4D2 signed16 bevel, CS:GO/Reactive Drop byte bevel/thin |
| Ambient | index4 bytes; legacy cube24, indexed positioned sample28 |
| Occlusion | revision1 object36/revision2 object40; counted polygon12 and vertex-index4 arrays; Reactive Drop accepts revision2 only |

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

## Spatial records and fixed-width safety

Use `Overlays`/`SetOverlays`, `Cubemaps`/`SetCubemaps` and
`Primitives`/`SetPrimitives` for owned snapshots and atomic edits:

| Record | Wire layout and retained data |
| --- | --- |
| `lumps.Overlay` | 352 bytes:64 face slots, U/V bounds, four UV vectors, origin/normal; packed lower14-bit count and upper2-bit order |
| `lumps.Cubemap` | 16 bytes: signed32 integer origin, original size code, three padding bytes |
| `lumps.Primitive` | ten bytes: type, byte1 padding, four uint16 index/vertex span fields |

Overlay `SetFaceCount`/`SetRenderOrder` reject >64/>3 without changing other bits.
Only active face slots and valid material texinfos are checked; inactive slots
remain unchanged. Primitive types0/1 are verified; unknown types remain raw.
Positive vertex-count primitives contain **absolute** ID38 indices within their
declared span. Zero vertex-count primitives instead use owning face-winding
indices; genuine maps may have an empty ID38. No automatic HDR migration occurs.

Every retained fixed-array codec rejects partial records before allocation;
untouched malformed payloads remain raw/exportable, not silently shortened.
Explicit widths and genuine consumer evidence are in
[task09](09-overlay-cubemap-and-primitive-layouts.md).


## Packed leaf and face fields

Leaf `Area()` reads the low nine bits (`0x01ff`); `Flags()` reads the high seven.
`SetArea(uint16)` and `SetFlags(uint16)` mask oversized inputs and preserve the
other subfield. Both 56-byte leaf0 and 32-byte leaf1 retain the exact wire word.

Face `PrimitiveCount()` masks `NumPrims` with `0x7fff`. The high bit **disables**
dynamic shadows: `DynamicShadowsEnabled()` tests it for zero, and
`SetDynamicShadowsEnabled(true)` clears it without changing the count.
`SetPrimitiveCount(uint16)` masks inputs and preserves shadow state.
Use `Faces`/`SetFaces`, `OriginalFaces`/`SetOriginalFaces` or `FacesHDR`/`SetFacesHDR`
to commit owned snapshots; all use the explicit 56-byte revision0/1 codec.
Primitive consumers use the logical count, including packed `0x8001` = one.

## Displacements and defined supplementary data

`Displacements`/`SetDisplacements` expose the complete176-byte record: all edge
and corner neighbors, explicit padding and ten allowed-vertex words. Powers2–4,
active neighbor enums and sibling face/vertex/triangle spans are checked before
commit. Inactive sentinel-neighbor bytes are retained, not normalized. Setters
never regenerate topology or remap indices.

`Flags` is the full marked byte24 word, **including** `0x80000000`; otherwise that
word is `MinTess`. This convention is compiler-proven in both SDK and Reactive
Drop, not a guessed mask in `Contents` or `Power`. `MinTess` and marked flags
cannot be supplied simultaneously.

`DisplacementSamples(index)` reads ID34 using the owning face's
`(lightmapWidth+1)*(lightmapHeight+1)` samples, with no styles multiplier.
Each index byte is followed by three barycentric bytes; marker255 adds one index
byte (`255+extra`, maximum510). `SetDisplacementSamples` preserves gaps/trailers
and rebases later starts atomically when encoded width changes. Same-width edits
leave ID26 untouched. Missing streams return `ErrLumpNotFound`; populated shared
starts require coordinated raw replacement, while zero-sample owners change no
bytes.

`WaterOverlays`, `OverlaySystemLevels` and `DisplacementMultiblends` have attached
atomic setters. Explicit level counts match overlays; multiblend records match
vertices of flagged displacements. Coordinated topology/channel cardinality
changes require explicit raw replacement rather than automatic repair.

Reactive Drop's `PropCollisions`, `PropHulls`, `PropHullVertices` and
`PropHullTriangles` expose verified wire fields; `PropBlob` is a byte vector with
an opaque interior. No public consumer establishes dictionary linkage or
triangle-index units, so field names never justify guessed blob/array bounds.
ID32 alpha data and proprietary physics remain opaque. Profile admission and
the source/authored-wire versus genuine-fixture distinction are recorded in
[task10](10-displacement-and-lump-coverage.md).

## Static-prop editing

`StaticProps()` returns fresh `lumps.StaticProps` arrays; `SetStaticProps(value)`
commits them into the existing unique sprp entry. Counts derive from slices.
Models/leaf spans and actual File leaf references are checked before commit.
Unknown versions and ambiguous sprp IDs return errors; raw preservation remains
available. Failed edits do not change bytes or metadata.

| Profile | Typed static-prop layouts |
| --- | --- |
| SDK2013 | common4–6; distinct SDK10 (72 bytes, uint32 primary flags, lightmap X/Y) |
| L4D2 | common4–9 |
| CS:GO | common4–11; v10 extended uint32 at72, v11 scale float at76 |
| Reactive Drop | common4–9; genuine publisher-distributed v9 independently consumed |

Records expose RGBA bytes, complete flags, explicit padding and optional fields.
`Layout` describes the selected dialect, not a request to convert it. Fields
absent from that dialect must remain zero; unsupported edits fail rather than
disappear. Original128-byte model slots retain bytes after the first NUL;
`Models[i].Name()` reads the name, and `SetName(name)` rejects embedded NUL or
more than127 encoded bytes, then deliberately clears the new slot's tail.

```go
props, err := file.StaticProps()
if err != nil {
	return err
}
if len(props.Props) == 0 {
	return nil
}
props.Props[0].Origin[0] += 16
if err := file.SetStaticProps(props); err != nil {
	return err
}
```

Modified compressed sprp emits uncompressed bytes, preserving the descriptor
version, other flags, game gaps and opaque siblings. Genuine v11 baseline scales
and independently decoded multi-instance edits are recorded in
[task08](08-static-prop-codecs-and-editing.md); no unverified v11 branch is inferred.

## Numeric lump support inventory

`Bsp.Support(id)` reports capabilities for the admitted file profile/BSP version and actual lump descriptor revision; it is not payload validation or proof of presence. SDK means SDK2013 BSP19/20; C means CS:GO BSP21; L means L4D2 BSP21; R means Reactive Drop BSP21. Map revision does not select a codec. Every valid numeric ID has independent `RawLump`/`SetRawLump` access and stored payload/metadata preservation, including unknown revisions and unsupported branch interpretations. Canonical export does not preserve physical ordering/gaps; game directory absolute offsets relocate.

Cells: `T` typed, `B` byte-stream, `O` opaque, `U` unsupported profile interpretation; following digits are accepted descriptor revisions. `r` means advertised getter-only; `rw` means attached persistent getter/setter. Unknown descriptor revisions become opaque without losing raw access. A zero width means variable/opaque/context-bypassed, never a guessed physics or alpha width. Getter-only legacy snapshots are not editable through mutation or codec `Marshall`. `RawReadable`, `RawEditable`, `Preserved` always remain true. `Readable`/`Editable` refer only to the decoded/attached contract listed below.

Evidence: [pinned SDK header](https://github.com/ValveSoftware/source-sdk-2013/blob/b8cfb12c0e083a2ef5b2f9f9b50f3902fa034474/src/public/bspfile.h), [pinned Reactive Drop header](https://github.com/ReactiveDrop/reactivedrop_public_src/blob/af3be02c1d7db2cd07c41105e7519cf8edd66ce3/src/public/bspfile.h), repository `README.md`, `internal/versions/{registry,versions}.go`, `lumps/fixed_arrays_test.go`, `internal/codec/{leaf,worldlight,occlusion,ambient,displacement,coverage_formats,game_directory,staticprops}.go`, and root attached APIs. Genuine Steam-distributed RD BSP21 maps establish profile admission, v1 worldlights/faces/leaves, sprp v9 and populated ID61/26/34. IDs22–25/49/50/63 were empty in the genuine RD fixtures: their populated layouts rely on source and authored wire evidence, not fabricated genuine coverage.

| ID | Name / branch aliases | SDK | C | L | R | Width bytes | Accessors | Practical contract / limitations |
| ---: | --- | --- | --- | --- | --- | --- | --- | --- |
| 0 | entities | B0rw | B0rw | B0rw | B0rw | 0 | `Entities`, `SetEntities` | Exact entity text; no key/value parser or reference rewriting. |
| 1 | planes | T0rw | T0rw | T0rw | T0rw | 20 | `Planes`, `SetPlanes` | Setter checks plane axis and count; no automatic geometry rebuild. |
| 2 | texdata | T0r | T0r | T0r | T0r | 32 | `Lump -> lumps.TexData.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 3 | vertexes | T0r | T0r | T0r | T0r | 12 | `Lump -> lumps.Vertex.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 4 | visibility | O | O | O | O | 0 | raw only | Legacy visibility decoder trusts signed cluster counts and offsets; bounded PVS/PAS contract pending task12. Use raw bytes, not the legacy decoder for untrusted data. |
| 5 | nodes | T0r | T0r | T0r | T0r | 32 | `Lump -> lumps.Node.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 6 | texinfo | T0r | T0r | T0r | T0r | 72 | `Lump -> lumps.TexInfo.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 7 | faces | T0,1rw | T0,1rw | T0,1rw | T0,1rw | 56 | `Faces`, `SetFaces` | Masked primitive spans and local winding references checked; high bit disables dynamic shadows. No topology regeneration. |
| 8 | lighting | T0,1r | T0,1r | T0,1r | T0,1r | 4 | `Lump -> lumps.Lighting.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 9 | occlusion | T1,2rw | T1,2rw | T1,2rw | T2rw | 0 | `Occlusion`, `SetOcclusion` | Counted occluders/polygons/indices, not a uniform array; plane/vertex references checked; v1 cannot encode area edits. Reactive Drop accepts only v2. |
| 10 | leafs | T0,1rw | T0,1rw | T0,1rw | T0,1rw | v0:56; v1:32 | `Leaves`, `SetLeaves` | Leaf descriptor selects inline lighting (v0) or compact layout (v1); no inferred conversion. |
| 11 | faceids | T0r | T0r | T0r | T0r | 2 | `Lump -> lumps.FaceId.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 12 | edges | T0r | T0r | T0r | T0r | 4 | `Lump -> lumps.Edge.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 13 | surfedges | T0r | T0r | T0r | T0r | 4 | `Lump -> lumps.Surfedge.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 14 | models | T0r | T0r | T0r | T0r | 48 | `Lump -> lumps.Model.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 15 | worldlights | T0rw | T0,1rw | T0rw | T0,1rw | v0:88; C/R v1:100 | `WorldLights`, `SetWorldLights` | v0=88 bytes; v1=100 only CS:GO/Reactive Drop. HDR is independent. |
| 16 | leaffaces | T0r | T0r | T0r | T0r | 2 | `Lump -> lumps.LeafFace.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 17 | leafbrushes | T0r | T0r | T0r | T0r | 2 | `Lump -> lumps.LeafBrush.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 18 | brushes | T0r | T0r | T0r | T0r | 12 | `Lump -> lumps.Brush.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 19 | brushsides | T0rw | T0rw | T0rw | T0rw | 8 | `BrushSides`, `SetBrushSides` | Eight bytes in all admitted profiles; branch selects bevel/thin fields. |
| 20 | areas | T0r | T0r | T0r | T0r | 8 | `Lump -> lumps.Area.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 21 | areaportals | T0r | T0r | T0r | T0r | 12 | `Lump -> lumps.AreaPortal.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 22 | portals / unused0, propcollision (R uses last alias) | U | U | U | T0rw | 8 | `PropCollisions`, `SetPropCollisions` | Reactive Drop only; negative/overflow scalar span checks, not guessed hull/dictionary resolution. Populated genuine evidence unavailable. |
| 23 | clusters / unused1, prophulls (R uses last alias) | U | U | U | T0rw | 16 | `PropHulls`, `SetPropHulls` | Reactive Drop only; negative/overflow scalar span checks, not guessed vertex/dictionary resolution; contents bits preserved. Populated genuine evidence unavailable. |
| 24 | portalverts / unused2, prophullverts (R uses last alias) | U | U | U | T0rw | 12 | `PropHullVertices`, `SetPropHullVertices` | Reactive Drop only; raw float bits preserved; no guessed collision dictionary association. |
| 25 | clusterportals / unused3, proptris (R uses last alias) | U | U | U | T0rw | 8 | `PropHullTriangles`, `SetPropHullTriangles` | Reactive Drop only; index units and cross-blob association are unproven; no semantic collision interpretation. |
| 26 | dispinfo | T0rw | T0rw | T0rw | T0rw | 176 | `Displacements`, `SetDisplacements` | Power and active neighbors checked; sibling face/vertex/triangle spans validated; no topology regeneration. |
| 27 | originalfaces | T0,1rw | T0,1rw | T0,1rw | T0,1rw | 56 | `OriginalFaces`, `SetOriginalFaces` | Independent original-face channel with exact packed words; no automatic split-face rebuild. |
| 28 | physdisp | O | O | O | O | 0 | raw only | Physics displacement interior is opaque; no proprietary physics codec. |
| 29 | physcollide | O | O | O | O | 0 | raw only | Physics collision interior is opaque; no proprietary physics codec. |
| 30 | vertnormals | T0r | T0r | T0r | T0r | 12 | `Lump -> lumps.VertNormal.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 31 | vertnormalindices | T0r | T0r | T0r | T0r | 2 | `Lump -> lumps.VertNormalIndice.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 32 | displightmapalphas | O | O | O | O | 0 | raw only | Stripped/opaque alpha data; no proven record width or typed codec. |
| 33 | dispverts | T0r | T0r | T0r | T0r | 20 | `Lump -> lumps.DispVert.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 34 | displightmapsamplepositions | B0rw | B0rw | B0rw | B0rw | 0 | `DisplacementSamples`, `SetDisplacementSamples` | Face-derived count and displacement start/span select variable-width sample stream; indexed getter/setter preserves unclaimed bytes and adjusts following starts. Missing face lighting is an error. |
| 35 | game | B0rw | B0rw | B0rw | B0rw | 0 | `GameLumps`, `GameLump`, `GameLumpByID`, `SetGameLumps`, `SetGameLump`, `StaticProps`, `SetStaticProps` | Counted 16-byte directory entries, not uniform game records; unknown children/gaps preserved. Static props require one unique sprp entry and explicit profile/revision selector; compressed outer directory unsupported. StaticPropStatus describes selection, not successful payload validation. |
| 36 | leafwaterdata | T0r | T0r | T0r | T0r | 12 | `Lump -> lumps.LeafWaterData.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 37 | primitives | T0rw | T0rw | T0rw | T0rw | 10 | `Primitives`, `SetPrimitives` | Primitive types 0 and 1 supported; sibling index/vertex spans checked; other types return unsupported. |
| 38 | primverts | T0r | T0r | T0r | T0r | 12 | `Lump -> lumps.PrimVert.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 39 | primindices | T0r | T0r | T0r | T0r | 2 | `Lump -> lumps.PrimIndice.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 40 | pakfile | B0(raw only) | B0(raw only) | B0(raw only) | B0(raw only) | 0 | raw only | Legacy ZIP reader exists via Lump(40). Empty ZIP handling, bounded resource reads, missing/ambiguous paths and reader closure are not repaired until task13; no attached safe resource API advertised. |
| 41 | clipportalverts | T0r | T0r | T0r | T0r | 12 | `Lump -> lumps.ClipPortalVerts.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 42 | cubemaps | T0rw | T0rw | T0rw | T0rw | 16 | `Cubemaps`, `SetCubemaps` | Integer origin, size code and padding preserved; no texture-resource synthesis. |
| 43 | texdatastringdata | B0r | B0r | B0r | B0r | 0 | `Lump -> lumps.TexDataStringData.GetData` | Exact NUL-separated string bytes; no table-reference-aware attached editor. |
| 44 | texdatastringtable | T0r | T0r | T0r | T0r | 4 | `Lump -> lumps.TexDataStringTable.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 45 | overlays | T0rw | T0rw | T0rw | T0rw | 352 | `Overlays`, `SetOverlays` | Active face/material references and packed count/order checked; inactive slots preserved. |
| 46 | leafmindisttowater | T0r | T0r | T0r | T0r | 2 | `Lump -> lumps.LeafMinDistToWater.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 47 | facemacrotextureinfo | T0r | T0r | T0r | T0r | 2 | `Lump -> lumps.FaceMacroTextureInfo.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 48 | disptris | T0r | T0r | T0r | T0r | 2 | `Lump -> lumps.DispTris.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 49 | physcollidesurface / propblob (R uses last alias) | U | U | U | B0rw | 0 | `PropBlob`, `SetPropBlob` | Reactive Drop defines a byte blob; interior/index units/dictionary association are opaque. Other profiles retain deprecated physcollidesurface alias without a codec. |
| 50 | wateroverlays | T0rw | T0rw | T0rw | T0rw | 1120 | `WaterOverlays`, `SetWaterOverlays` | 1120-byte water overlays; face/material references checked. RD evidence is source/authored wire, not a populated genuine lump. |
| 51 | leafambientindexhdr / lightmappages | T0rw | T0rw | T0rw | T0rw | indexed:4; bypass:0 | `AmbientIndicesHDR`, `SetAmbientIndicesHDR` | Leaf-v0 bypasses standalone lighting; absent LDR index selects legacy mode for both channels. Old lightmappages alias is not a separate codec. |
| 52 | leafambientindex / lightmappagesinfo | T0rw | T0rw | T0rw | T0rw | indexed:4; bypass:0 | `AmbientIndices`, `SetAmbientIndices` | Leaf-v0 bypasses standalone lighting; absent LDR index selects legacy mode for both channels. Old lightmappagesinfo alias is not a separate codec. |
| 53 | lightinghdr | T0,1r | T0,1r | T0,1r | T0,1r | 4 | `Lump -> lumps.Lighting.GetData` | Legacy Lump snapshot/GetData only; no attached persistent setter or automatic sibling-reference repair. |
| 54 | worldlightshdr | T0rw | T0,1rw | T0rw | T0,1rw | v0:88; C/R v1:100 | `WorldLightsHDR`, `SetWorldLightsHDR` | v0=88 bytes; v1=100 only CS:GO/Reactive Drop. LDR is independent. |
| 55 | leafambientlightinghdr | T0,1rw | T0,1rw | T0,1rw | T0,1rw | inline:0; legacy:24; indexed:28 | `AmbientSamplesHDR`, `SetAmbientSamplesHDR` | Leaf-v0 supplies inline lighting; legacy mode uses 24-byte cubes, indexed mode 28-byte samples; no implicit HDR fallback. |
| 56 | leafambientlighting | T0,1rw | T0,1rw | T0,1rw | T0,1rw | inline:0; legacy:24; indexed:28 | `AmbientSamples`, `SetAmbientSamples` | Leaf-v0 supplies inline lighting; legacy mode uses 24-byte cubes, indexed mode 28-byte samples; no implicit HDR fallback. |
| 57 | xzippakfile | O | O | O | O | 0 | raw only | Deprecated console XZIP payload; PC profiles preserve raw bytes only. |
| 58 | faceshdr | T0,1rw | T0,1rw | T0,1rw | T0,1rw | 56 | `FacesHDR`, `SetFacesHDR` | Independent HDR channel; no fallback to ordinary faces. |
| 59 | mapflags | T0r | T0r | T0r | T0r | 4 | `Lump -> lumps.MapFlags.GetData` | Legacy getter accepts absent data as zero and does not enforce exactly four bytes; absence/presence contract pending task12. No persistent typed setter. |
| 60 | overlayfades | T0r | T0r | T0r | T0r | 8 | `Lump -> lumps.OverlayFade.GetData` | Getter-only fade distances; no attached setter or overlay-reference-aware mutation. |
| 61 | overlaysystemlevels | U | T0rw | T0rw | T0rw | 4 | `OverlaySystemLevels`, `SetOverlaySystemLevels` | Only L4D2/CS:GO/Reactive Drop; four CPU/GPU level bytes preserved without guessing engine policy. |
| 62 | physlevel | O | O | O | O | 0 | raw only | Physics level interior opaque; no proprietary physics codec. |
| 63 | dispmultiblend | U | T0rw | U | T0rw | 80 | `DisplacementMultiblends`, `SetDisplacementMultiblends` | Only CS:GO/Reactive Drop; 80-byte weights/alpha/colors, reference spans checked; no invented SDK/L4D2 support. RD populated genuine evidence unavailable. |

Ambient metadata is conservative for unknown revisions even where the actual consumer bypasses an unused descriptor. Leaf-v0 supplies inline lighting; absent LDR index selects legacy cubes for both channels; neither channel implies HDR fallback.

ID35 is a counted directory, not a uniformly typed game-record array. `StaticPropVersions` comes from the shared static-prop profile selector; `StaticPropStatus` is typed only when the actual directory has one unique `sprp` child with an accepted revision (RD v9 is accepted). Missing/duplicate/unknown-revision children remain opaque while directory and other children remain accessible. This status is selector coverage, not successful validation of the child payload. The selector never infers versions from byte length.

ID4 intentionally advertises no safe decoded consumer until task12 repairs visibility count/offset/RLE bounds. ID40 advertises no attached safe resource API until task13 repairs ZIP empty/presence/path/read-limit/closure contracts. ID59 remains a legacy getter-only record whose empty-as-zero semantics do not distinguish absence and whose exact payload size is not enforced; task12 owns that repair. These limitations are not advertised as fully supported editing.

## Real World examples
* Proof of concept BSP viewer: [https://github.com/Galaco/kero](https://github.com/Galaco/kero)
* Insert game_text newline placeholder characters (avoids Hammer crash) as a compile step: [https://github.com/Galaco/CS-GO-game_text-newline-inserter/tree/golang](https://github.com/Galaco/CS-GO-game_text-newline-inserter/tree/golang)
* Bspzip filelist generator from a mountable resource directory: [https://github.com/Galaco/bspzip-traverser](https://github.com/Galaco/bspzip-traverser)


# Contributing
All contributions welcome. Known unsupported games/maps are especially useful.
