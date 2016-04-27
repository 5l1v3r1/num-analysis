package autodiff

// A DeepNum is a numeric with zero
// or more derivatives.
type DeepNum struct {
	Value float64
	Deriv *DeepNum
}

// NewDeepNum creates a DeepNum which
// represents an n-times differentiable
// constant.
func NewDeepNum(val float64, n int) *DeepNum {
	if n == 0 {
		return &DeepNum{Value: val}
	}
	return &DeepNum{Value: val, Deriv: NewDeepNum(0, n-1)}
}

// NewDeepNumVar creates a DeepNum
// which represents an n-times
// differentiable variable.
// In other words, the first derivative
// is 1 and subsequent derivates are 0.
func NewDeepNumVar(val float64, n int) *DeepNum {
	if n == 0 {
		return &DeepNum{Value: val}
	}
	return &DeepNum{Value: val, Deriv: NewDeepNum(1, n-1)}
}

func (d *DeepNum) Add(d1 *DeepNum) *DeepNum {
	sum := &DeepNum{Value: d.Value + d1.Value}
	if d.Deriv != nil && d1.Deriv != nil {
		sum.Deriv = d.Deriv.Add(d1.Deriv)
	}
	return sum
}

func (d *DeepNum) Sub(d1 *DeepNum) *DeepNum {
	diff := &DeepNum{Value: d.Value - d1.Value}
	if d.Deriv != nil && d1.Deriv != nil {
		diff.Deriv = d.Deriv.Sub(d1.Deriv)
	}
	return diff
}

func (d *DeepNum) Mul(d1 *DeepNum) *DeepNum {
	product := &DeepNum{Value: d.Value * d1.Value}
	if d.Deriv != nil && d1.Deriv != nil {
		product.Deriv = d.Mul(d1.Deriv).Add(d1.Mul(d.Deriv))
	}
	return product
}

func (d *DeepNum) Div(d1 *DeepNum) *DeepNum {
	return d.Mul(d1.Reciprocal())
}

func (d *DeepNum) Reciprocal() *DeepNum {
	recip := &DeepNum{Value: 1.0 / d.Value}
	if d.Deriv != nil {
		recip.Deriv = d.Deriv.Div(d.Mul(d).removeOneDerivative()).MulScaler(-1)
	}
	return recip
}

func (d *DeepNum) MulScaler(c float64) *DeepNum {
	res := &DeepNum{Value: d.Value * c}
	if d.Deriv != nil {
		res.Deriv = d.Deriv.MulScaler(c)
	}
	return res
}

func (d *DeepNum) AddScaler(c float64) *DeepNum {
	return &DeepNum{Value: d.Value + c, Deriv: d.Deriv}
}

func (d *DeepNum) removeOneDerivative() *DeepNum {
	res := &DeepNum{Value: d.Value}
	if d.Deriv != nil {
		if d.Deriv.Deriv == nil {
			return res
		}
		res.Deriv = d.Deriv.removeOneDerivative()
	}
	return res
}
