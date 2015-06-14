package svg

import "math"

const arcLengthApproximationInterval = 0.005

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
	params := a.Params()
	return params.Length()
}

// Params uses a bunch of math to generate ArcParams.
func (a *Arc) Params() ArcParams {
	rx, ry := math.Abs(a.XRadius), math.Abs(a.YRadius)
	if rx == 0 || ry == 0 {
		return lineToArcParams(Line{a.Start, a.End})
	}

	// Math from http://www.w3.org/TR/SVG/implnote.html#ArcImplementationNotes
	x1, y1 := a.Start.X, a.Start.Y
	x2, y2 := a.End.X, a.End.Y
	sin := math.Sin(math.Pi / 180 * a.Rotation)
	cos := math.Cos(math.Pi / 180 * a.Rotation)
	x1p := cos*(x1-x2)/2 + sin*(y1-y2)/2
	y1p := -sin*(x1-x2)/2 + cos*(y1-y2)/2

	// Canonicalize the radii
	lambda := math.Pow(x1p/rx, 2) + math.Pow(y1p/ry, 2)
	if lambda > 1 {
		sqrtLambda := math.Sqrt(lambda)
		rx = sqrtLambda * rx
		ry = sqrtLambda * ry
	}

	sqrtMe := (math.Pow(rx*ry, 2) - math.Pow(rx*y1p, 2) -
		math.Pow(ry*x1p, 2)) / (math.Pow(rx*y1p, 2) + math.Pow(ry*x1p, 2))
	if sqrtMe < 0 {
		sqrtMe = 0
	}
	coefficient := math.Sqrt(sqrtMe)
	if a.LargeArc == a.Sweep {
		coefficient *= -1
	}
	cxp := coefficient * rx * y1p / ry
	cyp := -coefficient * ry * x1p / rx
	center := Point{cos*cxp - sin*cyp + (x1+x2)/2,
		sin*cxp + cos*cyp + (y1+y2)/2}

	start := 180 / math.Pi * math.Atan2(y1-center.Y, x1-center.X)
	end := 180 / math.Pi * math.Atan2(y2-center.Y, x2-center.X)
	start -= a.Rotation
	end -= a.Rotation
	if start == 0 {
		start += 360
	}
	if end == 0 {
		end += 360
	}

	return ArcParams{center, start, end, a.Rotation, rx, ry, a.Sweep}
}

// From returns the arc's start point.
func (a *Arc) From() Point {
	return a.Start
}

// To returns the arc's end point.
func (a *Arc) To() Point {
	return a.End
}

// ArcParams contains the information needed to generate an arc parametrically.
type ArcParams struct {
	Center     Point
	StartAngle float64
	EndAngle   float64
	Rotation   float64
	XRadius    float64
	YRadius    float64
	Sweep      bool
}

// Length approximates the arc's length.
func (a *ArcParams) Length() float64 {
	var length float64
	for t := 0.0; t < 1; t += arcLengthApproximationInterval {
		length += Line{a.Evaluate(t), a.Evaluate(t +
			arcLengthApproximationInterval)}.Length()
	}
	return length
}

// Evaluate generates a point on the arc for a parameter between 0 and 1.
func (a *ArcParams) Evaluate(t float64) Point {
	var angle float64
	if a.Sweep == (a.EndAngle > a.StartAngle) {
		angle = a.StartAngle + t*(a.EndAngle-a.StartAngle)
	} else if !a.Sweep {
		end := a.EndAngle - 360
		angle = a.StartAngle + t*(end-a.StartAngle)
		if angle < 0 {
			angle += 360
		}
	} else {
		end := a.EndAngle + 360
		angle = a.StartAngle + t*(end-a.StartAngle)
		if angle > 360 {
			angle -= 360
		}
	}
	angle *= math.Pi / 180
	rotCos := math.Cos(math.Pi / 180 * a.Rotation)
	rotSin := math.Sin(math.Pi / 180 * a.Rotation)
	return Point{a.XRadius*math.Cos(angle)*rotCos -
		a.YRadius*rotSin*math.Sin(angle) + a.Center.X,
		a.YRadius*math.Cos(angle)*rotSin + a.YRadius*math.Sin(angle)*rotCos +
			a.Center.Y}
}

func lineToArcParams(l Line) ArcParams {
	center := Point{l.Start.X + (l.End.X-l.Start.X)/2,
		l.Start.Y + (l.End.Y-l.Start.Y)/2}
	radius := l.Length() / 2
	rotation := 180 / math.Pi * math.Atan2(l.End.Y-l.Start.Y,
		l.End.X-l.Start.X)
	return ArcParams{center, 0, 180, rotation, radius, 0, true}
}
