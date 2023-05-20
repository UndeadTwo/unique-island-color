package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

var gVisited = make([]bool, 1)
var gBuckets = make(map[color.RGBA][][]image.Point)

func floodFill(strata *image.RGBA, target color.RGBA, start image.Point, replacement color.RGBA) {
	activeSet := make([]image.Point, 1)
	activeSet = append(activeSet, start)
	processedSet := make([]image.Point, 1)

	//fmt.Printf("Strata image pixel count: %v \n", strata.Bounds().Size().X*strata.Bounds().Size().Y)

	width := strata.Bounds().Size().X

	neighbours := make([]image.Point, 1)
	for len(activeSet) > 0 {
		activePoint := activeSet[0]
		if strata.At(activePoint.X, activePoint.Y) == target {
			//strata.Set(activePoint.X, activePoint.Y, replacement)
			if activePoint.Y-1 >= 0 {
				neighbours = append(neighbours, image.Point{X: activePoint.X, Y: activePoint.Y - 1})
			}
			if activePoint.Y+1 < strata.Bounds().Size().Y {
				neighbours = append(neighbours, image.Point{X: activePoint.X, Y: activePoint.Y + 1})
			}
			if activePoint.X-1 >= 0 {
				neighbours = append(neighbours, image.Point{X: activePoint.X - 1, Y: activePoint.Y})
			}
			if activePoint.X+1 < strata.Bounds().Size().X {
				neighbours = append(neighbours, image.Point{X: activePoint.X + 1, Y: activePoint.Y})
			}

			for i := 0; i < len(neighbours); i++ {
				neighbour := neighbours[i]

				if strata.At(neighbour.X, neighbour.Y).(color.RGBA) != target {
					continue
				}

				skip := false

				if gVisited[neighbour.X+(neighbour.Y*width)] {
					continue
				}

				if skip {
					continue
				}

				gVisited[neighbour.X+(neighbour.Y*width)] = true
				activeSet = append(activeSet, neighbour)
			}
		}

		activeSet = activeSet[1:]
		processedSet = append(processedSet, activePoint)
	}

	for i := 0; i < len(processedSet); i++ {
		var x, y = processedSet[i].X, processedSet[i].Y
		strata.Set(x, y, replacement)
	}
	//fmt.Printf("Visited %v pixels while flood filling. \n", len(processedSet))
}

func floodSearch(strata *image.RGBA, target color.RGBA, start image.Point) []image.Point {
	activeSet := make([]image.Point, 1)
	activeSet = append(activeSet, start)
	processedSet := make([]image.Point, 1)

	width := strata.Bounds().Size().X

	neighbours := make([]image.Point, 1)
	for len(activeSet) > 0 {
		activePoint := activeSet[0]
		if strata.At(activePoint.X, activePoint.Y) == target {
			//strata.Set(activePoint.X, activePoint.Y, replacement)
			if activePoint.Y-1 >= 0 {
				neighbours = append(neighbours, image.Point{X: activePoint.X, Y: activePoint.Y - 1})
			}
			if activePoint.Y+1 < strata.Bounds().Size().Y {
				neighbours = append(neighbours, image.Point{X: activePoint.X, Y: activePoint.Y + 1})
			}
			if activePoint.X-1 >= 0 {
				neighbours = append(neighbours, image.Point{X: activePoint.X - 1, Y: activePoint.Y})
			}
			if activePoint.X+1 < strata.Bounds().Size().X {
				neighbours = append(neighbours, image.Point{X: activePoint.X + 1, Y: activePoint.Y})
			}

			for i := 0; i < len(neighbours); i++ {
				neighbour := neighbours[i]

				if strata.At(neighbour.X, neighbour.Y).(color.RGBA) != target {
					continue
				}

				skip := false

				if gVisited[neighbour.X+(neighbour.Y*width)] {
					continue
				}

				if skip {
					continue
				}

				gVisited[neighbour.X+(neighbour.Y*width)] = true
				activeSet = append(activeSet, neighbour)
			}
		}

		activeSet = activeSet[1:]
		processedSet = append(processedSet, activePoint)
	}

	return processedSet
}

func indexAsColor(index int) color.RGBA {
	colorR := uint8((index & 0xFF0000 >> 16) % 255)
	colorG := uint8((index & 0x00FF00 >> 8) % 256)
	colorB := uint8((index & 0x0000FF) % 256)

	return color.RGBA{colorR, colorG, colorB, 255}
}

var gColorFlag = color.RGBA{0, 0, 0, 0}
var gColorSpace = make([]color.RGBA, 0x1000000)
var gColorUsed = make([]bool, 0x1000000)

func setupColorSpace() {
	for r := 0; r < 256; r++ {
		for g := 0; g < 256; g++ {
			for b := 0; b < 256; b++ {
				i := (r + 256*(b%16)) + (g+256*int(math.Floor(float64(b)/16.0)))*4096
				gColorUsed[i] = false
				gColorSpace[i] = color.RGBA{uint8(r), uint8(g), uint8(b), 255}
			}
		}
	}
}

func getColorIndex(color color.RGBA) int {
	return (int(color.R) + 256*int(color.B%16)) + (int(color.G)+256*int(math.Floor(float64(color.B)/16.0)))*4096
}

func isColorUsed(color color.RGBA) bool {
	return gColorUsed[getColorIndex(color)]
}

func setColorUsed(color color.RGBA) {
	gUsedColorCount++
	gColorSpace[getColorIndex(color)] = color
	gColorUsed[getColorIndex(color)] = true
}

var gUsedColorCount = 0
var gFailedAtColors = -1

var gColorSpacing = 3
var gColorCount = 255 / gColorSpacing

