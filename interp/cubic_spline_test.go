package interp

import (
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/unixpickle/num-analysis/kahan"
)

type slopeComputer func(xs, ys []float64) float64

func TestMidArcSplineForwards(t *testing.T) {
	xs := []float64{1, 2, 3, 3.25, 5, 6.5, 7, 9}
	ys := []float64{2, 3, 5, 4, 6, 7, 7, 4}
	testMidArcSpline(t, xs, ys)
}

func TestMidArcSplineBackwards(t *testing.T) {
	xs := []float64{9, 7, 6.5, 5, 3.25, 3, 2, 1}
	ys := []float64{2, 3, 5, 4, 6, 7, 7, 4}
	testMidArcSpline(t, xs, ys)
}

func TestMidArcSplineUnordered(t *testing.T) {
	xs := []float64{7, 3.25, 1, 9, 3, 6, 5, 2}
	ys := []float64{2, 3, 5, 4, 6, 7, 7, 4}
	testMidArcSpline(t, xs, ys)
}

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

func testMidArcSpline(t *testing.T, xs, ys []float64) {
	testAnySpline(t, xs, ys, MidArcStyle, func(x, y []float64) float64 {
		return (y[2] - y[0]) / (x[2] - x[0])
	})
}

func testStandardSpline(t *testing.T, xs, ys []float64) {
	testAnySpline(t, xs, ys, StandardStyle, func(x, y []float64) float64 {
		m1 := (y[1] - y[0]) / (x[1] - x[0])
		m2 := (y[2] - y[1]) / (x[2] - x[1])
		return (m1 + m2) / 2
	})
}

func testAnySpline(t *testing.T, xs, ys []float64, style SplineStyle, slopeComp slopeComputer) {
	s := NewCubicSpline(style)
	for i, x := range xs {
		s.Add(x, ys[i])
		for j, x1 := range xs[:i+1] {
			actual := s.Eval(x1)
			expected := ys[j]
			if math.Abs(actual-expected) > 1e-5 || math.IsNaN(actual) {
				t.Errorf("got %f expected %f at %f", actual, expected, x1)
			}
		}
		sortedX, sortedY := sortedValues(xs[:i+1], ys[:i+1])
		for j := 1; j < i; j++ {
			slope := slopeComp(sortedX[j-1:j+2], sortedY[j-1:j+2])
			theX := sortedX[j]
			actual := s.Deriv(theX)
			if math.Abs(actual-slope) > 1e-5 || math.IsNaN(actual) {
				t.Errorf("got slope %f expected %f at %f", actual, slope, theX)
			}
		}
	}

	minX := xs[0]
	maxX := xs[0]
	for _, x := range xs {
		minX = math.Min(minX, x)
		maxX = math.Max(maxX, x)
	}
	for i := 0; i < 100; i++ {
		x1 := rand.Float64()*(maxX-minX) + minX
		x2 := rand.Float64()*(maxX-minX) + minX
		actual := s.Integ(x1, x2)
		expected := reimannSum(s, x1, x2)
		if math.IsNaN(actual) || math.Abs(actual-expected) > 1e-2 {
			t.Errorf("got integral %f when expected %f on interval [%f,%f]",
				actual, expected, x1, x2)
		}
	}
}

func sortedValues(x, y []float64) ([]float64, []float64) {
	res1 := make([]float64, len(x))
	res2 := make([]float64, len(y))
	copy(res1, x)
	copy(res2, y)
	p := &sortablePairs{res1, res2}
	sort.Sort(p)
	return p.X, p.Y
}

type sortablePairs struct {
	X []float64
	Y []float64
}

func (s *sortablePairs) Len() int {
	return len(s.X)
}

func (s *sortablePairs) Less(i, j int) bool {
	return s.X[i] < s.X[j]
}

func (s *sortablePairs) Swap(i, j int) {
	s.X[i], s.X[j] = s.X[j], s.X[i]
	s.Y[i], s.Y[j] = s.Y[j], s.Y[i]
}

func reimannSum(f *CubicSpline, x1, x2 float64) float64 {
	if x1 > x2 {
		return -reimannSum(f, x2, x1)
	}
	sum := kahan.NewSummer64()
	for x := x1; x < x2; x += 1e-3 {
		sum.Add(1e-3 * f.Eval(x))
	}
	return sum.Sum()
}
