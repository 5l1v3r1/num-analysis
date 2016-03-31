package leastsquares

import (
	"github.com/unixpickle/num-analysis/kahan"
	"github.com/unixpickle/num-analysis/linalg"
	"github.com/unixpickle/num-analysis/linalg/qrdecomp"
)

// A Solver solves equations of the form
// A'*A*x = A'*b for x given a vector b.
// Each Solver implicitly knows how to solve
// the above equation for a predefined A.
type Solver struct {
	upperTriangular *linalg.Matrix
	orthogonal      *linalg.Matrix
}

// NewSolver creates a Solver which solves
// equations of the form m'*m*x = m'*b.
//
// The matrix m must have independent columns.
func NewSolver(m *linalg.Matrix) *Solver {
	if m.Cols > m.Rows {
		panic("columns cannot be independent")
	}
	q, r := qrdecomp.Householder(m)
	return &Solver{r, q.Transpose()}
}

// Solve solves A'*A*x = A'*b for x given
// a vector b, where A is represented by s.
func (s *Solver) Solve(b linalg.Vector) linalg.Vector {
	// It can be shown that, if A has independent columns and
	// is decomposed as Q*R, then A'*A*x=A'*b <=> R*x = Q'*b.

	bMatrix := &linalg.Matrix{Rows: len(b), Cols: 1, Data: []float64(b)}
	orthoB := s.orthogonal.Mul(bMatrix)

	return s.backSubstitute(linalg.Vector(orthoB.Data))
}

func (s *Solver) backSubstitute(b linalg.Vector) linalg.Vector {
	res := make(linalg.Vector, s.upperTriangular.Rows)
	for row := s.upperTriangular.Rows - 1; row >= 0; row-- {
		value := kahan.NewSummer64()
		value.Add(b[row])
		for col := row + 1; col < s.upperTriangular.Rows; col++ {
			value.Add(-s.upperTriangular.Get(row, col) * res[col])
		}
		res[row] = value.Sum() / s.upperTriangular.Get(row, row)
	}
	return res
}
