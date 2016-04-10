package optimization

import (
	"math"
	"testing"

	"github.com/unixpickle/num-analysis/realroots"
)

func TestUnimodalPolynomial(t *testing.T) {
	poly := realroots.Polynomial{35, -17, 2}
	minimum := UnimodalMin(poly)
	if math.Abs(minimum-17.0/4.0) > 1e-8 {
		t.Error("invalid minimum:", minimum)
	}
}

func TestUnimodalNonPoly(t *testing.T) {
	minimum := UnimodalMin(xToTheX{})
	if math.Abs(minimum-1/math.E) > 1e-8 {
		t.Error("invalid minimum:", minimum)
	}
}

type xToTheX struct{}

func (_ xToTheX) Eval(x float64) float64 {
	if x < 0 {
		return -x + 1
	}
	return math.Pow(x, x)
}
