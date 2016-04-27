package autodiff

import (
	"math"
	"testing"
)

func TestDeepNumArithmetic(t *testing.T) {
	// We will compute the expression
	// x*(x+1) + 2x*(3-2x)/(x-15)
	// Where x=9

	x := NewDeepNumVar(9, 3)
	fifteen := NewDeepNum(15, 3)

	term1 := x.Mul(x.AddScaler(1))
	x.Div(x)
	term2 := x.MulScaler(2).Mul(x.MulScaler(-2).AddScaler(3)).Div(x.Sub(fifteen))
	sum := term1.Add(term2)

	testDeepNumValue(t, sum, []float64{135.0, 75.0 / 2.0, 19.0 / 2.0, 15.0 / 4.0})
}

func testDeepNumValue(t *testing.T, d *DeepNum, expected []float64) {
	for i, x := range expected {
		if d == nil {
			t.Error("not enough values: expected", len(expected), "but got", i)
			return
		}
		if math.Abs(d.Value-x)/x > 1e-5 {
			t.Errorf("invalid value %d: expected %f but got %f", i, x, d.Value)
		}
		d = d.Deriv
	}

	if d != nil {
		t.Error("too many derivatives: expected", len(expected), "but got more.")
	}
}
