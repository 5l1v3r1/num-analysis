package realroots

type Func interface {
	Eval(x float64) float64
}

type Interval struct {
	Start float64
	End   float64
}
