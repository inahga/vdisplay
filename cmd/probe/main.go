package main

import (
	"image"
	"log"

	"net/http"
	_ "net/http/pprof"

	"github.com/inahga/vdisplay/capture"
	"github.com/inahga/vdisplay/encoder"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	var encode encoder.H264
	pw, err := capture.NewPipewire()
	if err != nil {
		panic(err)
	}
	if err := pw.Start(60, image.Rectangle{}, func(img image.Image) {
		encode.Encode(img)
	}); err != nil {
		panic(err)
	}
	<-make(chan struct{})
}
