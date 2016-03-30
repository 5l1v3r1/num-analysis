package linalg

import "github.com/unixpickle/num-analysis/kahan"

// A Matrix is an MxN matrix with real entries.
type Matrix struct {
	Rows int
	Cols int

	// Data is ordered from left to right, top
	// to bottom.
	Data []float64
}

// Get returns the element at the i-th row and
// the j-th column, where i and j start at 0.
func (m *Matrix) Get(i, j int) float64 {
	if i < 0 || j < 0 || i >= m.Rows || j >= m.Cols {
		panic("index out of bounds")
	}
	return m.Data[i*m.Cols+j]
}

// Set updates the element referenced by i and
// j, as explained for Get().
func (m *Matrix) Set(i, j int, val float64) {
	if i < 0 || j < 0 || i >= m.Rows || j >= m.Cols {
		panic("index out of bounds")
	}
	m.Data[i*m.Cols+j] = val
}

// Copy returns a copy of this matrix.
func (m *Matrix) Copy() *Matrix {
	res := &Matrix{
		Rows: m.Rows,
		Cols: m.Cols,
		Data: make([]float64, len(m.Data)),
	}
	copy(res.Data, m.Data)
	return res
}

// Square returns true if and only if this matrix
// is square.
func (m *Matrix) Square() bool {
	return m.Rows == m.Cols
}

// Scale multiplies this matrix by c in place.
func (m *Matrix) Scale(c float64) {
	for i, d := range m.Data {
		m.Data[i] = d * c
	}
}

// Add performs matrix addition on m in place, adding
// the values from m1.
//
// The dimensions of m1 must match the dimensions of m.
func (m *Matrix) Add(m1 *Matrix) {
	if m.Rows != m1.Rows || m.Cols != m1.Cols {
		panic("dimension mismatch")
	}
	for i, d := range m1.Data {
		m.Data[i] += d
	}
}

// Mul performs matrix multiplication with m on the
// left and m1 on the right, returning the new matrix.
//
// Matrix multiplication requires that m.Cols == m1.Rows.
// The resulting matrix will have the size m.Rows by m1.Cols.
func (m *Matrix) Mul(m1 *Matrix) *Matrix {
	if m.Cols != m1.Rows {
		panic("dimension mismatch")
	}
	res := &Matrix{
		Rows: m.Rows,
		Cols: m1.Cols,
		Data: make([]float64, m.Rows*m1.Cols),
	}
	dataIdx := 0
	for i := 0; i < res.Rows; i++ {
		for j := 0; j < res.Cols; j++ {
			summer := kahan.NewSummer64()
			for k := 0; k < m.Rows; k++ {
				summer.Add(m.Get(i, k) * m1.Get(k, j))
			}
			res.Data[dataIdx] = summer.Sum()
			dataIdx++
		}
	}
	return res
}

// Transpose transposes m in place.
func (m *Matrix) Transpose() {
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < i; j++ {
			temp := m.Get(i, j)
			m.Set(i, j, m.Get(j, i))
			m.Set(j, i, temp)
		}
	}
}
