package eigen

import (
	"errors"
	"math"
	"math/rand"

	"github.com/unixpickle/num-analysis/kahan"
	"github.com/unixpickle/num-analysis/linalg"
	"github.com/unixpickle/num-analysis/linalg/ludecomp"
)

var ErrMaxSteps = errors.New("maximum steps exceeded")

// InverseIteration uses inverse iteration to approximate
// the eigenvalues and eigenvectors of a symmetric matrix m.
//
// The maxIters argument acts as a sort of "timeout".
// If the algorithm spends more than maxIters iterations
// looking for the next eigenvector without converging,
// then ErrMaxSteps is returned along with the eigenvalues
// and eigenvalues which were already found.
func InverseIteration(m *linalg.Matrix, maxIters int) ([]float64, []linalg.Vector, error) {
	if !m.Square() {
		panic("input matrix must be square")
	}
	values := make([]float64, 0, m.Rows)
	vectors := make([]linalg.Vector, 0, m.Rows)
	for i := 0; i < m.Rows; i++ {
		val, vec, remaining := inverseIterate(m, vectors, maxIters)
		if remaining < 0 {
			return values, vectors, ErrMaxSteps
		}
		val, vec, remaining = powerIterate(m, val, vec, vectors, remaining)
		if remaining < 0 {
			return values, vectors, ErrMaxSteps
		}
		vec.Scale(1 / math.Sqrt(vec.Dot(vec)))
		values = append(values, val)
		vectors = append(vectors, vec)
	}
	return values, vectors, nil
}

func inverseIterate(m *linalg.Matrix, vecs []linalg.Vector,
	maxIters int) (float64, linalg.Vector, int) {
	// Once the pivots differ by sqrt(epsilon), we may lose
	// half of our double's precision when computing A^-1*x.
	// This seems like a logical place to stop trying to
	// find a nearer approximation.
	pivotThreshold := math.Sqrt(math.Nextafter(1, 2) - 1)

	vec := randomVector(m.Rows)
	deleteProjections(vec, vecs)
	val := approxEigenvalue(m, vec)
	for i := 0; i < maxIters; i++ {
		mat := m.Copy()
		for j := 0; j < m.Rows; j++ {
			mat.Set(j, j, mat.Get(j, j)-val)
		}
		lu := ludecomp.Decompose(mat)
		if lu.PivotScale() < pivotThreshold {
			return val, vec, maxIters - (i + 1)
		}
		vec = lu.Solve(vec)
		deleteProjections(vec, vecs)
		normalizeVector(vec)
		val = approxEigenvalue(m, vec)
	}
	return 0, nil, -1
}

func powerIterate(m *linalg.Matrix, val float64, vec linalg.Vector, ortho []linalg.Vector,
	remaining int) (float64, linalg.Vector, int) {
	var lastError float64
	for i := 0; i < remaining; i++ {
		deleteProjections(vec, ortho)
		vec = m.Mul(linalg.NewMatrixColumn(vec)).Col(0)
		normalizeVector(vec)
		val = approxEigenvalue(m, vec)
		backError := backError(m, vec, val)
		if i == 0 {
			lastError = backError
		} else {
			if backError >= lastError {
				return val, vec, remaining - (i + 1)
			}
			lastError = backError
		}
	}
	return 0, nil, -1
}

func deleteProjections(v linalg.Vector, vecs []linalg.Vector) {
	for _, pv := range vecs {
		v.Add(pv.Copy().Scale(-pv.Dot(v)))
	}
}

func randomVector(size int) linalg.Vector {
	res := make(linalg.Vector, size)
	for i := range res {
		res[i] = rand.Float64()*2 - 1
	}
	return res
}

func approxEigenvalue(m *linalg.Matrix, vec linalg.Vector) float64 {
	colVec := linalg.NewMatrixColumn(vec)
	return vec.Dot(m.Mul(colVec).Col(0)) / vec.Dot(vec)
}

func normalizeVector(v linalg.Vector) {
	var mag float64
	for _, x := range v {
		mag = math.Max(mag, math.Abs(x))
	}
	if mag == 0 {
		for i := range v {
			v[i] = 1
		}
	} else {
		v.Scale(1 / mag)
	}
}

func oneNorm(v linalg.Vector) float64 {
	res := kahan.NewSummer64()
	for _, x := range v {
		res.Add(math.Abs(x))
	}
	return res.Sum()
}

func backError(m *linalg.Matrix, vec linalg.Vector, val float64) float64 {
	multiplied := m.Mul(linalg.NewMatrixColumn(vec)).Col(0)
	scaled := vec.Copy().Scale(-val)
	return oneNorm(multiplied.Add(scaled))
}
