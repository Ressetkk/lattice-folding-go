package main

import (
	"flag"
	"fmt"
)
var (
	protein = flag.String("protein", "", "character stream over the alphabet of {h, p}")
)
func main() {
	flag.Parse()
	fmt.Print(*protein)

	for i, char := range *protein {
		fmt.Println(i, string(char))
	}
}
