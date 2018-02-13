# bsp
Go library for manipulating Source Engine .bsp map files

Note:
At present only reading and even then only certain lumps are support, the rest will be added, as well as write support
as the library progresses.


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