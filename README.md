# bsp
Go library for manipulating Source Engine .bsp map files.

Note:
Read and write support is available, but only BSP version 20 is currently supported.
The following lumps currently have a full implementation:

```
0: Entdata
1: Planes
2: Texdata
3: Vertexes
4: Visibility
5: Nodes
6: texinfo
7: Faces
10: Leafs
12: Edges
13: Surfedges
14: Models
16: Leaffaces
18: Brushes
19: Brushsides
35: Game (partially implemented)
40: Pakfile
43: texdatastringdata
44: texdatastringtable
```

Lump 35 (Game) is not currently fully supported, however the Lump offset can be safely changed
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


# Contributing
If you know the underlying structure of a Lump that is currently implemented, please raise an issue or feel free to
implement it yourself. Adding a new Lump is as trivial as defining the struct, and referencing it at its index in
`lumps/lumps.go`.