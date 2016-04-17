package main

import "github.com/unixpickle/num-analysis/linalg"

// NeighborBlur is a linear transformation that
// averages pixels with their adjacent (but not
// diagonally adjacent) pixels.
type NeighborBlur struct {
	Width  int
	Height int
}

func (n NeighborBlur) Dim() int {
	return n.Width * n.Height
}

func (n NeighborBlur) Apply(v linalg.Vector) linalg.Vector {
	res := make(linalg.Vector, len(v))
	for y := 0; y < n.Height; y++ {
		yEdge := y == 0 || y == n.Height-1
		for x := 0; x < n.Width; x++ {
			xEdge := x == 0 || x == n.Width-1
			val := n.getPixel(v, x, y)
			if !xEdge && !yEdge {
				val += n.getPixel(v, x, y+1)
				val += n.getPixel(v, x, y-1)
				val += n.getPixel(v, x-1, y)
				val += n.getPixel(v, x+1, y)
				val /= 5
			} else if xEdge && yEdge {
				if x == 0 {
					val += n.getPixel(v, x+1, y)
				} else {
					val += n.getPixel(v, x-1, y)
				}
				if y == 0 {
					val += n.getPixel(v, x, y+1)
				} else {
					val += n.getPixel(v, x, y-1)
				}
				val /= 3
			} else if xEdge {
				val += n.getPixel(v, x, y-1)
				val += n.getPixel(v, x, y+1)
				if x == 0 {
					val += n.getPixel(v, x+1, y)
				} else {
					val += n.getPixel(v, x-1, y)
				}
				val /= 4
			} else if yEdge {
				val += n.getPixel(v, x-1, y)
				val += n.getPixel(v, x+1, y)
				if y == 0 {
					val += n.getPixel(v, x, y+1)
				} else {
					val += n.getPixel(v, x, y-1)
				}
				val /= 4
			}
			n.setPixel(res, x, y, val)
		}
	}
	return res
}

func (n NeighborBlur) getPixel(v linalg.Vector, x, y int) float64 {
	return v[x+y*n.Width]
}

func (n NeighborBlur) setPixel(v linalg.Vector, x, y int, f float64) {
	v[x+y*n.Width] = f
}
