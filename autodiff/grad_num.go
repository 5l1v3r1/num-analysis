package autodiff

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
	panic("TODO")
}

func (n GradNum) Sub(n1 Num) Num {
	panic("TODO")
}

func (n GradNum) Mul(n1 Num) Num {
	panic("TODO")
}

func (n GradNum) Div(n1 Num) Num {
	panic("TODO")
}

func (n GradNum) Reciprocal() Num {
	panic("TODO")
}

func (n GradNum) Sqrt() Num {
	panic("TODO")
}

func (n GradNum) Sin() Num {
	panic("TODO")
}

func (n GradNum) Cos() Num {
	panic("TODO")
}

func (n GradNum) Exp() Num {
	panic("TODO")
}

func (n GradNum) Pow(c float64) Num {
	panic("TODO")
}

func (n GradNum) PowNum(n1 Num) Num {
	panic("TODO")
}
