package main

import (
	"math"

	"github.com/unixpickle/num-analysis/kahan"
	"github.com/unixpickle/num-analysis/linalg"
)

// GaussianBlur is a linear transformation that
// uses gaussian distributions to average pixels
// with their adjacent (but not diagonally
// adjacent) pixels.
type GaussianBlur struct {
	Width  int
	Height int

	Radius   int
	Variance float64
}

func (g *GaussianBlur) Dim() int {
	return g.Width * g.Height
}

func (g *GaussianBlur) Apply(input linalg.Vector) linalg.Vector {
	res := make(linalg.Vector, len(input))
	v := make(linalg.Vector, len(input))
	copy(v, input)

	exps := g.cachedExponentials()

	for i := 0; i < 2; i++ {
		for y := 0; y < g.Height; y++ {
			for x := 0; x < g.Width; x++ {
				s := kahan.NewSummer64()
				expIdx := 0
				for k := -g.Radius; k <= g.Radius; k++ {
					for j := -g.Radius; j <= g.Radius; j++ {
						px := g.getPixel(v, x+k, y+j)
						exp := exps[expIdx]
						s.Add(px * exp)
						expIdx++
					}
				}
				g.setPixel(res, x, y, s.Sum())
			}
		}
		copy(v, res)
	}

	return res
}

func (g *GaussianBlur) getPixel(v linalg.Vector, x, y int) float64 {
	if x < 0 {
		x += g.Width
	} else if x >= g.Width {
		x = x - g.Width
	}
	if y < 0 {
		y += g.Height
	} else if y >= g.Height {
		y = y - g.Height
	}
	return v[x+y*g.Width]
}

func (g *GaussianBlur) setPixel(v linalg.Vector, x, y int, f float64) {
	v[x+y*g.Width] = f
}

func (g *GaussianBlur) cachedExponentials() []float64 {
	res := make([]float64, 0, g.Radius*g.Radius)
	s := kahan.NewSummer64()
	for x := -g.Radius; x <= g.Radius; x++ {
		for y := -g.Radius; y <= g.Radius; y++ {
			distSquared := math.Pow(float64(x), 2) + math.Pow(float64(y), 2)
			exp := math.Exp(-distSquared / (2 * g.Variance))
			res = append(res, exp)
			s.Add(exp)
		}
	}
	for i, x := range res {
		res[i] = x / s.Sum()
	}
	return res
}
