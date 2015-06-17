package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/unixpickle/svgdemos/svg"
)

func main() {
	fmt.Println("Enter path data, then deliver an EOF:")
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read data:", err)
		os.Exit(1)
	}
	path, err := svg.ParsePath(string(data))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse:", err)
		os.Exit(1)
	}
	var length float64
	for _, segment := range path.Segments() {
		length += segment.Length()
	}
	fmt.Println("length is approximately", length, "units")
}
