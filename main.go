package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

type RGBHistogram struct {
	R [256]int
	G [256]int
	B [256]int
}

func main() {
	// Opening an image file
	file, err := os.Open("./in.png")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	// Decoding the file into a easier to use format
	img, format, err := image.Decode(file)
	if err != nil {
		log.Println(err)
		return
	}

	// Calling the equalizeing function
	equalizeHistogram(img, format)
}

func equalizeHistogram(img image.Image, format string) {
	numberOfIntenistyLevels := 256
	totalPixels := img.Bounds().Dx() * img.Bounds().Dy()
	hMid := totalPixels / numberOfIntenistyLevels // Hmid
	coloursHistogram := countColours(img)
	newIndex := 0     // Rver
	histogramSum := 0 // Hsum
	fmt.Println(hMid)
	fmt.Println(format)
}

func countColours(img image.Image) RGBHistogram {
	result := RGBHistogram{}

	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// r >> 8 = r / 257 because the language scales the colour values to 2^16 which is the max value for a png
			result.R[r>>8]++
			result.G[g>>8]++
			result.B[b>>8]++
		}
	}

	return result
}
