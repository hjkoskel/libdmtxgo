package libdmtxgo

/*
#cgo LDFLAGS: -ldmtx
#include <stdlib.h>
#include <dmtx.h>
#include <stdio.h>
#include <string.h>

*/
import "C"
import (
	"fmt"
	"image"
	"image/color"
	"unsafe"
)

func Encode(code string) (image.Image, error) {
	str := (*C.uchar)(C.CBytes([]byte(code)))

	enc := C.dmtxEncodeCreate()
	C.dmtxEncodeDataMatrix(enc, C.int(len(code)), str)

	width := C.dmtxImageGetProp(enc.image, C.DmtxPropWidth)
	height := C.dmtxImageGetProp(enc.image, C.DmtxPropHeight)
	bytesPerPixel := C.dmtxImageGetProp(enc.image, C.DmtxPropBytesPerPixel)

	bytedata := C.GoBytes(unsafe.Pointer(enc.image.pxl), C.int(width*height*bytesPerPixel))

	C.dmtxEncodeDestroy(&enc)
	w := int(width)
	h := int(height)
	result := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			a := bytedata[(x+y*w)*3]
			if a > 0 {
				result.SetRGBA(x, y, color.RGBA{A: 255, R: 255, G: 255, B: 255})
			} else {
				result.SetRGBA(x, y, color.RGBA{A: 255, R: 0, G: 0, B: 0})
			}
		}
	}
	return result, nil
}

func Decode(picture image.Image) (string, error) {
	var result string
	w := picture.Bounds().Dx()
	h := picture.Bounds().Dy()

	rawdata := make([]byte, w*h*3) //24bits per pixel RGB
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ := picture.At(x, y).RGBA()
			rawdata[(x+y*w)*3] = byte(r)
			rawdata[(x+y*w)*3+1] = byte(g)
			rawdata[(x+y*w)*3+1] = byte(b)
		}
	}

	pxl := (*C.uchar)(C.CBytes(rawdata))
	img := C.dmtxImageCreate(pxl, C.int(w), C.int(h), C.DmtxPack24bppRGB)
	if img == nil {
		return "", fmt.Errorf("Image is null\n")
	}

	dec := C.dmtxDecodeCreate(img, 1)
	reg := C.dmtxRegionFindNext(dec, nil)

	codeDetected := false
	if reg != nil {
		msg := C.dmtxDecodeMatrixRegion(dec, reg, C.DmtxUndefined)
		result = C.GoString((*C.char)(unsafe.Pointer(msg.output)))
		if msg != nil {
			codeDetected = true
			C.dmtxMessageDestroy(&msg)
		}
		if int(C.dmtxRegionDestroy(&reg)) != C.DmtxPass {
			return result, fmt.Errorf("dtmxRegionDestroy failed")
		}
	}

	if int(C.dmtxDecodeDestroy(&dec)) != C.DmtxPass {
		return result, fmt.Errorf("dtmxDecodeDestroy failed")
	}

	if int(C.dmtxImageDestroy(&img)) != C.DmtxPass {
		return result, fmt.Errorf("dtmxImageDestroy failed")
	}

	C.free(unsafe.Pointer(pxl))

	if !codeDetected {
		return result, fmt.Errorf("No code detected")
	}
	return result, nil
}
