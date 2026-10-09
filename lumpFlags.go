package bsp

// LumpID is the numeric lump reference used by Source.
type LumpID int

// Lump Identifiers
const (
	// LumpEntities is Entity keyvalue data stored as string
	LumpEntities = LumpID(0)
	// LumpPlanes is bsp planes
	LumpPlanes = LumpID(1)
	// LumpTexData is texture data used by bsp faces
	LumpTexData = LumpID(2)
	// LumpVertexes is vertex data
	LumpVertexes = LumpID(3)
	// LumpVisibility is vvis calculated visibility pvs & pas information
	LumpVisibility = LumpID(4)
	// LumpNodes is bsp node tree entries
	LumpNodes = LumpID(5)
	// LumpTexInfo is face texture information
	LumpTexInfo = LumpID(6)
	// LumpFaces is  bsp faces
	LumpFaces = LumpID(7)
	// LumpLighting
	LumpLighting = LumpID(8)
	// LumpOcclusion
	LumpOcclusion = LumpID(9)
	// LumpLeafs
	LumpLeafs = LumpID(10)
	// LumpFaceIDs is normally stripped out by the compiler.
	LumpFaceIDs = LumpID(11)
	// LumpEdges is face edges. v1->v2, vertex order may be reversed
	LumpEdges = LumpID(12)
	// LumpSurfEdges
	LumpSurfEdges = LumpID(13)
	// LumpModels is models are root bsp nodes. m[0] = worldspawn. m[0+n] are brush entity data
	LumpModels = LumpID(14)
	// LumpWorldLights
	LumpWorldLights = LumpID(15)
	// LumpLeafFaces is faces that separate leaves.
	LumpLeafFaces = LumpID(16)
	// LumpLeafBrushes is brushes that define a leaf volume
	LumpLeafBrushes = LumpID(17)
	// LumpBrushes
	LumpBrushes = LumpID(18)
	// LumpBrushSides
	LumpBrushSides = LumpID(19)
	// LumpAreas
	LumpAreas = LumpID(20)
	// LumpAreaPortals
	LumpAreaPortals = LumpID(21)
	// LumpPortals
	LumpPortals = LumpID(22)
	// LumpUnused0
	LumpUnused0 = LumpID(22)
	// LumpPropCollision
	LumpPropCollision = LumpID(22)
	// LumpClusters
	LumpClusters = LumpID(23)
	// LumpUnused1
	LumpUnused1 = LumpID(23)
	// LumpPropHulls
	LumpPropHulls = LumpID(23)
	// LumpPortalVerts
	LumpPortalVerts = LumpID(24)
	// LumpUnused2
	LumpUnused2 = LumpID(24)
	// LumpPropHullVerts
	LumpPropHullVerts = LumpID(24)
	// LumpClusterPortals
	LumpClusterPortals = LumpID(25)
	// LumpUnused3
	LumpUnused3 = LumpID(25)
	// LumpPropTris
	LumpPropTris = LumpID(25)
	// LumpDispInfo
	LumpDispInfo = LumpID(26)
	// LumpOriginalFaces
	LumpOriginalFaces = LumpID(27)
	// LumpPhysDisp
	LumpPhysDisp = LumpID(28)
	// LumpPhysCollide
	LumpPhysCollide = LumpID(29)
	// LumpVertNormals
	LumpVertNormals = LumpID(30)
	// LumpVertNormalIndices
	LumpVertNormalIndices = LumpID(31)
	// LumpDispLightmapAlphas is contents is normally stripped out
	LumpDispLightmapAlphas = LumpID(32)
	// LumpDispVerts
	LumpDispVerts = LumpID(33)
	// LumpDispLightmapSamplePositions
	LumpDispLightmapSamplePositions = LumpID(34)
	// LumpGame is game specific data. includes staticprop data
	LumpGame = LumpID(35)
	// LumpLeafWaterData
	LumpLeafWaterData = LumpID(36)
	// LumpPrimitives
	LumpPrimitives = LumpID(37)
	// LumpPrimVerts
	LumpPrimVerts = LumpID(38)
	// LumpPrimIndices
	LumpPrimIndices = LumpID(39)
	// LumpPakfile is uncompressed zip of packed custom content
	LumpPakfile = LumpID(40)
	// LumpClipPortalVerts
	LumpClipPortalVerts = LumpID(41)
	// LumpCubemaps
	LumpCubemaps = LumpID(42)
	// LumpTexDataStringData is raw string data of material paths
	LumpTexDataStringData = LumpID(43)
	// LumpTexDataStringTable references entries in the string data lump
	LumpTexDataStringTable = LumpID(44)
	// LumpOverlays
	LumpOverlays = LumpID(45)
	// LumpLeafMinDistToWater
	LumpLeafMinDistToWater = LumpID(46)
	// LumpFaceMacroTextureInfo
	LumpFaceMacroTextureInfo = LumpID(47)
	// LumpDispTris
	LumpDispTris = LumpID(48)
	// LumpPhysCollideSurface is deprecated
	LumpPhysCollideSurface = LumpID(49)
	// LumpPropBlob
	LumpPropBlob = LumpID(49)
	// LumpWaterOverlays
	LumpWaterOverlays = LumpID(50)
	// LumpLightmapPages
	LumpLightmapPages = LumpID(51)
	// LumpLeafAmbientIndexHDR
	LumpLeafAmbientIndexHDR = LumpID(51)
	// LumpLightmapPagesInfo
	LumpLightmapPagesInfo = LumpID(52)
	// LumpLeafAmbientIndex
	LumpLeafAmbientIndex = LumpID(52)
	// LumpLightingHDR
	LumpLightingHDR = LumpID(53)
	// LumpWorldLightsHDR
	LumpWorldLightsHDR = LumpID(54)
	// LumpLeafAmbientLightingHDR
	LumpLeafAmbientLightingHDR = LumpID(55)
	// LumpLeafAmbientLighting
	LumpLeafAmbientLighting = LumpID(56)
	// LumpXZipPakfile is deprecated, and xbox specific
	LumpXZipPakfile = LumpID(57)
	// LumpFacesHDR
	LumpFacesHDR = LumpID(58)
	// LumpMapFlags
	LumpMapFlags = LumpID(59)
	// LumpOverlayFades
	LumpOverlayFades = LumpID(60)
	// LumpOverlaySystemLevels
	LumpOverlaySystemLevels = LumpID(61)
	// LumpPhysLevel
	LumpPhysLevel = LumpID(62)
	// LumpDispMultiBlend
	LumpDispMultiBlend = LumpID(63)
)
