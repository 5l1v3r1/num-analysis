package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/unixpickle/num-analysis/interp"
)

type Interpolation interface {
	Add(x, y float64)
	Eval(x float64) float64
}

func main() {
	interpolateFunc := js.MakeFunc(interpolate)
	js.Global.Get("window").Set("interpolate", interpolateFunc)
}

func interpolate(this *js.Object, args []*js.Object) interface{} {
	if len(args) != 6 {
		panic("expected 6 arguments")
	}
	method := args[0].String()
	xs := args[1]
	ys := args[2]
	if xs.Length() != ys.Length() {
		panic("expected xs and ys to be of same length")
	}

	interp := makeInterpolation(method, xs, ys)

	startX := args[3].Float()
	step := args[4].Float()
	count := args[5].Int()

	res := make([]float64, count)
	for i := 0; i < count; i++ {
		x := startX + step*float64(i)
		res[i] = interp.Eval(x)
	}

	return res
}

func makeInterpolation(method string, xs *js.Object, ys *js.Object) Interpolation {
	var res Interpolation
	switch method {
	case "poly":
		res = interp.NewPoly()
	case "std-spline":
		res = interp.NewCubicSpline(interp.StandardStyle)
	case "midarc-spline":
		res = interp.NewCubicSpline(interp.MidArcStyle)
	default:
		panic("unknown method: " + method)
	}

	for i := 0; i < xs.Length(); i++ {
		res.Add(xs.Index(i).Float(), ys.Index(i).Float())
	}

	return res
}
