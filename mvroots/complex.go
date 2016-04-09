package mvroots

import "github.com/unixpickle/num-analysis/linalg"

// A ComplexFunc is a function which takes complex
// numbers and returns complex numbers.
type ComplexFunc interface {
	Eval(c complex128) complex128
	Derivative(c complex128) complex128
}

// A ComplexAdapter turns a ComplexFunc F into a
// Func by vectorizing its inputs and outputs.
type ComplexAdapter struct {
	F ComplexFunc
}

// Dim returns 2, since there are two real
// components to each complex number.
func (c ComplexAdapter) Dim() int {
	return 2
}

// Eval treats the input vector as a complex
// number, evaluates the ComplexFunc, and then
// returns the result as a vector.
func (c ComplexAdapter) Eval(vec linalg.Vector) linalg.Vector {
	if len(vec) != 2 {
		panic("wrong dimensionality")
	}
	x := complex(vec[0], vec[1])
	res := c.F.Eval(x)
	return linalg.Vector{real(res), imag(res)}
}

// Jacobian computes the Jacobian of the function
// using the derivative of the underlying
// ComplexFunc.
func (c ComplexAdapter) Jacobian(vec linalg.Vector) *linalg.Matrix {
	if len(vec) != 2 {
		panic("wrong dimensionality")
	}

	x := complex(vec[0], vec[1])
	res := linalg.NewMatrix(2, 2)

	// Our function is p(x), where x = a + ib.
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

	realDeriv := c.F.Derivative(x)
	res.Set(0, 0, real(realDeriv))
	res.Set(1, 0, imag(realDeriv))

	imagDeriv := realDeriv * 1i
	res.Set(0, 1, real(imagDeriv))
	res.Set(1, 1, imag(imagDeriv))

	return res
}
