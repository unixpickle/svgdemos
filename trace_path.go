package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"

	"github.com/unixpickle/gogui"
	"github.com/unixpickle/svgdemos/svg"
)

const (
	Size   = 400
	Border = 1
	PathStep = 0.01
)

func main() {
	gogui.RunOnMain(setupEverything)
	gogui.Main(&gogui.AppInfo{Name: "Path Tracer"})
}

func setupEverything() {
	segments, bounds, err := readPathAndBounds()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	w, _ := gogui.NewWindow(gogui.Rect{0, 0, Size, Size})
	c, _ := gogui.NewCanvas(gogui.Rect{0, 0, Size, Size})
	w.Add(c)
	w.SetTitle("Path Tracer")
	w.Center()
	w.Show()
	w.SetCloseHandler(func() {
		os.Exit(0)
	})

	c.SetDrawHandler(func(ctx gogui.DrawContext) {
		tracePath(ctx, segments, bounds)
	})
	c.NeedsUpdate()
}

func tracePath(ctx gogui.DrawContext, segments []svg.PathSegment, bounds svg.Rect) {
	var scale float64
	var translateX, translateY float64
	if bounds.Width() > bounds.Height() {
		scale = (Size - Border*2) / bounds.Width()
		translateX = Border
		translateY = Border + (bounds.Width()-bounds.Height())*scale/2
	} else {
		scale = (Size - Border*2) / bounds.Height()
		translateX = Border + (bounds.Height()-bounds.Width())*scale/2
		translateY = Border
	}
	translateX -= bounds.Min.X * scale
	translateY -= bounds.Min.Y * scale
	ctx.SetStroke(gogui.Color{0, 0, 0, 1})
	for _, segment := range segments {
		startPoint := segment.To()
		ctx.MoveTo(startPoint.X*scale+translateX, startPoint.Y*scale+translateY)
		for t := PathStep; t < 1; t += PathStep {
			point := segment.Evaluate(t)
			ctx.LineTo(point.X*scale+translateX, point.Y*scale+translateY)
		}
	}
	ctx.StrokePath()
}

func readPathAndBounds() ([]svg.PathSegment, svg.Rect, error) {
	fmt.Println("Enter path data, then deliver an EOF:")
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, svg.Rect{}, err
	}
	path, err := svg.ParsePath(string(data))
	if err != nil {
		return nil, svg.Rect{}, err
	}

	segments := path.Segments()

	if len(segments) == 0 {
		return segments, svg.Rect{}, nil
	}

	bounds := segments[0].Bounds()
	for _, segment := range segments {
		b := segment.Bounds()
		bounds.Min = svg.Point{math.Min(bounds.Min.X, b.Min.X), math.Min(bounds.Min.Y, b.Min.Y)}
		bounds.Max = svg.Point{math.Max(bounds.Max.X, b.Max.X), math.Max(bounds.Max.Y, b.Max.Y)}
	}
	return segments, bounds, nil
}
