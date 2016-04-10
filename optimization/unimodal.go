package optimization

import "math"

// UnimodalMinPrec finds the minimum of a
// unimodal function up to a certain error
// bound, prec.
func UnimodalMinPrec(u UnimodalFunc, prec float64) float64 {
	before, after := unimodalMinBracket(u)
	return goldenSectionSearch(u, before, after, prec)
}

// UnimodalMin is like UnimodalMinPrec, but
// it tries to use a "reasonable" amount of
// precision.
func UnimodalMin(u UnimodalFunc) float64 {
	before, after := unimodalMinBracket(u)
	difference := after - before
	prec := math.Nextafter(difference, difference*2) - difference
	return goldenSectionSearch(u, before, after, prec)
}

func goldenSectionSearch(u UnimodalFunc, beforeX, afterX, prec float64) float64 {
	spacing := (3.0 - math.Sqrt(5)) / 2.0

	var x [4]float64
	var y [4]float64

	x[0], x[3] = beforeX, afterX
	gapSize := afterX - beforeX
	x[1], x[2] = x[0]+spacing*gapSize, x[3]-spacing*gapSize

	for i := 0; i < 4; i++ {
		y[i] = u.Eval(x[i])
	}

	stepCount := int(math.Ceil(math.Log(prec/gapSize) / math.Log(1-spacing)))
	for i := 0; i < stepCount; i++ {
		// If we have narrowed down the region to machine
		// precision, we cannot go any further.
		if x[1] <= x[0] || x[2] <= x[1] || x[3] <= x[2] {
			break
		}

		if y[1] < y[2] {
			x[3] = x[2]
			diff := x[3] - x[0]
			x[1] = x[0] + diff*spacing
			x[2] = x[3] - diff*spacing
			y[2], y[3] = y[1], y[2]
			y[1] = u.Eval(x[1])
		} else {
			x[0] = x[1]
			diff := x[3] - x[0]
			x[1] = x[0] + diff*spacing
			x[2] = x[3] - diff*spacing
			y[0], y[1] = y[1], y[2]
			y[2] = u.Eval(x[2])
		}
	}

	return (x[3] + x[0]) / 2
}

// unimodalMinBracket figures out an interval
// on u in which the global minimum must lie.
func unimodalMinBracket(u UnimodalFunc) (before, after float64) {
	after = 1
	for !unimodalIncreasingAfter(u, after) {
		after *= 2
		if math.IsInf(after, 0) {
			panic("function is not unimodal")
		}
	}

	before = -1
	for !unimodalIncreasingAfter(u, before) {
		before *= 2
		if math.IsInf(before, 0) {
			panic("function is not unimodal")
		}
	}

	return
}

func unimodalIncreasingAfter(u UnimodalFunc, x float64) bool {
	p2 := x * 0.75
	p3 := x * 0.5
	p2Val := u.Eval(p2)
	p3Val := u.Eval(p3)
	xVal := u.Eval(x)
	return xVal > p2Val && p2Val > p3Val
}
