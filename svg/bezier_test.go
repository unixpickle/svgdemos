package svg

import "testing"

func TestQuadBezierCurveBounds(t *testing.T) {
	l := Line{Point{10, 10}, Point{40, 40}}
	curves := []QuadraticBezier{
		{l.Start, l.Midpoint(), l.End},
		{Point{10, 50}, Point{20, 10}, Point{30, 50}},
		{Point{10, 10}, Point{40, 40}, Point{20, 50}},
	}
	expected := []Rect{
		{l.Start, l.End},
		{Point{10, 30}, Point{30, 50}},
		{Point{10, 10}, Point{28, 50}},
	}

	for i, curve := range curves {
		e := expected[i]
		bounds := curve.Bounds()
		if !bounds.approxEqual(e) {
			t.Error("expected bounds", e, "but got", bounds, "for case", i)
		}
	}
}

func TestCubicBezierCurveBounds(t *testing.T) {
	l := Line{Point{10, 10}, Point{40, 40}}
	curves := []CubicBezier{
		{l.Start, l.Midpoint(), l.Midpoint(), l.End},
		{Point{10, 50}, Point{40, 10}, Point{70, 90}, Point{100, 50}},
		{Point{96, 89}, Point{13, 46}, Point{14, 64}, Point{15, 91}},
		{Point{50, 10}, Point{10, 40}, Point{90, 70}, Point{50, 100}},
	}
	bounds := []Rect{
		{l.Start, l.End},
		{Point{10, 38.452994616207484}, Point{100, 61.547005383792516}},
		{Point{14.781782109764007, 63.23187046009383}, Point{96, 91}},
		{Point{38.452994616207484, 10}, Point{61.547005383792516, 100}},
	}

	for i, curve := range curves {
		e := bounds[i]
		bounds := curve.Bounds()
		if !bounds.approxEqual(e) {
			t.Error("expected bounds", e, "but got", bounds, "for case", i)
		}
	}
}
