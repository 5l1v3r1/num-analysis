package svd

import (
	"math"
	"testing"

	"github.com/unixpickle/num-analysis/linalg"
)

func TestDecompose3x3(t *testing.T) {
	mat := &linalg.Matrix{
		Rows: 3,
		Cols: 3,
		Data: []float64{
			1, 2, 3,
			4, 5, 6,
			7, 8, 10,
		},
	}
	verifySVD(t, mat)
}

func TestDecompose3x3Singular(t *testing.T) {
	mat := &linalg.Matrix{
		Rows: 3,
		Cols: 3,
		Data: []float64{
			1, 2, 3,
			4, 5, 6,
			7, 8, 9,
		},
	}
	verifySVD(t, mat)
}

func TestDecompose3x5(t *testing.T) {
	mat := &linalg.Matrix{
		Rows: 3,
		Cols: 5,
		Data: []float64{
			3.1356e-01, 1.6989e-02, 5.1117e-04, 1.6086e-01, 5.9027e-01,
			1.4946e-01, 5.8112e-01, 9.9693e-01, 6.5264e-01, 3.9784e-01,
			7.7404e-01, 7.1985e-01, 3.8669e-01, 6.9035e-01, 2.5369e-01,
		},
	}
	verifySVD(t, mat)
}

func TestDecompose5x3(t *testing.T) {
	mat := &linalg.Matrix{
		Rows: 5,
		Cols: 3,
		Data: []float64{
			3.1356e-01, 1.6989e-02, 5.1117e-04,
			1.6086e-01, 5.9027e-01, 1.4946e-01,
			5.8112e-01, 9.9693e-01, 6.5264e-01,
			3.9784e-01, 7.7404e-01, 7.1985e-01,
			3.8669e-01, 6.9035e-01, 2.5369e-01,
		},
	}
	verifySVD(t, mat)
}

func verifySVD(t *testing.T, m *linalg.Matrix) {
	v, d, u := Decompose(m)

	if v.Rows != v.Cols || v.Rows != m.Rows {
		t.Error("invalid dimensions for V", v.Rows, v.Cols)
	} else {
		if !isOrthogonal(v) {
			t.Error("V is not orthogonal")
		}
	}

	if d.Rows != m.Rows || d.Cols != m.Cols {
		t.Error("invalid dimensions for D", d.Rows, d.Cols)
	} else {
		if !isDiagonal(d) {
			t.Error("D is not diagonal")
		}
	}

	if u.Rows != u.Cols || u.Cols != m.Cols {
		t.Error("invalid dimensions for U", u.Rows, u.Cols)
	} else {
		if !isOrthogonal(u) {
			t.Error("U is not orthogonal")
		}
	}

	product := v.Mul(d).Mul(u)
	if product.Rows != m.Rows || product.Cols != m.Cols {
		t.Error("invalid product dimensions", product.Rows, product.Cols)
	} else {
		if !matricesClose(product, m) {
			t.Error("invalid product")
		}
	}
}

func isOrthogonal(q *linalg.Matrix) bool {
	qProduct := q.Transpose().Mul(q)
	for i := 0; i < qProduct.Rows; i++ {
		for j := 0; j < qProduct.Cols; j++ {
			entry := qProduct.Get(i, j)
			if i == j {
				if math.Abs(1-entry) > 1e-5 {
					return false
				}
			} else if math.Abs(entry) > 1e-5 {
				return false
			}
		}
	}
	return true
}

func isDiagonal(d *linalg.Matrix) bool {
	for i := 0; i < d.Rows; i++ {
		for j := 0; j < d.Cols; j++ {
			entry := math.Abs(d.Get(i, j))
			if i != j && entry > 1e-5 {
				return false
			}
		}
	}
	return true
}

func matricesClose(m1, m2 *linalg.Matrix) bool {
	for i := 0; i < m1.Rows; i++ {
		for j := 0; j < m1.Cols; j++ {
			actual := m1.Get(i, j)
			expected := m2.Get(i, j)
			if math.Abs(actual-expected) > 0.000001 {
				return false
			}
		}
	}
	return true
}
