package interp

import "github.com/unixpickle/num-analysis/kahan"

// Poly uses N-degree polynomials to interpolate
// between N data points.
type Poly struct {
	// coeffs are the coefficients of Newton
	// polynomials, going from the polynomial
	// 1 to (x-x0) to (x-x1) etc.
	coeffs []float64

	// roots are the data points, which represent
	// roots of Newton polynomials.
	roots []float64
}

// NewPoly creates a Poly interpolator which has
// not been given on any points.
// The interpolator will return a value of 0 for
// every input until it has been trained.
func NewPoly() *Poly {
	return &Poly{}
}

// Add adds a point to an interpolated polynomial
// for input x and output y.
func (p *Poly) Add(x, y float64) {
	value, product := p.evalAndProduct(x)
	p.roots = append(p.roots, x)
	p.coeffs = append(p.coeffs, (y-value)/product)
}

// Eval evaluates the interpolated function
// at a given point.
func (p *Poly) Eval(x float64) float64 {
	y, _ := p.evalAndProduct(x)
	return y
}

// Coefficients returns the coefficients of
// the interpolated polynomial's terms.
// For example, this will return [a b c]
// for a + bx + cx^2.
func (p *Poly) Coefficients() []float64 {
	basisPoly := make([]float64, len(p.coeffs))
	polySum := make([]*kahan.Summer64, len(p.coeffs))
	for i := range polySum {
		polySum[i] = kahan.NewSummer64()
	}

	if len(basisPoly) > 0 {
		basisPoly[0] = 1
	}

	for i, root := range p.roots {
		coeff := p.coeffs[i]
		for j, x := range basisPoly[:i+1] {
			polySum[j].Add(x * coeff)
		}

		if i != len(p.roots)-1 {
			for j := i + 1; j > 0; j-- {
				basisPoly[j] = basisPoly[j-1] - root*basisPoly[j]
			}
			basisPoly[0] = -root * basisPoly[0]
		}
	}

	res := make([]float64, len(polySum))
	for i, x := range polySum {
		res[i] = x.Sum()
	}
	return res
}

// evalAndProduct evaluates the function at a
// point and returns the product of all the
// terms in the polynomial
// (i.e. (x-x0)(x-x1)...).
func (p *Poly) evalAndProduct(x float64) (y, prod float64) {
	res := kahan.NewSummer64()
	product := 1.0
	for i, coeff := range p.coeffs {
		res.Add(product * coeff)
		product *= (x - p.roots[i])
	}
	return res.Sum(), product
}
