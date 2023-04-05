
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"path/filepath"
	"image/png"
)

func main() {
	//get the input file path from the command-line arguments - binary image
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s [input_filename]\n", os.Args[0])
		os.Exit(1)
	}
	inputFilePath := os.Args[1]

	//open the input file path
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer inputFile.Close()

	//decode the input file as an image

	img, _, err := image.Decode(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	grayImg := image.NewGray(img.Bounds())
	for x := grayImg.Bounds().Min.X; x < grayImg.Bounds().Max.X; x++ {
		for y := grayImg.Bounds().Min.Y; y < grayImg.Bounds().Max.Y; y++ {
			grayImg.Set(x, y, img.At(x, y))
		}
	}

	// Create a new grayscale image to hold the edges
	edgeImg := image.NewGray(grayImg.Bounds())

	// Apply the Sobel operator to the image
	width := grayImg.Bounds().Max.X
	height := grayImg.Bounds().Max.Y
	for x := 1; x < width-1; x++ {
		for y := 1; y < height-1; y++ {
			// Compute the Gx and Gy values using the Sobel operator - rmove placeholder 0 * values
			Gx := -1*int(grayImg.GrayAt(x-1, y-1).Y) + 0*int(grayImg.GrayAt(x, y-1).Y) + 1*int(grayImg.GrayAt(x+1, y-1).Y) +
				-2*int(grayImg.GrayAt(x-1, y).Y) + 0*int(grayImg.GrayAt(x, y).Y) + 2*int(grayImg.GrayAt(x+1, y).Y) +
				-1*int(grayImg.GrayAt(x-1, y+1).Y) + 0*int(grayImg.GrayAt(x, y+1).Y) + 1*int(grayImg.GrayAt(x+1, y+1).Y)
			Gy := -1*int(grayImg.GrayAt(x-1, y-1).Y) - 2*int(grayImg.GrayAt(x, y-1).Y) - 1*int(grayImg.GrayAt(x+1, y-1).Y) +
				0*int(grayImg.GrayAt(x-1, y).Y) + 0*int(grayImg.GrayAt(x, y).Y) + 0*int(grayImg.GrayAt(x+1, y).Y) +
				1*int(grayImg.GrayAt(x-1, y+1).Y) + 2*int(grayImg.GrayAt(x, y+1).Y) + 1*int(grayImg.GrayAt(x+1, y+1).Y)

			// Compute the magnitude of the gradient
			magnitude := math.Sqrt(float64(Gx*Gx + Gy*Gy))

			// Set the pixel value in the edge image
			if magnitude > 128 {
				edgeImg.SetGray(x, y, color.Gray{Y: 255})
			} else {
				edgeImg.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

	outputFilePath := filepath.Join(filepath.Dir(inputFilePath), fmt.Sprintf("%s_outline%s", filepath.Base(inputFilePath[:len(inputFilePath)-len(filepath.Ext(inputFilePath))]), filepath.Ext(inputFilePath)))

	//save binary image to a file
	outFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Print("Error creating binary output file: ",err)
		return
	}
	defer outFile.Close()
	if filepath.Ext(inputFilePath) == ".png"{
		err = png.Encode(outFile, edgeImg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if filepath.Ext(inputFilePath) == ".jpg" || filepath.Ext(inputFilePath) == ".jpeg" {
		err = jpeg.Encode(outFile,edgeImg, &jpeg.Options{Quality: 100})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Unsupported File type: %s\n",filepath.Ext(inputFilePath))
		os.Exit(1)
	}
	fmt.Printf("outline image saved as %s\n",outputFilePath)
}
