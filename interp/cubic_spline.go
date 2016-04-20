package interp

import (
	"fmt"
	"sort"
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

type CubicSpline struct {
	style  SplineStyle
	x      []float64
	y      []float64
	slopes []float64
	funcs  []CubicFunc
}

func NewCubicSpline(style SplineStyle) *CubicSpline {
	return &CubicSpline{style: style}
}

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
	c.updateFunc(idx)
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

func (c *CubicSpline) updateFunc(idx int) float64 {
	// TODO: do some algebra here to generate a cubic polynomial
	// that goes through two points with two slopes.
	panic("todo: this.")
}
