package linalg

import (
	"math"

	"github.com/unixpickle/num-analysis/kahan"
)

// Vector is an ordered list of floating points
// which can be manipulated like a vector.
type Vector []float64

// Dot returns the dot product of two vectors.
// The dimensions of v and v1 must match.
func (v Vector) Dot(v1 Vector) float64 {
	if len(v) != len(v1) {
		panic("dimension mismatch")
	}
	summer := kahan.NewSummer64()
	for i, x := range v {
		summer.Add(x * v1[i])
	}
	return summer.Sum()
}

// Copy returns a copy of this vector.
func (v Vector) Copy() Vector {
	res := make(Vector, len(v))
	copy(res, v)
	return res
}

// Scale scales v in place and returns v.
func (v Vector) Scale(c float64) Vector {
	for i, x := range v {
		v[i] = x * c
	}
	return v
}

// Add adds v1 to v in place and returns v.
func (v Vector) Add(v1 Vector) Vector {
	for i, x := range v1 {
		v[i] += x
	}
	return v
}

// Mag returns the magnitude of this vector using
// a 2-norm.
func (v Vector) Mag() float64 {
	return math.Sqrt(v.Dot(v))
}

// MaxAbs returns the max of the absolute values
// of every component in the vector.
func (v Vector) MaxAbs() float64 {
	var res float64
	for _, x := range v {
		res = math.Max(res, math.Abs(x))
	}
	return res
}
