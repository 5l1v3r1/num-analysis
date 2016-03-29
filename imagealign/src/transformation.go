package main

import "github.com/unixpickle/num-analysis/ludecomp"

type Transformation struct {
	Matrix     [4]float64
	TranslateX float64
	TranslateY float64
}

func ApproxTransformation(source, destination []Point) *Transformation {
	columnVectors := sourceColumnEquations(source)
	normalMatrix := ludecomp.NewMatrix(6)
	for i, v1 := range columnVectors {
		for j, v2 := range columnVectors {
			normalMatrix.Set(i, j, v1.Dot(v2))
		}
	}

	results := make(ludecomp.Vector, len(destination)*2)
	for i, point := range destination {
		twoI := i
		results[twoI] = point.X
		results[twoI+1] = point.Y
	}
	normalResult := make(ludecomp.Vector, 6)
	for i, vec := range columnVectors {
		normalResult[i] = vec.Dot(results)
	}

	solution := ludecomp.Decompose(normalMatrix).Solve(normalResult)
	return &Transformation{
		Matrix:     [4]float64{solution[0], solution[1], solution[3], solution[4]},
		TranslateX: solution[2],
		TranslateY: solution[5],
	}
}

func (t *Transformation) Apply(p Point) Point {
	return Point{
		X: p.X*t.Matrix[0] + p.Y*t.Matrix[1] + t.TranslateX,
		Y: p.X*t.Matrix[2] + p.Y*t.Matrix[3] + t.TranslateY,
	}
}

func sourceColumnEquations(points []Point) [6]ludecomp.Vector {
	var columnVectors [6]ludecomp.Vector
	for i := range columnVectors {
		columnVectors[i] = make(ludecomp.Vector, len(points)*2)
	}
	for i, point := range points {
		twoI := i * 2
		columnVectors[0][twoI] = point.X
		columnVectors[1][twoI] = point.Y
		columnVectors[2][twoI] = 1
		columnVectors[3][twoI+1] = point.X
		columnVectors[4][twoI+1] = point.Y
		columnVectors[5][twoI+1] = 1
	}
	return columnVectors
}
