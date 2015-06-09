package svg

import "math"

type CubicBezier struct {
	Start    Point
	Control1 Point
	Control2 Point
	End      Point
}

// Bounds computes the bounding box for the Bezier curve.
func (c *CubicBezier) Bounds() Rect {
	minX := math.Min(c.Start.X, c.End.X)
	maxX := math.Max(c.Start.X, c.End.X)
	minY := math.Min(c.Start.Y, c.End.Y)
	maxY := math.Max(c.Start.Y, c.End.Y)

	xExtrema := cubicBezierExtrema(c.Start.X, c.Control1.X, c.Control2.X,
		c.End.X)
	yExtrema := cubicBezierExtrema(c.Start.Y, c.Control1.Y, c.Control2.Y,
		c.End.Y)
	for _, xValue := range xExtrema {
		minX = math.Min(minX, xValue)
		maxX = math.Max(maxX, xValue)
	}
	for _, yValue := range yExtrema {
		minY = math.Min(minY, yValue)
		maxY = math.Max(maxY, yValue)
	}

	return Rect{Point{minX, minY}, Point{maxX, maxY}}
}

// Evaluate gets a point on the bezier curve for a parameter between 0 and 1.
func (c *CubicBezier) Evaluate(t float64) Point {
	x := cubicBezierPolynomial(c.Start.X, c.Control1.X, c.Control2.X, c.End.X,
		t)
	y := cubicBezierPolynomial(c.Start.Y, c.Control1.Y, c.Control2.Y, c.End.Y,
		t)
	return Point{x, y}
}

func cubicBezierExtrema(A, B, C, D float64) []float64 {
	// These coefficients result from taking the derivative of the cubic bezier
	// polynomial.
	a := 3*D - 9*C + 9*B - 3*A
	b := 6*A - 12*B + 6*C
	c := 3 * (B - A)
	discriminant := math.Pow(b, 2) - 4*a*c
	if discriminant < 0 {
		return []float64{}
	}

	solution1 := (-b + math.Sqrt(discriminant)) / (2 * a)
	solution2 := (-b - math.Sqrt(discriminant)) / (2 * a)
	result := make([]float64, 0, 2)
	if solution1 >= 0 && solution1 <= 1 {
		result = append(result, cubicBezierPolynomial(A, B, C, D, solution1))
	}
	if solution2 >= 0 && solution2 <= 1 {
		result = append(result, cubicBezierPolynomial(A, B, C, D, solution2))
	}
	return result
}

func cubicBezierPolynomial(A, B, C, D, t float64) float64 {
	return A*math.Pow(1-t, 3) + 3*B*t*math.Pow(1-2, 2) + 3*C*(1-t)*t*t + D*t*t*t
}
