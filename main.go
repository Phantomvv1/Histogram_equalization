package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
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

	CnR := equalizeColour(coloursHistogram.R, hMid)
	CnG := equalizeColour(coloursHistogram.G, hMid)
	CnB := equalizeColour(coloursHistogram.B, hMid)

	newImg := makeNewImage(CnR, CnG, CnB, img)

	file, err := os.Create("out." + format)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	switch format {
	case "png":
		err = png.Encode(file, newImg)
		if err != nil {
			log.Println(err)
			return
		}
	case "jpg":
		err = jpeg.Encode(file, newImg, nil)
		if err != nil {
			log.Println(err)
			return
		}
	default:
		log.Println("Unsupported format")
	}
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

func equalizeColour(histogram [256]int, hMid int) [256]int {
	var Cn [256]int
	newIndex := 0     // Rver
	histogramSum := 0 // Hsum

	for i := range 256 {
		l := newIndex
		histogramSum += histogram[i]

		if histogramSum > hMid {
			histogramSum -= hMid
			newIndex++
		}

		r := newIndex
		Cn[i] = (l + r) / 2
	}

	return Cn
}

func makeNewImage(CnR, CnG, CnB [256]int, img image.Image) image.Image {
	newImg := image.NewRGBA(img.Bounds())

	// Mapping the old colours to the new colours and setting them in the new image
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			newR := uint8(CnR[r>>8])
			newG := uint8(CnG[g>>8])
			newB := uint8(CnB[b>>8])

			newImg.Set(x, y, color.RGBA{R: newR, G: newG, B: newB, A: uint8(a >> 8)})
		}
	}

	return newImg
}
