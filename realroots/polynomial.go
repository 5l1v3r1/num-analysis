package realroots

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

	res := p[0]

	varTerm := x
	for _, coeff := range p[1:] {
		res += varTerm * coeff
		varTerm *= x
	}

	return res
}
