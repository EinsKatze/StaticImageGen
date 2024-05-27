package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"sync"
	"time"
)

var progress int

func genPixels(start int, end int, width int, img *image.NRGBA, maxProg int, wg *sync.WaitGroup) {
	defer wg.Done()
	randGen := rand.New(rand.NewSource(time.Now().Unix())) // Seed the random generator, to speed it up when calling it often. (We call it often)
	for y := start; y < end; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{
				uint8(randGen.Intn(255)),
				uint8(randGen.Intn(255)),
				uint8(randGen.Intn(255)),
				255})
		}
		progress += width
		fmt.Printf("%d/%d pixels generated.\r", progress, maxProg)
	}
}

func main() {
	var wg sync.WaitGroup
	var width int
	var height int
	var fileName string

	flag.IntVar(&width, "w", 64, "Specify width. Default: 64")                       // Console arg for width
	flag.IntVar(&height, "h", 64, "Specify height. Default: 64")                     // Console arg for height
	flag.StringVar(&fileName, "f", "img.png", "Specify file name. Default: img.png") // Console arg for file name
	flag.Parse()
	maxProg := width * height // Max Progress

	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	// Multi-threading the number generation, might reduce speed on old cpus (?), but overall should improve speed.
	// We could do this on all available threads, but in my case, too many threads reduce the speed instead of increasing it
	rowsPerWorker := height / 4
	for i := 0; i < 4; i++ {
		startY := i * rowsPerWorker
		endY := startY + rowsPerWorker
		if i == 3 { // Special case for the last worker to cover any remaining rows
			endY = height
		}
		wg.Add(1)                                            // Increment the wait group to track the number of active workers
		go genPixels(startY, endY, width, img, maxProg, &wg) // Launch a goroutine to generate pixels for the specified rows
	}

	f, err := os.Create(fileName) // Create the file, where the img data should be written to
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	wg.Wait() // Wait for all goroutines to finish before writing content to file

	enc := &png.Encoder{ // Disabling compression in the encoder to speed up encoding
		CompressionLevel: png.NoCompression,
	}
	err = enc.Encode(f, img) // Write image data to file
	if err != nil {
		fmt.Println(err)
		return
	}
}
