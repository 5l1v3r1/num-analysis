package main

import (
	"github.com/unixpickle/num-analysis/conjgrad"
	"github.com/unixpickle/num-analysis/linalg"
)

func Equilibriate(nodes []*MeshNode) {
	tran := newMeshLinTran(nodes)
	residual := residualForceVector(tran)
	solution := conjgrad.Solve(tran, residual)

	componentIdx := 0
	for _, x := range nodes {
		if x.Fixed {
			continue
		}
		x.Position.X = solution[componentIdx]
		x.Position.Y = solution[componentIdx+1]
		componentIdx += 2
	}

	updateRestingForces(nodes)
}

func residualForceVector(m meshLinTran) linalg.Vector {
	res := make(linalg.Vector, m.variableNodeCount*2)

	componentIdx := 0
	for _, x := range m.nodes {
		if x.Fixed {
			continue
		}
		res[componentIdx] = x.RestingForce.X
		res[componentIdx+1] = x.RestingForce.Y
		for _, neighbor := range x.Neighbors {
			if !neighbor.Fixed {
				continue
			}
			res[componentIdx] += neighbor.Position.X
			res[componentIdx+1] += neighbor.Position.Y
		}
		componentIdx += 2
	}

	return res
}

func updateRestingForces(n []*MeshNode) {
	for _, x := range n {
		if !x.Fixed {
			continue
		}
		neighborCount := float64(len(x.Neighbors))
		x.RestingForce = Vec2{x.Position.X * neighborCount,
			x.Position.Y * neighborCount}
		for _, neighbor := range x.Neighbors {
			x.RestingForce.X -= neighbor.Position.X
			x.RestingForce.Y -= neighbor.Position.Y
		}
	}
}

type meshLinTran struct {
	nodes             []*MeshNode
	variableNodeCount int
}

func newMeshLinTran(nodes []*MeshNode) meshLinTran {
	res := meshLinTran{nodes: nodes}
	for _, n := range nodes {
		if !n.Fixed {
			res.variableNodeCount++
		}
	}
	return res
}

func (m meshLinTran) Dim() int {
	return m.variableNodeCount * 2
}

func (m meshLinTran) Apply(v linalg.Vector) linalg.Vector {
	res := make(linalg.Vector, len(v))

	newPositions := map[*MeshNode]Vec2{}
	componentIdx := 0
	for _, x := range m.nodes {
		if x.Fixed {
			continue
		}
		newPositions[x] = Vec2{v[componentIdx], v[componentIdx+1]}
		componentIdx += 2
	}

	componentIdx = 0
	for _, x := range m.nodes {
		if x.Fixed {
			continue
		}
		pos := newPositions[x]
		res[componentIdx] = float64(len(x.Neighbors)) * pos.X
		res[componentIdx+1] = float64(len(x.Neighbors)) * pos.Y
		for _, neighbor := range x.Neighbors {
			if neighbor.Fixed {
				continue
			}
			neighborPos := newPositions[neighbor]
			res[componentIdx] -= neighborPos.X
			res[componentIdx+1] -= neighborPos.Y
		}
		componentIdx += 2
	}

	return res
}
