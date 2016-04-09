package realroots

import (
	"math"
	"sort"

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
// that pass through the x-axis without
// "bouncing off" of it.
// In other words, it returns the roots of odd
// multiplicity.
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
		if len(trimmed)%2 == 1 {
			return nil
		} else {
			// TODO: use the end behavior points of the function
			// to bracket the root, then solve.
			panic("not yet implemented.")
		}
	}

	evaluatedPoints := make([]float64, len(criticalPoints))
	for i, x := range criticalPoints {
		evaluatedPoints[i] = trimmed.Eval(x)
	}

	var roots []float64
	for i := 0; i < len(criticalPoints)-1; i++ {
		val1 := evaluatedPoints[i]
		val2 := evaluatedPoints[i+1]
		if val1 == 0 || val2 == 0 {
			continue
		}
		if (val1 < 0) != (val2 < 0) {
			roots = append(roots, Root(trimmed, Interval{val1, val2}))
		}
	}

	// TODO: look for a root before the first critical point.
	// TODO: look for a root after the last critical point.
	panic("not yet implemented.")

	sort.Float64s(roots)
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
