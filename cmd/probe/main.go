package main

import (
	"image"
	"image/png"
	"os"

	"github.com/inahga/vdisplay/capture"
	_ "github.com/inahga/vdisplay/capture"
	_ "github.com/inahga/vdisplay/vdisplay"
)

func main() {
	pw, err := capture.NewPipewire()
	if err != nil {
		panic(err)
	}

	pw.SetMaxFramerate(60)

	pw.Register(func(img image.Image) {
		out, err := os.Create("output.png")
		if err != nil {
			panic(err)
		}
		defer out.Close()
		png.Encode(out, img)
		os.Exit(0)
	})

	if err := pw.Start(); err != nil {
		panic(err)
	}
}
