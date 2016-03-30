package qrdecomp

import (
	"math"

	"github.com/unixpickle/num-analysis/kahan"
	"github.com/unixpickle/num-analysis/linalg"
)

// GramSchmidt decomposes an MxN matrix m into a
// product q*r where q is an MxN matrix with
// orthonormal columns and Q is an NxN upper-triangular
// matrix.
func GramSchmidt(m *linalg.Matrix) (q, r *linalg.Matrix) {
	q = m.Copy()
	r = linalg.NewMatrix(m.Cols, m.Cols)

	for step := 0; step < q.Cols; step++ {
		for projCol := 0; projCol < step; projCol++ {
			dot := dotColumns(q, step, projCol)
			for row := 0; row < q.Rows; row++ {
				val := q.Get(row, step)
				val -= dot * q.Get(row, projCol)
				q.Set(row, step, val)
			}
			r.Set(projCol, step, dot)
		}
		mag := math.Sqrt(dotColumns(q, step, step))
		scaler := 1.0 / mag
		r.Set(step, step, mag)
		for row := 0; row < q.Rows; row++ {
			val := q.Get(row, step)
			q.Set(row, step, val*scaler)
		}
	}

	return
}

func dotColumns(m *linalg.Matrix, i, j int) float64 {
	s := kahan.NewSummer64()
	for k := 0; k < m.Rows; k++ {
		s.Add(m.Get(k, i) * m.Get(k, j))
	}
	return s.Sum()
}
