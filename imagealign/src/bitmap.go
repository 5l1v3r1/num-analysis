package main

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"
)

// A Bitmap is a <canvas> element whose pixels
// can be manipulated.
type Bitmap struct {
	Canvas  *js.Object
	DataArr *js.Object
	ImgData *js.Object
	Width   int
	Height  int
}

// NewBitmap creates a new square canvas and a
// corresponding Bitmap for it.
func NewBitmap(size int) *Bitmap {
	canvas := js.Global.Get("document").Call("createElement", "canvas")
	canvas.Set("width", size)
	canvas.Set("height", size)
	canvas.Get("style").Set("width", fmt.Sprintf("%dpx", size))
	canvas.Get("style").Set("height", fmt.Sprintf("%dpx", size))
	return NewBitmapCanvas(canvas)
}

// NewBitmapCanvas wraps an existing canvas in a
// Bitmap object so that its pixels can be read
// or written.
func NewBitmapCanvas(c *js.Object) *Bitmap {
	ctx := c.Call("getContext", "2d")
	data := ctx.Call("getImageData", 0, 0, c.Get("width"), c.Get("height"))
	return &Bitmap{
		Canvas:  c,
		ImgData: data,
		DataArr: data.Get("data"),
		Width:   c.Get("width").Int(),
		Height:  c.Get("height").Int(),
	}
}

// Set changes the red, green, and blue values at a
// given coordinate in the bitmap.
//
// The color components range from 0 to 255.
//
// This will not update the actual canvas until you
// call Flush().
//
// If the coordinates are out of bounds, this has no
// effect.
func (b *Bitmap) Set(x, y int, r, g, blue int) {
	if x < 0 || y < 0 || x >= b.Width || y >= b.Height {
		return
	}
	baseIdx := 4 * (y*b.Width + x)
	b.DataArr.SetIndex(baseIdx, r)
	b.DataArr.SetIndex(baseIdx+1, g)
	b.DataArr.SetIndex(baseIdx+2, blue)
	b.DataArr.SetIndex(baseIdx+3, 0xff)
}

// Get returns the pixel at the given coordinates.
//
// If the coordinates are out of bounds, this returns
// a completely black pixel.
func (b *Bitmap) Get(x, y int) (r, g, blue int) {
	if x < 0 || y < 0 || x >= b.Width || y >= b.Height {
		return
	}
	baseIdx := 4 * (y*b.Width + x)
	r = b.DataArr.Index(baseIdx).Int()
	g = b.DataArr.Index(baseIdx + 1).Int()
	blue = b.DataArr.Index(baseIdx + 2).Int()
	return
}

// Flush updates the canvas to reflect the current
// bitmap data.
func (b *Bitmap) Flush() {
	b.Canvas.Call("getContext", "2d").Call("putImageData", b.ImgData, 0, 0)
}
