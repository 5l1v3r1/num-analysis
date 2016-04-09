package mvroots

import (
	"math/cmplx"
	"testing"
)

func TestPolynomialRootsRealCoeffs(t *testing.T) {
	poly := Polynomial{3, -1, 1, -3, 4, 5}
	roots := []complex128{
		-1.4818635,
		-0.3241806 + 0.7199509i,
		-0.3241806 - 0.7199509i,
		0.6651123 + 0.4550801i,
		0.6651123 - 0.4550801i,
	}
	generalPolynomialTest(t, poly, roots, 1e-5)
}

func TestPolynomialRootsComplexCoeffs(t *testing.T) {
	poly := Polynomial{(3 + 2i), (1 - 3i), (15 - 0.5i), -(1 + 4i), (0.1 + 5i), (0.1 + 0.2i)}
	roots := []complex128{
		-20.938812713787172 - 9.322132020458058i,
		1.674619725670770 + 0.983776748528890i,
		-0.883651066344169 - 1.382478799221257i,
		-0.182747574812801 + 0.508194313349658i,
		0.130591629273395 - 0.387360242199228i,
	}
	generalPolynomialTest(t, poly, roots, 1e-7)
}

func TestPolynomialRootsSparse(t *testing.T) {
	poly := Polynomial{7, 0, 0, -1}
	roots := []complex128{
		-0.956465591386194 + 1.656646999972302i,
		-0.956465591386194 - 1.656646999972302i,
		1.912931182772389,
	}
	generalPolynomialTest(t, poly, roots, 1e-7)
}

func generalPolynomialTest(t *testing.T, p Polynomial, roots []complex128, prec float64) {
	actual := p.Roots()
	expected := map[complex128]bool{}
	for _, x := range roots {
		expected[x] = true
	}

ActualLoop:
	for _, x := range actual {
		for exp := range expected {
			if cmplx.Abs(exp-x) < prec {
				delete(expected, exp)
				continue ActualLoop
			}
		}
		t.Error("incorrect root:", x)
	}
	if len(expected) > 0 {
		t.Error("missing roots:", expected)
	}
}
