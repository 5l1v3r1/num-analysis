package realroots

import "math"

type dekker struct {
	guessX float64
	guessY float64

	lastGuessX float64
	lastGuessY float64

	counterPointX float64
	counterPointY float64

	done bool

	function Func
}

func newDekker(f Func, i Interval) *dekker {
	startY := f.Eval(i.Start)
	if startY == 0 {
		return &dekker{done: true, guessX: i.Start}
	}
	endY := f.Eval(i.End)
	if endY == 0 {
		return &dekker{done: true, guessX: i.End}
	}
	if math.Abs(startY) > math.Abs(endY) {
		startY, endY = endY, startY
		i.Start, i.End = i.End, i.Start
	}
	return &dekker{
		guessX:        i.Start,
		guessY:        startY,
		lastGuessX:    i.End,
		lastGuessY:    endY,
		counterPointX: i.End,
		counterPointY: endY,
		function:      f,
	}
}

func (d *dekker) Step() {
	if d.done {
		return
	}

	bisectionPoint := (d.guessX + d.counterPointX) / 2
	if bisectionPoint == d.guessX || bisectionPoint == d.counterPointX {
		d.done = true
		return
	}

	secantPoint, valid := d.secantRoot()
	if !valid || !valueBetween(secantPoint, bisectionPoint, d.guessX) {
		d.updateGuess(bisectionPoint)
	} else {
		d.updateGuess(secantPoint)
	}
}

func (d *dekker) Done() bool {
	return d.done
}

func (d *dekker) Root() float64 {
	return d.guessX
}

func (d *dekker) secantRoot() (root float64, valid bool) {
	dY := (d.guessY - d.lastGuessY)
	if dY == 0 {
		return 0, false
	}
	dX := (d.guessX - d.lastGuessX)
	return d.guessX - (dX/dY)*d.guessY, true
}

func (d *dekker) updateGuess(x float64) {
	d.lastGuessX = d.guessX
	d.lastGuessY = d.guessY

	val := d.function.Eval(x)
	if val == 0 {
		d.guessX = x
		d.guessY = val
		d.done = true
		return
	}

	if (val < 0) == (d.counterPointY < 0) {
		d.counterPointX = d.guessX
		d.counterPointY = d.guessY
	}

	d.guessX = x
	d.guessY = val

	if math.Abs(val) > math.Abs(d.counterPointY) {
		d.guessX, d.counterPointX = d.counterPointX, d.guessX
		d.guessY, d.counterPointY = d.counterPointY, d.guessY
	}
}

func valueBetween(x, a, b float64) bool {
	if a < b {
		return a < x && x < b
	} else {
		return b < x && x < a
	}
}
