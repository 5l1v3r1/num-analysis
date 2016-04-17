package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"

	"github.com/unixpickle/num-analysis/linalg"
)

type Image struct {
	Width  int
	Height int
	RGB    [3]linalg.Vector
}

func ReadImage(path string) (*Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	res := &Image{
		Width:  img.Bounds().Dx(),
		Height: img.Bounds().Dy(),
	}

	for i := 0; i < 3; i++ {
		res.RGB[i] = make(linalg.Vector, res.Width*res.Height)
	}

	idx := 0
	for y := 0; y < res.Height; y++ {
		for x := 0; x < res.Width; x++ {
			pixel := img.At(x+img.Bounds().Min.X, y+img.Bounds().Min.Y)
			r, g, b, _ := pixel.RGBA()
			res.RGB[0][idx] = float64(r) / 0xffff
			res.RGB[1][idx] = float64(g) / 0xffff
			res.RGB[2][idx] = float64(b) / 0xffff
			idx++
		}
	}

	return res, nil
}

func (i *Image) Write(path string) error {
	img := image.NewRGBA(image.Rect(0, 0, i.Width, i.Height))
	idx := 0
	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			r, g, b := i.RGB[0][idx], i.RGB[1][idx], i.RGB[2][idx]
			if r > 1 {
				r = 1
			} else if r < 0 {
				r = 0
			}
			if g > 1 {
				g = 1
			} else if g < 0 {
				g = 0
			}
			if b > 1 {
				b = 1
			} else if b < 0 {
				b = 0
			}
			c := color.RGBA{
				R: uint8(int(r*0xff + 0.5)),
				G: uint8(int(g*0xff + 0.5)),
				B: uint8(int(b*0xff + 0.5)),
				A: 0xff,
			}
			img.Set(x, y, c)
			idx++
		}
	}
	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()
	return png.Encode(w, img)
}
