package mvroots

import "github.com/unixpickle/num-analysis/kahan"

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
