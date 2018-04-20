# Bsp
Go library for manipulating Source Engine .bsp map files.

### Features:
* Read support for (probably) all documented bsp formats
* Write support for (probably) all documented bsp formats.
* Freely modify and resize any Lump data.

##### Not all lumps are current supported, but can be freely read and modified, as they are treated as `[]byte`

The following lumps currently have a full implementation for v20 bsp's:

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
27: OriginalFaces
28: PhysDisp
30: VertNormals
33: DispVerts
34: DispLightmapSamplePosition
36: LeafWaterData
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
blocks of a specified .bsp to terminal.

```go
package main

import (
	"fmt"
	"github.com/galaco/bsp"
	"log"
	"os"
)


func main() {
	f,_ := os.Open("de_dust2.bsp")
	fi,_ := f.Stat()
	
	fileData := make([]byte, fi.Size())
	f.Read(fileData)
	f.Close()

	// Create a new bsp reader
	reader := bsp.NewReader()
	reader.SetBuffer(fileData)
	
	// Read buffer
	file := reader.Read()
    
	lump := file.GetLump(0)
	fmt.Println(lump.GetContents().GetData().(string))
}
```

## Real World examples
Replace game_text newline placeholder characters (avoids Hammer crash) as a compile step: [https://github.com/Galaco/CS-GO-game_text-newline-inserter/tree/golang](https://github.com/Galaco/CS-GO-game_text-newline-inserter/tree/golang)


# Contributing
If you want to contribute, feel free to fork and raise a Pull Request for new additions.