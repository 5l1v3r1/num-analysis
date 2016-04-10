package optimization

import (
	"math"
	"testing"

	"github.com/unixpickle/num-analysis/linalg"
)

func TestGradientDescentSolvableSystem(t *testing.T) {
	matrix := &linalg.Matrix{
		Rows: 4,
		Cols: 4,
		Data: []float64{
			1, 2, 3, 4,
			5, 6, 7, 8,
			3, 4, 1, 2,
			8, 9, 10, 12,
		},
	}
	product := linalg.Vector{400, 300, 20, -30.5}
	sys := NewLinSysFunc(matrix, product)
	actual := GradientDescent(sys, 1e-8)
	expected := linalg.Vector{-623, 515.5, 338, -255.5}
	diff := actual.Copy().Add(expected.Scale(-1))
	if math.Sqrt(diff.Dot(diff)) > 1e-3 {
		t.Error("expected", expected, "but got", actual)
	}
}
