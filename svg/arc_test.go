package svg

import "testing"

func TestArcBounds(t *testing.T) {
	arcs := []Arc{
		{Point{10, 10}, Point{30, 30}, 0, 50, 0, true, true},
		{Point{10, 10}, Point{30, 30}, 50, 0, 0, false, false},
		{Point{50, 10}, Point{60, 20}, 10, 10, 30, true, false},
		{Point{50, 10}, Point{60, 20}, 10, 20, 30, true, false},
		{Point{60, 20}, Point{50, 10}, 10, 20, 30, true, true},
		{Point{10, 10}, Point{30, 10}, 20, 10, 90, false, false},
		{Point{10, 10}, Point{20, 20}, 10, 10, 0, false, false},
		{Point{10, 10}, Point{20, 20}, 10, 10, 0, true, false},
		{Point{10, 10}, Point{20, 25}, 10, 15, 0, true, false},
	}
	expected := []Rect{
		{Point{10, 10}, Point{30, 30}},
		{Point{10, 10}, Point{30, 30}},
		{Point{40, 10}, Point{60, 30}},
		{Point{33.755562782270175, 10}, Point{60, 44.98681871515589}},
		{Point{33.755562782270175, 10}, Point{60, 44.98681871515589}},
		{Point{10, 10}, Point{30, 30}},
		{Point{10, 10}, Point{20, 20}},
		{Point{0, 10}, Point{20, 30}},
		{Point{0, 10}, Point{20, 40}},
	}
	for i, arc := range arcs {
		e := expected[i]
		var bounds Rect
		if params, line := arc.Params(); params != nil {
			bounds = params.Bounds()
		} else {
			bounds = line.Bounds()
		}
		if !bounds.approxEqual(e) {
			t.Error("expected bounds", e, "but got", bounds, "for case", i)
		}
	}
}
