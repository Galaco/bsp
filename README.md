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
5: Nodes
6: Texinfo
7: Faces
8: Lighting
9: Occlusion
10: Leafs
12: Edges
13: Surfedges
14: Models
16: Leaffaces
18: Brushes
19: Brushsides
20: Areas
21: AreaPortals
40: Pakfile
42: Cubemaps
43: Texdatastringdata
44: Texdatastringtable
45: Overlays
46: OverlayFades
```

##### This library may reorganise lump order during the first export. This is intentional to handle lump resizing, but will change your checksum if you export without changes.

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

	fmt.Println(file.GetLump(0).GetContents().GetData().(string))
	fmt.Println(file.GetLump(43).GetContents().GetData().(string))
}
```

## Real World examples
Replace game_text newline placeholder characters (avoids Hammer crash) as a compile step: [https://github.com/Galaco/CS-GO-game_text-newline-inserter/tree/golang](https://github.com/Galaco/CS-GO-game_text-newline-inserter/tree/golang)


# Contributing
If you want to contribute, feel free to fork and raise a Pull Request for new additions.