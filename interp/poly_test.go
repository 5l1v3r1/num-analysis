package interp

import (
	"math"
	"math/rand"
	"testing"

	"github.com/unixpickle/num-analysis/kahan"
)

const TestSmallPolyDegree = 3
const TestLargePolyDegree = 30
const TestCoefficientsPolyDegree = 5

func TestPolySmall(t *testing.T) {
	testPolyDegree(t, TestSmallPolyDegree)
}

func TestPolyLarge(t *testing.T) {
	testPolyDegree(t, TestLargePolyDegree)
}

func TestPolyCoefficients(t *testing.T) {
	polynomial := make([]float64, TestCoefficientsPolyDegree+1)
	for i := range polynomial {
		polynomial[i] = float64(i + 2)
	}

	p := NewPoly()
	for i := 0; i < TestCoefficientsPolyDegree+1; i++ {
		x := float64(i) / float64(TestCoefficientsPolyDegree+1)
		p.Add(x, evaluatePolynomial(polynomial, x))
	}

	coeffs := p.Coefficients()
	if len(coeffs) != len(polynomial) {
		t.Fatal("invalid number of coefficients:", len(coeffs))
	}

	for i, actual := range coeffs {
		expected := polynomial[i]
		if math.Abs(actual-expected) > 1e-5 {
			t.Error("expected", expected, "but got", actual, "for x ^", i)
		}
	}
}

func testPolyDegree(t *testing.T, deg int) {
	testPoly := make([]float64, deg+1)
	for i := range testPoly {
		testPoly[i] = float64(i + 2)
	}

	p := NewPoly()
	for i := 0; i < deg+1; i++ {
		xValue := float64(i) / float64(deg+1)
		p.Add(xValue, evaluatePolynomial(testPoly, xValue))
	}

	for i := 0; i < 100; i++ {
		xValue := rand.Float64()
		expected := evaluatePolynomial(testPoly, xValue)
		actual := p.Eval(xValue)
		if math.Abs((actual-expected)/expected) > 1e-5 {
			t.Fatal("expected", expected, "but got", actual, "for x", xValue)
		}
	}
}

func evaluatePolynomial(p []float64, x float64) float64 {
	sum := kahan.NewSummer64()
	term := 1.0
	for _, coeff := range p {
		sum.Add(coeff * term)
		term *= x
	}
	return sum.Sum()
}
