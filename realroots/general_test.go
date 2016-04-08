package realroots

import (
	"math"
	"testing"
)

func TestRootWellBehaved(t *testing.T) {
	root := Root(sineFunction{}, Interval{3, 4})
	if math.Abs(root-math.Pi) > 1e-11 {
		t.Error("invalid root:", root, "(expected pi)")
	}
	root = RootPrec(sineFunction{}, Interval{3, 4}, 1e-4)
	if math.Abs(root-math.Pi) > 1.5e-4 {
		t.Error("invalid root:", root, "(expected pi)")
	}
}

func TestRootPoorlyBehaved(t *testing.T) {
	f := &badDekker{}
	root := RootPrec(f, Interval{0, 61}, 1e-11)
	if math.Abs(root-60) > 1.5e-10 {
		t.Error("invalid root", root, "(expected 60)")
	}
	if f.callCount > 140 {
		t.Errorf("too many iterations (%d); bisection not used.", f.callCount)
	}
	f.callCount = 0

	root = Root(f, Interval{0, 61})
	if math.Abs(root-60) > 1e-11 {
		t.Error("invalid root", root, "(expected 60)")
	}
	if f.callCount > 160 {
		t.Errorf("too many iterations (%d); bisection not used.", f.callCount)
	}
	f.callCount = 0
}
