package ludecomp

import (
	"math"
	"testing"
)

func TestBackSubstitution(t *testing.T) {
	m1 := Matrix{
		N: 4,
		V: []float64{
			3, 5, 7, 2,
			1, 9.5, 2, 4,
			3, 18, 3.6, 5,
			2, 8, 1, 3.14,
		},
	}
	lowerSolution := m1.SolveLowerTriangular(Vector{1, 2, 3, 4})
	upperSolution := m1.SolveUpperTriangular(Vector{1, 2, 3, 4})

	expectedLower := Vector{1 / 3, 10 / 57, -55 / 171, 6317 / 8810}
	expectedUpper := Vector{12, 20, -17, 4}
	if math.Abs(cosTheta(lowerSolution, expectedLower)-1) > 0.0001 {
		t.Error("incorrect lower triangular solution")
	}
	if math.Abs(cosTheta(upperSolution, expectedUpper)-1) > 0.0001 {
		t.Error("incorrect upper triangular solution", upperSolution)
	}
}

func cosTheta(v1, v2 Vector) float64 {
	return v1.Dot(v2) / math.Sqrt(v1.Dot(v1)*v2.Dot(v2))
}
