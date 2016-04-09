package mvroots

import "github.com/unixpickle/num-analysis/linalg"

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
