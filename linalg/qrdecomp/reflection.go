package qrdecomp

import (
	"github.com/unixpickle/num-analysis/kahan"
	"github.com/unixpickle/num-analysis/linalg"
)

// A Reflection represents a linear transformation
// which reflects vectors across a given vector.
type Reflection struct {
	V linalg.Vector
}

// NewReflection creates a Reflection which
// reflects vectors across v, which needn't
// be normalized.
func NewReflection(v linalg.Vector) *Reflection {
	return &Reflection{v.Copy().Scale(1 / v.Mag())}
}

// Apply returns the Reflection of the given vector
// over r.V.
//
// If v has more components than r.V, then the first
// components of v are not modified and the reflection
// is applied to the last components of v.
func (r *Reflection) Apply(v linalg.Vector) linalg.Vector {
	if len(r.V) > len(v) {
		panic("dimension is too low")
	}
	ignoreCount := len(v) - len(r.V)

	s := kahan.NewSummer64()
	for i, x := range r.V {
		s.Add(x * v[i+ignoreCount])
	}
	dot := s.Sum()

	res := make(linalg.Vector, len(v))
	copy(res, v[:ignoreCount])
	for i, x := range r.V {
		resIdx := i + ignoreCount
		res[resIdx] = v[resIdx] - 2*x*dot
	}
	return res
}

// applyColumn applies the Reflection to a column of
// a matrix in place.
func (r *Reflection) applyColumn(m *linalg.Matrix, col int) {
	// TODO: perhaps make this more efficient.
	v := make(linalg.Vector, m.Rows)
	for i := range v {
		v[i] = m.Get(i, col)
	}
	v = r.Apply(v)
	for i, x := range v {
		m.Set(i, col, x)
	}
}
