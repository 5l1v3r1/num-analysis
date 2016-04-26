package autodiff

import "math"

// Num is a float64-backed numeric that
// stores a value and the value's gradient
// with respect to any number of variables.
type Num struct {
	Value    float64
	Gradient []float64
}

// NewNum creates a Num given a value and
// the dimension of a gradient.
// The gradient will be set to 0, making
// the resulting Num a "constant" value.
func NewNum(val float64, gradSize int) Num {
	res := Num{
		Value:    val,
		Gradient: make([]float64, gradSize),
	}
	return res
}

// NewNumVar is like NewNum, but it
// sets one element of the gradient to 1,
// thus representing one of the variables in
// the gradient.
func NewNumVar(val float64, gradSize int, varIdx int) Num {
	res := NewNum(val, gradSize)
	res.Gradient[varIdx] = 1
	return res
}

func (n Num) Add(n1 Num) Num {
	res := NewNum(n.Value+n1.Value, len(n.Gradient))
	for i, x := range n1.Gradient {
		res.Gradient[i] = x + n.Gradient[i]
	}
	return res
}

func (n Num) Sub(n1 Num) Num {
	res := NewNum(n.Value-n1.Value, len(n.Gradient))
	for i, x := range n1.Gradient {
		res.Gradient[i] = -x + n.Gradient[i]
	}
	return res
}

func (n Num) Mul(n1 Num) Num {
	res := NewNum(n.Value*n1.Value, len(n.Gradient))
	for i, x := range n1.Gradient {
		res.Gradient[i] = x*n.Value + n1.Value*n.Gradient[i]
	}
	return res
}

func (n Num) Div(n1 Num) Num {
	return n.Mul(n1.Reciprocal())
}

func (n Num) Pow(n1 Num) Num {
	res := NewNum(math.Pow(n.Value, n1.Value), len(n.Gradient))

	basePart := res.Value * (n1.Value / n.Value)
	powerPart := res.Value * math.Log(n.Value)

	for i, x := range n1.Gradient {
		res.Gradient[i] = basePart*n.Gradient[i] + powerPart*x
	}

	return res
}

func (n Num) Reciprocal() Num {
	return n.chainRule(1/n.Value, -1.0/(n.Value*n.Value))
}

func (n Num) Sqrt() Num {
	sqrt := math.Sqrt(n.Value)
	return n.chainRule(sqrt, 1/(2*sqrt))
}

func (n Num) Sin() Num {
	return n.chainRule(math.Sin(n.Value), math.Cos(n.Value))
}

func (n Num) Cos() Num {
	return n.chainRule(math.Cos(n.Value), -math.Sin(n.Value))
}

func (n Num) Exp() Num {
	exp := math.Exp(n.Value)
	return n.chainRule(exp, exp)
}

func (n Num) PowScaler(c float64) Num {
	return n.chainRule(math.Pow(n.Value, c), c*math.Pow(n.Value, c-1))
}

func (n Num) chainRule(newVal, opDerivative float64) Num {
	res := NewNum(newVal, len(n.Gradient))
	for i, fPrime := range n.Gradient {
		res.Gradient[i] = opDerivative * fPrime
	}
	return res
}
