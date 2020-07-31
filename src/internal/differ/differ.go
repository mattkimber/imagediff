package differ

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
)

func GetDiff(original image.Image, new image.Image) (image.Image, error) {
	if original.Bounds() != new.Bounds() {
		return nil, fmt.Errorf("image bounds do not match")
	}

	indexed := false
	var pal color.Palette

	if p, ok := original.ColorModel().(color.Palette); ok {
		pal = p
		indexed = true
	}

	var output draw.Image

	if indexed {
		output = image.NewPaletted(original.Bounds(), pal)
	} else {
		output = image.NewRGBA(original.Bounds())
	}

	for x := original.Bounds().Min.X; x < original.Bounds().Max.X; x++ {
		for y := original.Bounds().Min.Y; y < original.Bounds().Max.Y; y++ {
			if original.At(x,y) != new.At(x,y) {
				output.Set(x, y, new.At(x,y))
			}
		}
	}

	return output, nil
}


func ApplyDiff(original image.Image, diff image.Image) (image.Image, error) {
	if original.Bounds() != diff.Bounds() {
		return nil, fmt.Errorf("image bounds do not match")
	}

	indexed := false
	var pal color.Palette

	if p, ok := original.ColorModel().(color.Palette); ok {
		pal = p
		indexed = true
	}

	var output draw.Image
	var transparent color.Color

	if indexed {
		output = image.NewPaletted(original.Bounds(), pal)
		transparent = pal[0]
	} else {
		output = image.NewRGBA(original.Bounds())
		transparent = color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 0,
		}
	}

	for x := original.Bounds().Min.X; x < original.Bounds().Max.X; x++ {
		for y := original.Bounds().Min.Y; y < original.Bounds().Max.Y; y++ {
			if diff.At(x,y) != transparent {
				output.Set(x, y, diff.At(x,y))
			} else {
				output.Set(x, y, original.At(x,y))
			}
		}
	}

	return output, nil
}

func Apply(inputDir string, diffDir string, outputDir string) error{
	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return fmt.Errorf("could not stat input directory: %v", err)
	}

	for _, file := range files {
		compareFile, err := os.Stat(diffDir + "/" + file.Name())

		// Skip files which don't exist
		if err != nil {
			continue
		}

		fmt.Printf("Processing %s\n", compareFile.Name())

		input, err := ReadImageFile(inputDir + "/" + file.Name())
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		diff, err := ReadImageFile(diffDir + "/" + file.Name())
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		result, err := ApplyDiff(input, diff)
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		err = WriteImageFile(result, outputDir + "/" + file.Name())
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

	}

	return nil
}

func Generate(inputDir string, compareDir string, outputDir string) error {
	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return fmt.Errorf("could not stat input directory: %v", err)
	}

	for _, file := range files {
		compareFile, err := os.Stat(compareDir + "/" + file.Name())

		// Skip files which don't exist
		if err != nil {
			continue
		}

		fmt.Printf("Processing %s\n", compareFile.Name())

		input, err := ReadImageFile(inputDir + "/" + file.Name())
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		edited, err := ReadImageFile(compareDir + "/" + file.Name())
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		result, err := GetDiff(input, edited)
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		err = WriteImageFile(result, outputDir + "/" + file.Name())
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

	}

	return nil
}

func WriteImageFile(image image.Image, filename string) error {
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("could not open output image file: %v", err)
	}

	return png.Encode(file, image)
}

func ReadImageFile(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil, fmt.Errorf("could not open image file: %v", err)
	}

	result, _, err := image.Decode(file)
	return result, err
}