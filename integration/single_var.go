package integration

import (
	"math"

	"github.com/unixpickle/num-analysis/kahan"
	"github.com/unixpickle/num-analysis/linalg"
	"github.com/unixpickle/num-analysis/linalg/ludecomp"
)

const defaultDivisionCount = 1e7

// IntegrateDegree computes the integral of a
// f by interpolating f as a piecewise
// polynomial function and integrating the
// resulting interpolation.
//
// The spacing argument specifies the how wide
// each interpolated polynomial should be.
// The smaller the spacing, the more accurate
// the result (to a point).
func IntegrateDegree(f Func, i Interval, spacing float64, degree int) float64 {
	if i.Length() < 0 {
		return -IntegrateDegree(f, i.Reverse(), spacing, degree)
	} else if i.Length() == 0 {
		return 0
	} else if spacing == 0 {
		spacing = i.Length()
	}

	if degree == 0 {
		return midpointIntegral(f, i, spacing)
	}

	weights := polynomialIntegralWeights(spacing, degree)

	res := kahan.NewSummer64()
	firstSample := f(i.Start)

	for k := 0; true; k++ {
		endX := i.Start + spacing*float64(k+1)
		if endX >= i.End {
			break
		}

		res.Add(firstSample * weights[0])

		startX := i.Start + spacing*float64(k)
		for i := 1; i <= degree; i++ {
			x := startX + spacing*float64(i)/float64(degree)
			y := f(x)
			if i == degree {
				firstSample = y
			}
			res.Add(y * weights[i])
		}
	}

	return res.Sum()
}

// IntegrateReimann computes the integral of a
// function using trapazoidal Reimann sums.
//
// The spacing argument specifies how wide each
// approximated trapazoid should be.
// For an accurate result, spacing should be small.
func IntegrateReimann(f Func, i Interval, spacing float64) float64 {
	return IntegrateDegree(f, i, spacing, 1)
}

// Integrate approximates the integral of a
// function along an interval.
func Integrate(f Func, i Interval) float64 {
	return IntegrateReimann(f, i, i.Length()/defaultDivisionCount)
}

func midpointIntegral(f Func, i Interval, spacing float64) float64 {
	res := kahan.NewSummer64()
	for k := 0; true; k++ {
		x := i.Start + spacing/2 + spacing*float64(k)
		if x >= i.End {
			break
		}
		res.Add(spacing * f(x))
	}
	return res.Sum()
}

// polynomialIntegralWeights returns the coefficients
// for f(x0), f(x0+spacing/d), f(x0+2*spacing/d), etc.,
// where d is the degree of the interpolation used for
// an integral.
// For a d-degree interpolation, the function must be
// sampled d+1 times, so there are d+1 coefficients.
func polynomialIntegralWeights(spacing float64, degree int) []float64 {
	numSamples := degree + 1
	system := &linalg.Matrix{
		Rows: numSamples,
		Cols: numSamples,
		Data: make([]float64, numSamples*numSamples),
	}
	for d := 0; d <= degree; d++ {
		for n := 0; n < numSamples; n++ {
			numArg := float64(n) * spacing / float64(degree)
			arg := math.Pow(numArg, float64(d))
			system.Set(d, n, arg)
		}
	}
	integrals := make(linalg.Vector, numSamples)
	for d := 1; d <= degree; d++ {
		integrals[d] = 1 / float64(d+1) * math.Pow(spacing, float64(d+1))
	}
	coefficients := ludecomp.Decompose(system).Solve(integrals)

	return []float64(coefficients)
}
