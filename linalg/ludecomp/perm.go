package ludecomp

import "github.com/unixpickle/num-analysis/linalg"

// A Perm represents a the result of a permutation on an arbitrary list.
// A Perm is encoded as the result of applying a permutation to the list
// []int{0, 1, 2, ...}.
type Perm []int

// IdentityPerm returns a permutation of size n that does nothing.
func IdentityPerm(n int) Perm {
	res := make(Perm, n)
	for i := range res {
		res[i] = i
	}
	return res
}

// Apply generates a new Vector a Perm to a Vector.
// The Perm must be the same size as the vector.
func (p Perm) Apply(vec linalg.Vector) linalg.Vector {
	if len(p) != len(vec) {
		panic("dimension mismatch")
	}
	res := make(linalg.Vector, len(vec))
	for i, x := range p {
		res[i] = vec[x]
	}
	return res
}

// Swap applies a swap to this permutation.
// The permutation p becomes a new permutation which is
// equivalent to applying the old one, then applying a swap.
func (p Perm) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Inverse returns the inverse of this permutation.
func (p Perm) Inverse() Perm {
	res := make(Perm, len(p))
	for i, x := range p {
		res[x] = i
	}
	return res
}
