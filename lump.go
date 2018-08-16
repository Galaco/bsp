package bsp

import (
	"github.com/galaco/bsp/lumps"
	"log"
	"github.com/galaco/bsp/versions"
)

const LUMP_ENTITIES = 0
const LUMP_PLANES = 1
const LUMP_TEXDATA = 2
const LUMP_VERTEXES = 3
const LUMP_VISIBILITY = 4
const LUMP_NODES = 5
const LUMP_TEXINFO = 6
const LUMP_FACES = 7
const LUMP_LIGHTING = 8
const LUMP_OCCLUSION = 9
const LUMP_LEAFS = 10
const LUMP_FACEIDS = 11 // contents is normally stripped out
const LUMP_EDGES = 12
const LUMP_SURFEDGES = 13
const LUMP_MODELS = 14
const LUMP_WORLDLIGHTS = 15
const LUMP_LEAFFACES = 16
const LUMP_LEAFBRUSHES = 17
const LUMP_BRUSHES = 18
const LUMP_BRUSHSIDES = 19
const LUMP_AREAS = 20
const LUMP_AREAPORTALS = 21
const LUMP_PORTALS = 22
const LUMP_UNUSED0 = 22
const LUMP_PROPCOLLISION = 22
const LUMP_CLUSTERS = 23
const LUMP_UNUSED1 = 23
const LUMP_PROPHULLS = 23
const LUMP_PORTALVERTS = 24
const LUMP_UNUSED2 = 24
const LUMP_PROPHULLVERTS = 24
const LUMP_CLUSTERPORTALS = 25
const LUMP_UNUSED3 = 25
const LUMP_PROPTRIS = 25
const LUMP_DISPINFO = 26
const LUMP_ORIGINALFACES = 27
const LUMP_PHYSDISP = 28
const LUMP_PHYSCOLLIDE = 29
const LUMP_VERTNORMALS = 30
const LUMP_VERTNORMALINDICES = 31
const LUMP_DISP_LIGHTMAP_ALPHAS = 32 // contents is normally stripped out
const LUMP_DISP_VERTS = 33
const LUMP_DISP_LIGHTMAP_SAMPLE_POSITIONS = 34
const LUMP_GAME_LUMP = 35
const LUMP_LEAFWATERDATA = 36
const LUMP_PRIMITIVES = 37
const LUMP_PRIMVERTS = 38
const LUMP_PRIMINDICES = 39
const LUMP_PAKFILE = 40
const LUMP_CLIPPORTALVERTS = 41
const LUMP_CUBEMAPS = 42
const LUMP_TEXDATA_STRING_DATA = 43
const LUMP_TEXDATA_STRING_TABLE = 44
const LUMP_OVERLAYS = 45
const LUMP_LEAFMINDISTTOWATER = 46
const LUMP_FACE_MACRO_TEXTURE_INFO = 47
const LUMP_DISP_TRIS = 48
const LUMP_PHYSCOLLIDESURFACE = 49 // deprecated
const LUMP_PROP_BLOB = 49
const LUMP_WATEROVERLAYS = 50
const LUMP_LIGHTMAPPAGES = 51
const LUMP_LEAF_AMBIENT_INDEX_HDR = 51
const LUMP_LIGHTMAPPAGESINFOS = 52
const LUMP_LEAF_AMBIENT_INDEX = 52
const LUMP_LIGHTING_HDR = 53
const LUMP_WORLDLIGHTS_HDR = 54
const LUMP_LEAF_AMBIENT_LIGHTING_HDR = 55
const LUMP_LEAF_AMBIENT_LIGHTING = 56
const LUMP_XZIPPAKFILE = 57 // deprecated, and xbox specific
const LUMP_FACES_HDR = 58
const LUMP_MAP_FLAGS = 59
const LUMP_OVERLAY_FADES = 60
const LUMP_OVERLAY_SYSTEM_LEVELS = 61
const LUMP_PHYSLEVEL = 62
const LUMP_DISP_MULTIBLEND = 63

// Container for a lump. Also includes metadata about the lump.
// N.B. Some information mirrors the header's lump descriptor, but header information should not be trusted after
// import completion.
type Lump struct {
	raw []byte
	data lumps.ILump
	length int32
	index int
	loaded bool
}

// Get lump identifier
// Id is the lump type index (not the index for the order the lumps are stored)
func (l *Lump) SetId(index int) {
	l.index = index
}

// Get the contents of a lump.
// NOTE: Will need to be cast to the relevant lumps
func (l *Lump) GetContents() lumps.ILump {
	if l.loaded == false {
		l.data = l.data.FromBytes(l.raw, int32(len(l.raw)))
		l.loaded = true
	}
	return l.data
}

// Set content type of a lump.
func (l *Lump) SetContents(data lumps.ILump) {
	l.data = data
	l.loaded = false
}

// Get the raw []byte contents of a lump.
// N.B. This is the raw imported value. To get the raw value of a modified lump, use GetContents().ToBytes()
func (l *Lump) GetRawContents() []byte {
	return l.raw
}

// Set raw []byte contents of a lump.
func (l *Lump) SetRawContents(raw []byte) {
	l.raw = raw
}

// Get length of a lump in bytes.
func (l *Lump) GetLength() int32 {
	return l.length
}

// Return an instance of a Lump for a given offset.
func getReferenceLumpByIndex(index int, version int32) (lumps.ILump,error) {
	if index < 0 || index > 63 {
		log.Fatalf("Invalid lump index: %d provided\n", index)
	}

	return versions.GetLumpForVersion(int(version), index)
}