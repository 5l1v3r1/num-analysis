package mvroots

import (
	"github.com/unixpickle/num-analysis/kahan"
	"github.com/unixpickle/num-analysis/linalg"
)

// Func is a multivariable, differentiable function
// with the same number of inputs and outputs.
type Func interface {
	// Dim returns the function's dimensionality.
	// Number of inputs = number of outputs = Dim().
	Dim() int

	// Eval evaluates the function at a given
	// input vector.
	Eval(input linalg.Vector) linalg.Vector

	// Jacobian computes the Jacobian at a given
	// input vector.
	// Each row of the Jacobian corresponds to a
	// component of the output, and each column
	// corresponds to a component of the input.
	Jacobian(input linalg.Vector) *linalg.Matrix
}

// Polynomial is a complex-valued polynomial which
// is treated as a multivariable function.
// The input and output vectors are of the form
// [real, imaginary].
type Polynomial []complex128

// Dim returns 2, since there are two real
// components to each complex number.
func (p Polynomial) Dim() int {
	return 2
}

// Eval treats the input vector as a complex
// number, evaluates the polynomial, and then
// returns the result as a vector.
func (p Polynomial) Eval(vec linalg.Vector) linalg.Vector {
	if len(vec) != 2 {
		panic("wrong dimensionality")
	}

	x := complex(vec[0], vec[1])

	xPower := complex128(1)
	s := kahan.NewComplexSummer128()
	for _, a := range p {
		s.Add(a * xPower)
		xPower *= x
	}

	sum := s.Sum()
	return linalg.Vector{real(sum), imag(sum)}
}

// Jacobian computes the partial derivatives of
// the real and complex parts of the polynomial
// with respect to the real and complex values
// of the input.
func (p Polynomial) Jacobian(vec linalg.Vector) *linalg.Matrix {
	if len(vec) != 2 {
		panic("wrong dimensionality")
	}

	x := complex(vec[0], vec[1])
	res := linalg.NewMatrix(2, 2)

	// Our polynomial is p(x), where x = a + ib.
	// Thus, we can construct g(a,b) = p(a+ib).
	// We now can compute the partials:
	// - (dg)/(da) = p'(a+ib)
	// - (dg)/(db) = i*p'(a+ib)

	// Since p(a+ib) can be split into two real
	// functions h and k like:
	// p(a+ib)=h(a,b) + i*k(a,b)
	// We know that p'(a+ib) = h'(a,b) + i*k'(a,b).
	// Since the real and complex parts of p(x) and
	// p'(x) can be split up into functions like
	// this, we know that it is okay to use real(p')
	// and imag(p') as partials.

	realDeriv := p.evaluateDerivative(x)
	res.Set(0, 0, real(realDeriv))
	res.Set(1, 0, imag(realDeriv))

	imagDeriv := realDeriv * 1i
	res.Set(0, 1, real(imagDeriv))
	res.Set(1, 1, imag(imagDeriv))

	return res
}

func (p Polynomial) evaluateDerivative(x complex128) complex128 {
	xPower := complex128(1)
	s := kahan.NewComplexSummer128()
	for i, coeff := range p[1:] {
		diffCoeff := complex(float64(i+1), 0) * coeff
		s.Add(diffCoeff * xPower)
		xPower *= x
	}
	return s.Sum()
}
