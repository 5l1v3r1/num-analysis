package ludecomp

import (
	"github.com/unixpickle/num-analysis/linalg"
)

// solveLowerTriangular solves the system Lx = b for x,
// where L is the lower triangular part of m and b is given.
// The lower triangular part of m refers to all the entries
// below and including the diagonal of m.
//
// This runs in O(m.N^2) time.
func solveLowerTriangular(m *linalg.Matrix, b linalg.Vector) linalg.Vector {
	if len(b) != m.Rows || !m.Square() {
		panic("dimension mismatch")
	}
	solution := make(linalg.Vector, len(b))
	for i := 0; i < m.Rows; i++ {
		answer := b[i]
		for j := 0; j < i; j++ {
			answer -= m.Get(i, j) * solution[j]
		}
		answer /= m.Get(i, i)
		solution[i] = answer
	}
	return solution
}

// solveUpperTriangular solves the system Ux = b for x,
// where U is the upper triangular part of m and b is given.
// The upper triangular part of m refers to all the entries
// below the diagonal of m, and an assumed diagonal of all 1's.
//
// This runs in O(m.N^2) time.
func solveUpperTriangular(m *linalg.Matrix, b linalg.Vector) linalg.Vector {
	if len(b) != m.Rows || !m.Square() {
		panic("dimension mismatch")
	}
	solution := make(linalg.Vector, len(b))
	for i := m.Rows - 1; i >= 0; i-- {
		answer := b[i]
		for j := m.Rows - 1; j > i; j-- {
			answer -= m.Get(i, j) * solution[j]
		}
		solution[i] = answer
	}
	return solution
}
