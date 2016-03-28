package main

import "github.com/gopherjs/gopherjs/js"

func main() {
	js.Global.Get("window").Call("addEventListener", "load", func() {
		straightCell := NewImageCell("Straight Image")
		crookedCell := NewImageCell("Crooked Image")
		body := js.Global.Get("document").Get("body")
		body.Call("appendChild", straightCell.Element)
		body.Call("appendChild", crookedCell.Element)
	})
}
