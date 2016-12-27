package conjgrad

import "github.com/unixpickle/num-analysis/linalg"

const residualUpdateFrequency = 20

// SolveStoppable solves a system of linear equations
// t*x = b for x, where t is a symmetric positive-definite
// linear operator.
//
// If precond is nil, then no preconditioning is used.
//
// The prec argument specifies a bound on the
// residual error of the solution. If the largest
// element of (Ax-b) has an absolute value less than
// prec, then the current x is returned.
//
// The cancelChan argument is a channel which you
// can close to stop the solve early.
// If the solve is cancelled, an approximate
// solution is returned.
func SolveStoppable(t, precond LinTran, b linalg.Vector, prec float64,
	cancelChan <-chan struct{}) linalg.Vector {
	if precond == nil {
		precond = identity{}
	}

	var conjVec linalg.Vector
	var residual linalg.Vector
	var solution linalg.Vector

	var lastResidualDot float64

	residual = b.Copy()
	solution = make(linalg.Vector, t.Dim())

	for i := 0; residual.MaxAbs() > prec; i++ {
		z := precond.Apply(residual)
		if i == 0 {
			conjVec = z.Copy()
			lastResidualDot = z.Dot(residual)
		} else {
			residualDot := z.Dot(residual)
			projAmount := -residualDot / lastResidualDot
			lastResidualDot = residualDot
			conjVec = z.Copy().Add(conjVec.Scale(-projAmount))
		}
		if allZero(conjVec) {
			break
		}
		optimalDistance := z.Dot(residual) / conjVec.Dot(t.Apply(conjVec))

		solution.Add(conjVec.Copy().Scale(optimalDistance))
		if i != 0 && (i%residualUpdateFrequency) == 0 {
			residual = t.Apply(solution).Scale(-1).Add(b)
		} else {
			residual.Add(t.Apply(conjVec).Scale(-optimalDistance))
		}

		select {
		case <-cancelChan:
			return solution
		default:
		}
	}

	return solution
}

// SolvePrec is like SolveStoppable, but it does not
// give you the option to cancel the solve early.
func SolvePrec(t, precond LinTran, b linalg.Vector, prec float64) linalg.Vector {
	return SolveStoppable(t, precond, b, prec, nil)
}

func allZero(v linalg.Vector) bool {
	for _, x := range v {
		if x != 0 {
			return false
		}
	}
	return true
}
