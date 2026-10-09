package bsp

import (
	"github.com/galaco/bsp/internal/codec"
)

// LumpSupportStatus separates decoded records, byte contracts and opaque data.
type LumpSupportStatus string

const (
	SupportTyped              LumpSupportStatus = "typed"
	SupportByteStream         LumpSupportStatus = "byte-stream"
	SupportOpaque             LumpSupportStatus = "opaque"
	SupportUnsupportedProfile LumpSupportStatus = "unsupported-profile"
)

// LumpSupport describes capabilities, not payload validity or presence. Readable
// and Editable refer to the advertised decoded/attached contract; RawReadable,
// RawEditable and Preserved always hold for valid numeric IDs. Preserved means
// stored payload/metadata retention, not preservation of physical file layout.
// Versions are accepted lump descriptor revisions for this profile, not BSP
// versions. RecordWidth is zero for variable, bypassed or opaque representations.
// Returned slices are independent; callers cannot modify the static inventory.
type LumpSupport struct {
	ID             LumpID
	Name           string
	Aliases        []string
	Versions       []int32
	Representation string
	RecordWidth    int
	Status         LumpSupportStatus
	Readable       bool
	Editable       bool
	RawReadable    bool
	RawEditable    bool
	Preserved      bool
	Accessors      []string
	Evidence       string
	Limitations    string
	// The directory is supported independently of its sprp child. Status here
	// reflects the actual unique child's selector, never a payload-size guess.
	StaticPropStatus   LumpSupportStatus
	StaticPropVersions []int32
}

const supportSDKSource = "https://github.com/ValveSoftware/source-sdk-2013/blob/b8cfb12c0e083a2ef5b2f9f9b50f3902fa034474/src/public/bspfile.h"
const supportRDSource = "https://github.com/ReactiveDrop/reactivedrop_public_src/blob/af3be02c1d7db2cd07c41105e7519cf8edd66ce3/src/public/bspfile.h"

type lumpSupportEntry struct {
	name                  string
	aliases               []string
	versions              []int32
	representation        string
	width                 int
	status                LumpSupportStatus
	readable, editable    bool
	accessors             []string
	evidence, limitations string
}

