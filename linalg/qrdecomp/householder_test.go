package qrdecomp

import "testing"

func TestHouseholderTall(t *testing.T) {
	testDecomposer(t, test5x3Matrix, Householder, 3)
}

func TestHouseholderWide(t *testing.T) {
	testDecomposer(t, test5x3Matrix.Transpose(), Householder, 3)
}

func TestHouseholderSquare(t *testing.T) {
	testDecomposer(t, test4x4Matrix, Householder, 4)
}

func TestHouseholderSingular(t *testing.T) {
	testDecomposer(t, testSingularMatrix, Householder, 4)
}
