package optimization

import (
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
	actual := GradientDescent(sys, 1e-12)
	expected := linalg.Vector{-623, 515.5, 338, -255.5}
	diff := actual.Copy().Add(expected.Scale(-1))
	if diff.Dot(diff) > 1e-10 {
		t.Error("expected", expected, "but got", actual)
	}
}

func TestGradientDescentProjection(t *testing.T) {
	matrix := &linalg.Matrix{
		Rows: 4,
		Cols: 3,
		Data: []float64{
			1, 2, 3,
			5, 6, 7,
			3, 4, 1,
			8, 9, 10,
		},
	}
	product := linalg.Vector{1, 2, 3, 4}
	sys := NewLinSysFunc(matrix, product)
	actual := GradientDescent(sys, 0)
	expected := linalg.Vector{-0.300675675675677, 1.060810810810812, -0.341216216216216}
	diff := actual.Copy().Add(expected.Copy().Scale(-1))
	if diff.Dot(diff) > 1e-10 {
		t.Error("expected", expected, "but got", actual)
	}
}

func TestGradientDescentProjectionMany(t *testing.T) {
	matrix := &linalg.Matrix{
		Rows: 4,
		Cols: 4,
		Data: []float64{
			1, 2, 3, 1,
			5, 6, 7, 5,
			3, 4, 1, 3,
			8, 9, 10, 8,
		},
	}
	product := linalg.Vector{1, 2, 3, 4}
	sys := NewLinSysFunc(matrix, product)
	answer := GradientDescent(sys, 0)
	actualProduct := linalg.Vector(matrix.Mul(linalg.NewMatrixColumn(answer)).Data)
	expectedProduct := linalg.Vector{0.797297297297297, 2.472972972972972, 3, 3.729729729729727}
	diff := actualProduct.Copy().Add(expectedProduct.Copy().Scale(-1))
	if diff.Dot(diff) > 1e-10 {
		t.Error("expected", expectedProduct, "but got", actualProduct)
	}
}
