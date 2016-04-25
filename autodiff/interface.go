package autodiff

// Num is a generic numeric type that
// supports basic arithmetic operations.
type Num interface {
	Less(Num) bool
	Equal(Num) bool
	Greater(Num) bool

	Add(Num) Num
	Sub(Num) Num
	Mul(Num) Num
	Div(Num) Num

	Reciprocal() Num
	Sqrt() Num
	Sin() Num
	Cos() Num
	Exp() Num

	Pow(c float64) Num
	PowNum(Num) Num
}
