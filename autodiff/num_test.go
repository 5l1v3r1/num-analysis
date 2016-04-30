package autodiff

import (
	"math"
	"testing"
)

func TestNumArithmetic(t *testing.T) {
	// There are three variables, x0, x1, and x2.
	// We will compute (x0*x1)^2 + (x1)/(x0-x2) - x0^3
	// Where x0=4, x1 = 10, and x2 = 15.

	x0 := NewNumVar(4, 3, 0)
	x1 := NewNumVar(10, 3, 1)
	x2 := NewNumVar(15, 3, 2)

	value := (x0.Mul(x1)).PowScaler(2)
	value = value.Add(x1.Div(x0.Sub(x2)))
	value = value.Sub(x0.PowScaler(3))

	testNumericalValue(t, value, 16886.0/11.0, []float64{751.917, 319.909, 0.0826446})
}

func TestNumPow(t *testing.T) {
	// There are two variables, x0 and x1.
	// We will compute (x0+x1)^(pi*x0-x1^2)
	// Where x0=5 and x1=10.

	x0 := NewNumVar(5, 2, 0)
	x1 := NewNumVar(10, 2, 1)
	const2 := NewNum(math.Pi, 2)

	value := x0.Add(x1).Pow(x0.Mul(const2).Sub(x1.PowScaler(2)))

	testNumericalValue(t, value, 7.32609e-100, []float64{2.11586e-99, -4.37957e-98})

	// We will compute x0^0 x0=0.
	x0 = NewNumVar(0, 2, 0)
	value = x0.PowScaler(0)
	testNumericalValue(t, value, 1, []float64{0, 0})
}

func TestNumFuncs(t *testing.T) {
	// There are two variables, x0 and x1.
	// We will compute x0*sqrt(exp(x1*x0*sin(2cos(x1))))
	// Where x0=2 and x1=3.

	x0 := NewNumVar(2, 2, 0)
	x1 := NewNumVar(3, 2, 1)
	const2 := NewNum(2, 2)

	value := x0.Mul(x1.Cos().Mul(const2).Sin().Mul(x0).Mul(x1).Exp().Sqrt())

	testNumericalValue(t, value, 0.127558, []float64{-0.111762, -0.0740555})
}

func testNumericalValue(t *testing.T, v Num, expected float64, grad []float64) {
	if math.Abs((v.Value-expected)/expected) > 1e-5 {
		t.Error("value should be", expected, "but got", v.Value)
	}

	for i, ex := range grad {
		actual := v.Gradient[i]
		if math.IsNaN(actual) || math.Abs((actual-ex)/ex) > 1e-5 {
			t.Error("partial", i, "should be", ex, "but got", actual)
		}
	}
}
