package svg

import "math"

type Point struct {
	X float64
	Y float64
}

func (p Point) approxEqual(p1 Point) bool {
	return Line{p, p1}.Length() < 0.00001
}

type Rect struct {
	Min Point
	Max Point
}

func (r Rect) approxEqual(r1 Rect) bool {
	return r.Min.approxEqual(r1.Min) && r.Max.approxEqual(r1.Max)
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

// Length returns the length of the line.
func (l Line) Length() float64 {
	return math.Sqrt(math.Pow(l.End.X-l.Start.X, 2) + math.Pow(l.End.Y-l.Start.Y, 2))
}

// Midpoint returns the midpoint of the line.
func (l Line) Midpoint() Point {
	return Point{l.Start.X + (l.End.X-l.Start.X)*0.5, l.Start.Y + (l.End.Y-l.Start.Y)*0.5}
}

// From returns the line's start point.
func (l Line) From() Point {
	return l.Start
}

// To returns the line's end point.
func (l Line) To() Point {
	return l.End
}
