package mvroots

import (
	"math"
	"math/cmplx"

	"github.com/unixpickle/num-analysis/kahan"
)

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
