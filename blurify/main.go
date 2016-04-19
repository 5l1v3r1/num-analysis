package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"time"

	"github.com/unixpickle/num-analysis/conjgrad"
	"github.com/unixpickle/num-analysis/linalg"
)

const DescentThreshold = 1e-5
const ConjGradThreshold = 1e-2
const DescentTimeout = time.Second * 5

type BlurGen func(w, h int) conjgrad.LinTran

var Blurers = map[string]BlurGen{
	"neighbor": func(w, h int) conjgrad.LinTran {
		return NeighborBlur{w, h, 2}
	},
	"gaussian": func(w, h int) conjgrad.LinTran {
		size := math.Max(float64(w), float64(h))
		variance := size / 2000
		return &GaussianBlur{w, h, int(math.Sqrt(variance)*10 + 0.5), variance}
	},
}

func main() {
	if len(os.Args) != 5 {
		dieUsage()
	}
	op := os.Args[1]
	if op != "blur" && op != "unblur" {
		dieUsage()
	}
	blurer, ok := Blurers[os.Args[2]]
	if !ok {
		dieUsage()
	}

	inFile, outFile := os.Args[3], os.Args[4]

	var err error
	if op == "blur" {
		err = applyBlur(inFile, outFile, blurer)
	} else {
		err = applyUnblur(inFile, outFile, blurer)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func dieUsage() {
	fmt.Fprintln(os.Stderr, "Usage: blurify <blur|unblur> <algo> <in_img> <out_img>")
	fmt.Fprintln(os.Stderr, "Algorithms:")
	for name := range Blurers {
		fmt.Fprintln(os.Stderr, " "+name)
	}
	os.Exit(1)
}

func applyBlur(in, out string, blurer BlurGen) error {
	img, err := ReadImage(in)
	if err != nil {
		return err
	}
	algo := blurer(img.Width, img.Height)
	img.RGB[0] = algo.Apply(img.RGB[0])
	img.RGB[1] = algo.Apply(img.RGB[1])
	img.RGB[2] = algo.Apply(img.RGB[2])
	return img.Write(out)
}

func applyUnblur(in, out string, blurer BlurGen) error {
	img, err := ReadImage(in)
	if err != nil {
		return err
	}
	algo := blurer(img.Width, img.Height)
	for i := 0; i < 3; i++ {
		img.RGB[i] = optimize(img.RGB[i], algo)
		fmt.Println("done", i+1, "out of 3")
	}
	return img.Write(out)
}

func optimize(data linalg.Vector, algo conjgrad.LinTran) linalg.Vector {
	done := make(chan struct{}, 0)
	solution := make(chan linalg.Vector, 2)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		timeout := time.After(DescentTimeout)
		d := NewDescender(data, algo)
		defer func() {
			solution <- d.Guess()
			wg.Done()
		}()
		for {
			if d.Step() < DescentThreshold*float64(len(data)) {
				break
			}
			select {
			case <-timeout:
				return
			case <-done:
				return
			default:
			}
		}
	}()

	go func() {
		defer wg.Done()
		solution <- conjgrad.SolveStoppable(algo, data, ConjGradThreshold, done)
	}()

	s := <-solution
	close(done)
	wg.Wait()
	s1 := <-solution

	for _, sol := range []linalg.Vector{s, s1} {
		for i, x := range sol {
			sol[i] = math.Min(1, math.Max(0, x))
		}
	}

	diff1 := algo.Apply(s).Scale(-1).Add(data)
	diff2 := algo.Apply(s1).Scale(-1).Add(data)
	if diff1.Dot(diff1) < diff2.Dot(diff2) {
		return s
	} else {
		return s1
	}
}
