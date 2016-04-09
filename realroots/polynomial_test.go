package realroots

import (
	"math"
	"testing"
)

func TestPolynomialRootsEven(t *testing.T) {
	poly := Polynomial{20, 0, -30, 0.3, 0.1}
	generalRootTest(t, poly, []float64{-18.8692, -0.814086, 0.820798, 15.8625}, 1e-4)
}

func TestPolynomialRootsOdd(t *testing.T) {
	poly := Polynomial{-1, 2, 10, 4, 5, 6}
	generalRootTest(t, poly, []float64{-1.17799, -0.477396, 0.222548}, 1e-4)
}

func TestPolynomialRootsFancy(t *testing.T) {
	poly := Polynomial{-1, 1, 1.0 / 2.0, 1.0 / 6.0, 1.0 / 24.0, 1.0 / 120.0, 1.0 / 720.0,
		1.0 / 5040.0}
	generalRootTest(t, poly, []float64{math.Log(2)}, 1e-4)
}

func generalRootTest(t *testing.T, poly Polynomial, roots []float64, prec float64) {
	actual := poly.OddRoots()
	if len(actual) != len(roots) {
		t.Error("unexpected root count:", actual)
	} else {
		for i, x := range roots {
			if math.Abs(x-actual[i]) > prec {
				t.Error("invalid root:", actual[i], "expected", x)
			}
		}
	}
}
