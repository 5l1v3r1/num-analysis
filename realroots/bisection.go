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
	b := newBisector(f, i)
	for i := 0; i < steps && !b.Done(); i++ {
		b.Step()
	}
	return b.Root()
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

type bisector struct {
	interval Interval
	function Func
	done     bool

	startPos bool
}

func newBisector(f Func, i Interval) *bisector {
	startVal := f.Eval(i.Start)
	if startVal == 0 {
		return &bisector{
			interval: Interval{i.Start, i.Start},
			function: f,
			done:     true,
		}
	}

	endVal := f.Eval(i.End)
	if endVal == 0 {
		return &bisector{
			interval: Interval{i.End, i.End},
			function: f,
			done:     true,
		}
	}

	startPos := f.Eval(i.Start) > 0
	return &bisector{
		interval: i,
		function: f,
		startPos: startPos,
	}
}

func (b *bisector) Step() {
	if b.done {
		return
	}

	x := b.Root()
	val := b.function.Eval(x)

	if val == 0 {
		b.done = true
		b.interval.Start = val
		b.interval.End = val
		return
	}

	valPos := val > 0
	if valPos == b.startPos {
		b.interval.Start = x
	} else {
		b.interval.End = x
	}
}

func (b *bisector) Done() bool {
	return b.done
}

func (b *bisector) Root() float64 {
	if b.done {
		return b.interval.Start
	}
	return (b.interval.Start + b.interval.End) / 2
}
