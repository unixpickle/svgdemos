package svg

import "testing"

func TestParsePath(t *testing.T) {
	path, err := ParsePath("M600,350 l 50,-25 a25,25 -30 0,1 50,-25 l50-25" +
		"a25,50-30 0,1 50,-25 c 0 -10 25 25 50 -25 0 -10 25 25 50 -25" +
		"zL 100 100 H 200 V 200H100Z M 200,200C300,300 0,0 400 400 S 425,375 " +
		"450,350 Q 400 400 500,230 T 450 230")
	if err != nil {
		t.Fatal(err)
	}
	expected := Path{
		PathCmd{"M", []float64{600, 350}},
		PathCmd{"l", []float64{50, -25}},
		PathCmd{"a", []float64{25, 25, -30, 0, 1, 50, -25}},
		PathCmd{"l", []float64{50, -25}},
		PathCmd{"a", []float64{25, 50, -30, 0, 1, 50, -25}},
		PathCmd{"c", []float64{0, -10, 25, 25, 50, -25,
			0, -10, 25, 25, 50, -25}},
		PathCmd{"z", []float64{}},
		PathCmd{"L", []float64{100, 100}},
		PathCmd{"H", []float64{200}},
		PathCmd{"V", []float64{200}},
		PathCmd{"H", []float64{100}},
		PathCmd{"Z", []float64{}},
		PathCmd{"M", []float64{200, 200}},
		PathCmd{"C", []float64{300, 300, 0, 0, 400, 400}},
		PathCmd{"S", []float64{425, 375, 450, 350}},
		PathCmd{"Q", []float64{400, 400, 500, 230}},
		PathCmd{"T", []float64{450, 230}},
	}
	if len(expected) != len(path) {
		t.Fatal("invalid length:", len(expected))
	}
	for i, expect := range expected {
		actual := path[i]
		if !expect.Equals(actual) {
			t.Error("path command", i, "should be", expect, "but is", actual)
		}
	}
}
