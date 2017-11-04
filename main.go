package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Provide two command line arguments")
		fmt.Println("one for the input and one for the output file)")
		os.Exit(1)
	}

	inFileName := os.Args[1]
	outFileName := os.Args[2]

	imgfile, err := os.Open(inFileName)
	if err != nil {
		fmt.Println("File not found!")
		os.Exit(1)
	}
	defer imgfile.Close()

	img, _, err := image.Decode(imgfile)

	outimg := ColorEqualize(img)

	outimgfile, err := os.Create(outFileName)
	if err != nil {
		fmt.Println("Can not create output file!")
		os.Exit(1)
	}
	defer outimgfile.Close()

	png.Encode(outimgfile, outimg)
}

func ColorEqualize(img image.Image) image.Image {

	bounds := img.Bounds()

	Npixels := (bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y)
	outimg := image.NewRGBA(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y))

	var hist [256][3]int
	var cumhist [256][3]int

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgb := PixelRGB(img, x, y)
			for c := 0; c < 3; c++ {
				hist[rgb[c]][c]++
			}
		}
	}

	for c := 0; c < 3; c++ {
		sum := 0
		for n := 0; n < 256; n++ {
			sum += hist[n][c]
			cumhist[n][c] = sum
		}
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgb := PixelRGB(img, x, y)
			var rgb2 [3]uint8
			for c := 0; c < 3; c++ {
				rgb2[c] = uint8((255 * cumhist[rgb[c]][c]) / Npixels)
			}
			outimg.Set(x, y, color.RGBA{rgb2[0], rgb2[1], rgb2[2], 255})
		}
	}

	return outimg
}

func PixelRGB(img image.Image, x, y int) (rgb [3]uint8) {
	r, g, b, _ := img.At(x, y).RGBA()
	rgb[0] = uint8(r)
	rgb[1] = uint8(g)
	rgb[2] = uint8(b)
	return
}
