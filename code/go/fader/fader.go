package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
)

var (
	inProgFolder         = "in_progress"
	opacity_step_procent = 3
)

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
	images := []string{"1.png", "2.png", "1.png"}
	create_sequence(images)
}

func create_sequence(images []string) {
	numOfImages := len(images)

	numToBeDrawn := (numOfImages - 1) * 100 / opacity_step_procent

	numOfDrawn := 0
	counter := 0
	for {

		firstImgIndex := counter / 100

		opacity := counter % 100

		numOfDrawn++
		create_image(images[firstImgIndex], images[firstImgIndex+1], uint(opacity), uint(numOfDrawn))
		fmt.Printf("Created %d. image out of %d.\n", numOfDrawn, numToBeDrawn)

		counter = counter + opacity_step_procent
		if counter >= (numOfImages-1)*100 {
			break
		}

	}

}

func create_image(firstImage string, secondImage string, opacityPercent uint, count uint) {
	perc := float64(opacityPercent) / float64(100)

	//Background image
	fImg1, _ := os.Open(firstImage)
	defer fImg1.Close()
	img1, _, _ := image.Decode(fImg1)

	//Logo to stick over background image
	fImg2, _ := os.Open(secondImage)
	defer fImg2.Close()
	img2, _, _ := image.Decode(fImg2)

	mask := image.NewUniform(color.Alpha{uint8(perc * 256)})

	//Create a new blank image m
	m := image.NewRGBA(image.Rect(0, 0, img1.Bounds().Dx(), img1.Bounds().Dy()))

	//Paste background image over m
	draw.Draw(m, m.Bounds(), img1, image.Point{0, 0}, draw.Src)

	//Now paste logo image over m using a mask (ref. http://golang.org/doc/articles/image_draw.html )

	//******Goal is to have opacity value 50 of logo image, when we paste it****
	draw.DrawMask(m, m.Bounds(), img2, image.Point{0, 0}, mask, image.Point{0, 0}, draw.Over)

	toimg, _ := os.Create(fmt.Sprintf("%s/%04d.png", inProgFolder, count))
	defer toimg.Close()

	png.Encode(toimg, m)

}
