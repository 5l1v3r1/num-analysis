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
