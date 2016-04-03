package eigen

import (
	"time"

	"github.com/unixpickle/num-analysis/linalg"
)

type convergenceCriteria interface {
	Step(err float64, val float64, vec linalg.Vector)
	Converging() bool
	Best() (float64, linalg.Vector)
}

// backErrorCriteria is a convergenceCriteria that
// assumes a solution is converging if the backwards
// error goes below a certain value.
type backErrorCriteria struct {
	threshold float64

	first     bool
	bestError float64
	bestVal   float64
	bestVec   linalg.Vector
}

func newBackErrorCriteria(threshold float64) *backErrorCriteria {
	return &backErrorCriteria{
		threshold: threshold,
		first:     true,
	}
}

func (b *backErrorCriteria) Step(err float64, val float64, vec linalg.Vector) {
	if b.first || b.bestError > err {
		b.first = false
		b.bestError = err
		b.bestVal = val
		b.bestVec = vec
	}
}

func (b *backErrorCriteria) Converging() bool {
	return b.bestError <= b.threshold && !b.first
}

func (b *backErrorCriteria) Best() (float64, linalg.Vector) {
	return b.bestVal, b.bestVec
}

// oscillationCriteria is a convergenceCriteria that
// assumes a solution is converging if the backwards
// error oscillates or stays the same enough times.
type oscillationCriteria struct {
	first bool

	increasing       bool
	remainingChanges int
	lastError        float64

	bestError float64
	bestVal   float64
	bestVec   linalg.Vector
}

func newOscillationCriteria(m *linalg.Matrix) *oscillationCriteria {
	return &oscillationCriteria{
		first:            true,
		remainingChanges: m.Rows*2 + 1,
	}
}

func (c *oscillationCriteria) Step(err float64, val float64, vec linalg.Vector) {
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

func (c *oscillationCriteria) Converging() bool {
	return c.remainingChanges <= 0
}

func (c *oscillationCriteria) Best() (float64, linalg.Vector) {
	return c.bestVal, c.bestVec
}

// timeCriteria is a convergenceCriteria that
// claims the approximation is converging
// when a timeout is elapsed.
type timeCriteria struct {
	*backErrorCriteria
	timeout <-chan struct{}
}

func newTimeCriteria(t time.Duration) *timeCriteria {
	ch := make(chan struct{}, 0)
	go func() {
		time.Sleep(t)
		close(ch)
	}()
	return &timeCriteria{
		backErrorCriteria: newBackErrorCriteria(0),
		timeout:           ch,
	}
}

func (t *timeCriteria) Converging() bool {
	select {
	case <-t.timeout:
		return true
	default:
		return false
	}
}
