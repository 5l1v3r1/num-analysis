package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"math/rand"
	"os"

	"github.com/unixpickle/num-analysis/linalg"
	"github.com/unixpickle/num-analysis/mvroots"
)

const ImageSize = 1024
const NewtonIterations = 20
const Epsilon = 1e-8
const MinColorDiff = 0.05

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: newton-basins <polynomial> <output.png>\n"+
			"Example polynomial format: -3x^2 + 2 - 5x + 7x^3")
		os.Exit(1)
	}

	poly, err := ParsePolynomial(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse polynomial:", err)
		os.Exit(1)
	}

	roots := poly.Roots()
	colors := generateColors(len(roots))

	boundMag := poly.RootBound()
	fmt.Printf("Graphing from -%f to %f\n", boundMag, boundMag)
	img := image.NewRGBA(image.Rect(0, 0, ImageSize, ImageSize))

	for y := 0; y < ImageSize; y++ {
		imagPart := ((float64(y) / ImageSize) - 0.5) * boundMag * 2
		for x := 0; x < ImageSize; x++ {
			realPart := ((float64(x) / ImageSize) - 0.5) * boundMag * 2
			s := complex(realPart, imagPart)
			convergeTarget := runNewtonFromPoint(poly, s)
			if rootIdx, ok := closeRoot(roots, convergeTarget); ok {
				img.Set(x, y, colors[rootIdx])
			} else {
				img.Set(x, y, color.RGBA{0xff, 0xff, 0xff, 0xff})
			}
		}
	}

	out, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not create output file:", err)
		os.Exit(1)
	}
	defer out.Close()
	png.Encode(out, img)
}

func runNewtonFromPoint(p mvroots.Polynomial, s complex128) complex128 {
	r := real(s)
	i := imag(s)
	iterator := mvroots.NewIterator(mvroots.ComplexAdapter{p}, linalg.Vector{r, i})

	var smallestVal float64
	var bestRoot complex128

	for i := 0; i < NewtonIterations; i++ {
		iterator.Step()
		guess := iterator.Guess()
		argument := complex(guess[0], guess[1])
		funcVal := cmplx.Abs(p.Eval(argument))
		if i == 0 || funcVal < smallestVal {
			smallestVal = funcVal
			bestRoot = argument
		}
	}

	return bestRoot
}

func closeRoot(roots []complex128, value complex128) (int, bool) {
	for i, x := range roots {
		if cmplx.Abs(x-value) < Epsilon {
			return i, true
		}
	}
	return 0, false
}

func generateColors(n int) []color.Color {
	colors := make([]color.Color, n+2)
	colors[0] = color.RGBA{}
	colors[1] = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	for i := range colors {
	ColorLoop:
		for {
			colors[i] = color.RGBA{
				R: uint8(rand.Intn(256)),
				G: uint8(rand.Intn(256)),
				B: uint8(rand.Intn(256)),
				A: uint8(rand.Intn(256)),
			}
			for _, x := range colors[:i] {
				if colorDiff(x, colors[i]) < MinColorDiff {
					continue ColorLoop
				}
			}
			break ColorLoop
		}
	}
	return colors
}

func colorDiff(c1, c2 color.Color) float64 {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	x1 := []uint32{r1, g1, b1, a1}
	x2 := []uint32{r2, g2, b2, a2}
	var diff uint32
	for i := 0; i < 4; i++ {
		if x1[i] < x2[i] {
			diff += x2[i] - x1[i]
		} else {
			diff += x1[i] - x2[i]
		}
	}
	return float64(diff) / float64(0xffff) * 4
}
