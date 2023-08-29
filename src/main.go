package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
)

func main() {
	var width int
	var height int

	flag.IntVar(&width, "w", 64, "Specify width. Default: 64")
	flag.IntVar(&height, "h", 64, "Specify height. Default: 64")
	flag.Parse()

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	progress := 0
	maxProg := width * height
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{
				uint8(rand.Intn(255)),
				uint8(rand.Intn(255)),
				uint8(rand.Intn(255)),
				255})
			progress += 1
		}
		fmt.Printf("%d/%d pixels generated.\r", progress, maxProg)
	}

	f, err := os.Create("img.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		fmt.Println(err)
		return
	}
}
