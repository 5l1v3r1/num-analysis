package optimization

import "math"

type UnimodalFunc interface {
	Eval(x float64) float64
}

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
	prec := math.Nextafter(after-before, (after-before)*2)
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

	stepCount := int(math.Ceil(math.Log(gapSize/prec) / math.Log(1-spacing)))
	for i := 0; i < stepCount; i++ {
		if y[1] < y[2] {
			x[2], x[3] = x[1], x[2]
			y[2], y[3] = y[1], y[2]
			x[1] = x[0] + (x[3] - x[2])
			y[1] = u.Eval(x[1])
		} else {
			x[0], x[1] = x[1], x[2]
			y[0], y[1] = y[1], y[2]
			x[2] = x[3] - (x[1] - x[0])
			y[2] = u.Eval(x[2])
		}
	}

	return 0
}

// unimodalMinBracket figures out an interval
// on u in which the global minimum must lie.
func unimodalMinBracket(u UnimodalFunc) (before, after float64) {
	before = -1
	for !unimodalDecreasingBefore(u, before) {
		before *= 2
	}

	after = 1
	for !unimodalIncreasingAfter(u, after) {
		after *= 2
	}

	return
}

func unimodalIncreasingAfter(u UnimodalFunc, x float64) bool {
	p3 := x * 2
	p2 := (x + p3) / 2
	p3Val := u.Eval(p3)
	p2Val := u.Eval(p2)
	xVal := u.Eval(x)
	return xVal < p2Val && p2Val < p3Val
}

func unimodalDecreasingBefore(u UnimodalFunc, x float64) bool {
	p3 := x * 2
	p2 := (x + p3) / 2
	p3Val := u.Eval(p3)
	p2Val := u.Eval(p2)
	xVal := u.Eval(x)
	return xVal > p2Val && p2Val > p3Val
}
