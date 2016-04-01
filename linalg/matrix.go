package linalg

import (
	"bytes"
	"strconv"

	"github.com/unixpickle/num-analysis/kahan"
)

const outputPrecision = 10

// A Matrix is an MxN matrix with real entries.
type Matrix struct {
	Rows int
	Cols int

	// Data is ordered from left to right, top
	// to bottom.
	Data []float64
}

// NewMatrix creates a matrix of a given size with
// zeroes in every cell.
func NewMatrix(rows, cols int) *Matrix {
	res := &Matrix{
		Rows: rows,
		Cols: cols,
		Data: make([]float64, rows*cols),
	}
	return res
}

// NewMatrixIdentity returns an identity matrix of
// the given size.
func NewMatrixIdentity(size int) *Matrix {
	res := NewMatrix(size, size)
	for i := 0; i < size; i++ {
		res.Set(i, i, 1)
	}
	return res
}

// NewMatrixColumn creates a column matrix using a
// vector's values.
func NewMatrixColumn(v Vector) *Matrix {
	res := NewMatrix(len(v), 1)
	copy(res.Data, v)
	return res
}

// Get returns the element at the i-th row and
// the j-th column, where i and j start at 0.
func (m *Matrix) Get(i, j int) float64 {
	return m.Data[i*m.Cols+j]
}

// Set updates the element referenced by i and
// j, as explained for Get().
func (m *Matrix) Set(i, j int, val float64) {
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

// Scale multiplies m by c in place and returns m.
func (m *Matrix) Scale(c float64) *Matrix {
	for i, d := range m.Data {
		m.Data[i] = d * c
	}
	return m
}

// Add performs matrix addition on m in place and
// returns m.
//
// The dimensions of m1 must match the dimensions of m.
func (m *Matrix) Add(m1 *Matrix) *Matrix {
	if m.Rows != m1.Rows || m.Cols != m1.Cols {
		panic("dimension mismatch")
	}
	for i, d := range m1.Data {
		m.Data[i] += d
	}
	return m
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
			for k := 0; k < m.Cols; k++ {
				summer.Add(m.Get(i, k) * m1.Get(k, j))
			}
			res.Data[dataIdx] = summer.Sum()
			dataIdx++
		}
	}
	return res
}

// Transpose returns a new matrix which represents
// the transpose of m.
func (m *Matrix) Transpose() *Matrix {
	res := NewMatrix(m.Cols, m.Rows)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			res.Set(j, i, m.Get(i, j))
		}
	}
	return res
}

// Col gets the column vector at the given column index.
func (m *Matrix) Col(col int) Vector {
	res := make(Vector, m.Rows)
	for i := range res {
		res[i] = m.Get(i, col)
	}
	return res
}

// String returns a human-readable, row-by-row string
// representation of this matrix.
// Each row in m is separated by semicolons in m's
// string representation.
func (m *Matrix) String() string {
	var res bytes.Buffer
	res.WriteRune('[')
	for row := 0; row < m.Rows; row++ {
		if row != 0 {
			res.WriteString("; ")
		}
		for col := 0; col < m.Cols; col++ {
			if col != 0 {
				res.WriteRune(' ')
			}
			val := m.Get(row, col)
			res.WriteString(strconv.FormatFloat(val, 'g', outputPrecision, 64))
		}
	}
	res.WriteRune(']')
	return res.String()
}
