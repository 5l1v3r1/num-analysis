package autodiff

import "math"

// GradNum is a float64-backed Num that
// stores a value and the value's gradient
// with respect to any number of variables.
type GradNum struct {
	Value    float64
	Gradient []float64
}

// NewGradNum creates a GradNum given a
// value and the dimension of the gradient.
// The gradient will be set to 0, making
// the resulting GradNum a constant value.
func NewGradNum(val float64, gradSize int) GradNum {
	res := GradNum{
		Value:    val,
		Gradient: make([]float64, gradSize),
	}
	return res
}

// NewGradNumVar is like NewGradNum, but it
// sets one element of the gradient to 1,
// thus representing one of the variables in
// the gradient.
func NewGradNumVar(val float64, gradSize int, varIdx int) GradNum {
	res := NewGradNum(val, gradSize)
	res.Gradient[varIdx] = 1
	return res
}

func (n GradNum) Less(n1 Num) bool {
	return n.Value < n1.(GradNum).Value
}

func (n GradNum) Equal(n1 Num) bool {
	return n.Value == n1.(GradNum).Value
}

func (n GradNum) Greater(n1 Num) bool {
	return n.Value > n1.(GradNum).Value
}

func (n GradNum) Add(n1 Num) Num {
	res := NewGradNum(n.Value+n1.(GradNum).Value, len(n.Gradient))
	for i, x := range n1.(GradNum).Gradient {
		res.Gradient[i] = x + n.Gradient[i]
	}
	return res
}

func (n GradNum) Sub(n1 Num) Num {
	res := NewGradNum(n.Value-n1.(GradNum).Value, len(n.Gradient))
	for i, x := range n1.(GradNum).Gradient {
		res.Gradient[i] = -x + n.Gradient[i]
	}
	return res
}

func (n GradNum) Mul(n1 Num) Num {
	gn1 := n1.(GradNum)
	res := NewGradNum(n.Value*gn1.Value, len(n.Gradient))
	for i, x := range gn1.Gradient {
		res.Gradient[i] = x*n.Value + gn1.Value*n.Gradient[i]
	}
	return res
}

func (n GradNum) Div(n1 Num) Num {
	return n.Mul(n1.Reciprocal())
}

func (n GradNum) Pow(n1 Num) Num {
	gn1 := n1.(GradNum)
	res := NewGradNum(math.Pow(n.Value, gn1.Value), len(n.Gradient))

	basePart := res.Value * (gn1.Value / n.Value)
	powerPart := res.Value * math.Log(n.Value)

	for i, x := range gn1.Gradient {
		res.Gradient[i] = basePart*n.Gradient[i] + powerPart*x
	}

	return res
}

func (n GradNum) Reciprocal() Num {
	return n.chainRule(1/n.Value, -1.0/(n.Value*n.Value))
}

func (n GradNum) Sqrt() Num {
	sqrt := math.Sqrt(n.Value)
	return n.chainRule(sqrt, 1/(2*sqrt))
}

func (n GradNum) Sin() Num {
	return n.chainRule(math.Sin(n.Value), math.Cos(n.Value))
}

func (n GradNum) Cos() Num {
	return n.chainRule(math.Cos(n.Value), -math.Sin(n.Value))
}

func (n GradNum) Exp() Num {
	exp := math.Exp(n.Value)
	return n.chainRule(exp, exp)
}

func (n GradNum) PowScaler(c float64) Num {
	return n.chainRule(math.Pow(n.Value, c), c*math.Pow(n.Value, c-1))
}

func (n GradNum) chainRule(newVal, opDerivative float64) Num {
	res := NewGradNum(newVal, len(n.Gradient))
	for i, fPrime := range n.Gradient {
		res.Gradient[i] = opDerivative * fPrime
	}
	return res
}
