package qrdecomp

import (
	"math"
	"math/rand"
	"testing"

	"github.com/unixpickle/num-analysis/kahan"
	"github.com/unixpickle/num-analysis/linalg"
)

const smallValue = 0.00001

type decomposer func(m *linalg.Matrix) (q, r *linalg.Matrix)

var test5x3Matrix = &linalg.Matrix{
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

var test4x4Matrix = &linalg.Matrix{
	Rows: 4,
	Cols: 4,
	Data: []float64{
		0.196693, 0.972318, 0.839954, 0.950959,
		0.977388, 0.528317, 0.676612, 0.486393,
		0.188801, 0.199278, 0.805243, 0.606177,
		0.112958, 0.665220, 0.082632, 0.123123,
	},
}

var testSingularMatrix = &linalg.Matrix{
	Rows: 4,
	Cols: 4,
	Data: []float64{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	},
}

func testDecomposer(t *testing.T, m *linalg.Matrix, d decomposer, qCols int) {
	q, r := d(m)

	if q.Cols != qCols {
		t.Error("unexpected number of columns in Q:", q.Cols)
	}

	product := q.Mul(r)
	if product.Rows != m.Rows || product.Cols != m.Cols {
		t.Error("dimension dismatch for QR:", product.Rows, "by", product.Cols)
	} else if matrixDifference(m, product) > smallValue {
		t.Error("bad matrix product:", product)
	}

	qTransposeQ := q.Transpose().Mul(q)
	identity := linalg.NewMatrixIdentity(qTransposeQ.Cols)
	if matrixDifference(qTransposeQ, identity) > smallValue {
		t.Error("Q is not orthogonal:", q)
	}

	for i := 0; i < r.Rows; i++ {
		for j := 0; j < i; j++ {
			if math.Abs(r.Get(i, j)) > smallValue {
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

func randomMatrix(size int) *linalg.Matrix {
	res := &linalg.Matrix{
		Rows: size,
		Cols: size,
		Data: make([]float64, size*size),
	}
	for i := range res.Data {
		res.Data[i] = rand.Float64()
	}
	return res
}
