package autodiff

import "math"

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

// Depth returns the number of derivatives
// stored in d.
func (d *DeepNum) Depth() int {
	res := 0
	for d.Deriv != nil {
		res++
		d = d.Deriv
	}
	return res
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

func (d *DeepNum) Pow(d1 *DeepNum) *DeepNum {
	if d.Deriv == nil || d1.Deriv == nil {
		return &DeepNum{Value: math.Pow(d.Value, d1.Value)}
	}

	lessD := d.removeOneDerivative()
	lessD1 := d1.removeOneDerivative()
	value := lessD.Pow(lessD1)

	res := &DeepNum{Value: value.Value}
	if d.Deriv != nil && d1.Deriv != nil {
		basePart := value.Mul(lessD1.Div(lessD))
		powerPart := value.Mul(lessD.Log())
		res.Deriv = d.Deriv.Mul(basePart).Add(d1.Deriv.Mul(powerPart))
	}

	return res
}

func (d *DeepNum) Reciprocal() *DeepNum {
	recip := &DeepNum{Value: 1.0 / d.Value}
	if d.Deriv != nil {
		lessD := d.removeOneDerivative()
		recip.Deriv = d.Deriv.Div(lessD.Mul(lessD)).MulScaler(-1)
	}
	return recip
}

func (d *DeepNum) Log() *DeepNum {
	res := &DeepNum{Value: math.Log(d.Value)}
	if d.Deriv != nil {
		res.Deriv = d.Deriv.Div(d)
	}
	return res
}

func (d *DeepNum) Sqrt() *DeepNum {
	if d.Deriv == nil {
		return &DeepNum{Value: math.Sqrt(d.Value)}
	}
	subSqrt := d.removeOneDerivative().Sqrt()
	res := &DeepNum{Value: subSqrt.Value}
	res.Deriv = d.Deriv.Div(subSqrt.MulScaler(2))
	return res
}

func (d *DeepNum) Sin() *DeepNum {
	res := &DeepNum{Value: math.Sin(d.Value)}
	if d.Deriv == nil {
		return res
	}
	// TODO: optimize this by looping and alternating
	// between sin, cos, -sin, -cos, since all the
	// derivatives only require two computations in
	// total.
	res.Deriv = d.Deriv.Mul(res.removeOneDerivative().Cos())
	return res
}

func (d *DeepNum) Cos() *DeepNum {
	res := &DeepNum{Value: math.Sin(d.Value)}
	if d.Deriv == nil {
		return res
	}
	// TODO: see comment in Sin().
	res.Deriv = d.Deriv.Mul(res.removeOneDerivative().Cos()).MulScaler(-1)
	return res
}

func (d *DeepNum) Exp() *DeepNum {
	if d.Deriv == nil {
		return &DeepNum{Value: math.Exp(d.Value)}
	}
	subExp := d.removeOneDerivative().Exp()
	return &DeepNum{
		Value: subExp.Value,
		Deriv: subExp.Mul(d.Deriv),
	}
}

func (d *DeepNum) PowScaler(c float64) *DeepNum {
	res := &DeepNum{Value: math.Pow(d.Value, c)}
	if d.Deriv != nil {
		res.Deriv = d.Deriv.Mul(d.removeOneDerivative().PowScaler(c - 1).MulScaler(c))
	}
	return res
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
