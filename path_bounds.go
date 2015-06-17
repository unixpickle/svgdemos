package main

import (
	"fmt"
	"io/ioutil"
	"math"
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
	segments := path.Segments()
	if len(segments) == 0 {
		fmt.Println("path has no bounds.")
		return
	}
	bounds := segments[0].Bounds()
	for _, segment := range segments {
		b := segment.Bounds()
		bounds.Min = svg.Point{math.Min(bounds.Min.X, b.Min.X), math.Min(bounds.Min.Y, b.Min.Y)}
		bounds.Max = svg.Point{math.Max(bounds.Max.X, b.Max.X), math.Max(bounds.Max.Y, b.Max.Y)}
	}
	fmt.Println("x =", bounds.Min.X, "y =", bounds.Min.Y, "width =", bounds.Max.X-bounds.Min.X,
		"height =", bounds.Max.Y-bounds.Min.Y)
}
