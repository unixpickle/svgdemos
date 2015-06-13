package svg

import "math"

type Point struct {
	X float64
	Y float64
}

type Rect struct {
	Min Point
	Max Point
}

type Line struct {
	Start Point
	End   Point
}

// Bounds returns the bounding box of the line.
func (l Line) Bounds() Rect {
	minX := math.Min(l.Start.X, l.End.X)
	maxX := math.Max(l.Start.X, l.End.X)
	minY := math.Min(l.Start.Y, l.End.Y)
	maxY := math.Max(l.Start.Y, l.End.Y)
	return Rect{Point{minX, minY}, Point{maxX, maxY}}
}

// From returns the line's start point.
func (l Line) From() Point {
	return l.Start
}

// To returns the line's end point.
func (l Line) To() Point {
	return l.End
}
