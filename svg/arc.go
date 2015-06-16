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
	if params, line := a.Params(); params != nil {
		return params.Bounds()
	} else {
		return line.Bounds()
	}
}

// Length computes the length of the arc.
func (a *Arc) Length() float64 {
	if params, line := a.Params(); params != nil {
		return params.Length()
	} else {
		return line.Length()
	}
}

// Params uses a bunch of math to generate ArcParams. In some cases, an arc is
// treated like a line. In these cases, the ArcParams will be nil and the Line
// will be non-nil.
func (a *Arc) Params() (*ArcParams, *Line) {
	rx, ry := math.Abs(a.XRadius), math.Abs(a.YRadius)
	if rx == 0 || ry == 0 {
		return nil, &Line{a.Start, a.End}
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

	sqrtMe := (math.Pow(rx*ry, 2) - math.Pow(rx*y1p, 2) - math.Pow(ry*x1p, 2)) /
		(math.Pow(rx*y1p, 2) + math.Pow(ry*x1p, 2))
	if sqrtMe < 0 {
		sqrtMe = 0
	}
	coefficient := math.Sqrt(sqrtMe)
	if a.LargeArc == a.Sweep {
		coefficient *= -1
	}
	cxp := coefficient * rx * y1p / ry
	cyp := -coefficient * ry * x1p / rx
	center := Point{cos*cxp - sin*cyp + (x1+x2)/2, sin*cxp + cos*cyp + (y1+y2)/2}

	startAngle := angleBetween(1, 0, (x1p-cxp)/rx, (y1p-cyp)/ry)
	angleDelta := clipDegreesTo360(angleBetween((x1p-cxp)/rx, (y1p-cyp)/ry, (-x1p-cxp)/rx,
		(-y1p-cyp)/ry))
	endAngle := clipDegreesTo360(startAngle + angleDelta)
	startAngle = clipDegreesTo360(startAngle)

	return &ArcParams{center, startAngle, endAngle, a.Rotation, rx, ry, a.Sweep}, nil
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
	Center Point

	// StartAngle must be between 0 and 360
	StartAngle float64

	// EndAngle must be between 0 and 360
	EndAngle float64

	// Rotation must be between 0 and 360
	Rotation float64

	// XRadius must be greater than 0
	XRadius float64

	// YRadius must be greater than 0
	YRadius float64

	Sweep bool
}

// Bounds computes the bounding box of the arc.
func (a *ArcParams) Bounds() Rect {
	minX, maxX := a.minMaxX()
	minY, maxY := a.minMaxY()
	return Rect{Point{minX, minY}, Point{maxX, maxY}}
}

// Length approximates the arc's length.
func (a *ArcParams) Length() float64 {
	var length float64
	for t := 0.0; t < 1; t += arcLengthApproximationInterval {
		length += Line{a.Evaluate(t), a.Evaluate(t + arcLengthApproximationInterval)}.Length()
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
	return a.evaluateAngle(angle)
}

func (a *ArcParams) minMaxX() (min, max float64) {
	x1, x2 := a.Evaluate(0).X, a.Evaluate(1).X
	min = math.Min(x1, x2)
	max = math.Max(x1, x2)

	tanRot := math.Tan(a.Rotation * math.Pi / 180)
	angle1 := clipDegreesTo360(180 / math.Pi * math.Atan(tanRot*-a.YRadius/a.XRadius))
	angle2 := clipDegreesTo360(180 + angle1)

	for _, angle := range []float64{angle1, angle2} {
		if a.includesAngle(angle) {
			xValue := a.evaluateAngle(angle).X
			min = math.Min(min, xValue)
			max = math.Max(max, xValue)
		}
	}

	return
}

func (a *ArcParams) minMaxY() (min, max float64) {
	y1, y2 := a.Evaluate(0).Y, a.Evaluate(1).Y
	min = math.Min(y1, y2)
	max = math.Max(y1, y2)

	cotanRot := 1 / math.Tan(a.Rotation*math.Pi/180)
	angle1 := clipDegreesTo360(180 / math.Pi * math.Atan(cotanRot*a.YRadius/a.XRadius))
	angle2 := clipDegreesTo360(180 + angle1)

	for _, angle := range []float64{angle1, angle2} {
		if a.includesAngle(angle) {
			yValue := a.evaluateAngle(angle).Y
			min = math.Min(min, yValue)
			max = math.Max(max, yValue)
		}
	}

	return
}

func (a *ArcParams) includesAngle(angle float64) bool {
	assertClippedTo360(angle)
	if a.StartAngle < a.EndAngle {
		if a.Sweep {
			return angle >= a.StartAngle && angle <= a.EndAngle
		} else {
			return angle <= a.StartAngle || angle >= a.EndAngle
		}
	} else {
		if a.Sweep {
			return angle <= a.EndAngle || angle >= a.StartAngle
		} else {
			return angle >= a.EndAngle && angle <= a.StartAngle
		}
	}
}

func (a *ArcParams) evaluateAngle(angle float64) Point {
	angle *= math.Pi / 180
	rotCos := math.Cos(math.Pi / 180 * a.Rotation)
	rotSin := math.Sin(math.Pi / 180 * a.Rotation)
	return Point{a.XRadius*math.Cos(angle)*rotCos - a.YRadius*rotSin*math.Sin(angle) + a.Center.X,
		a.XRadius*math.Cos(angle)*rotSin + a.YRadius*math.Sin(angle)*rotCos + a.Center.Y}
}

func clipDegreesTo360(angle float64) float64 {
	for angle < 0 {
		angle += 360
	}
	for angle > 360 {
		angle -= 360
	}
	return angle
}

func assertClippedTo360(angle float64) {
	if angle < 0 || angle > 360 {
		panic("angle is not clipped between 0 and 360")
	}
}

func angleBetween(v1x, v1y, v2x, v2y float64) float64 {
	dot := v1x*v2x + v1y*v2y
	magProduct := math.Sqrt(math.Pow(v1x, 2)+math.Pow(v1y, 2)) *
		math.Sqrt(math.Pow(v2x, 2)+math.Pow(v2y, 2))
	angle := math.Acos(dot / magProduct)
	if v1x*v2y-v1y*v2x < 0 {
		angle *= -1
	}
	return 180 / math.Pi * angle
}
