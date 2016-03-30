package cholesky

import (
	"math"
	"math/rand"
	"testing"

	"github.com/unixpickle/num-analysis/linalg"
)

func TestSolve3x3(t *testing.T) {
	mat := &linalg.Matrix{
		Rows: 3,
		Cols: 3,
		Data: []float64{
			14, 26, 17,
			26, 57, 32,
			17, 32, 25,
		},
	}
	dec := Decompose(mat)

	problems := []linalg.Vector{
		{1, 2, 3},
	}
	solutions := []linalg.Vector{
		{-222.0 / 529.0, -2.0 / 529.0, 217.0 / 529.0},
	}
	for i, problem := range problems {
		solution := dec.Solve(linalg.Vector{1, 2, 3})
		if math.IsNaN(solution[0]) {
			t.Error("NaN's in solution", solution, "for", problem)
			continue
		}
		if solutionDiff(solution, solutions[i]) > 0.000001 {
			t.Error("wrong solution for", problem, "got", solution, "expected", solutions[i])
		}
	}

}

func BenchmarkDecompose200x200(b *testing.B) {
	matrix := randMatrix(200)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decompose(matrix)
	}
}

func BenchmarkDecompose100x100(b *testing.B) {
	matrix := randMatrix(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decompose(matrix)
	}
}

func BenchmarkSolve800x800(b *testing.B) {
	dec := Decompose(randMatrix(800))
	answer := randVec(800)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dec.Solve(answer)
	}
}

func BenchmarkSolve400x400(b *testing.B) {
	dec := Decompose(randMatrix(400))
	answer := randVec(400)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dec.Solve(answer)
	}
}

func solutionDiff(s1, s2 linalg.Vector) float64 {
	var diff float64
	for i, x := range s1 {
		diff += math.Pow(x-s2[i], 2)
	}
	return math.Sqrt(diff)
}

func randMatrix(size int) *linalg.Matrix {
	vecs := make([]linalg.Vector, size)
	for i := range vecs {
		vecs[i] = randVec(size)
	}
	res := linalg.NewMatrix(size, size)
	for i, v1 := range vecs {
		for j, v2 := range vecs {
			res.Set(i, j, v1.Dot(v2))
		}
	}
	return res
}

func randVec(size int) linalg.Vector {
	res := make(linalg.Vector, size)
	for i := 0; i < size; i++ {
		res[i] = rand.Float64()
	}
	return res
}
