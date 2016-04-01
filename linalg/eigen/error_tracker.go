package eigen

import "github.com/unixpickle/num-analysis/linalg"

type errorDirection int

type convergeTracker struct {
	first bool

	increasing       bool
	remainingChanges int
	lastError        float64

	bestError float64
	bestVal   float64
	bestVec   linalg.Vector
}

func newConvergeTracker(m *linalg.Matrix) *convergeTracker {
	return &convergeTracker{
		first:            true,
		remainingChanges: m.Rows*2 + 1,
	}
}

func (c *convergeTracker) Step(err float64, val float64, vec linalg.Vector) {
	if c.first {
		c.first = false
		c.bestError = err
		c.bestVal = val
		c.bestVec = vec
		c.lastError = err
		return
	}

	increasing := err > c.lastError
	if increasing != c.increasing {
		c.increasing = increasing
		c.remainingChanges--
	} else if err == c.lastError {
		c.remainingChanges--
	}
	c.lastError = err

	if err < c.bestError {
		c.bestError = err
		c.bestVal = val
		c.bestVec = vec
	}
}

func (c *convergeTracker) Converging() bool {
	return c.remainingChanges <= 0
}

func (c *convergeTracker) Best() (float64, linalg.Vector) {
	return c.bestVal, c.bestVec
}
