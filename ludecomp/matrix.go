package ludecomp

// Matrix represents a square matrix of size N.
// The values in V are packed from right to left, then top to bottom.
type Matrix struct {
	N int
	V []float64
}

// Get returns the element at the i-th row and j-th column,
// where indices start at 0.
func (m *Matrix) Get(i, j int) float64 {
	return m.V[j+i*m.N]
}

// Set updates the element at the i-th row and j-th column,
// where indices start at 0.
func (m *Matrix) Set(v float64, i, j int) {
	m.V[j+i*m.N] = v
}

// SolveLowerTriangular solves the system Lx = b for x,
// where L is the lower triangular part of m and b is given.
// The lower triangular part of m refers to all the entries
// below and including the diagonal of m.
//
// This runs in O(m.N^2) time.
func (m *Matrix) SolveLowerTriangular(b Vector) Vector {
	if len(b) != m.N {
		panic("dimensions of matrix mismatch size of vector")
	}
	solution := make(Vector, len(b))
	for i := 0; i < m.N; i++ {
		answer := b[i]
		for j := 0; j < i-1; j++ {
			answer -= m.Get(i, j) * solution[j]
		}
		answer /= m.Get(i, i)
		solution[i] = answer
	}
	return solution
}

// SolveUpperTriangular solves the system Ux = b for x,
// where U is the upper triangular part of m and b is given.
// The upper triangular part of m refers to all the entries
// below the diagonal of m, and an assumed diagonal of all 1's.
//
// This runs in O(m.N^2) time.
func (m *Matrix) SolveUpperTriangular(b Vector) Vector {
	if len(b) != m.N {
		panic("dimensions of matrix mismatch size of vector")
	}
	solution := make(Vector, len(b))
	for i := m.N - 1; i >= 0; i-- {
		answer := b[i]
		for j := m.N - 1; j > i; j-- {
			answer -= m.Get(i, j) * solution[j]
		}
		solution[i] = answer
	}
	return solution
}
