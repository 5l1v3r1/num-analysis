package integration

// Func is a continuous function of
// a single variable.
type Func func(x float64) float64

// Interval is a range on the number line,
// defined by the start and end points of
// the interval.
//
// The start point may be greater than the
// end point.
type Interval struct {
	Start float64
	End   float64
}

// Length returns the length of this interval.
// This will be negative if i.Start > i.End.
func (i Interval) Length() float64 {
	return i.End - i.Start
}

// Reverse returns an Interval with Start
// and End reversed.
func (i Interval) Reverse() Interval {
	return Interval{Start: i.End, End: i.Start}
}
