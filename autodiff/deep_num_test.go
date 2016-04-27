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

func TestDeepNumExponentials(t *testing.T) {
	// We will compute the expression
	// (x+2)^(-x^3-3x^2+x+1) where x=1.

	x := NewDeepNumVar(1, 3)

	base := x.AddScaler(2)
	exponent := x.PowScaler(3).MulScaler(-1).Sub(x.PowScaler(2).MulScaler(3)).Add(x).AddScaler(1)
	value := base.Pow(exponent)

	testDeepNumValue(t, value, []float64{1.0 / 9.0, -1.05062, 7.90147, -38.0578})
}

func TestDeepNumFunctions(t *testing.T) {
	// We will compute the expression
	// sin(sqrt(x^3 + 2x)) where x=3.

	x := NewDeepNumVar(3, 3)

	value := x.PowScaler(3).Add(x.MulScaler(2)).Sqrt().Sin()
	testDeepNumValue(t, value, []float64{-0.512954, 2.16675, 3.66096, -12.0968})
}

func testDeepNumValue(t *testing.T, d *DeepNum, expected []float64) {
	for i, x := range expected {
		if d == nil {
			t.Error("not enough values: expected", len(expected), "but got", i)
			return
		}
		if math.IsNaN(d.Value) || math.Abs((d.Value-x)/x) > 1e-5 {
			t.Errorf("invalid value %d: expected %f but got %f", i, x, d.Value)
		}
		d = d.Deriv
	}

	if d != nil {
		t.Error("too many derivatives: expected", len(expected), "but got more.")
	}
}
