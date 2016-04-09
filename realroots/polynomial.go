package realroots

import (
	"math"

	"github.com/unixpickle/num-analysis/kahan"
)

// A Polynomial represents a polynomial with
// real coefficients.
//
// For a Polynomial p, the x^0 coefficient is
// p[0], the x^1 coefficient is p[1], etc.
// The degree of p is len(p).
type Polynomial []float64

// Eval evaluates the Polynomial at a point x.
func (p Polynomial) Eval(x float64) float64 {
	if len(p) == 0 {
		return 0
	}

	s := kahan.NewSummer64()
	s.Add(p[0])

	varTerm := x
	for _, coeff := range p[1:] {
		s.Add(varTerm * coeff)
		varTerm *= x
	}

	return s.Sum()
}

// OddRoots approximates the real roots of p
// that pass through the x-axis.
// It is not guaranteed to return the even roots
// of the function, since said roots "bounce off"
// the axis. However, it may return some even
// roots if it finds them.
//
// The roots are sorted in ascending order.
func (p Polynomial) OddRoots() []float64 {
	trimmed := p.trim()
	if len(trimmed) < 2 {
		return nil
	} else if len(trimmed) == 2 {
		return []float64{-trimmed[0] / trimmed[1]}
	}

	// Even roots of the derivative correspond to inflection
	// points, so using OddRoots() still gives us all the
	// potential mins and maxes of p.
	criticalPoints := trimmed.Derivative().OddRoots()

	if len(criticalPoints) == 0 {
		if trimmed.even() {
			return nil
		} else {
			return []float64{trimmed.singleOddRoot()}
		}
	}

	evaluatedPoints := make([]float64, len(criticalPoints))
	for i, x := range criticalPoints {
		evaluatedPoints[i] = trimmed.Eval(x)
	}

	var roots []float64

	if r, ok := trimmed.initialRoot(criticalPoints[0], evaluatedPoints[0]); ok {
		roots = []float64{r}
	}

	for i := 0; i < len(criticalPoints)-1; i++ {
		val1 := evaluatedPoints[i]
		if val1 == 0 {
			roots = append(roots, criticalPoints[i])
			continue
		}
		val2 := evaluatedPoints[i+1]
		if val2 == 0 {
			continue
		}
		if (val1 < 0) != (val2 < 0) {
			p1, p2 := criticalPoints[i], criticalPoints[i+1]
			roots = append(roots, Root(trimmed, Interval{p1, p2}))
		}
	}

	k := len(criticalPoints) - 1
	if evaluatedPoints[k] == 0 {
		roots = append(roots, criticalPoints[k])
	} else if r, ok := trimmed.finalRoot(criticalPoints[k], evaluatedPoints[k]); ok {
		roots = append(roots, r)
	}

	return roots
}

// Derivative returns the derivative of p.
func (p Polynomial) Derivative() Polynomial {
	if len(p) < 2 {
		return Polynomial{}
	}
	res := make(Polynomial, len(p)-1)
	for i := range res {
		res[i] = p[i+1] * float64(i+1)
	}
	return res
}

// epsilon returns a number which is small
// relative to the coefficients of this
// polynomial.
func (p Polynomial) epsilon() float64 {
	var maxCoeff float64
	for _, x := range p {
		maxCoeff = math.Max(maxCoeff, math.Abs(x))
	}
	return math.Nextafter(maxCoeff, maxCoeff*2) - maxCoeff
}

// trim returns a polynomial which is closely
// equivalent to p, but without negligible
// leading coefficients.
func (p Polynomial) trim() Polynomial {
	res := make(Polynomial, len(p))
	copy(res, p)
	for len(res) > 0 {
		eps := p.epsilon()
		if math.Abs(res[len(res)-1]) > eps {
			break
		} else {
			res = res[:len(res)-1]
		}
	}
	return res
}

// even returns true if this polynomial's
// highest term is of an even power.
func (p Polynomial) even() bool {
	return len(p)%2 == 1
}

// singleOddRoot returns the root for odd
// polynomials which have no local maxes
// or mins.
func (p Polynomial) singleOddRoot() float64 {
	zeroVal := p.Eval(0)
	if zeroVal == 0 {
		return zeroVal
	}
	if zeroVal > 0 {
		negativePoint := p.endBehaviorLeft(-p.epsilon())
		return Root(p, Interval{negativePoint, 0})
	} else {
		positivePoint := p.endBehaviorRight(p.epsilon())
		return Root(p, Interval{0, positivePoint})
	}
}

// endBehaviorRight returns an x value after
// a given x at which the function is positive.
func (p Polynomial) endBehaviorRight(x float64) float64 {
	epsilon := math.Abs(math.Nextafter(x, x*2) - x)
	for p.Eval(x+epsilon) <= 0 {
		epsilon *= 2
	}
	return x + epsilon
}

// endBehaviorLeft returns an x value before
// a given x value at which the function is
// negative or positive, depending on the
// degree of p.
func (p Polynomial) endBehaviorLeft(x float64) float64 {
	epsilon := math.Abs(math.Nextafter(x, x*2) - x)
	if p.even() {
		for p.Eval(x-epsilon) <= 0 {
			epsilon *= 2
		}
	} else {
		for p.Eval(x-epsilon) >= 0 {
			epsilon *= 2
		}
	}
	return x - epsilon
}

// initialRoot returns the root before any of
// the critical points, if one exists.
// The x and y arguments are the coordinates of
// the leftmost critical point of the function.
func (p Polynomial) initialRoot(x, y float64) (float64, bool) {
	if (p.even() && y < 0) || (!p.even() && y > 0) {
		bracketLeft := p.endBehaviorLeft(x)
		return Root(p, Interval{bracketLeft, x}), true
	} else {
		return 0, false
	}
}

// finalRoot returns the root after any of the
// critical points, if one exists.
// This is like initialRoot, except for the right
// of the function instead of the left of it.
func (p Polynomial) finalRoot(x, y float64) (float64, bool) {
	if y >= 0 {
		return 0, false
	}
	bracketRight := p.endBehaviorRight(x)
	return Root(p, Interval{x, bracketRight}), true
}
