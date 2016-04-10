package optimization

import "github.com/unixpickle/num-analysis/linalg"

type UnimodalFunc interface {
	Eval(x float64) float64
}

// GradFunc is a multivariable function with a
// computable gradient.
type GradFunc interface {
	// Dim returns the number of input arguments
	// this function takes.
	Dim() int

	// Eval evaluates the function for the given
	// argument vector.
	Eval(vec linalg.Vector) float64

	// Gradient evaluates the gradient of the
	// function at a given argument vector.
	Gradient(vec linalg.Vector) linalg.Vector
}

// LinSysFunc is a GradFunc which computes
// ||Ax-b||^2 where A is a pre-determined
// matrix and b is a pre-determined vector.
type LinSysFunc struct {
	// matrix represents A in Ax=b.
	matrix *linalg.Matrix

	// normal represents 2*A'*A.
	normal *linalg.Matrix

	// transpose represents 2*A'*b.
	constTerm *linalg.Matrix

	// product represents b in Ax=b.
	product linalg.Vector
}

// NewLinSysFunc generates a LinSysFunc which
// represents the equation a*x=b.
func NewLinSysFunc(a *linalg.Matrix, b linalg.Vector) *LinSysFunc {
	return &LinSysFunc{
		matrix:    a,
		normal:    a.Transpose().Mul(a).Scale(2),
		constTerm: a.Transpose().Mul(linalg.NewMatrixColumn(b)).Scale(-2),
		product:   b,
	}
}

// Dim returns the number of unknowns in
// the system.
func (l *LinSysFunc) Dim() int {
	return l.matrix.Cols
}

// Eval evaluates ||Ax-b||^2 for the given
// value x.
func (l *LinSysFunc) Eval(x linalg.Vector) float64 {
	mat := linalg.NewMatrixColumn(x)
	product := l.matrix.Mul(mat)
	res := linalg.Vector(product.Data)
	diff := res.Scale(-1).Add(l.product)
	return diff.Dot(diff)
}

// Gradient returns the gradient of ||Ax-b||^2
// with respect to x.
func (l *LinSysFunc) Gradient(x linalg.Vector) linalg.Vector {
	variableTerm := l.normal.Mul(linalg.NewMatrixColumn(x))
	total := variableTerm.Add(l.constTerm)
	return linalg.Vector(total.Data)
}
