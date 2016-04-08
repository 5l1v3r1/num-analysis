package realroots

import "math"

// defaultEpsilonScale is the number of epsilons
// off a solution to Root() can be before it is
// returned without hesitation.
const defaultEpsilonScale = 4

// RootPrec finds a root of f on the interval i.
// f must have opposite signs at i.Start and i.End.
//
// The prec argument bounds the margin of error
// of the returned solution. The possible error
// will be no greater than prec unless better
// precision is impossible using floating points.
func RootPrec(f Func, i Interval, prec float64) float64 {
	remainingSteps := bisectionSteps(i, prec)
	b := newBisector(f, i)
	d := newDekker(f, i)

	for !b.Done() && !d.Done() && remainingSteps > 0 {
		remainingSteps--
		b.Step()
		d.Step()
		if d.Bounded(prec) {
			return d.Root()
		}
	}

	if remainingSteps == 0 || b.Done() {
		return b.Root()
	} else {
		return d.Root()
	}
}

// Root is like RootPrec except that it chooses a
// reasonable default precision based on the
// magnitude of the endpoints of i.
func Root(f Func, i Interval) float64 {
	mag := math.Max(math.Abs(i.Start), math.Abs(i.End))
	if mag == 0 {
		return 0
	}
	epsilon := defaultEpsilonScale * (math.Nextafter(mag, mag*2) - mag)
	return RootPrec(f, i, epsilon)
}
