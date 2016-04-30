package autodiff

type DeepFunc func(n *DeepNum) *DeepNum

// PowerSeries generates the coefficients
// for a power series of a DeepFunc centered
// around a point x.
//
// The result is an array of coefficients,
// where the n-th element corresponds to the
// (x-center)^n term.
func PowerSeries(d DeepFunc, center float64, degree int) []float64 {
	res := make([]float64, degree+1)
	value := NewDeepNumVar(center, degree)
	output := d(value)

	factorial := 1.0
	for i := range res {
		res[i] = output.Value / factorial
		output = output.Deriv
		factorial *= float64(i + 1)
	}

	return res
}
