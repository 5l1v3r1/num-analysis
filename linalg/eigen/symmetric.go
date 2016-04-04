package eigen

import (
	"errors"
	"math"
	"math/rand"
	"time"

	"github.com/unixpickle/num-analysis/kahan"
	"github.com/unixpickle/num-analysis/linalg"
	"github.com/unixpickle/num-analysis/linalg/ludecomp"
)

var ErrTimeout = errors.New("timeout exceeded")

type EigenChan struct {
	// Values is a channel of eigenvalues.
	// It will be closed once all eigenvalues
	// have been read.
	Values <-chan float64

	// Vectors is a channel of eigenvectors.
	// It will be closed once all eigenvectors
	// have been read.
	Vectors <-chan linalg.Vector

	// Cancel is a channel which you may close to
	// terminate the algorithm early.
	Cancel chan<- struct{}
}

// Symmetric approximates the eigenvalues and eigenvectors
// of a symmetric matrix m.
func Symmetric(m *linalg.Matrix) ([]float64, []linalg.Vector) {
	vals := make([]float64, 0, m.Rows)
	vecs := make([]linalg.Vector, 0, m.Rows)
	ch := SymmetricAsync(m)
	for val := range ch.Values {
		vals = append(vals, val)
		vecs = append(vecs, <-ch.Vectors)
	}
	return vals, vecs
}

// SymmetricTimeout is like Symmetric, but terminates early
// if the supplied timeout is exceeded.
//
// If the timeout expires, the eigenpairs which were already
// found are returned along with ErrTimeout.
func SymmetricTimeout(m *linalg.Matrix, t time.Duration) ([]float64, []linalg.Vector, error) {
	return SymmetricPrec(m, t, 0)
}

// SymmetricAsync is like Symmetric, but it runs in the
// background and reports eigenpairs as it finds them.
// It returns a channel of eigenvalues and eigenvectors,
// as well as a cancel channel which the caller may close
// to terminate the algorithm early.
func SymmetricAsync(m *linalg.Matrix) *EigenChan {
	return SymmetricPrecAsync(m, 0)
}

// SymmetricPrec is like SymmetricTimeout, but it attempts
// to achieve a given level of precision.
//
// Specifically, it stops making better approximations of
// each eigenvalue after a certain backwards error is
// achieved, where backwards error is defined as norm(Av-xv)
// where v is the approximate eigenvector and x is the
// approximate eigenvalue.
//
// It may be impossible to achieve the given level of
// precision. To address this, you must specify a timeout
// after which the algorithm gives up. If the timeout is
// exceeded, the eigenpairs which were found are returned
// along with ErrTimeout.
func SymmetricPrec(m *linalg.Matrix, t time.Duration,
	p float64) ([]float64, []linalg.Vector, error) {
	vals := make([]float64, 0, m.Rows)
	vecs := make([]linalg.Vector, 0, m.Rows)
	ch := SymmetricPrecAsync(m, p)
	go func() {
		time.Sleep(t)
		close(ch.Cancel)
	}()
	for val := range ch.Values {
		vals = append(vals, val)
		vecs = append(vecs, <-ch.Vectors)
	}
	if len(vals) < m.Rows {
		return vals, vecs, ErrTimeout
	} else {
		return vals, vecs, nil
	}
}

// SymmetricPrecAsync is a combination of SymmetricPrec
// and SymmetricAsync.
// It allows you to compute the eigenvalues and
// eigenvectors up to a certain precision, but also
// allows you to terminate the process early.
func SymmetricPrecAsync(m *linalg.Matrix, p float64) *EigenChan {
	valChan := make(chan float64, m.Rows)
	vecChan := make(chan linalg.Vector, m.Rows)
	cancelChan := make(chan struct{}, 0)
	go func() {
		defer close(valChan)
		defer close(vecChan)
		iterator := symmetricIterator{
			matrix:       m,
			cancelChan:   cancelChan,
			eigenVectors: make([]linalg.Vector, 0, m.Rows),
			eigenValues:  make([]float64, 0, m.Rows),
			precision:    p,
		}
		for i := 0; i < m.Rows; i++ {
			if !iterator.findNextVector() {
				return
			}
			valChan <- iterator.eigenValues[i]
			vecChan <- iterator.eigenVectors[i]
		}
	}()
	return &EigenChan{
		Values:  valChan,
		Vectors: vecChan,
		Cancel:  cancelChan,
	}
}

// SymmetricFixedTime is like Symmetric, but it only
// spends a certain amount of time converging to the
// answer.
// In other words, it will find an answer with the
// appropriate precision for the given time constraint.
func SymmetricFixedTime(m *linalg.Matrix, t time.Duration) ([]float64, []linalg.Vector) {
	iterator := symmetricIterator{
		matrix:        m,
		cancelChan:    nil,
		eigenVectors:  make([]linalg.Vector, 0, m.Rows),
		eigenValues:   make([]float64, 0, m.Rows),
		fixedTime:     true,
		iterationTime: t / time.Duration(m.Rows*2),
	}
	for i := 0; i < m.Rows; i++ {
		iterator.findNextVector()
	}
	return iterator.eigenValues, iterator.eigenVectors
}

