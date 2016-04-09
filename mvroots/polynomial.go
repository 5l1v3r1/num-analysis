package mvroots

import (
	"math"
	"math/cmplx"
	"math/rand"

	"github.com/unixpickle/num-analysis/kahan"
	"github.com/unixpickle/num-analysis/linalg"
)

// polyRootThreshold is a scaling factor by which a
// function's value must be shrunken by Newton's
// method in order to be considered a root.
// For instance, if f(x0) = 10, then you want
// f(xn) = polyRootThreshold*10 in order for xn to
// be considered a root.
const polyRootThreshold = 1e-11

const polyRootIterationSteps = 40

// Polynomial is a ComplexFunc which represents
// a polynomial.
// The n-th element in a polynomial is the
// coefficient for the x^n term.
type Polynomial []complex128

// Eval evaluates the polynomial for a given x.
func (p Polynomial) Eval(x complex128) complex128 {
	xPower := complex128(1)
	s := kahan.NewComplexSummer128()
	for _, a := range p {
		s.Add(a * xPower)
		xPower *= x
	}

	return s.Sum()
}

// Derivative evaluates the derivative of the
// polynomial at a given x.
func (p Polynomial) Derivative(x complex128) complex128 {
	xPower := complex128(1)
	s := kahan.NewComplexSummer128()
	for i, coeff := range p[1:] {
		diffCoeff := complex(float64(i+1), 0) * coeff
		s.Add(diffCoeff * xPower)
		xPower *= x
	}
	return s.Sum()
}

// RootBound returns a "bound" for the roots of
// the polynomial. No roots will have a magnitude
// greater than the returned bound.
func (p Polynomial) RootBound() float64 {
	if len(p) < 2 {
		return 0
	}

	leadingMag := cmplx.Abs(p[len(p)-1])

	// This is the Fujiwara bound, which, believe me,
	// is very very clever.
	var maxX float64
	for i, x := range p[:len(p)-1] {
		rootNum := float64(len(p) - (i + 1))
		coeffMag := cmplx.Abs(x)
		quotient := coeffMag / leadingMag
		if i == 0 {
			quotient /= 2
		}
		x := math.Pow(quotient, 1/rootNum)
		maxX = math.Max(maxX, x)
	}

	return 2 * maxX
}

// Quotient divides this polynomial by a binomial
// (x-b) and returns the new polynomial, discarding
// the remainder.
// The binomial that this divides by is (x - b).
func (p Polynomial) Quotient(b complex128) Polynomial {
	if len(p) <= 1 {
		return Polynomial{}
	}
	leadingCoefficient := p[len(p)-1]
	res := make(Polynomial, len(p)-1)
	for i := len(p) - 1; i > 0; i-- {
		res[i-1] = leadingCoefficient
		leadingCoefficient = b*leadingCoefficient + p[i-1]
	}
	return res
}

// OneRoot approximates a single root of this
// polynomial.
func (p Polynomial) OneRoot() complex128 {
	if len(p) <= 1 {
		return cmplx.NaN()
	}
	bound := p.RootBound()

	for {
		r := rand.Float64()*bound*2 - bound
		i := rand.Float64()*bound*2 - bound

		startMagnitude := cmplx.Abs(p.Eval(complex(r, i)))
		if startMagnitude == 0 {
			return complex(r, i)
		}

		iterator := NewIterator(ComplexAdapter{p}, linalg.Vector{r, i})

		var smallestVal float64
		var bestRoot complex128

		for i := 0; i < polyRootIterationSteps; i++ {
			iterator.Step()
			guess := iterator.Guess()
			argument := complex(guess[0], guess[1])
			funcVal := cmplx.Abs(p.Eval(argument))
			if i == 0 || funcVal < smallestVal {
				smallestVal = funcVal
				bestRoot = argument
			}
		}

		if smallestVal < polyRootThreshold*startMagnitude {
			return bestRoot
		}
	}
}

// Roots returns the complex roots of this
// polynomial.
func (p Polynomial) Roots() []complex128 {
	root := p.OneRoot()
	if cmplx.IsNaN(root) {
		return nil
	}
	return append(p.Quotient(root).Roots(), root)
}