func getUniqueColor(originalColor color.RGBA) color.RGBA {
	attemptCount := 0
	if isColorUsed(originalColor) {
		rearSpiral := originalColor
		foreSpiral := originalColor
		for attemptCount < 614125 {
			if attemptCount%3 == 0 {
				rearSpiral.R -= 3
				foreSpiral.R += 3
			} else if attemptCount%3 == 1 {
				rearSpiral.G -= 3
				foreSpiral.G += 3
			} else if attemptCount%3 == 2 {
				rearSpiral.B -= 3
				foreSpiral.B += 3
			}

			if !isColorUsed(rearSpiral) {
				setColorUsed(rearSpiral)
				return rearSpiral
			} else if !isColorUsed(foreSpiral) {
				setColorUsed(foreSpiral)
				return foreSpiral
			}

			attemptCount++
		}

		return color.RGBA{0, 0, 0, 0}
	} else {
		setColorUsed(originalColor)
		return originalColor
	}
}

func mainf() {
	if len(os.Args) < 3 {
		os.Exit(0)
	}

	for i := 0; i < 16777216; i++ {
		gColorSpace[i] = gColorFlag
	}

	infile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err.Error())
	}
	defer infile.Close()

	channelColors := 24
	channelGulf := 255 / channelColors
	imageSize := int(math.Pow(2, math.Ceil(math.Log2(math.Sqrt(float64(channelColors*channelColors*channelColors))))))
	sectors := int(math.Floor(float64(imageSize) / float64(channelColors)))

	sortedImage := image.NewRGBA(image.Rect(0, 0, imageSize, imageSize))

	for r := 0; r < channelColors; r++ {
		for g := 0; g < channelColors; g++ {
			for b := 0; b < channelColors; b++ {
				sortedImage.Set(
					g+channelColors*(r%sectors),
					b+channelColors*int(math.Floor(float64(r)/float64(sectors))),
					color.RGBA{
						R: uint8(r * channelGulf),
						G: uint8(g * channelGulf),
						B: uint8(b * channelGulf),
						A: 255})
			}
		}
	}

	outfile, err := os.OpenFile(os.Args[2], os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		print(os.Args[2], '\n')
		print(err)
		panic(err.Error())
	}

	png.Encode(outfile, sortedImage)
	outfile.Close()
}

func main() {
	if len(os.Args) < 3 {
		os.Exit(0)
	}

	setupColorSpace()

	infile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err.Error())
	}
	defer infile.Close()

	src, _, err := image.Decode(infile)
	if err != nil {
		panic(err.Error())
	}

	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	gVisited = make([]bool, w*h)
	sortedImage := image.NewRGBA(bounds)
	usedImage := image.NewRGBA(image.Rect(0, 0, 4096, 4096))
	pixelList := make([]color.RGBA, w*h)

	//floodFill(src.(*image.RGBA), color.RGBA{0, 0, 0, 255}, image.Point{X: 0, Y: 0}, color.RGBA{255, 128, 128, 255})

	maxIndex := w * h
	count := 0
	threshold := float64(maxIndex) * 0.1
	for i := 0; i < w*h; i++ {
		var x, y = i % w, int(math.Floor(float64(i) / float64(w)))
		if gVisited[i] {
			continue
		}

		if float64(i) > threshold {
			threshold += float64(maxIndex) * 0.1
			fmt.Printf("Flooding at index: %v of maximum %v \n", i, maxIndex)
		}

		gVisited[i] = true
		currentColor := src.At(x, y).(color.RGBA)
		//gBuckets[currentColor] = append(gBuckets[currentColor], floodSearch(src.(*image.RGBA), currentColor, image.Point{X: x, Y: y}))
		floodFill(src.(*image.RGBA), currentColor, image.Point{X: x, Y: y}, getUniqueColor(currentColor))
		count++
	}
	fmt.Printf("Estimated island count: %v \n", count)

	//colorTable := make([]color.RGBA, 24169)

	/*for i := 0; i < 24169; i++ {
		colorFraction := int(math.Floor(float64(0xFFFFFF) * (1.0 - (float64(i) / 24169.00))))
		colorR := uint8((colorFraction & 0xFF0000 >> 16) % 255)
		colorG := uint8((colorFraction & 0x00FF00 >> 8) % 256)
		colorB := uint8((colorFraction & 0x0000FF) % 256)
		colorTable[i] = color.RGBA{R: colorR, G: colorG, B: colorB, A: 255}
	}*/

	fmt.Printf("Copying pixels for output... \n")
	for i := 0; i < w*h; i++ {
		var x, y = i % w, int(math.Floor(float64(i) / float64(w)))

		pixelList[i] = src.At(x, y).(color.RGBA)
	}
	fmt.Printf("Finished copying... \n")

	for i := 0; i < 4096*4096; i++ {
		var x, y = i % 4096, int(math.Floor(float64(i) / 4096.0))
		usedImage.Set(x, y, gColorSpace[i])
	}

	//sort.SliceStable(pixelList, func(i, j int) bool {
	//	if pixelList[i].R > pixelList[j].R {
	//		return true
	//	}
	//	if pixelList[i].G > pixelList[j].G {
	//		return true
	//	}
	//	if pixelList[i].B > pixelList[j].B {
	//		return true
	//	}
	//	return false
	//})

	for i := 0; i < w*h; i++ {
		var x, y = i % w, int(math.Floor(float64(i) / float64(w)))
		sortedImage.Set(x, y, pixelList[i])
	}

	outfile, err := os.OpenFile(os.Args[2], os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		print(os.Args[2], '\n')
		print(err)
		panic(err.Error())
	}

	outfile2, err := os.OpenFile(os.Args[3], os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		print(os.Args[3], '\n')
		print(err)
		panic(err.Error())
	}

	png.Encode(outfile, sortedImage)
	png.Encode(outfile2, usedImage)
	outfile.Close()
}
