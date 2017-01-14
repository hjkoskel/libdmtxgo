/*
Simple use example. Code and decode images
*/
package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/hjkoskel/libdmtxgo"
)

//For creating "compression artefacts"
func saveAsJPEG(img *image.RGBA, quality int, filename string) {
	out, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	q := jpeg.Options{Quality: quality}
	err = jpeg.Encode(out, img, &q)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func saveAsPNG(img *image.RGBA, filename string) {
	out, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = png.Encode(out, img)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readPictureFile(filename string) (image.Image, error) {
	imgfile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("File %v not found\n", filename)
	}
	defer imgfile.Close()
	img, _, err := image.Decode(imgfile)
	return img, err
}

func main() {
	//TODO I have to create test from this. But actual functionality is in C-library so there are not much to test
	fmt.Printf("Coding and decoding same code\n")
	img, err := libdmtxgo.Encode("0Q324343430794<OQQ")
	if err != nil {
		fmt.Printf("VIRHE %#v", err)
	}
	saveAsJPEG(img.(*image.RGBA), 0, "testEncode.jpeg")
	saveAsPNG(img.(*image.RGBA), "testEncode.png")

	code, codeerr := libdmtxgo.Decode(img)
	if codeerr == nil {
		fmt.Printf("Code =%#v\n", code)
	} else {
		fmt.Printf("Code error=%#v\n", codeerr)
	}

	cleanIMG, _ := readPictureFile("testEncode.png")
	code, codeerr = libdmtxgo.Decode(cleanIMG)
	if codeerr == nil {
		fmt.Printf("Code =%#v\n", code)
	} else {
		fmt.Printf("Code error=%#v\n", codeerr)
	}

	dirtyIMG, _ := readPictureFile("testEncode.jpeg")
	code, codeerr = libdmtxgo.Decode(dirtyIMG)
	if codeerr == nil {
		fmt.Printf("Code =%#v\n", code)
	} else {
		fmt.Printf("Code error=%#v\n", codeerr)
	}

}
