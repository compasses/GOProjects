package main

import (
	"fmt"
	"flag"
)

var infile *string = flag.String("i", "infile", "Files contains value to sorting")
var outfile *string = flag.String("o", "outfile", "Files save sorting results")
var algorithems *string = flag.String("a", "qsort", "Sort algorithem")


func main() {
	flag.Parse()

	if infile != nil {
		fmt.Println("infile = ", *infile, "outfile = ", *outfile, "algorithem = ", *algorithems)
	}
}
