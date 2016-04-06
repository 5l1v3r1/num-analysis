package realroots

import (
	"math"
	"testing"
)

func TestBisection(t *testing.T) {
	p := Polynomial{75.8782, -45.3416, -15.3639, 5.29236}
	root := Bisection(p, Interval{4, 7}, 10)
	if math.Abs(root-4.1363) > 0.004 {
		t.Error("expected", 4.1363, "but got", root)
	}
	root = Bisection(p, Interval{1, 2}, 10)
	if math.Abs(root-1.34461) > 0.002 {
		t.Error("expected", 1.34461, "but got", root)
	}
}

func TestBisectionPrec(t *testing.T) {
	p := Polynomial{75.8782, -45.3416, -15.3639, 5.29236}
	root := BisectionPrec(p, Interval{4, 7}, 1)
	if math.Abs(root-4.1363) > 0.502 {
		t.Error("expected roughly", 4.1363, "but got", root)
	}
	root = BisectionPrec(p, Interval{4, 7}, 1e-7)
	if math.Abs(root-4.13629866) > 1e-7 {
		t.Error("expected roughly", 4.13629866, "but got", root)
	}
}
