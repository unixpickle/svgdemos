package svg

import "testing"

func TestParsePath(t *testing.T) {
	path, err := ParsePath(`M600,350 l 50,-25 a25,25 -30 0,1 50,-25 l50-25
		a25,50-30 0,1 50,-25c 0 -10 25 25 50 -25 0 -10 25 25 50 -25
		zL 100 100 H 200 V 200H100Z M 200,200C300,300 0,0 400 400
		S 425,375 450,350 Q 400 400 500,230 T 450 230`)
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
	errorPaths := []string{
		"M Z", "M", "Z 1", "M 1", "M 1 2 3", "L", "L 1", "L 1 2 3",
		"m Z", "m", "z 1", "m 1", "m 1 2 3", "l", "l 1", "l 1 2 3",
		"p 1 2 3", "x y z", "b 1 2",
		"1", "1 2", "1 2 3", ",1",
	}
	for _, pathStr := range errorPaths {
		if _, err := ParsePath(pathStr); err == nil {
			t.Error("expected path to trigger error:", pathStr)
		}
	}
}

func TestAbsolutePath(t *testing.T) {
	path, err := ParsePath(`m 10,10,10-10 l 20,20 h 10-20 v 30
		c 10,10 20-20 -20,30 s 10 10 20-20 q 10 10 20 0 t 20 0
		a 10,20 20.5 1 1 10,-20 l -40-40`)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := ParsePath(`M 10,10 20,0 L 40,20 H 50 30 V 50
		C 40,60 50,30 10,80 S 20,90 30,60 Q 40,70 50,60 T 70,60
		A 10,20 20.5 1 1 80,40 L 40,0`)
	if err != nil {
		t.Fatal(err)
	}
	actual := path.Absolute()
	for i, x := range expected {
		a := actual[i]
		if !a.Equals(x) {
			t.Error("command", i, "should be", x, "but it is", a)
		}
	}
}
