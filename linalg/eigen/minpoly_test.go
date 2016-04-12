package eigen

import (
	"math"
	"math/cmplx"
	"math/rand"
	"testing"

	"github.com/unixpickle/num-analysis/linalg"
)

func TestMinPoly(t *testing.T) {
	matrix := &linalg.Matrix{
		Rows: 4,
		Cols: 4,
		Data: []float64{
			0.4067747645043207, 0.6164628072332874, 0.4354664525880786, 0.6401143239610350,
			0.6672015179143060, 0.6948118476963832, 0.2866760212120346, 0.2382317902035142,
			0.1964834003600978, 0.1599047271415556, 0.5562498223962068, 0.6483002730071004,
			0.0257814064728408, 0.6624187228821636, 0.5832017567431756, 0.6189453299949805,
		},
	}
	actual := MinPoly(matrix)
	expected := []float64{0.0443497, -0.0733366, 0.82639, -2.27678, 1}
	if len(actual) != len(expected) {
		t.Error("invalid polynomial length for:", actual)
	} else {
		for i, x := range expected {
			if a := actual[i]; math.Abs(a-x) > 1e-5 {
				t.Error("expected term", i, "to be", x, "but it's", a)
			}
		}
	}
}

func TestMinEigs(t *testing.T) {
	matrix := &linalg.Matrix{
		Rows: 4,
		Cols: 4,
		Data: []float64{
			0.4067747645043207, 0.6164628072332874, 0.4354664525880786, 0.6401143239610350,
			0.6672015179143060, 0.6948118476963832, 0.2866760212120346, 0.2382317902035142,
			0.1964834003600978, 0.1599047271415556, 0.5562498223962068, 0.6483002730071004,
			0.0257814064728408, 0.6624187228821636, 0.5832017567431756, 0.6189453299949805,
		},
	}
	actual := MinEigs(matrix)
	expected := map[complex128]bool{
		1.842862536469372 + 0.000000000000000i:  true,
		0.482027937707537 + 0.000000000000000i:  true,
		-0.024054354792509 + 0.222142382259406i: true,
		-0.024054354792509 - 0.222142382259406i: true,
	}
ActualLoop:
	for _, act := range actual {
		for exp := range expected {
			if cmplx.Abs(exp-act) < 1e-6 {
				delete(expected, exp)
				continue ActualLoop
			}
		}
		t.Error("incorrect eigenvalue:", act)
	}
	for exp := range expected {
		t.Error("missing eigenvalue:", exp)
	}
}

func BenchmarkMinEigs10x10(b *testing.B) {
	matrix := linalg.NewMatrix(10, 10)
	for i := range matrix.Data {
		matrix.Data[i] = (rand.Float64() * 2) - 1
	}
	for i := 0; i < b.N; i++ {
		MinEigs(matrix)
	}
}

func BenchmarkMinEigs50x50(b *testing.B) {
	matrix := linalg.NewMatrix(50, 50)
	for i := range matrix.Data {
		matrix.Data[i] = (rand.Float64() * 2) - 1
	}
	for i := 0; i < b.N; i++ {
		MinEigs(matrix)
	}
}
