package qrdecomp

import "testing"

func TestGramSchmidtTall(t *testing.T) {
	testDecomposer(t, test5x3Matrix, GramSchmidt, 3)
}

func TestGramSchmidtSquare(t *testing.T) {
	testDecomposer(t, test4x4Matrix, GramSchmidt, 4)
}

func BenchmarkGramSchmidt100x100(b *testing.B) {
	mat := randomMatrix(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GramSchmidt(mat)
	}
}

func BenchmarkGramSchmidt50x50(b *testing.B) {
	mat := randomMatrix(50)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GramSchmidt(mat)
	}
}
