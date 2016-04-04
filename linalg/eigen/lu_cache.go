package eigen

import "github.com/unixpickle/num-analysis/linalg/ludecomp"

const luCacheSize = 4

// luCache caches a few LU decompositions for matrices
// (A - v) where v is an approximated eigenvalue.
// This is useful when an approximation is converging
// and the same few eigenvalues come up again and
// again which are equal down to machine precision.
type luCache struct {
	vals []float64
	lu   []*ludecomp.LU
}

func newLUCache() *luCache {
	return &luCache{
		vals: make([]float64, 0, luCacheSize),
		lu:   make([]*ludecomp.LU, 0, luCacheSize),
	}
}

func (l *luCache) get(val float64) *ludecomp.LU {
	for i, v := range l.vals {
		if v == val {
			return l.lu[i]
		}
	}
	return nil
}

func (l *luCache) set(val float64, lu *ludecomp.LU) {
	if len(l.vals) < luCacheSize {
		l.vals = append(l.vals, val)
		l.lu = append(l.lu, lu)
	} else {
		copy(l.vals, l.vals[1:])
		copy(l.lu, l.lu[1:])
		l.vals[len(l.vals)-1] = val
		l.lu[len(l.lu)-1] = lu
	}
}
