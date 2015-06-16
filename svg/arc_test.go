package svg

import "testing"

func TestArcBounds(t *testing.T) {
	arcs := []Arc{
		{Point{10, 10}, Point{30, 30}, 0, 50, 0, true, true},
		{Point{10, 10}, Point{30, 30}, 50, 0, 0, false, false},
	}
	expected := []Rect{
		{Point{10, 10}, Point{30, 30}},
		{Point{10, 10}, Point{30, 30}},
	}
	for i, arc := range arcs {
		e := expected[i]
		bounds := arc.Bounds()
		if !bounds.approxEqual(e) {
			t.Error("expected bounds", e, "but got", bounds, "for case", i)
		}
	}
}
