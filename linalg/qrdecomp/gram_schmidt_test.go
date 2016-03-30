package qrdecomp

import (
	"math"
	"testing"

	"github.com/unixpickle/num-analysis/kahan"
	"github.com/unixpickle/num-analysis/linalg"
)

func TestGramSchmidt(t *testing.T) {
	matrix := &linalg.Matrix{
		Rows: 5,
		Cols: 3,
		Data: []float64{
			0.196693, 0.972318, 0.839954,
			0.950959, 0.977388, 0.528317,
			0.676612, 0.486393, 0.188801,
			0.199278, 0.805243, 0.606177,
			0.112958, 0.665220, 0.082632,
		},
	}
	q, r := GramSchmidt(matrix)
	newMat := q.Mul(r)
	if matrixDifference(newMat, matrix) > 0.00001 {
		t.Error("matrix product is incorrect:", newMat)
	}
	if matrixDifference(q.Transpose().Mul(q), linalg.NewMatrixIdentity(3)) > 0.00001 {
		t.Error("Q is not orthogonal:", q)
	}
	for i := 0; i < r.Rows; i++ {
		for j := 0; j < i; j++ {
			if math.Abs(r.Get(i, j)) > 0.00001 {
				t.Error("R is not upper-triangular:", r)
			}
		}
	}
}

func matrixDifference(m1, m2 *linalg.Matrix) float64 {
	res := kahan.NewSummer64()
	for i, x := range m1.Data {
		res.Add(math.Abs(x - m2.Data[i]))
	}
	return res.Sum()
}
