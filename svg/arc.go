package svg

type Arc struct {
	From     Point
	To       Point
	XRadius  float64
	YRadius  float64
	Rotation float64
	LargeArc bool
	Sweep    bool
}
