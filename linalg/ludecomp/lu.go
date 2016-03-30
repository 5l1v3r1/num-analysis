package ludecomp

import (
	"math"

	"github.com/unixpickle/num-analysis/linalg"
)

// LU stores all of the information about a matrix
// that has been decomposed into LU form.
//
// More specifically, the matrix is decomposed such that
// PAQ = LU, where P and Q are permutation matrices, L is
// lower-triangular, U is upper-triangular, and A is the
// original matrix.
type LU struct {
	// LU is a matrix which stores both L and U.
	// The lower part of this matrix stores L, and
	// the upper part stores U.
	LU *linalg.Matrix

	// InPerm is the permutation that should be applied
	// to the input linalg.Vector before solving.
	InPerm Perm

	// OutPerm is the permutation that should be applied
	// to the solution linalg.Vector after solving (LU)x = Pb.
	OutPerm Perm
}

// Decompose generates the LU decomposition for a
// square, invertible matrix m.
func Decompose(m *linalg.Matrix) *LU {
	if !m.Square() {
		panic("dimension mismatch")
	}
	res := &LU{
		LU:      m.Copy(),
		InPerm:  IdentityPerm(m.Rows),
		OutPerm: IdentityPerm(m.Rows),
	}
	for i := 0; i < m.Rows; i++ {
		pivotRow, pivotCol := res.bestPivot(i)
		if pivotCol != i {
			res.swapColumns(i, pivotCol)
		}
		if pivotRow != i {
			res.swapRows(i, pivotRow)
		}
		pivot := res.LU.Get(i, i)
		res.upperTriangularElimination(i, pivot)

		// As doing some math will show, the entries of L are the same
		// as the lower-triangular entries of the original matrix, so
		// no further computation needs to be done to compute these entries.
	}
	res.OutPerm = res.OutPerm.Inverse()
	return res
}

// Solve computes the linalg.Vector x such that Ax=v, where A is the
// decomposed matrix represented by l.
func (l *LU) Solve(v linalg.Vector) linalg.Vector {
	in := l.InPerm.Apply(v)
	sol1 := solveLowerTriangular(l.LU, in)
	sol2 := solveUpperTriangular(l.LU, sol1)
	return l.OutPerm.Apply(sol2)
}

func (l *LU) bestPivot(stepsDone int) (row, col int) {
	var biggestValue float64
	row = stepsDone
	col = stepsDone
	for i := stepsDone; i < l.LU.Rows; i++ {
		for j := stepsDone; j < l.LU.Rows; j++ {
			x := math.Abs(l.LU.Get(i, j))
			if x > biggestValue {
				biggestValue = x
				row = i
				col = j
			}
		}
	}
	return
}

func (l *LU) swapColumns(i, j int) {
	l.OutPerm.Swap(i, j)
	for k := 0; k < l.LU.Rows; k++ {
		v1 := l.LU.Get(k, i)
		v2 := l.LU.Get(k, j)
		l.LU.Set(k, i, v2)
		l.LU.Set(k, j, v1)
	}
}

func (l *LU) swapRows(step, pivotRow int) {
	// By swapping the rows of the upper and lower triangular
	// matrices simultaneously, we make up for the fact that we
	// cannot entirely swap rows of either of the two matrices
	// without losing triangularity.

	// It can be shown that the resulting product LU is equivalent
	// to fully swapping the rows of the lower-triangular matrix,
	// making it non-triangular.
	l.InPerm.Swap(step, pivotRow)
	for i := 0; i < l.LU.Rows; i++ {
		val1 := l.LU.Get(step, i)
		val2 := l.LU.Get(pivotRow, i)
		l.LU.Set(pivotRow, i, val1)
		l.LU.Set(step, i, val2)
	}
}

func (l *LU) upperTriangularElimination(step int, pivot float64) {
	// Divide the current row by the pivot.
	invPivot := 1 / pivot
	for col := step + 1; col < l.LU.Rows; col++ {
		v := l.LU.Get(step, col)
		l.LU.Set(step, col, v*invPivot)
	}

	// Subtract the current row from subsequent rows.
	for row := step + 1; row < l.LU.Rows; row++ {
		subScale := l.LU.Get(row, step)
		for col := step + 1; col < l.LU.Rows; col++ {
			val := l.LU.Get(row, col)
			subVal := l.LU.Get(step, col)
			l.LU.Set(row, col, val-subVal*subScale)
		}
	}
}
