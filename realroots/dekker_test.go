package realroots

import (
	"math"
	"testing"
)

type sineFunction struct{}

func (_ sineFunction) Eval(x float64) float64 {
	return math.Sin(x)
}

func TestDekker(t *testing.T) {
	d := newDekker(sineFunction{}, Interval{1, 4})
	bound := 1e-11
	for {
		d.Step()
		if d.Bounded(bound) {
			break
		}
	}
	root := d.Root()
	if math.Abs(root-math.Pi) > 1.5*bound {
		t.Error("expected", math.Pi, "got", root, "(using bound ", bound, ")")
	}
}
