package qrdecomp

import "github.com/unixpickle/num-analysis/linalg"

// A ReflectionChain represents a series of reflections
// which have increasing dimensionalities, starting at 2.
//
// In other words, a ReflectionChain represents a
// product of householder matrices.
type ReflectionChain struct {
	// Reflections is a list of reflections.
	// The n-th reflection must have dimension n+2.
	Reflections []Reflection
}

// Apply is the equivalent of multiplying a vector on the
// right of the matrix represented by r.
func (r *ReflectionChain) Apply(v linalg.Vector) linalg.Vector {
	for _, x := range r.Reflections {
		v = x.Apply(v)
	}
	return v
}

// Dim returns the size of the vectors that r operates on,
// which can also be thought of as the size of the matrix
// represented by r.
func (r *ReflectionChain) Dim() int {
	return len(r.Reflections) + 1
}

// Matrix generates a matrix which represents the linear
// transformation performed by r.
//
// Sometimes, you may not want a completely square matrix.
// The cols argument specifies how many columns of the
// matrix to generate.
//
// If cols == r.Dim(), then the resulting matrix is square.
func (r *ReflectionChain) Matrix(cols int) *linalg.Matrix {
	dim := r.Dim()
	if cols > dim {
		panic("too many columns requested")
	}
	res := linalg.NewMatrix(dim, cols)

	for i := 0; i < cols; i++ {
		// Since the unit vector will have zeroes below the i-th
		// value, some reflections have no effect.
		skipReflections := len(r.Reflections) - (i + 1)
		if skipReflections < 0 {
			skipReflections = 0
		}
		v := make(linalg.Vector, dim)
		v[i] = 1
		for k := skipReflections; k < len(r.Reflections); k++ {
			v = r.Reflections[i].Apply(v)
		}
		for k, x := range v {
			res.Set(k, i, x)
		}
	}

	return res
}
