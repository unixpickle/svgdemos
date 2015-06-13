package svg

type Arc struct {
	Start    Point
	End      Point
	XRadius  float64
	YRadius  float64
	Rotation float64
	LargeArc bool
	Sweep    bool
}

// Bounds returns the bounding box of the arc.
func (a *Arc) Bounds() Rect {
	panic("not yet implemented")
	// TODO: this
	return Rect{}
}

// Length computes the length of the arc.
func (a *Arc) Length() float64 {
	panic("not yet implemented")
	return 0
}

// From returns the arc's start point.
func (a *Arc) From() Point {
	return a.Start
}

// To returns the arc's end point.
func (a *Arc) To() Point {
	return a.End
}
