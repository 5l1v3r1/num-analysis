package autodiff

import (
	"math"
	"testing"
)

func TestPowerSeriesSine(t *testing.T) {
	actual := PowerSeries(differentiableSine, 0, 8)
	expected := []float64{0, 1, 0, -1.0 / 6.0, 0, 1.0 / 120.0, 0, -1.0 / 5040.0, 0}
	for i, x := range expected {
		a := actual[i]
		if math.Abs(x-a) > 1e-8 {
			t.Error("expected", x, "but got", a, "for term", i)
		}
	}

	actual = PowerSeries(differentiableExp, 0, 6)
	expected = []float64{1, 1, 1.0 / 2.0, 1.0 / 6.0, 1.0 / 24.0, 1.0 / 120.0, 1.0 / 720.0}
	for i, x := range expected {
		a := actual[i]
		if math.Abs(x-a) > 1e-8 {
			t.Error("expected", x, "but got", a, "for term", i)
		}
	}
}

func differentiableSine(d *DeepNum) *DeepNum {
	return d.Sin()
}

func differentiableExp(d *DeepNum) *DeepNum {
	return d.Exp()
}
