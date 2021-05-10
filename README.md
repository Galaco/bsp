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
	"github.com/galaco/bsp/lumps"
	"log"
	"os"
)

func main() {
	f,_ := os.Open("de_dust2.bsp")

	// Create a new bsp reader
	reader := bsp.NewReader(f)
	
	// Read buffer
	file,err := reader.Read()
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
    
	lump := file.Lump(bsp.LumpEntities).(*lumps.Entities)
	log.Println(lump.GetData())
}
```

## Real World examples
* Proof of concept BSP viewer: [https://github.com/Galaco/Lambda-Client](https://github.com/Galaco/Lambda-Client)
* Insert game_text newline placeholder characters (avoids Hammer crash) as a compile step: [https://github.com/Galaco/CS-GO-game_text-newline-inserter/tree/golang](https://github.com/Galaco/CS-GO-game_text-newline-inserter/tree/golang)
* Bspzip filelist generator from a mountable resource directory: [https://github.com/Galaco/bspzip-traverser](https://github.com/Galaco/bspzip-traverser)


# Contributing
All contributions welcome. Known unsupported games/maps are especially useful.
