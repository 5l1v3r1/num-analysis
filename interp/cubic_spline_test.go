package interp

import (
	"math"
	"testing"
)

func TestStandardSplineForwards(t *testing.T) {
	xs := []float64{1, 2, 3, 3.25, 5, 6.5, 7, 9}
	ys := []float64{2, 3, 5, 4, 6, 7, 7, 4}
	testStandardSpline(t, xs, ys)
}

func TestStandardSplineBackwards(t *testing.T) {
	xs := []float64{9, 7, 6.5, 5, 3.25, 3, 2, 1}
	ys := []float64{2, 3, 5, 4, 6, 7, 7, 4}
	testStandardSpline(t, xs, ys)
}

func TestStandardSplineUnordered(t *testing.T) {
	xs := []float64{7, 3.25, 1, 9, 3, 6, 5, 2}
	ys := []float64{2, 3, 5, 4, 6, 7, 7, 4}
	testStandardSpline(t, xs, ys)
}

func testStandardSpline(t *testing.T, xs, ys []float64) {
	s := NewCubicSpline(StandardStyle)
	for i, x := range xs {
		s.Add(x, ys[i])
		for j, x1 := range xs[:i+1] {
			actual := s.Eval(x1)
			expected := ys[j]
			if math.Abs(actual-expected) > 1e-5 || math.IsNaN(actual) {
				t.Errorf("got %f expected %f at %f", actual, expected, x1)
			}
		}
		for j := 1; j < i; j++ {
			slope := (ys[j+1] - ys[j-1]) / (xs[j+1] - xs[j-1])
			actual := s.Deriv(xs[j])
			if math.Abs(actual - slope) > 1e-5 || math.IsNaN(actual) {
				t.Errorf("got slope %f expected %f at %f", actual, slope, xs[j])
			}
		}
	}
}
