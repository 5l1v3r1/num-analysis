package ludecomp

import "github.com/unixpickle/num-analysis/kahan"

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
