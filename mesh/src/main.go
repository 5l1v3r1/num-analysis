package main

import "github.com/gopherjs/gopherjs/js"

const MeshDimensions = 5

var GlobalMesh []*MeshNode

func main() {
	GlobalMesh = NewMeshGrid(MeshDimensions)
	window := js.Global.Get("window")
	window.Set("getMeshDimensions", js.MakeFunc(GetMeshDimensions))
	window.Set("getMeshCoords", js.MakeFunc(GetMeshCoords))
	window.Set("getMeshEdges", js.MakeFunc(GetMeshEdges))
	window.Set("moveMeshNode", js.MakeFunc(MoveMeshNode))
}

func GetMeshDimensions(this *js.Object, args []*js.Object) interface{} {
	return MeshDimensions
}

func GetMeshCoords(this *js.Object, args []*js.Object) interface{} {
	coords := make([]float64, len(GlobalMesh)*2)
	for i, x := range GlobalMesh {
		coords[i*2] = x.Position.X
		coords[i*2+1] = x.Position.Y
	}
	return coords
}

func GetMeshEdges(this *js.Object, args []*js.Object) interface{} {
	var edges []int

	nodeIdx := map[*MeshNode]int{}
	for i, node := range GlobalMesh {
		nodeIdx[node] = i
	}

	seenNodes := map[*MeshNode]bool{}
	for i, node := range GlobalMesh {
		seenNodes[node] = true
		for _, neighbor := range node.Neighbors {
			if seenNodes[neighbor] {
				continue
			}
			edges = append(edges, i, nodeIdx[neighbor])
		}
	}

	return edges
}

func MoveMeshNode(this *js.Object, args []*js.Object) interface{} {
	if len(args) != 3 {
		panic("expected 3 arguments")
	}
	nodeIdx := args[0].Int()
	x := args[1].Float()
	y := args[2].Float()

	node := GlobalMesh[nodeIdx]
	wasFixed := node.Fixed
	node.Fixed = true
	node.Position = Vec2{x, y}
	Equilibriate(GlobalMesh)
	node.Fixed = wasFixed

	return nil
}
