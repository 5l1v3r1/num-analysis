package main

import "github.com/unixpickle/num-analysis/linalg"

// NeighborBlur is a linear transformation that
// averages pixels with their adjacent (but not
// diagonally adjacent) pixels.
type NeighborBlur struct {
	Width        int
	Height       int
	Applications int
}

func (n NeighborBlur) Dim() int {
	return n.Width * n.Height
}

func (n NeighborBlur) Apply(input linalg.Vector) linalg.Vector {
	res := make(linalg.Vector, len(input))
	v := make(linalg.Vector, len(input))
	copy(v, input)
	for i := 0; i < n.Applications; i++ {
		for y := 0; y < n.Height; y++ {
			for x := 0; x < n.Width; x++ {
				val := n.getPixel(v, x, y)
				val += n.getPixel(v, x, y+1)
				val += n.getPixel(v, x, y-1)
				val += n.getPixel(v, x-1, y)
				val += n.getPixel(v, x+1, y)
				val /= 5
				n.setPixel(res, x, y, val)
			}
		}
		copy(v, res)
	}
	return res
}

func (n NeighborBlur) getPixel(v linalg.Vector, x, y int) float64 {
	if x < 0 {
		x += n.Width
	} else if x >= n.Width {
		x = x - n.Width
	}
	if y < 0 {
		y += n.Height
	} else if y >= n.Height {
		y = y - n.Height
	}
	return v[x+y*n.Width]
}

func (n NeighborBlur) setPixel(v linalg.Vector, x, y int, f float64) {
	v[x+y*n.Width] = f
}
