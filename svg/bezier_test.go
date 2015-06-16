package svg

import "testing"

func TestQuadBezierCurveBounds(t *testing.T) {
	l := Line{Point{10, 10}, Point{40, 40}}
	curve := QuadraticBezier{l.Start, l.Midpoint(), l.End}
	bounds := curve.Bounds()
	if !pointsClose(bounds.Min, l.Start) || !pointsClose(bounds.Max, l.End) {
		t.Error("bad bounds:", bounds)
	}
	
	curve = QuadraticBezier{Point{10, 50}, Point{20, 10}, Point{30, 50}}
	bounds = curve.Bounds()
	if !pointsClose(bounds.Min, Point{10, 30}) ||
		!pointsClose(bounds.Max, Point{30, 50}) {
		t.Error("bad bounds:", bounds)
	}
}

func TestCubicBezierCurveBounds(t *testing.T) {
	l := Line{Point{10, 10}, Point{40, 40}}
	curve := CubicBezier{l.Start, l.Midpoint(), l.Midpoint(), l.End}
	bounds := curve.Bounds()
	if !pointsClose(bounds.Min, l.Start) || !pointsClose(bounds.Max, l.End) {
		t.Error("bad bounds:", bounds)
	}
	
	// TODO: write more tests here
}

func pointsClose(p1, p2 Point) bool {
	return Line{p1, p2}.Length() < 0.00001
}
