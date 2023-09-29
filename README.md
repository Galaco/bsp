[![GoDoc](https://godoc.org/github.com/Galaco/bsp?status.svg)](https://godoc.org/github.com/Galaco/bsp)
[![Go report card](https://goreportcard.com/badge/github.com/galaco/bsp)](https://goreportcard.com/badge/github.com/galaco/bsp)
[![GolangCI](https://golangci.com/badges/github.com/galaco/bsp.svg)](https://golangci.com)
[![CircleCI](https://circleci.com/gh/Galaco/bsp/tree/master.svg?style=svg)](https://circleci.com/gh/Galaco/bsp/tree/master)
[![codecov](https://codecov.io/gh/Galaco/bsp/branch/master/graph/badge.svg)](https://codecov.io/gh/Galaco/bsp)

# Bsp
The most comprehensive library for reading and writing Source Engine .bsp map files.

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

Lumps not listed here are parsed and available as `[]byte` format.

Note: Some lumps in some BSP versions have data with unidentified purpose. These fields are available as byte arrays. 
Please submit an issue or a PR if you can help fill in any of these fields.


# Usage

Minimal example of obtaining entdata from a BSP. The following will print the entdata
lump (entdata is a single json-like string) of a specified .bsp to console.

```go
package main

import (
	"github.com/galaco/bsp"
	"github.com/galaco/bsp/lump"
	"log"
	"os"
)

func main() {
	f, err := os.Open("de_dust2.bsp")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create a new bsp reader
	reader := bsp.NewReader()

	// Read buffer
	file, err := reader.Read(f)
	if err != nil {
		log.Fatal(err)
	}

	lump := file.Lump(bsp.LumpEntities).(*lumps.Entities)
	log.Println(lump.Contents())
}
```

There are more usage examples available in the `examples/` directory.

## Exporting BSPs

This library supports writing BSPs. It aims to preserve identical binaries where possible, but this is not guaranteed
due to wide-ranging difference in format across games (and even within the same game!). 
For example:
* Counterstrike: Source
  * de_dust2 Lump 59 (MapFlags) has 0 flags set, a the 4byte lump is written. Format is BSP v20.
  * de_nuke Lump 59 (MapFlags) has 0 flags set, but the lump is not written. Format is BSP v20.

There are plenty of other scenarios where this can occur, and in a way that we cannot guess with certainty what the 
expected behaviour should be. By default, this library assumes that structures that contain > 0 bytes are written, 
but this behaviour can be overridden (see examples).


## Real World examples
* Proof of concept BSP viewer: [https://github.com/Galaco/kero](https://github.com/Galaco/kero)
* Insert game_text newline placeholder characters (avoids Hammer crash) as a compile step: [https://github.com/Galaco/CS-GO-game_text-newline-inserter/tree/golang](https://github.com/Galaco/CS-GO-game_text-newline-inserter/tree/golang)
* Bspzip filelist generator from a mountable resource directory: [https://github.com/Galaco/bspzip-traverser](https://github.com/Galaco/bspzip-traverser)


# Contributing
All contributions welcome, in particular any maps that are found to be incompatible.
