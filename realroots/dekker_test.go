package realroots

import (
	"math"
	"testing"
)

type sineFunction struct{}

func (_ sineFunction) Eval(x float64) float64 {
	return math.Sin(x)
}

type badDekker struct {
	callCount int
}

func (b *badDekker) Eval(x float64) float64 {
	b.callCount++
	if x <= 60 {
		return math.Pow(5, -x)
	} else {
		return math.Pow(5, -60) + -60*(x-60)
	}
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

	bad := &badDekker{}
	d = newDekker(bad, Interval{0, 61})
	for {
		d.Step()
		if d.Bounded(1e-7) {
			break
		}
	}
	if bad.callCount < 320 {
		t.Error("unexpected call count", bad.callCount)
	}
	if math.Abs(d.Root()-60) > 1e-6 {
		t.Error("expected root of ~60 but got", d.Root())
	}
}
