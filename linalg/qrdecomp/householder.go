package qrdecomp

import (
	"math"

	"github.com/unixpickle/num-analysis/kahan"
	"github.com/unixpickle/num-analysis/linalg"
)

// Householder decomposes an MxN matrix m into a
// product q*r where q is an MxN matrix with
// orthonormal columns and Q is an NxN
// upper-triangular matrix.
func Householder(m *linalg.Matrix) (q, r *linalg.Matrix) {
	var qRef *ReflectionChain
	qRef, r = HouseholderReflections(m)
	q = qRef.Matrix(r.Rows)
	return
}

// HouseholderReflections is like Householder, except
// that q is represented as a ReflectionChain rather
// than as a numerical matrix.
func HouseholderReflections(m *linalg.Matrix) (q *ReflectionChain, r *linalg.Matrix) {
	r = m.Copy()

	size := r.Cols
	if r.Rows < r.Cols {
		size = r.Rows
	}

	if r.Cols < r.Rows {
		q = &ReflectionChain{
			Reflections: make([]Reflection, size),
		}
	} else {
		q = &ReflectionChain{
			Reflections: make([]Reflection, size-1),
		}
	}

	for col := 0; col < len(q.Reflections); col++ {
		ref := eliminationReflection(col, r)
		q.Reflections[len(q.Reflections)-(col+1)] = *ref
	}

	if r.Cols < r.Rows {
		trimmedR := linalg.NewMatrix(r.Cols, r.Cols)
		for i := 0; i < r.Cols; i++ {
			for j := i; j < r.Cols; j++ {
				trimmedR.Set(i, j, r.Get(i, j))
			}
		}
		r = trimmedR
	}

	return
}

// eliminationReflection figures out a reflection to
// eliminate the sub-diagonal of the given column.
// It applies this reflection to the columns of the
// matrix and returns the reflection.
func eliminationReflection(col int, r *linalg.Matrix) *Reflection {
	magSum := kahan.NewSummer64()
	for i := col; i < r.Rows; i++ {
		magSum.Add(r.Get(i, col) * r.Get(i, col))
	}
	mag := math.Sqrt(magSum.Sum())

	var reflection *Reflection
	if mag == 0 {
		basisVec := make(linalg.Vector, r.Rows-col)
		basisVec[0] = 1
		reflection = NewReflection(basisVec)
	} else {
		refVec := make(linalg.Vector, r.Rows-col)
		newPivot := mag

		// Avoid cancelling firstComp and +/-mag for
		// better numerical accuracy.
		if firstComp := r.Get(col, col); firstComp < 0 {
			refVec[0] = firstComp - mag
		} else {
			refVec[0] = firstComp + mag
			newPivot = -newPivot
		}

		for i := 1; i < len(refVec); i++ {
			refVec[i] = r.Get(col+i, col)
		}
		reflection = NewReflection(refVec)
		r.Set(col, col, newPivot)
		for i := col + 1; i < r.Rows; i++ {
			r.Set(i, col, 0)
		}
	}

	for c := col + 1; c < r.Cols; c++ {
		reflection.applyColumn(r, c)
	}

	return reflection
}
