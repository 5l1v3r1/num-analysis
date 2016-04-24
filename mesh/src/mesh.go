package main

type Vec2 struct {
	X float64
	Y float64
}

type MeshNode struct {
	Neighbors []*MeshNode
	Fixed     bool

	Position     Vec2
	RestingForce Vec2
}

func NewMeshGrid(size int) []*MeshNode {
	res := make([]*MeshNode, 0, size*size)
	for i := 0; i < size; i++ {
		iEdge := (i == 0) || (i == size-1)
		for j := 0; j < size; j++ {
			jEdge := (j == 0) || (j == size-1)
			count := 2
			if !iEdge {
				count++
			}
			if !jEdge {
				count++
			}
			res = append(res, &MeshNode{
				Neighbors: make([]*MeshNode, 0, count),
				Fixed:     jEdge || iEdge,
			})
		}
	}

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			node := res[i*size+j]

			posX := float64(j) / float64(size-1)
			posY := float64(i) / float64(size-1)
			node.Position = Vec2{posX, posY}

			for y := i - 1; y <= i+1; y++ {
				if y < 0 || y == size {
					continue
				}
				for x := j - 1; x <= j+1; x++ {
					if x < 0 || x == size {
						continue
					}
					neighbor := res[y*size+x]
					node.Neighbors = append(node.Neighbors, neighbor)
				}
			}
		}
	}

	return res
}
