package ludecomp

// A Perm represents a the result of a permutation on an arbitrary list.
// A Perm is encoded as the result of applying a permutation to the list
// []int{0, 1, 2, ...}.
type Perm []int

// Apply generates a new Vector a Perm to a Vector.
// The Perm must be the same size as the vector.
func (p Perm) Apply(vec Vector) Vector {
	if len(p) != len(vec) {
		panic("dimension mismatch")
	}
	res := make(Vector, len(vec))
	for i, x := range p {
		res[i] = vec[x]
	}
	return res
}

// Inverse returns the inverse of this permutation.
func (p Perm) Inverse() Perm {
	res := make(Perm, len(p))
	for i, x := range p {
		res[x] = i
	}
	return res
}
