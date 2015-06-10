package svg

type Point struct {
	X float64
	Y float64
}

type Rect struct {
	Min Point
	Max Point
}

// Bounds returns a copy of a rectangle.
func (r Rect) Bounds() Rect {
	return r
}
