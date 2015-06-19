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
	Size     = 400
	Border   = 10
	PathStep = 0.01
)

var Segments []svg.PathSegment
var Bounds svg.Rect
var Selected svg.PathSegment

func main() {
	gogui.RunOnMain(setupEverything)
	gogui.Main(&gogui.AppInfo{Name: "Path Tracer"})
}

func setupEverything() {
	if err := readPathAndBounds(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	w, _ := gogui.NewWindow(gogui.Rect{0, 0, Size, Size})
	c, _ := gogui.NewCanvas(gogui.Rect{0, 0, Size, Size})
	w.Add(c)
	w.SetTitle("Path Inspector")
	w.Center()
	w.Show()
	w.SetCloseHandler(func() {
		os.Exit(0)
	})
	w.SetMouseMoveHandler(func(e gogui.MouseEvent) {
		updateSelected(e)
		c.NeedsUpdate()
	})

	c.SetDrawHandler(tracePath)
	c.NeedsUpdate()
}

func updateSelected(e gogui.MouseEvent) {
	translateX, translateY, scale := transformation()
	mousePoint := svg.Point{e.X, e.Y}
	Selected = nil

SegmentLoop:
	for _, segment := range Segments {
		for t := 0.0; t <= 1; t += PathStep {
			point := segment.Evaluate(t)
			screenPoint := svg.Point{point.X*scale + translateX, point.Y*scale + translateY}
			if (svg.Line{mousePoint, screenPoint}).Length() < 5 {
				Selected = segment
				break SegmentLoop
			}
		}
	}
}

func tracePath(ctx gogui.DrawContext) {
	translateX, translateY, scale := transformation()
	ctx.SetStroke(gogui.Color{0, 0, 0, 1})
	ctx.SetThickness(1)
	for i, segment := range Segments {
		if segment == Selected {
			ctx.MoveTo(segment.To().X*scale+translateX, segment.To().Y*scale+translateY)
			continue
		}
		if i == 0 {
			startPoint := segment.From()
			ctx.MoveTo(startPoint.X*scale+translateX, startPoint.Y*scale+translateY)
		}
		for t := PathStep; t < 1; t += PathStep {
			point := segment.Evaluate(t)
			ctx.LineTo(point.X*scale+translateX, point.Y*scale+translateY)
		}
	}
	ctx.StrokePath()
	if Selected != nil {
		ctx.SetStroke(gogui.Color{0, 1, 0, 1})
		ctx.SetThickness(4)
		ctx.MoveTo(Selected.From().X*scale+translateX, Selected.From().Y*scale+translateY)
		for t := PathStep; t <= 1; t += PathStep {
			point := Selected.Evaluate(t)
			ctx.LineTo(point.X*scale+translateX, point.Y*scale+translateY)
		}
		ctx.StrokePath()
	}
}

func transformation() (translateX, translateY, scale float64) {
	if Bounds.Width() > Bounds.Height() {
		scale = (Size - Border*2) / Bounds.Width()
		translateX = Border
		translateY = Border + (Bounds.Width()-Bounds.Height())*scale/2
	} else {
		scale = (Size - Border*2) / Bounds.Height()
		translateX = Border + (Bounds.Height()-Bounds.Width())*scale/2
		translateY = Border
	}
	translateX -= Bounds.Min.X * scale
	translateY -= Bounds.Min.Y * scale
	return
}

func readPathAndBounds() error {
	fmt.Println("Enter path data, then deliver an EOF:")
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	path, err := svg.ParsePath(string(data))
	if err != nil {
		return err
	}

	Segments = path.Segments()

	if len(Segments) == 0 {
		return nil
	}

	Bounds = Segments[0].Bounds()
	for _, segment := range Segments {
		b := segment.Bounds()
		Bounds.Min = svg.Point{math.Min(Bounds.Min.X, b.Min.X), math.Min(Bounds.Min.Y, b.Min.Y)}
		Bounds.Max = svg.Point{math.Max(Bounds.Max.X, b.Max.X), math.Max(Bounds.Max.Y, b.Max.Y)}
	}
	return nil
}
