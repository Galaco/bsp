package bsp

// LumpId is the lump reference used by Source
type LumpId int

// Lump Identifiers
const (
	// LumpEntities is Entity keyvalue data stored as string
	LumpEntities = LumpId(0)
	// LumpPlanes is bsp planes
	LumpPlanes = LumpId(1)
	// LumpTexData is texture data used by bsp faces
	LumpTexData = LumpId(2)
	// LumpVertexes is vertex data
	LumpVertexes = LumpId(3)
	// LumpVisibility is vvis calculated visibility pvs & pas information
	LumpVisibility = LumpId(4)
	// LumpNodes is bsp node tree entries
	LumpNodes = LumpId(5)
	// LumpTexInfo is face texture information
	LumpTexInfo = LumpId(6)
	// LumpFaces is  bsp faces
	LumpFaces = LumpId(7)
	// LumpLighting
	LumpLighting = LumpId(8)
	// LumpOcclusion
	LumpOcclusion = LumpId(9)
	// LumpLeafs
	LumpLeafs = LumpId(10)
	// LumpFaceIds is contents is normally stripped out by compiler
	LumpFaceIds = LumpId(11)
	// LumpEdges is face edges. v1->v2, vertex order may be reversed
	LumpEdges = LumpId(12)
	// LumpSurfEdges
	LumpSurfEdges = LumpId(13)
	// LumpModels is models are root bsp nodes. m[0] = worldspawn. m[0+n] are brush entity data
	LumpModels = LumpId(14)
	// LumpWorldLights
	LumpWorldLights = LumpId(15)
	// LumpLeafFaces is faces that separate leaves.
	LumpLeafFaces = LumpId(16)
	// LumpLeafBrushes is brushes that define a leaf volume
	LumpLeafBrushes = LumpId(17)
	// LumpBrushes
	LumpBrushes = LumpId(18)
	// LumpBrushSides
	LumpBrushSides = LumpId(19)
	// LumpAreas
	LumpAreas = LumpId(20)
	// LumpAreaPortals
	LumpAreaPortals = LumpId(21)
	// LumpPortals
	LumpPortals = LumpId(22)
	// LumpUnused0
	LumpUnused0 = LumpId(22)
	// LumpPropCollision
	LumpPropCollision = LumpId(22)
	// LumpClusters
	LumpClusters = LumpId(23)
	// LumpUnused1
	LumpUnused1 = LumpId(23)
	// LumpPropHulls
	LumpPropHulls = LumpId(23)
	// LumpPortalVerts
	LumpPortalVerts = LumpId(24)
	// LumpUnused2
	LumpUnused2 = LumpId(24)
	// LumpPropHullVerts
	LumpPropHullVerts = LumpId(24)
	// LumpClusterPortals
	LumpClusterPortals = LumpId(25)
	// LumpUnused3
	LumpUnused3 = LumpId(25)
	// LumpPropTris
	LumpPropTris = LumpId(25)
	// LumpDispInfo
	LumpDispInfo = LumpId(26)
	// LumpOriginalFaces
	LumpOriginalFaces = LumpId(27)
	// LumpPhysDisp
	LumpPhysDisp = LumpId(28)
	// LumpPhysCollide
	LumpPhysCollide = LumpId(29)
	// LumpVertNormals
	LumpVertNormals = LumpId(30)
	// LumpVertNormalIndices
	LumpVertNormalIndices = LumpId(31)
	// LumpDispLightmapAlphas is contents is normally stripped out
	LumpDispLightmapAlphas = LumpId(32)
	// LumpDispVerts
	LumpDispVerts = LumpId(33)
	// LumpDispLightmapSamplePositions
	LumpDispLightmapSamplePositions = LumpId(34)
	// LumpGame is game specific data. includes staticprop data
	LumpGame = LumpId(35)
	// LumpLeafWaterData
	LumpLeafWaterData = LumpId(36)
	// LumpPrimitives
	LumpPrimitives = LumpId(37)
	// LumpPrimVerts
	LumpPrimVerts = LumpId(38)
	// LumpPrimIndices
	LumpPrimIndices = LumpId(39)
	// LumpPakfile is uncompressed zip of packed custom content
	LumpPakfile = LumpId(40)
	// LumpClipPortalVerts
	LumpClipPortalVerts = LumpId(41)
	// LumpCubemaps
	LumpCubemaps = LumpId(42)
	// LumpTexDataStringData is raw string data of material paths
	LumpTexDataStringData = LumpId(43)
	// LumpTexDataStringTable references entries in the string data lump
	LumpTexDataStringTable = LumpId(44)
	// LumpOverlays
	LumpOverlays = LumpId(45)
	// LumpLeafMinDistToWater
	LumpLeafMinDistToWater = LumpId(46)
	// LumpFaceMacroTextureInfo
	LumpFaceMacroTextureInfo = LumpId(47)
	// LumpDispTris
	LumpDispTris = LumpId(48)
	// LumpPhysCollideSurface is deprecated
	LumpPhysCollideSurface = LumpId(49)
	// LumpPropBlob
	LumpPropBlob = LumpId(49)
	// LumpWaterOverlays
	LumpWaterOverlays = LumpId(50)
	// LumpLightmapPages
	LumpLightmapPages = LumpId(51)
	// LumpLeafAmbientIndexHDR
	LumpLeafAmbientIndexHDR = LumpId(51)
	// LumpLightmapPagesInfo
	LumpLightmapPagesInfo = LumpId(52)
	// LumpLeafAmbientIndex
	LumpLeafAmbientIndex = LumpId(52)
	// LumpLightingHDR
	LumpLightingHDR = LumpId(53)
	// LumpWorldLightsHDR
	LumpWorldLightsHDR = LumpId(54)
	// LumpLeafAmbientLightingHDR
	LumpLeafAmbientLightingHDR = LumpId(55)
	// LumpLeafAmbientLighting
	LumpLeafAmbientLighting = LumpId(56)
	// LumpXZipPakfile is deprecated, and xbox specific
	LumpXZipPakfile = LumpId(57)
	// LumpFacesHDR
	LumpFacesHDR = LumpId(58)
	// LumpMapFlags
	LumpMapFlags = LumpId(59)
	// LumpOverlayFades
	LumpOverlayFades = LumpId(60)
	// LumpOverlaySystemLevels
	LumpOverlaySystemLevels = LumpId(61)
	// LumpPhysLevel
	LumpPhysLevel = LumpId(62)
	// LumpDispMultiBlend
	LumpDispMultiBlend = LumpId(63)
)
