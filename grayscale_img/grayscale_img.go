//convert a colored image to grayscale

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

	//create a new grayscale image with the same dimensions as the input image
	grayImg := image.NewGray(img.Bounds())

	//set every pixel of the gray image to the r,g,b avg of input image
	for x := grayImg.Bounds().Min.X; x < grayImg.Bounds().Max.X; x++ {
		for y := grayImg.Bounds().Min.Y; y < grayImg.Bounds().Max.Y; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			grayValue := uint8((r + g + b) / 3 >> 8) //this might need to change - check SO
			grayImg.Set(x, y, color.Gray{grayValue})

		}
	}
	//create the output file path
	outputFilePath := filepath.Join(filepath.Dir(inputFilePath), fmt.Sprintf("%s_gray%s", filepath.Base(inputFilePath[:len(inputFilePath)-len(filepath.Ext(inputFilePath))]), filepath.Ext(inputFilePath)))

	//open the output file for writing
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer outputFile.Close()

	//write the grayscale image to the output file as a png or jpeg
	if filepath.Ext(inputFilePath) == ".png"{
		err = png.Encode(outputFile, grayImg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if filepath.Ext(inputFilePath) == ".jpg" || filepath.Ext(inputFilePath) == ".jpeg" {
		err = jpeg.Encode(outputFile,grayImg, &jpeg.Options{Quality: 100})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Unsupported File type: %s\n",filepath.Ext(inputFilePath))
		os.Exit(1)
	}
	fmt.Printf("Grayscale image saved as %s\n",outputFilePath)
}
