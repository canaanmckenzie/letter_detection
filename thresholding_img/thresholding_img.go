//apply thresholding to convert grayscale to binary img

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

const THRESHOLD = 200

func main() {
	//get the input file path from the command-line arguments
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

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	//create a new binary image with the same dimensions as the grayscale image
	binaryImg := image.NewGray(image.Rect(0, 0, width, height))

	//apply thresholding to the grayscale image to convert it into a binary image
	threshold := THRESHOLD

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gray := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			if gray.Y < uint8(threshold) {
				binaryImg.SetGray(x, y, color.Gray{Y: 0})
			} else {
				binaryImg.SetGray(x, y, color.Gray{Y: 225})
			}
		}
	}

	outputFilePath := filepath.Join(filepath.Dir(inputFilePath), fmt.Sprintf("%s_binary%s", filepath.Base(inputFilePath[:len(inputFilePath)-len(filepath.Ext(inputFilePath))]), filepath.Ext(inputFilePath)))

	//save binary image to a file
	outFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Print("Error creating binary output file: ",err)
		return
	}
	defer outFile.Close()
	if filepath.Ext(inputFilePath) == ".png"{
		err = png.Encode(outFile, binaryImg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if filepath.Ext(inputFilePath) == ".jpg" || filepath.Ext(inputFilePath) == ".jpeg" {
		err = jpeg.Encode(outFile,binaryImg, &jpeg.Options{Quality: 100})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Unsupported File type: %s\n",filepath.Ext(inputFilePath))
		os.Exit(1)
	}
	fmt.Printf("Binary image saved as %s\n",outputFilePath)
}
