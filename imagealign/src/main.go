package main

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"
)

var StraightCell *ImageCell
var CrookedCell *ImageCell

const minPointCount = 6

func main() {
	js.Global.Get("window").Call("addEventListener", "load", func() {
		StraightCell = NewImageCell("Straight Image")
		CrookedCell = NewImageCell("Crooked Image")
		body := js.Global.Get("document").Get("body")
		body.Call("appendChild", StraightCell.Element)
		body.Call("appendChild", CrookedCell.Element)

		js.Global.Get("align-button").Call("addEventListener", "click", alignImage)
	})
}

func alignImage() {
	if len(StraightCell.Points) != len(CrookedCell.Points) {
		alert("mismatching number of points")
		return
	}
	if len(StraightCell.Points) < minPointCount {
		alert(fmt.Sprintf("need at least %d points", minPointCount))
	}
	trans := ApproxTransformation(StraightCell.ScaledPoints(), CrookedCell.ScaledPoints())
	crookedBitmap := NewBitmapCanvas(CrookedCell.Canvas)
	newBitmap := NewBitmap(crookedBitmap.Width)
	for x := 0; x < newBitmap.Width; x++ {
		for y := 0; y < newBitmap.Height; y++ {
			sourcePoint := trans.Apply(Point{float64(x), float64(y)})
			r, g, b := crookedBitmap.Get(int(sourcePoint.X), int(sourcePoint.Y))
			newBitmap.Set(x, y, r, g, b)
		}
	}
	newBitmap.Flush()
	js.Global.Get("document").Get("body").Call("appendChild", newBitmap.Canvas)
}

func alert(msg string) {
	js.Global.Call("alert", msg)
}
