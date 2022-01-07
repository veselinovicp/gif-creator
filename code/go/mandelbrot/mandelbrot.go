package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/esimov/gobrot/palette"
)

var (
	colorPalette    string
	colorStep       float64
	xpos, ypos      float64
	width, height   int
	imageSmoothness int
	maxIteration    int
	escapeRadius    float64
	outputFile      string

	numOfFades = 100
	radiusStep = 1.1

	iterFactor = 1

	inProgFolder = "in_progress"
)

// var waitGroup sync.WaitGroup

func init() {
	flag.Float64Var(&colorStep, "step", 6000, "Color smooth step. Value should be greater than iteration count, otherwise the value will be adjusted to the iteration count.")
	flag.IntVar(&width, "width", 1024, "Rendered image width")
	flag.IntVar(&height, "height", 768, "Rendered image height")
	flag.Float64Var(&xpos, "xpos", -0.00275, "Point position on the real axis (defined on `x` axis)")
	flag.Float64Var(&ypos, "ypos", 0.78912, "Point position on the imaginary axis (defined on `y` axis)")
	flag.Float64Var(&escapeRadius, "radius", .125689, "Escape Radius")
	flag.IntVar(&maxIteration, "iteration", 800, "Iteration count")
	flag.IntVar(&imageSmoothness, "smoothness", 8, "The rendered mandelbrot set smoothness. For a more detailded and clear image use higher numbers. For 4xAA (AA = antialiasing) use -smoothness 4")
	flag.StringVar(&colorPalette, "palette", "Hippi", "Hippi | Plan9 | AfternoonBlue | SummerBeach | Biochimist | Fiesta")
	flag.StringVar(&outputFile, "file", "mandelbrot.png", "The rendered mandelbrot image filname")
	flag.Parse()
}

func initialise() {
	path := filepath.Join(".", inProgFolder)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Println("create in progress folder error")
		}
	}
}

func main() {
	initialise()

	width = width * imageSmoothness
	height = height * imageSmoothness

	numOfDrawn := 0
	for {
		numOfDrawn++
		if numOfDrawn > numOfFades {
			break
		}

		escapeRadius = escapeRadius / radiusStep

		maxIteration = maxIteration - iterFactor

		outputFile = fmt.Sprintf("%s/%04d.png", inProgFolder, numOfDrawn)
		createOne()
	}

	for {
		if numOfDrawn != numOfFades {
			numOfDrawn++
		}
		if numOfDrawn > numOfFades*2 {
			break
		}

		escapeRadius = escapeRadius * radiusStep

		maxIteration = maxIteration + iterFactor

		outputFile = fmt.Sprintf("%s/%04d.png", inProgFolder, numOfDrawn)

		createOne()
	}
}

func createOne() {

	if colorStep < float64(maxIteration) {
		colorStep = float64(maxIteration)
	}

	colors := interpolateColors(&colorPalette, colorStep)
	if len(colors) > 0 {

		render(maxIteration, colors)
		fmt.Printf("\n\nMandelbrot set rendered into `%s`\n", outputFile)
	}
	time.Sleep(time.Second)
}

