package qrdecomp

import "testing"

func TestGramSchmidtTall(t *testing.T) {
	testDecomposer(t, test5x3Matrix, GramSchmidt, 3)
}

func TestGramSchmidtSquare(t *testing.T) {
	testDecomposer(t, test4x4Matrix, GramSchmidt, 4)
}
