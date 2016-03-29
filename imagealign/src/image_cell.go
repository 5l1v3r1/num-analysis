package main

import (
	"fmt"
	"math"

	"github.com/gopherjs/gopherjs/js"
)

const imageSize = 300
const dotRadius = 4

var dotColors = []string{"#65bcd4", "#e36a8f", "#32be6e", "#f19e4d", "#9b59b6", "#814938"}

type Point struct {
	X float64
	Y float64
}

type ImageCell struct {
	Element       *js.Object
	Input         *js.Object
	DotCountLabel *js.Object
	Canvas        *js.Object
	Image         *js.Object

	Points []Point
}

func NewImageCell(label string) *ImageCell {
	document := js.Global.Get("document")

	cell := document.Call("createElement", "div")
	cell.Set("className", "image-picker")
	title := document.Call("createElement", "h1")
	title.Set("textContent", label)
	input := document.Call("createElement", "input")
	input.Set("type", "file")

	cell.Call("appendChild", title)
	cell.Call("appendChild", input)

	controls := document.Call("createElement", "div")
	clearButton := document.Call("createElement", "button")
	clearButton.Set("textContent", "Reset Points")
	dotCountLabel := document.Call("createElement", "label")
	controls.Call("appendChild", clearButton)
	controls.Call("appendChild", dotCountLabel)
	cell.Call("appendChild", controls)

	devicePixelRatio := js.Global.Get("window").Get("devicePixelRatio").Float()

	canvas := document.Call("createElement", "canvas")
	canvas.Set("width", imageSize*devicePixelRatio)
	canvas.Set("height", imageSize*devicePixelRatio)
	canvas.Get("style").Set("width", fmt.Sprintf("%dpx", imageSize))
	canvas.Get("style").Set("height", fmt.Sprintf("%dpx", imageSize))
	cell.Call("appendChild", canvas)

	res := &ImageCell{
		Element:       cell,
		Input:         input,
		DotCountLabel: dotCountLabel,
		Canvas:        canvas,
	}

	input.Call("addEventListener", "change", res.inputChanged)
	clearButton.Call("addEventListener", "click", res.resetPoints)
	canvas.Call("addEventListener", "click", res.addPoint)

	res.updateUI()

	return res
}

func (i *ImageCell) ScaledPoints() []Point {
	scale := i.Canvas.Get("width").Float() / imageSize
	res := make([]Point, len(i.Points))
	for i, p := range i.Points {
		res[i] = Point{p.X * scale, p.Y * scale}
	}
	return res
}

func (i *ImageCell) inputChanged(event *js.Object) {
	file := event.Get("target").Get("files").Index(0)
	if file == nil {
		return
	}
	reader := js.Global.Get("FileReader").New()
	reader.Call("addEventListener", "load", func(e *js.Object) {
		i.Image = js.Global.Get("document").Call("createElement", "img")
		i.Image.Set("src", reader.Get("result"))
		i.Image.Call("addEventListener", "load", func() {
			i.resetPoints()
		})
	})
	reader.Call("readAsDataURL", file)
}

func (i *ImageCell) resetPoints() {
	i.Points = nil
	i.updateUI()
}

func (i *ImageCell) updateUI() {
	t := fmt.Sprintf("Label count: %d", len(i.Points))
	i.DotCountLabel.Set("textContent", t)

	ctx := i.Canvas.Call("getContext", "2d")

	ctx.Call("save")
	defer ctx.Call("restore")

	scale := i.Canvas.Get("width").Float() / imageSize
	ctx.Call("scale", scale, scale)

	ctx.Set("fillStyle", "black")
	ctx.Call("fillRect", 0, 0, imageSize, imageSize)

	if i.Image != nil {
		width := i.Image.Get("width").Float()
		height := i.Image.Get("height").Float()
		if width > height {
			relHeight := (height / width) * imageSize
			ctx.Call("drawImage", i.Image, 0, (imageSize-relHeight)/2, imageSize, relHeight)
		} else {
			relWidth := (width / height) * imageSize
			ctx.Call("drawImage", i.Image, (imageSize-relWidth)/2, 0, relWidth, imageSize)
		}
	}

	for idx, point := range i.Points {
		color := dotColors[idx%len(dotColors)]
		ctx.Set("fillStyle", color)
		ctx.Call("beginPath")
		ctx.Call("arc", point.X, point.Y, dotRadius, 0, math.Pi*2, false)
		ctx.Call("closePath")
		ctx.Call("fill")
	}
}

func (i *ImageCell) addPoint(e *js.Object) {
	rect := i.Canvas.Call("getBoundingClientRect")
	left := rect.Get("left").Float()
	top := rect.Get("top").Float()
	pointLeft := e.Get("clientX").Float() - left
	pointTop := e.Get("clientY").Float() - top
	i.Points = append(i.Points, Point{pointLeft, pointTop})
	i.updateUI()
}
