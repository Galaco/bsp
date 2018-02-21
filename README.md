# bsp
Go library for manipulating Source Engine .bsp map files.

Note:
Read and write support is available, but only BSP version 20 is currently supported. Many lumps are not fully mapped, although support is increasing all the time. 


The following lumps currently have a full implementation:

```
0: Entdata
1: Planes
2: Texdata
3: Vertexes
5: Nodes
6: Texinfo
7: Faces
8: Lighting
10: Leafs
12: Edges
13: Surfedges
14: Models
16: Leaffaces
18: Brushes
19: Brushsides
20: Areas
21: AreaPortals
36: LeafWaterData
40: Pakfile
42: Cubemaps
43: Texdatastringdata
44: Texdatastringtable
45: Overlays
46: OverlayFades
```

The following lumps have some progress made towards them:

```
4: Visibility
9: Occlusion
15: WorldLights
35: Game
50: WaterOverlays
51: LightMapPages | Leaf Ambient Index HDR
52: LightMapPageInfos | Lead Ambient Index
``` 

All Lumps not listed above can be read and modified, but Lump data is treated as a `[]byte`.

# Usage

Minimal example of obtaining entdata and texdata from a BSP. The following will print both the entdata and texdata
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
	f,err := os.Open("de_dust2.bsp")
	if err!= nil {
		log.Fatal(err)
	}

	file := bsp.Parse(f)

	fmt.Println(file.Lumps[0].GetData().(string))
	fmt.Println(file.Lumps[43].GetData().(string))
}

```

## Real World usage examples
Replace game_text newline placeholder characters (avoids Hammer crash) as a compile step: [https://github.com/Galaco/CS-GO-game_text-newline-inserter/tree/golang](https://github.com/Galaco/CS-GO-game_text-newline-inserter/tree/golang)


# Contributing
If you want to contribute, feel free to fork and raise a Pull Request for new additions.