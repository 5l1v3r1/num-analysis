package ludecomp

import "testing"

func TestSolve4x4(t *testing.T) {
	m := &Matrix{
		N: 4,
		V: []float64{
			1, 2, 3, 7,
			4, 5, 6, 9.5,
			7, 8, 10, 3.2,
			1.4, 1.5, 7.9, 2.1,
		},
	}
	dec := Decompose(m)

	inputs := []Vector{
		Vector{1, 2, 3, 7},
		Vector{1, 2, 5, 2},
	}
	expected := []Vector{
		Vector{2075.0 / 1589.0, -512.0 / 245.0, 7769.0 / 7612.0, 397.0 / 3423.0},
		Vector{-223.0 / 193.0, 1329.0 / 941.0, 1925.0 / 7927.0, -371.0 / 1858.0},
	}
	for i, in := range inputs {
		out := dec.Solve(in)
		if vectorDiff(out, expected[i]) > 0.0001 {
			t.Error("Test", i, "expected", expected[i], "but got", out)
		}
	}
}
