package eigen

import (
	"math"
	"testing"

	"github.com/unixpickle/num-analysis/linalg"
)

type eigenSolver func(m *linalg.Matrix) ([]float64, []linalg.Vector)

func testEigenSolver(t *testing.T, solver eigenSolver, mat *linalg.Matrix, expected []float64) {
	vals, vecs := solver(mat)
	if vecs == nil || vals == nil {
		return
	}

	verifyEigs(t, vals, expected)

	for i, vec := range vecs {
		multVec := mat.Mul(linalg.NewMatrixColumn(vec)).Col(0)
		scaledVec := vec.Copy().Scale(-vals[i])
		errorVec := scaledVec.Add(multVec)
		if errorVec.Dot(errorVec) > 1e-10 {
			t.Error("bad eigenvector", vec, "for eigenvalue", vals[i])
		}
	}
}

func verifyEigs(t *testing.T, actual, expected []float64) {
	expectedRemaining := make([]float64, len(expected))
	copy(expectedRemaining, expected)

ActualLoop:
	for _, v := range actual {
		for i, x := range expectedRemaining {
			if math.Abs(x-v) < 0.00001 {
				expectedRemaining[i] = expectedRemaining[len(expectedRemaining)-1]
				expectedRemaining = expectedRemaining[:len(expectedRemaining)-1]
				continue ActualLoop
			}
		}
		t.Error("incorrect or duplicated eigenvalue:", v, "expected one of", expectedRemaining)
	}
}
