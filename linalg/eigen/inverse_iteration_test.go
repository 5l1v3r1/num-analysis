package eigen

import (
	"math"
	"testing"

	"github.com/unixpickle/num-analysis/linalg"
)

func TestInverseIterationBasic(t *testing.T) {
	mat := &linalg.Matrix{
		Rows: 3,
		Cols: 3,
		Data: []float64{
			66, 78, 76,
			78, 93, 92,
			76, 92, 94,
		},
	}

	vals, vecs, err := InverseIteration(mat, 10000)
	if err != nil {
		t.Error(err)
	}

	verifyEigs(t, vals, []float64{4.81397359013199e-02, 2.99176945337813e+00,
		2.49960090810721e+02})

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
	expectedRemaining := map[float64]struct{}{}
	for _, x := range expected {
		expectedRemaining[x] = struct{}{}
	}

ActualLoop:
	for _, v := range actual {
		for x := range expectedRemaining {
			if math.Abs(x-v) < 0.00001 {
				delete(expectedRemaining, x)
				continue ActualLoop
			}
		}
		t.Error("incorrect or duplicated eigenvalue:", v)
	}
}
