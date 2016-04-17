package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/unixpickle/num-analysis/conjgrad"
)

type BlurGen func(w, h int) conjgrad.LinTran

var Blurers = map[string]BlurGen{
	"neighbor": func(w, h int) conjgrad.LinTran {
		return NeighborBlur{w, h}
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
	// TODO: this.
	return errors.New("unblur not yet implemented")
}