type symmetricIterator struct {
	matrix *linalg.Matrix

	eigenVectors []linalg.Vector
	eigenValues  []float64

	// cancelChan is non-nil if this computation can be
	// cancelled, and will be closed if said computation
	// is cancelled.
	cancelChan <-chan struct{}

	// precision is non-zero if backErrorCriteria should
	// be used instead of oscillationCriteria.
	precision float64

	// fixedTime is used if iteration for each eigenvector
	// should be performed for a fixed amount of time.
	fixedTime     bool
	iterationTime time.Duration
}

func (i *symmetricIterator) findNextVector() bool {
	val, vec := i.inverseIterate()
	if vec == nil {
		return false
	}
	val, vec = i.powerIterate(val, vec)
	if vec == nil {
		return false
	}
	normalizeTwoNorm(vec)
	i.eigenVectors = append(i.eigenVectors, vec)
	i.eigenValues = append(i.eigenValues, val)
	return true
}

func (i *symmetricIterator) inverseIterate() (float64, linalg.Vector) {
	epsilon := math.Nextafter(1, 2) - 1

	vec := i.randomStart()
	i.deleteProjections(vec)
	val := i.scaleFactor(vec)

	criteria := i.convergenceCriteria()
	criteria.Step(i.backError(val, vec), val, vec)

	cache := newLUCache()

	for !i.cancelled() {
		lu := cache.get(val)
		if lu == nil {
			mat := i.shiftedMatrix(val)
			lu = ludecomp.Decompose(mat)
			cache.set(val, lu)
			if lu.PivotScale() < epsilon {
				return val, vec
			}
		}
		vec = lu.Solve(vec)
		normalizeMaxElement(vec)
		i.deleteProjections(vec)
		normalizeMaxElement(vec)
		val = i.scaleFactor(vec)

		criteria.Step(i.backError(val, vec), val, vec)
		if criteria.Converging() {
			return criteria.Best()
		}
	}

	return 0, nil
}

func (i *symmetricIterator) powerIterate(val float64, vec linalg.Vector) (float64, linalg.Vector) {
	criteria := i.convergenceCriteria()
	criteria.Step(i.backError(val, vec), val, vec)

	for !i.cancelled() {
		vec = i.matrix.Mul(linalg.NewMatrixColumn(vec)).Col(0)
		normalizeMaxElement(vec)
		i.deleteProjections(vec)
		normalizeMaxElement(vec)
		val = i.scaleFactor(vec)

		criteria.Step(i.backError(val, vec), val, vec)
		if criteria.Converging() {
			return criteria.Best()
		}
	}

	return 0, nil
}

func (i *symmetricIterator) deleteProjections(vec linalg.Vector) {
	for _, eigVec := range i.eigenVectors {
		projMag := eigVec.Dot(vec)
		for i, x := range eigVec {
			vec[i] -= projMag * x
		}
	}
}

func (i *symmetricIterator) randomStart() linalg.Vector {
	res := make(linalg.Vector, i.matrix.Rows)
	for i := range res {
		res[i] = rand.Float64()*2 - 1
	}
	return res
}

func (i *symmetricIterator) scaleFactor(v linalg.Vector) float64 {
	colVec := linalg.NewMatrixColumn(v)
	return v.Dot(i.matrix.Mul(colVec).Col(0)) / v.Dot(v)
}

func (i *symmetricIterator) shiftedMatrix(s float64) *linalg.Matrix {
	mat := i.matrix.Copy()
	for j := 0; j < mat.Rows; j++ {
		mat.Set(j, j, mat.Get(j, j)-s)
	}
	return mat
}

func (i *symmetricIterator) backError(val float64, vec linalg.Vector) float64 {
	multiplied := i.matrix.Mul(linalg.NewMatrixColumn(vec))
	errorSum := kahan.NewSummer64()
	for i, x := range vec {
		productVal := multiplied.Get(i, 0)
		errorSum.Add(math.Pow(productVal-val*x, 2))
	}
	return math.Sqrt(errorSum.Sum())
}

func (i *symmetricIterator) convergenceCriteria() convergenceCriteria {
	if i.precision != 0 {
		return newBackErrorCriteria(i.precision)
	} else if i.fixedTime {
		return newTimeCriteria(i.iterationTime)
	} else {
		return newOscillationCriteria(i.matrix)
	}
}

func (i *symmetricIterator) cancelled() bool {
	if i.cancelChan == nil {
		return false
	}
	select {
	case <-i.cancelChan:
		return true
	default:
		return false
	}
}

// normalizeMaxElement normalizes the given vector using
// the infinity norm (i.e. the norm which returns the
// maximum component of the vector).
func normalizeMaxElement(v linalg.Vector) {
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

// normalizeTwoNorm normalizes the given vector using
// the standard two-norm (a.k.a. the Euclidean norm).
func normalizeTwoNorm(v linalg.Vector) {
	v.Scale(1 / math.Sqrt(v.Dot(v)))
}
