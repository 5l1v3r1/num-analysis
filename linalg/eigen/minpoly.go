package eigen

import (
	"math"

	"github.com/unixpickle/num-analysis/kahan"
	"github.com/unixpickle/num-analysis/linalg"
	"github.com/unixpickle/num-analysis/linalg/qrdecomp"
	"github.com/unixpickle/num-analysis/mvroots"
)

// defaultMinPolyPrecision is the default factor
// by which the 1-norm of a matrix must be scaled
// down for the scaled down matrix to be "0".
const defaultMinPolyPrecision = 1e-11

// MinPolyPrec computes the polynomial P of least
// degree satisfying P(m)=0.
//
// The returned polynomial is represented as an array
// of coefficients, where the n-th element is the
// coefficient for the x^n term.
//
// The prec parameter specifies how to tell if a
// matrix is 0. If all of the entries have an
// absolute value <= prec, then the matrix is
// considered 0.
func MinPolyPrec(m *linalg.Matrix, prec float64) []float64 {
	if m.Rows != m.Cols {
		panic("matrix must be square")
	}
	if isZeroMatrix(m, prec) {
		return []float64{1}
	}
	for degree := 1; degree < m.Rows; degree++ {
		poly, m := minimizingPolynomial(m, degree)
		if isZeroMatrix(m, prec) {
			return poly
		}
	}
	poly, _ := minimizingPolynomial(m, m.Rows)
	return poly
}

// MinPoly is like MinPolyPrec, but it uses a
// reasonable default precision.
func MinPoly(m *linalg.Matrix) []float64 {
	var maxEntry float64
	for _, x := range m.Data {
		maxEntry = math.Max(maxEntry, math.Abs(x))
	}
	return MinPolyPrec(m, maxEntry*defaultMinPolyPrecision)
}

// MinEigs computes the roots of the minimal
// polynomial of a matrix.
//
// This is equivalent to finding the eigenvalues
// of the matrix and removing some (but not all)
// of the repeated ones.
func MinEigs(m *linalg.Matrix) []complex128 {
	poly := MinPoly(m)
	complexPoly := make(mvroots.Polynomial, len(poly))
	for i, x := range poly {
		complexPoly[i] = complex(x, 0)
	}
	return complexPoly.Roots()
}

func isZeroMatrix(m *linalg.Matrix, prec float64) bool {
	for _, x := range m.Data {
		if math.Abs(x) > prec {
			return false
		}
	}
	return true
}

// minimizingPolynomial finds a monic polynomial
// of degree deg which is the minimal polynomial
// of P (if the provided degree is correct).
// If deg is less than the degree of the minimal
// polynomial, then the result is meaningless.
// It returns the polynomial P and P(m).
func minimizingPolynomial(m *linalg.Matrix, deg int) ([]float64, *linalg.Matrix) {
	equations := matrixPowerColumnEquations(m, deg)
	_, r := qrdecomp.HouseholderReflections(equations)

	// This is a general technique for finding the
	// nullspace of an upper-triangular matrix when
	// you know that the last pivot is 0.
	// If r is invertible, then this will still give
	// a vector, but it will be meaningless.
	nullspace := make([]float64, deg+1)
	nullspace[deg] = 1
	for i := deg - 1; i >= 0; i-- {
		s := kahan.NewSummer64()
		s.Add(r.Get(i, r.Cols-1))
		for j := r.Cols - 2; j > i; j-- {
			s.Add(r.Get(i, j) * nullspace[j])
		}
		nullspace[i] = -s.Sum() / r.Get(i, i)
	}

	product := equations.Mul(linalg.NewMatrixColumn(linalg.Vector(nullspace)))
	product.Rows = m.Rows
	product.Cols = m.Cols

	return nullspace, product
}

// matrixPowerColumnEquations computes the powers of m
// like m^0, m^1, ..., m^deg, then generates a matrix
// where the i-th column is the "flattened"
// representation of m^i.
func matrixPowerColumnEquations(m *linalg.Matrix, deg int) *linalg.Matrix {
	columnEquations := linalg.NewMatrix(m.Rows*m.Cols, deg+1)
	powMat := linalg.NewMatrixIdentity(m.Rows)
	for i := 0; i <= deg; i++ {
		for j, x := range powMat.Data {
			columnEquations.Set(j, i, x)
		}
		if i < deg {
			powMat = m.Mul(powMat)
		}
	}
	return columnEquations
}
