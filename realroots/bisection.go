package realroots

import "math"

// Bisection approximates a real root on a
// given interval of a continuous function f,
// provided that the sign of f at i.Start
// differs from the sign of f at i.End.
//
// The steps argument specifies the number
// of bisections to performed before an
// answer is returned.
//
// If f is exactly zero at either end of the
// start interval, or at any step during the
// procedure, then the perfect root will be
// returned immediately.
func Bisection(f Func, i Interval, steps int) float64 {
	startVal := f.Eval(i.Start)
	if startVal == 0 {
		return i.Start
	}

	endVal := f.Eval(i.End)
	if endVal == 0 {
		return endVal
	}

	startPos := f.Eval(i.Start) > 0

	for k := 0; k < steps; k++ {
		x := (i.Start + i.End) / 2
		val := f.Eval(x)
		if val == 0 {
			return x
		}
		valPos := val > 0
		if valPos == startPos {
			i.Start = x
		} else {
			i.End = x
		}
	}

	return (i.Start + i.End) / 2
}

// BisectionPrec is like Bisection, but it
// runs enough steps so that, theoretically,
// the size of the narrowed-down interval is
// less than prec.
func BisectionPrec(f Func, i Interval, prec float64) float64 {
	currentSpace := i.End - i.Start
	ratio := currentSpace / prec
	stepCount := int(math.Ceil(math.Log2(ratio)))
	return Bisection(f, i, stepCount)
}
