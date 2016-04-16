package conjgrad

import "github.com/unixpickle/num-analysis/linalg"

// A LinTran is a square linear transformation.
type LinTran interface {
	// Dim returns the number of dimensions in
	// the input and output vectors of this
	// linear transformation.
	Dim() int

	// Apply applies this linear transformation
	// to a vector and returns the result.
	Apply(v linalg.Vector) linalg.Vector
}

// MatLinTran is a LinTran which is defined as
// L such that L(v) = M*v for a matrix M.
//
// This is generally only useful for relatively
// dense linear transformations; for a sparse
// matrix, it might be better to implement your
// own LinTran.
type MatLinTran struct {
	M *linalg.Matrix
}

func (m MatLinTran) Dim() int {
	return m.M.Rows
}

func (m MatLinTran) Apply(v linalg.Vector) linalg.Vector {
	column := &linalg.Matrix{Rows: len(v), Cols: 1, Data: v}
	res := m.M.Mul(column)
	return linalg.Vector(res.Data)
}
