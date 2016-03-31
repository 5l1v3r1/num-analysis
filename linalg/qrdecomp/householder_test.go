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

func BenchmarkHouseholder100x100(b *testing.B) {
	mat := randomMatrix(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Householder(mat)
	}
}

func BenchmarkHouseholderReflections100x100(b *testing.B) {
	mat := randomMatrix(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HouseholderReflections(mat)
	}
}

func BenchmarkHouseholder50x50(b *testing.B) {
	mat := randomMatrix(50)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Householder(mat)
	}
}

func BenchmarkHouseholderReflections50x50(b *testing.B) {
	mat := randomMatrix(50)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HouseholderReflections(mat)
	}
}
