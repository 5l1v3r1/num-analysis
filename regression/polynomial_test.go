package regression

import (
	"math"
	"testing"
)

func TestFitPolynomial(t *testing.T) {
	points := []Point{
		Point{-2, 1},
		Point{0, 1},
		Point{2, 25},
		Point{1, 7},
	}
	poly := FitPolynomial(3, points)
	expected := Polynomial{1, 3, 2, 1}
	for i, x := range poly {
		if math.Abs(x-expected[i]) > 0.0001 {
			t.Errorf("bad x^%d term; got %f expected %f", i, x, expected[i])
		}
	}
}
