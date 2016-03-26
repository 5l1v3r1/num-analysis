package ludecomp

// LU stores all of the information about a matrix
// that has been decomposed into LU form.
//
// More specifically, the matrix is decomposed such that
// PAQ = LU, where P and Q are permutation matrices, L is
// lower-triangular, U is upper-triangular, and A is the
// original matrix.
type LU struct {
	// LU is a matrix which stores both L and U.
	// The lower part of this matrix stores L, and
	// the upper part stores U.
	LU *Matrix

	// InPerm is the permutation that should be applied
	// to the input vector before solving.
	InPerm Perm

	// OutPerm is the permutation that should be applied
	// to the solution vector after solving (LU)x = Pb.
	OutPerm Perm
}

// Solve computes the vector x such that Ax=v, where A is the
// decomposed matrix represented by l.
func (l *LU) Solve(v Vector) Vector {
	in := l.InPerm.Apply(v)
	sol1 := l.LU.SolveUpperTriangular(in)
	sol2 := l.LU.SolveLowerTriangular(sol1)
	return l.OutPerm.Apply(sol2)
}
