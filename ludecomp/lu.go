package ludecomp

import "math"

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
	LU *Matrix

	// InPerm is the permutation that should be applied
	// to the input vector before solving.
	InPerm Perm

	// OutPerm is the permutation that should be applied
	// to the solution vector after solving (LU)x = Pb.
	OutPerm Perm
}

// Decompose generates the LU decomposition for an invertible matrix.
func Decompose(m *Matrix) *LU {
	res := &LU{
		LU:      &Matrix{N: m.N, V: make([]float64, m.N*m.N)},
		InPerm:  IdentityPerm(m.N),
		OutPerm: IdentityPerm(m.N),
	}
	copy(res.LU.V, m.V)
	for i := 0; i < m.N; i++ {
		pivotRow, pivotCol := res.bestPivot(i)
		if pivotCol != i {
			res.swapColumns(i, pivotCol)
		}
		if pivotRow != i {
			res.swapRows(i, pivotRow)
		}
		scaler := 1.0 / res.LU.Get(i, i)
		res.LU.Set(i, i, scaler)
		for k := i + 1; k < m.N; k++ {
			v := res.LU.Get(i, k)
			res.LU.Set(i, k, v*scaler)
		}
		for k := i + 1; k < m.N; k++ {
			v := res.LU.Get(k, i)
			res.LU.Set(k, i, scaler*v)
		}
	}
	res.OutPerm = res.OutPerm.Inverse()
	return res
}

// Solve computes the vector x such that Ax=v, where A is the
// decomposed matrix represented by l.
func (l *LU) Solve(v Vector) Vector {
	in := l.InPerm.Apply(v)
	sol1 := l.LU.SolveUpperTriangular(in)
	sol2 := l.LU.SolveLowerTriangular(sol1)
	return l.OutPerm.Apply(sol2)
}

func (l *LU) bestPivot(stepsDone int) (row, col int) {
	var biggestValue float64
	row = stepsDone
	col = stepsDone
	for i := stepsDone; i < l.LU.N; i++ {
		for j := stepsDone; j < l.LU.N; j++ {
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
	for k := 0; k < l.LU.N; k++ {
		v1 := l.LU.Get(k, i)
		v2 := l.LU.Get(k, j)
		l.LU.Set(k, i, v2)
		l.LU.Set(k, j, v1)
	}
}

func (l *LU) swapRows(step, pivotRow int) {
	l.InPerm.Swap(step, pivotRow)
	for i := 0; i < step; i++ {
		val1 := l.LU.Get(step, i)
		val2 := l.LU.Get(pivotRow, i)
		l.LU.Set(pivotRow, i, val1)
		l.LU.Set(step, i, val2)
	}
}