func interpolateColors(paletteCode *string, numberOfColors float64) []color.RGBA {
	var factor float64
	steps := []float64{}
	cols := []uint32{}
	interpolated := []uint32{}
	interpolatedColors := []color.RGBA{}

	for _, v := range palette.ColorPalettes {
		factor = 1.0 / numberOfColors
		switch v.Keyword {
		case *paletteCode:
			if paletteCode != nil {
				for index, col := range v.Colors {
					if col.Step == 0.0 && index != 0 {
						stepRatio := float64(index+1) / float64(len(v.Colors))
						step := float64(int(stepRatio*100)) / 100 // truncate to 2 decimal precision
						steps = append(steps, step)
					} else {
						steps = append(steps, col.Step)
					}
					r, g, b, a := col.Color.RGBA()
					r /= 0xff
					g /= 0xff
					b /= 0xff
					a /= 0xff
					uintColor := uint32(r)<<24 | uint32(g)<<16 | uint32(b)<<8 | uint32(a)
					cols = append(cols, uintColor)
				}

				var min, max, minColor, maxColor float64
				if len(v.Colors) == len(steps) && len(v.Colors) == len(cols) {
					for i := 0.0; i <= 1; i += factor {
						for j := 0; j < len(v.Colors)-1; j++ {
							if i >= steps[j] && i < steps[j+1] {
								min = steps[j]
								max = steps[j+1]
								minColor = float64(cols[j])
								maxColor = float64(cols[j+1])
								uintColor := cosineInterpolation(maxColor, minColor, (i-min)/(max-min))
								interpolated = append(interpolated, uint32(uintColor))
							}
						}
					}
				}

				for _, pixelValue := range interpolated {
					r := pixelValue >> 24 & 0xff
					g := pixelValue >> 16 & 0xff
					b := pixelValue >> 8 & 0xff
					a := 0xff

					interpolatedColors = append(interpolatedColors, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
				}
			}
		}
	}

	return interpolatedColors
}

func render(maxIteration int, colors []color.RGBA) {
	var waitGroup sync.WaitGroup

	ratio := float64(height) / float64(width)
	xmin, xmax := xpos-escapeRadius/2.0, math.Abs(xpos+escapeRadius/2.0)
	ymin, ymax := ypos-escapeRadius*ratio/2.0, math.Abs(ypos+escapeRadius*ratio/2.0)

	fmt.Printf("y min: %f, y max: %f \n", ymin, ymax)
	fmt.Printf("x min: %f, x max: %f \n", xmin, xmax)

	image := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

	var xMiddle []int
	var yMiddle []int

	for iy := 0; iy < height; iy++ {
		waitGroup.Add(1)
		go func(iy int) {
			defer waitGroup.Done()
			var xs []int
			var ys []int
			var previousColor uint32
			for ix := 0; ix < width; ix++ {

				var x = xmin + (xmax-xmin)*float64(ix)/float64(width-1)
				var y = ymin + (ymax-ymin)*float64(iy)/float64(height-1)
				norm, it := mandelIteration(x, y, maxIteration)
				iteration := float64(maxIteration-it) + math.Log(norm)

				if int(math.Abs(iteration)) < len(colors)-1 {
					color1 := colors[int(math.Abs(iteration))]
					color2 := colors[int(math.Abs(iteration))+1]
					color := linearInterpolation(rgbaToUint(color1), rgbaToUint(color2), uint32(iteration))
					if previousColor != color {
						xs = append(xs, ix)
						ys = append(ys, iy)

					}
					previousColor = color

					image.Set(ix, iy, uint32ToRgba(color))

				}
			}
			xMiddle = append(xMiddle, sum(xs)/len(xs))
			yMiddle = append(yMiddle, sum(ys)/len(ys))
		}(iy)
	}

	waitGroup.Wait()
	middX := sum(xMiddle) / len(xMiddle)
	middY := sum(yMiddle) / len(yMiddle)

	if middX < width && middX > 0 {
		xpos = convertToRelative(xmin, xmax, width, middX)
	}

	if middY < height && middY > 0 {
		ypos = convertToRelative(ymin, ymax, height, middY)
	}

	fmt.Printf("middle x: %d, pos x: %f, middle y: %d, pos y: %f, width: %d, height: %d \n", middX, xpos, middY, ypos, width, height)

	output, _ := os.Create(outputFile)
	defer output.Close()

	png.Encode(output, image)
}

func convertToRelative(min float64, max float64, absoluteScale int, absoluteNumber int) float64 {
	rel := float64(absoluteNumber) / float64(absoluteScale)
	return (1-rel)*min + rel*max
}

func sum(arr []int) int {
	result := 0
	for _, i := range arr {
		result += i
	}
	return result
}

func cosineInterpolation(c1, c2, mu float64) float64 {
	mu2 := (1 - math.Cos(mu*math.Pi)) / 2.0
	return c1*(1-mu2) + c2*mu2
}

func linearInterpolation(c1, c2, mu uint32) uint32 {
	return c1*(1-mu) + c2*mu
}

func mandelIteration(cx, cy float64, maxIter int) (float64, int) {
	var x, y, xx, yy float64 = 0.0, 0.0, 0.0, 0.0

	for i := 0; i < maxIter; i++ {
		xy := x * y
		xx = x * x
		yy = y * y
		if xx+yy > 4 {
			return xx + yy, i
		}
		x = xx - yy + cx
		y = 2*xy + cy
	}

	logZn := (x*x + y*y) / 2
	return logZn, maxIter
}

func rgbaToUint(color color.RGBA) uint32 {
	r, g, b, a := color.RGBA()
	r /= 0xff
	g /= 0xff
	b /= 0xff
	a /= 0xff
	return uint32(r)<<24 | uint32(g)<<16 | uint32(b)<<8 | uint32(a)
}

func uint32ToRgba(col uint32) color.RGBA {
	r := col >> 24 & 0xff
	g := col >> 16 & 0xff
	b := col >> 8 & 0xff
	a := 0xff
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}