// Numeric positions are the wire IDs; branch aliases never allocate new IDs.
var lumpSupportEntries = [64]lumpSupportEntry{
	0:  {name: "entities", aliases: []string{}, versions: []int32{0}, representation: "entity text bytes", width: 0, status: SupportByteStream, readable: true, editable: true, accessors: []string{"Entities", "SetEntities"}, evidence: supportSDKSource + "; README.md; internal/codec/registry.go", limitations: "Exact entity text; no key/value parser or reference rewriting."},
	1:  {name: "planes", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 20, status: SupportTyped, readable: true, editable: true, accessors: []string{"Planes", "SetPlanes"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Setter checks plane axis and count; no automatic geometry rebuild."},
	2:  {name: "texdata", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 32, status: SupportTyped, readable: true, editable: true, accessors: []string{"TextureData", "SetTextureData"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	3:  {name: "vertexes", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 12, status: SupportTyped, readable: true, editable: true, accessors: []string{"Vertices", "SetVertices"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	4:  {name: "visibility", aliases: []string{}, versions: []int32{0}, representation: "counted offset pairs and RLE bytes", width: 0, status: SupportTyped, readable: true, editable: true, accessors: []string{"Visibility", "SetVisibility"}, evidence: supportSDKSource + "; README.md; internal/codec/visibility.go", limitations: "At most65536 clusters; offsets are lump-relative, shared rows allowed; -1 means missing. Bounded PVS/PAS rows, no geometry-based recomputation."},
	5:  {name: "nodes", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 32, status: SupportTyped, readable: true, editable: true, accessors: []string{"Nodes", "SetNodes"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	6:  {name: "texinfo", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 72, status: SupportTyped, readable: true, editable: true, accessors: []string{"TextureInfo", "SetTextureInfo"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	7:  {name: "faces", aliases: []string{}, versions: []int32{0, 1}, representation: "fixed records", width: 56, status: SupportTyped, readable: true, editable: true, accessors: []string{"Faces", "SetFaces"}, evidence: supportSDKSource + "; README.md; internal/codec/face_wire_test.go", limitations: "Low15-bit primitive count excludes the high dynamic-shadow disable bit; primitive ranges and local winding references checked. No topology regeneration."},
	8:  {name: "lighting", aliases: []string{}, versions: []int32{0, 1}, representation: "fixed records", width: 4, status: SupportTyped, readable: true, editable: true, accessors: []string{"Lighting", "SetLighting"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	9:  {name: "occlusion", aliases: []string{}, versions: []int32{1, 2}, representation: "counted occlusion arrays", width: 0, status: SupportTyped, readable: true, editable: true, accessors: []string{"Occlusion", "SetOcclusion"}, evidence: supportSDKSource + "; README.md; internal/codec/occlusion.go", limitations: "Counted occluders/polygons/indices, not a uniform array; plane/vertex references checked; v1 cannot encode area edits. Reactive Drop accepts only v2."},
	10: {name: "leafs", aliases: []string{}, versions: []int32{0, 1}, representation: "descriptor-selected fixed leaf records", width: 0, status: SupportTyped, readable: true, editable: true, accessors: []string{"Leaves", "SetLeaves"}, evidence: supportSDKSource + "; README.md; internal/codec/leaf.go", limitations: "Leaf descriptor selects inline lighting (v0) or compact layout (v1); no inferred conversion."},
	11: {name: "faceids", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 2, status: SupportTyped, readable: true, editable: true, accessors: []string{"FaceIDs", "SetFaceIDs"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	12: {name: "edges", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 4, status: SupportTyped, readable: true, editable: true, accessors: []string{"Edges", "SetEdges"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	13: {name: "surfedges", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 4, status: SupportTyped, readable: true, editable: true, accessors: []string{"SurfaceEdges", "SetSurfaceEdges"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	14: {name: "models", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 48, status: SupportTyped, readable: true, editable: true, accessors: []string{"Models", "SetModels"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	15: {name: "worldlights", aliases: []string{}, versions: []int32{0, 1}, representation: "fixed records", width: 88, status: SupportTyped, readable: true, editable: true, accessors: []string{"WorldLights", "SetWorldLights"}, evidence: supportSDKSource + "; " + supportRDSource + "; README.md; internal/codec/worldlight.go", limitations: "v0=88 bytes; v1=100 only CS:GO/Reactive Drop. HDR is independent."},
	16: {name: "leaffaces", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 2, status: SupportTyped, readable: true, editable: true, accessors: []string{"LeafFaces", "SetLeafFaces"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	17: {name: "leafbrushes", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 2, status: SupportTyped, readable: true, editable: true, accessors: []string{"LeafBrushes", "SetLeafBrushes"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	18: {name: "brushes", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 12, status: SupportTyped, readable: true, editable: true, accessors: []string{"Brushes", "SetBrushes"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	19: {name: "brushsides", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 8, status: SupportTyped, readable: true, editable: true, accessors: []string{"BrushSides", "SetBrushSides"}, evidence: supportSDKSource + "; README.md; internal/codec/registry.go", limitations: "Eight bytes in all admitted profiles; branch selects bevel/thin fields."},
	20: {name: "areas", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 8, status: SupportTyped, readable: true, editable: true, accessors: []string{"Areas", "SetAreas"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	21: {name: "areaportals", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 12, status: SupportTyped, readable: true, editable: true, accessors: []string{"AreaPortals", "SetAreaPortals"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	22: {name: "portals", aliases: []string{"unused0", "propcollision"}, versions: []int32{0}, representation: "fixed records", width: 8, status: SupportTyped, readable: true, editable: true, accessors: []string{"PropCollisions", "SetPropCollisions"}, evidence: supportRDSource + "; README.md; internal/codec/coverage_formats.go", limitations: "Reactive Drop only; negative/overflow scalar span checks, not guessed hull/dictionary resolution. Populated genuine evidence unavailable."},
	23: {name: "clusters", aliases: []string{"unused1", "prophulls"}, versions: []int32{0}, representation: "fixed records", width: 16, status: SupportTyped, readable: true, editable: true, accessors: []string{"PropHulls", "SetPropHulls"}, evidence: supportRDSource + "; README.md; internal/codec/coverage_formats.go", limitations: "Reactive Drop only; negative/overflow scalar span checks, not guessed vertex/dictionary resolution; contents bits preserved. Populated genuine evidence unavailable."},
	24: {name: "portalverts", aliases: []string{"unused2", "prophullverts"}, versions: []int32{0}, representation: "fixed records", width: 12, status: SupportTyped, readable: true, editable: true, accessors: []string{"PropHullVertices", "SetPropHullVertices"}, evidence: supportRDSource + "; README.md; internal/codec/coverage_formats.go", limitations: "Reactive Drop only; raw float bits preserved; no guessed collision dictionary association."},
	25: {name: "clusterportals", aliases: []string{"unused3", "proptris"}, versions: []int32{0}, representation: "fixed records", width: 8, status: SupportTyped, readable: true, editable: true, accessors: []string{"PropHullTriangles", "SetPropHullTriangles"}, evidence: supportRDSource + "; README.md; internal/codec/coverage_formats.go", limitations: "Reactive Drop only; index units and cross-blob association are unproven; no semantic collision interpretation."},
	26: {name: "dispinfo", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 176, status: SupportTyped, readable: true, editable: true, accessors: []string{"Displacements", "SetDisplacements"}, evidence: supportSDKSource + "; README.md; internal/codec/displacement.go", limitations: "Power and active neighbors checked; sibling face/vertex/triangle spans validated; no topology regeneration."},
	27: {name: "originalfaces", aliases: []string{}, versions: []int32{0, 1}, representation: "fixed records", width: 56, status: SupportTyped, readable: true, editable: true, accessors: []string{"OriginalFaces", "SetOriginalFaces"}, evidence: supportSDKSource + "; README.md; internal/codec/face_wire_test.go", limitations: "Independent original-face channel with exact packed words; no automatic split-face/topology rebuild."},
	28: {name: "physdisp", aliases: []string{}, versions: []int32{0}, representation: "opaque bytes", width: 0, status: SupportOpaque, readable: false, editable: false, accessors: []string{}, evidence: supportSDKSource + "; README.md; internal/codec/registry.go", limitations: "Physics displacement interior is opaque; no proprietary physics codec."},
	29: {name: "physcollide", aliases: []string{}, versions: []int32{0}, representation: "opaque bytes", width: 0, status: SupportOpaque, readable: false, editable: false, accessors: []string{}, evidence: supportSDKSource + "; README.md; internal/codec/registry.go", limitations: "Physics collision interior is opaque; no proprietary physics codec."},
	30: {name: "vertnormals", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 12, status: SupportTyped, readable: true, editable: true, accessors: []string{"VertexNormals", "SetVertexNormals"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	31: {name: "vertnormalindices", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 2, status: SupportTyped, readable: true, editable: true, accessors: []string{"VertexNormalIndices", "SetVertexNormalIndices"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	32: {name: "displightmapalphas", aliases: []string{}, versions: []int32{0}, representation: "opaque bytes", width: 0, status: SupportOpaque, readable: false, editable: false, accessors: []string{}, evidence: supportSDKSource + "; README.md; internal/codec/registry.go", limitations: "Stripped/opaque alpha data; no proven record width or typed codec."},
	33: {name: "dispverts", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 20, status: SupportTyped, readable: true, editable: true, accessors: []string{"DisplacementVertices", "SetDisplacementVertices"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	34: {name: "displightmapsamplepositions", aliases: []string{}, versions: []int32{0}, representation: "face-derived sample byte stream", width: 0, status: SupportByteStream, readable: true, editable: true, accessors: []string{"DisplacementSamples", "SetDisplacementSamples"}, evidence: supportSDKSource + "; README.md; internal/codec/displacement.go", limitations: "Face-derived count and displacement start/span select variable-width sample stream; indexed getter/setter preserves unclaimed bytes and adjusts following starts. Missing face lighting is an error."},
	35: {name: "game", aliases: []string{}, versions: []int32{0}, representation: "counted game directory and opaque children", width: 0, status: SupportByteStream, readable: true, editable: true, accessors: []string{"GameLumps", "GameLump", "GameLumpByID", "SetGameLumps", "SetGameLump", "StaticProps", "SetStaticProps"}, evidence: supportSDKSource + "; " + supportRDSource + "; README.md#game-directory-ownership; internal/codec/game_directory.go; internal/codec/staticprops.go; staticprops.go", limitations: "Counted 16-byte directory entries, not uniform game records; unknown children/gaps preserved. Static props require one unique sprp entry and explicit profile/revision selector; compressed outer directory unsupported. StaticPropStatus describes selection, not successful payload validation."},
	36: {name: "leafwaterdata", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 12, status: SupportTyped, readable: true, editable: true, accessors: []string{"LeafWaterData", "SetLeafWaterData"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	37: {name: "primitives", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 10, status: SupportTyped, readable: true, editable: true, accessors: []string{"Primitives", "SetPrimitives"}, evidence: supportSDKSource + "; README.md; internal/codec/spatial.go", limitations: "Primitive types 0 and 1 supported; sibling index/vertex spans checked; other types return unsupported."},
	38: {name: "primverts", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 12, status: SupportTyped, readable: true, editable: true, accessors: []string{"PrimitiveVertices", "SetPrimitiveVertices"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	39: {name: "primindices", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 2, status: SupportTyped, readable: true, editable: true, accessors: []string{"PrimitiveIndices", "SetPrimitiveIndices"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	40: {name: "pakfile", aliases: []string{}, versions: []int32{0}, representation: "ZIP resource byte stream", width: 0, status: SupportByteStream, readable: true, editable: false, accessors: []string{"Pakfile", "Pakfile.Entries", "Pakfile.ReadFile"}, evidence: supportSDKSource + "; README.md; pakfile_integration_test.go", limitations: "Read-only ZIP snapshots; normalized first-match lookup, explicit not-found and entry closure. Whole ZIP changes use RawLump/SetRawLump; no extraction, per-entry editor or ZIP expansion cap."},
	41: {name: "clipportalverts", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 12, status: SupportTyped, readable: true, editable: true, accessors: []string{"ClipPortalVertices", "SetClipPortalVertices"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	42: {name: "cubemaps", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 16, status: SupportTyped, readable: true, editable: true, accessors: []string{"Cubemaps", "SetCubemaps"}, evidence: supportSDKSource + "; README.md; internal/codec/registry.go", limitations: "Integer origin, size code and padding preserved; no texture-resource synthesis."},
	43: {name: "texdatastringdata", aliases: []string{}, versions: []int32{0}, representation: "NUL-separated texture string bytes", width: 0, status: SupportByteStream, readable: true, editable: true, accessors: []string{"TextureStringData", "SetTextureStringData"}, evidence: supportSDKSource + "; README.md; internal/codec/registry.go", limitations: "Exact NUL-separated bytes; no automatic string-table rebasing."},
	44: {name: "texdatastringtable", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 4, status: SupportTyped, readable: true, editable: true, accessors: []string{"TextureStringTable", "SetTextureStringTable"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	45: {name: "overlays", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 352, status: SupportTyped, readable: true, editable: true, accessors: []string{"Overlays", "SetOverlays"}, evidence: supportSDKSource + "; README.md; internal/codec/registry.go", limitations: "Active face/material references and packed count/order checked; inactive slots preserved."},
	46: {name: "leafmindisttowater", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 2, status: SupportTyped, readable: true, editable: true, accessors: []string{"LeafMinDistancesToWater", "SetLeafMinDistancesToWater"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	47: {name: "facemacrotextureinfo", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 2, status: SupportTyped, readable: true, editable: true, accessors: []string{"FaceMacroTextureInfo", "SetFaceMacroTextureInfo"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	48: {name: "disptris", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 2, status: SupportTyped, readable: true, editable: true, accessors: []string{"DisplacementTriangles", "SetDisplacementTriangles"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	49: {name: "physcollidesurface", aliases: []string{"propblob"}, versions: []int32{0}, representation: "opaque interior byte blob", width: 0, status: SupportByteStream, readable: true, editable: true, accessors: []string{"PropBlob", "SetPropBlob"}, evidence: supportRDSource + "; README.md; internal/codec/coverage_formats.go", limitations: "Reactive Drop defines a byte blob; interior/index units/dictionary association are opaque. Other profiles retain deprecated physcollidesurface alias without a codec."},
	50: {name: "wateroverlays", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 1120, status: SupportTyped, readable: true, editable: true, accessors: []string{"WaterOverlays", "SetWaterOverlays"}, evidence: supportSDKSource + "; README.md; internal/codec/coverage_formats.go", limitations: "1120-byte water overlays; face/material references checked. RD evidence is source/authored wire, not a populated genuine lump."},
	51: {name: "leafambientindexhdr", aliases: []string{"lightmappages"}, versions: []int32{0}, representation: "fixed records", width: 4, status: SupportTyped, readable: true, editable: true, accessors: []string{"AmbientIndicesHDR", "SetAmbientIndicesHDR"}, evidence: supportSDKSource + "; README.md; internal/codec/ambient.go", limitations: "Leaf-v0 bypasses standalone lighting; absent LDR index selects legacy mode for both channels. Old lightmappages alias is not a separate codec."},
	52: {name: "leafambientindex", aliases: []string{"lightmappagesinfo"}, versions: []int32{0}, representation: "fixed records", width: 4, status: SupportTyped, readable: true, editable: true, accessors: []string{"AmbientIndices", "SetAmbientIndices"}, evidence: supportSDKSource + "; README.md; internal/codec/ambient.go", limitations: "Leaf-v0 bypasses standalone lighting; absent LDR index selects legacy mode for both channels. Old lightmappagesinfo alias is not a separate codec."},
	53: {name: "lightinghdr", aliases: []string{}, versions: []int32{0, 1}, representation: "fixed records", width: 4, status: SupportTyped, readable: true, editable: true, accessors: []string{"LightingHDR", "SetLightingHDR"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	54: {name: "worldlightshdr", aliases: []string{}, versions: []int32{0, 1}, representation: "fixed records", width: 88, status: SupportTyped, readable: true, editable: true, accessors: []string{"WorldLightsHDR", "SetWorldLightsHDR"}, evidence: supportSDKSource + "; " + supportRDSource + "; README.md; internal/codec/worldlight.go", limitations: "v0=88 bytes; v1=100 only CS:GO/Reactive Drop. LDR is independent."},
	55: {name: "leafambientlightinghdr", aliases: []string{}, versions: []int32{0, 1}, representation: "fixed records", width: 28, status: SupportTyped, readable: true, editable: true, accessors: []string{"AmbientSamplesHDR", "SetAmbientSamplesHDR"}, evidence: supportSDKSource + "; README.md; internal/codec/ambient.go", limitations: "Leaf-v0 supplies inline lighting; legacy mode uses 24-byte cubes, indexed mode 28-byte samples; no implicit HDR fallback."},
	56: {name: "leafambientlighting", aliases: []string{}, versions: []int32{0, 1}, representation: "fixed records", width: 28, status: SupportTyped, readable: true, editable: true, accessors: []string{"AmbientSamples", "SetAmbientSamples"}, evidence: supportSDKSource + "; README.md; internal/codec/ambient.go", limitations: "Leaf-v0 supplies inline lighting; legacy mode uses 24-byte cubes, indexed mode 28-byte samples; no implicit HDR fallback."},
	57: {name: "xzippakfile", aliases: []string{}, versions: []int32{0}, representation: "opaque bytes", width: 0, status: SupportOpaque, readable: false, editable: false, accessors: []string{}, evidence: supportSDKSource + "; README.md; internal/codec/registry.go", limitations: "Deprecated console XZIP payload; PC profiles preserve raw bytes only."},
	58: {name: "faceshdr", aliases: []string{}, versions: []int32{0, 1}, representation: "fixed records", width: 56, status: SupportTyped, readable: true, editable: true, accessors: []string{"FacesHDR", "SetFacesHDR"}, evidence: supportSDKSource + "; README.md; internal/codec/face_wire_test.go", limitations: "Independent HDR channel with exact packed words; no fallback to ordinary faces."},
	59: {name: "mapflags", aliases: []string{}, versions: []int32{0}, representation: "optional uint32", width: 4, status: SupportTyped, readable: true, editable: true, accessors: []string{"MapFlags", "SetMapFlags", "RemoveMapFlags"}, evidence: supportSDKSource + "; README.md; internal/codec/mapflags.go", limitations: "Empty means absent; exactly four bytes means present, including zero. Presence is explicit and removal is deliberate."},
	60: {name: "overlayfades", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 8, status: SupportTyped, readable: true, editable: true, accessors: []string{"OverlayFades", "SetOverlayFades"}, evidence: supportSDKSource + "; README.md; internal/codec/fixed_arrays_test.go", limitations: "Owned records and explicit persistent setter; no automatic geometry or sibling-reference rebuilding."},
	61: {name: "overlaysystemlevels", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 4, status: SupportTyped, readable: true, editable: true, accessors: []string{"OverlaySystemLevels", "SetOverlaySystemLevels"}, evidence: supportRDSource + "; README.md; internal/codec/coverage_formats.go", limitations: "Only L4D2/CS:GO/Reactive Drop; four CPU/GPU level bytes preserved without guessing engine policy."},
	62: {name: "physlevel", aliases: []string{}, versions: []int32{0}, representation: "opaque bytes", width: 0, status: SupportOpaque, readable: false, editable: false, accessors: []string{}, evidence: supportSDKSource + "; README.md; internal/codec/registry.go", limitations: "Physics level interior opaque; no proprietary physics codec."},
	63: {name: "dispmultiblend", aliases: []string{}, versions: []int32{0}, representation: "fixed records", width: 80, status: SupportTyped, readable: true, editable: true, accessors: []string{"DisplacementMultiblends", "SetDisplacementMultiblends"}, evidence: supportRDSource + "; README.md; internal/codec/coverage_formats.go", limitations: "Only CS:GO/Reactive Drop; 80-byte weights/alpha/colors, reference spans checked; no invented SDK/L4D2 support. RD populated genuine evidence unavailable."},
}

// Support reports the current file's profile and descriptor-selected contract.
// It does not decode ordinary payloads: malformed data can still fail its getter.
// Game-directory inspection only selects an optional unique static-prop child.
func (bsp *File) Support(id LumpID) (LumpSupport, error) {
	if err := validateID(id); err != nil {
		return LumpSupport{}, err
	}
	entry := lumpSupportEntries[id]
	result := LumpSupport{
		ID: id, Name: entry.name,
		Aliases:        append([]string(nil), entry.aliases...),
		Representation: entry.representation, RecordWidth: entry.width,
		Status: entry.status, Readable: entry.readable, Editable: entry.editable,
		RawReadable: true, RawEditable: true, Preserved: true,
		Accessors: append([]string(nil), entry.accessors...),
		Evidence:  entry.evidence, Limitations: entry.limitations,
	}
	rd := bsp.profile == ProfileReactiveDrop
	// These historical names share numeric IDs but not wire representations.
	if rd {
		switch id {
		case 22, 23, 24, 25, 49:
			result.Name = entry.aliases[len(entry.aliases)-1]
			result.Aliases = append([]string{entry.name}, entry.aliases[:len(entry.aliases)-1]...)
		}
	}
	if _, err := profileLayout(bsp.profile, bsp.header.Version); err != nil {
		result.unsupported(SupportUnsupportedProfile, "File profile/BSP version is not admitted.")
		return result, nil
	}
	switch id {
	case 22, 23, 24, 25, 49:
		if !rd {
			result.unsupported(SupportUnsupportedProfile, "This profile has no verified codec for this numeric ID/alias.")
			return result, nil
		}
	case 61:
		if bsp.profile == ProfileSDK2013 {
			result.unsupported(SupportUnsupportedProfile, "Overlay system levels are not verified for SDK2013.")
			return result, nil
		}
	case 63:
		if bsp.profile != ProfileCSGO && !rd {
			result.unsupported(SupportUnsupportedProfile, "Displacement multiblends are not verified for this profile.")
			return result, nil
		}
	}
	context := bsp.codecContext(id)
	// Use the same authoritative selector as attached consumers. The finite
	// inventory intentionally makes no typed claim for unknown revisions, even
	// when an ambient consumer bypasses an unused descriptor.
	accepted := false
	for _, revision := range entry.versions {
		candidate := context
		candidate.LumpVersion = revision
		if codec.ValidateRevision(int(id), candidate) == nil {
			result.Versions = append(result.Versions, revision)
			accepted = accepted || revision == context.LumpVersion
		}
	}
	if !accepted || codec.ValidateRevision(int(id), context) != nil {
		result.unsupported(SupportOpaque, "Unknown or unsupported descriptor revision; stored bytes remain available.")
		return result, nil
	}
	if result.Status == SupportOpaque {
		result.Versions = nil
		return result, nil
	}
	switch id {
	case 10:
		width, err := codec.LeafRecordSize(context)
		if err != nil {
			result.unsupported(SupportOpaque, "Unsupported leaf context.")
		} else {
			result.RecordWidth = width
		}
	case 15, 54:
		width, err := codec.WorldLightRecordSize(context)
		if err != nil {
			result.unsupported(SupportOpaque, "Unsupported worldlight context.")
		} else {
			result.RecordWidth = width
		}
	case 51, 52, 55, 56:
		if context.LeafVersion != 0 && context.LeafVersion != 1 {
			result.unsupported(SupportOpaque, "Unsupported companion leaf revision.")
		} else if context.LeafVersion == 0 || context.LegacyAmbient {
			result.RecordWidth = 0
			if context.LeafVersion == 1 && (id == 55 || id == 56) {
				result.RecordWidth = 24
			}
		}
	case 35:
		result.StaticPropStatus = SupportOpaque
		for revision := int32(4); revision <= 11; revision++ {
			candidate := context
			candidate.LumpVersion = revision
			if _, err := codec.StaticPropLayoutFor(candidate); err == nil {
				result.StaticPropVersions = append(result.StaticPropVersions, revision)
			}
		}
		if _, _, _, err := bsp.staticPropEntry(); err == nil {
			result.StaticPropStatus = SupportTyped
		}
	}
	return result, nil
}

func (support *LumpSupport) unsupported(status LumpSupportStatus, reason string) {
	support.Status = status
	support.Readable = false
	support.Editable = false
	support.RecordWidth = 0
	support.Representation = "opaque stored bytes"
	support.Accessors = nil
	support.Limitations += " " + reason
	if status == SupportUnsupportedProfile {
		support.Versions = nil
	}
}
