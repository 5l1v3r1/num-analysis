package leastsquares

import (
	"math"
	"testing"

	"github.com/unixpickle/num-analysis/linalg"
)

func TestSolver(t *testing.T) {
	matrix := &linalg.Matrix{
		Rows: 5,
		Cols: 3,
		Data: []float64{
			5.95913322176935e-01, 2.08633691137558e-01, 7.65776841603420e-01,
			1.17606644889278e-02, 6.51852164874548e-01, 7.99626569156839e-01,
			2.75872166958871e-01, 3.15046281477938e-01, 2.43200532424610e-01,
			5.95581731565972e-01, 5.06533069056311e-01, 8.38506134994520e-01,
			6.20826184766993e-01, 8.23262230426032e-01, 9.02147303422546e-01,
		},
	}
	solver := NewSolver(matrix)
	problems := []linalg.Vector{
		{5.2951620e-01, 1.4021439e-01, 2.1952707e-01, 3.7669667e-01, 2.7811324e-02},
		{6.7811274e-01, 9.0022400e-01, 4.2592947e-01, 8.2595138e-01, 2.7262794e-02},
		{0, 0, 0, 0, 0},
	}
	solutions := []linalg.Vector{
		{2.1177021e-01, -7.1687047e-01, 6.9765115e-01},
		{-6.8761059e-01, -8.6266337e-01, 1.7186178e+00},
		{0, 0, 0},
	}
	for i, problem := range problems {
		actual := solver.Solve(problem)
		expected := solutions[i]
		if vectorDiff(actual, expected) > 0.000001 {
			t.Error("got", actual, "expected", expected, "for problem", problem)
		}
	}
}

func vectorDiff(v1, v2 linalg.Vector) float64 {
	var sum float64
	for i, x := range v1 {
		sum += math.Abs(v2[i] - x)
	}
	return sum
}
