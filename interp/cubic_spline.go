package interp

import (
	"fmt"
	"sort"

	"github.com/unixpickle/num-analysis/kahan"
	"github.com/unixpickle/num-analysis/linalg"
	"github.com/unixpickle/num-analysis/linalg/ludecomp"
)

// A SplineStyle determines how the shape of
// a spline interpolation is computed.
// Different styles ensure different properties
// of the resulting spline.
type SplineStyle int

// These are the valid SplineStyle values.
const (
	// StandardStyle generates the default type
	// of spline without any unusual properties.
	StandardStyle SplineStyle = iota

	// MonotoneStyle forces the generated spline
	// to retain the original data's monotonicity.
	MonotoneStyle
)

// A CubicFunc represents a cubic function
// of a single variable.
//
// The four numbers [a b c d] in a CubicFunc
// correspond to a + bx + cx^2 + dx^3.
type CubicFunc [4]float64

// Eval evaluates the cubic function at a
// point x.
func (f CubicFunc) Eval(x float64) float64 {
	sum := kahan.NewSummer64()
	prod := 1.0
	for _, coeff := range f[:] {
		sum.Add(coeff * prod)
		prod *= x
	}
	return sum.Sum()
}

// A CubicSpline is a piecewise function made up of
// cubic components, designed to have be continuous
// up to the first derivative.
type CubicSpline struct {
	style  SplineStyle
	x      []float64
	y      []float64
	slopes []float64
	funcs  []CubicFunc
}

// NewCubicSpline creates a CubicSpline with no pieces.
func NewCubicSpline(style SplineStyle) *CubicSpline {
	return &CubicSpline{style: style}
}

// Add adds a point to the cubic spline, generating
// a new piece and affecting old pieces in the process.
func (c *CubicSpline) Add(x, y float64) {
	idx := sort.SearchFloat64s(c.x, x)

	c.x = append(c.x, 0)
	copy(c.x[idx+1:], c.x[idx:])
	c.y = append(c.y, 0)
	copy(c.y[idx+1:], c.y[idx:])
	c.slopes = append(c.slopes, 0)
	copy(c.slopes[idx+1:], c.slopes[idx:])

	c.x[idx] = x
	c.y[idx] = y

	if len(c.x) > 1 {
		c.funcs = append(c.funcs, CubicFunc{})
		copy(c.funcs[:idx+1], c.funcs[:idx])
	}

	c.updateSlope(idx)
	if idx > 0 {
		c.updateSlope(idx - 1)
		c.updateFunc(idx - 1)
		if idx > 1 {
			c.updateFunc(idx - 2)
		}
	}
	if idx < len(c.slopes)-1 {
		c.updateSlope(idx)
		c.updateFunc(idx)
		if idx < len(c.slopes)-2 {
			c.updateFunc(idx + 1)
		}
	}
}

// Eval evaluates the cubic spline at a given point.
func (c *CubicSpline) Eval(x float64) float64 {
	if yCount := len(c.y); yCount == 1 {
		return c.y[0]
	} else if yCount == 0 {
		return 0
	}

	idx := sort.SearchFloat64s(c.x, x) - 1
	if idx < 0 {
		idx = 0
	} else if idx >= len(c.funcs) {
		idx = len(c.funcs) - 1
	}

	return c.funcs[idx].Eval(x)
}

func (c *CubicSpline) updateSlope(idx int) float64 {
	switch c.style {
	case StandardStyle:
		c.slopes[idx] = c.computeStandardSlope(idx)
	case MonotoneStyle:
		c.slopes[idx] = c.computeMonotoneSlope(idx)
	}
	panic(fmt.Sprintf("unknown style: %d", c.style))
}

func (c *CubicSpline) computeStandardSlope(idx int) float64 {
	if len(c.x) < 2 {
		return 0
	}
	if idx == 0 {
		return (c.y[1] - c.y[0]) / (c.x[1] - c.x[0])
	} else if last := len(c.x) - 1; idx == last {
		return (c.y[last] - c.y[last-1]) / (c.x[last] - c.x[last-1])
	}
	// TODO: simplify this formula.
	mx1 := (c.x[idx-1] + c.x[idx]) / 2
	mx2 := (c.x[idx] + c.x[idx+1]) / 2
	my1 := (c.y[idx-1] + c.y[idx]) / 2
	my2 := (c.y[idx] + c.y[idx+1]) / 2
	return (my2 - my1) / (mx2 - mx1)
}

func (c *CubicSpline) computeMonotoneSlope(idx int) float64 {
	// TODO: this.
	panic("monotone cubic splines not yet implemented.")
}

func (c *CubicSpline) updateFunc(idx int) {
	// TODO: use a closed-form solution to this
	// system to improve performance.
	x0 := c.x[idx]
	x1 := c.x[idx+1]
	system := &linalg.Matrix{
		Rows: 4,
		Cols: 4,
		Data: []float64{
			1, x0, x0 * x0, x0 * x0 * x0,
			1, x1, x1 * x1, x1 * x1 * x1,
			0, 1, 2 * x0, 3 * x0 * x0,
			0, 1, 2 * x1, 3 * x1 * x1,
		},
	}
	solutions := linalg.Vector{c.y[idx], c.y[idx+1], c.slopes[idx], c.slopes[idx+1]}
	lu := ludecomp.Decompose(system)
	solution := lu.Solve(solutions)
	c.funcs[idx] = CubicFunc{solution[0], solution[1], solution[2], solution[3]}
}
