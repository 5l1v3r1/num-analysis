package svd

import (
	"math"
	"math/rand"

	"github.com/unixpickle/num-analysis/linalg"
	"github.com/unixpickle/num-analysis/linalg/eigen"
)

// Decompose returns the singular value decomposition
// of a matrix m.
// The MxN matrix m is decomposed into v*d*u where v
// is an MxM orthogonal matrix, d is an MxN diagonal
// matrix, and u is an NxN orthogonal matrix.
func Decompose(m *linalg.Matrix) (v, d, u *linalg.Matrix) {
	if m.Cols > m.Rows {
		u, d, v = Decompose(m.Transpose())
		u = u.Transpose()
		v = v.Transpose()
		d = d.Transpose()
		return
	}

	rightVecs := make([]linalg.Vector, 0, m.Cols)
	leftVecs := make([]linalg.Vector, 0, m.Rows)
	singularValues := make([]float64, 0, m.Cols)

	normalMat := m.Transpose().Mul(m)
	res := eigen.SymmetricAsync(normalMat)
	for vec := range res.Vectors {
		<-res.Values
		rightVecs = append(rightVecs, vec)
		leftVec := m.Mul(linalg.NewMatrixColumn(vec)).Col(0)

		// Computing the singular value using matrix multiplication
		// is more accurate than using the square root of an
		// eigenvalue of the normal matrix.
		val := math.Sqrt(leftVec.Dot(leftVec))

		singularValues = append(singularValues, val)

		projectOut(leftVec, leftVecs)
		if normalize(leftVec) {
			leftVecs = append(leftVecs, leftVec)
		} else {
			leftVecs = addRandomVector(leftVecs, m.Rows)
		}
	}

	leftVecs = completeOrthoBasis(leftVecs, m.Rows)

	return svdMatrices(leftVecs, rightVecs, singularValues)
}

func projectOut(v linalg.Vector, vecs []linalg.Vector) {
	for _, x := range vecs {
		projScale := x.Dot(v)
		v.Add(x.Copy().Scale(-projScale))
	}
}

func normalize(v linalg.Vector) bool {
	mag := math.Sqrt(v.Dot(v))
	if mag == 0 {
		return false
	}
	v.Scale(1 / mag)
	return true
}

func completeOrthoBasis(basis []linalg.Vector, size int) []linalg.Vector {
	for len(basis) < size {
		basis = addRandomVector(basis, size)
	}
	return basis
}

func addRandomVector(basis []linalg.Vector, size int) []linalg.Vector {
	randVec := make(linalg.Vector, size)
	for i := range randVec {
		randVec[i] = (rand.Float64() * 2) - 1
	}
	projectOut(randVec, basis)
	normalize(randVec)
	return append(basis, randVec)
}

func svdMatrices(leftVecs, rightVecs []linalg.Vector, vals []float64) (v, d, u *linalg.Matrix) {
	v = linalg.NewMatrix(len(leftVecs), len(leftVecs))
	for j, leftVec := range leftVecs {
		for i, val := range leftVec {
			v.Set(i, j, val)
		}
	}

	d = linalg.NewMatrix(len(leftVecs), len(rightVecs))
	for i, x := range vals {
		d.Set(i, i, x)
	}

	u = linalg.NewMatrix(len(rightVecs), len(rightVecs))
	for i, rightVec := range rightVecs {
		for j, val := range rightVec {
			u.Set(i, j, val)
		}
	}

	return
}
