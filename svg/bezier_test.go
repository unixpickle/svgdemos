package svg

import "testing"

func TestQuadBezierCurveBounds(t *testing.T) {
	l := Line{Point{10, 10}, Point{40, 40}}
	curve := QuadraticBezier{l.Start, l.Midpoint(), l.End}
	bounds := curve.Bounds()
	expected := Rect{l.Start, l.End}
	if !rectsClose(bounds, expected) {
		t.Error("expected bounds", expected, "but got", bounds)
	}

	curve = QuadraticBezier{Point{10, 50}, Point{20, 10}, Point{30, 50}}
	bounds = curve.Bounds()
	expected = Rect{Point{10, 30}, Point{30, 50}}
	if !rectsClose(bounds, expected) {
		t.Error("expected bounds", expected, "but got", bounds)
	}

	curve = QuadraticBezier{Point{10, 10}, Point{40, 40}, Point{20, 50}}
	bounds = curve.Bounds()
	expected = Rect{Point{10, 10}, Point{28, 50}}
	if !rectsClose(bounds, expected) {
		t.Error("expected bounds", expected, "but got", bounds)
	}
}

func TestCubicBezierCurveBounds(t *testing.T) {
	l := Line{Point{10, 10}, Point{40, 40}}
	curve := CubicBezier{l.Start, l.Midpoint(), l.Midpoint(), l.End}
	bounds := curve.Bounds()
	expected := Rect{l.Start, l.End}
	if !rectsClose(bounds, expected) {
		t.Error("expected bounds", expected, "but got", bounds)
	}

	curve = CubicBezier{Point{10, 50}, Point{40, 10}, Point{70, 90}, Point{100, 50}}
	bounds = curve.Bounds()
	expected = Rect{Point{10, 38.452994616207484}, Point{100, 61.547005383792516}}
	if !rectsClose(bounds, expected) {
		t.Error("expected bounds", expected, "but got", bounds)
	}

	curve = CubicBezier{Point{96, 89}, Point{13, 46}, Point{14, 64}, Point{15, 91}}
	bounds = curve.Bounds()
	expected = Rect{Point{14.781782109764007, 63.23187046009383}, Point{96, 91}}
	if !rectsClose(bounds, expected) {
		t.Error("expected bounds", expected, "but got", bounds)
	}

	curve = CubicBezier{Point{50, 10}, Point{10, 40}, Point{90, 70}, Point{50, 100}}
	bounds = curve.Bounds()
	expected = Rect{Point{38.452994616207484, 10}, Point{61.547005383792516, 100}}
	if !rectsClose(bounds, expected) {
		t.Error("expected bounds", expected, "but got", bounds)
	}
}

func rectsClose(b1, b2 Rect) bool {
	return pointsClose(b1.Min, b2.Min) && pointsClose(b1.Max, b2.Max)
}

func pointsClose(p1, p2 Point) bool {
	return Line{p1, p2}.Length() < 0.00001
}
