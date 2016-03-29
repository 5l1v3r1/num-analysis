package cholesky

import (
	"math"
	"testing"

	"github.com/unixpickle/num-analysis/ludecomp"
)

func TestSolve3x3(t *testing.T) {
	mat := &ludecomp.Matrix{
		N: 3,
		V: []float64{
			14, 26, 17,
			26, 57, 32,
			17, 32, 25,
		},
	}
	dec := Decompose(mat)

	problems := []ludecomp.Vector{
		{1, 2, 3},
	}
	solutions := []ludecomp.Vector{
		{-222.0 / 529.0, -2.0 / 529.0, 217.0 / 529.0},
	}
	for i, problem := range problems {
		solution := dec.Solve(ludecomp.Vector{1, 2, 3})
		if math.IsNaN(solution[0]) {
			t.Error("NaN's in solution", solution, "for", problem)
			continue
		}
		if solutionDiff(solution, solutions[i]) > 0.000001 {
			t.Error("wrong solution for", problem, "got", solution, "expected", solutions[i])
		}
	}

}

func solutionDiff(s1, s2 ludecomp.Vector) float64 {
	var diff float64
	for i, x := range s1 {
		diff += math.Pow(x-s2[i], 2)
	}
	return math.Sqrt(diff)
}
