package app

import (
	"image"
	. "image/color"
	"os"
)

type StegOption struct {
	channel string
	bits    []uint
	order   string
	xy      string
}

func getPixels(img image.Image, opt StegOption) (result []Color) {
	var bounds = img.Bounds().Max
	result = make([]Color, bounds.X*bounds.Y)
	if opt.xy == "xy" {
		for i := range bounds.X {
			for j := range bounds.Y {
				result = append(result, img.At(i, j))
			}
		}
	} else {
		for i := range bounds.Y {
			for j := range bounds.X {
				result = append(result, img.At(j, i))
			}
		}
	}
	return
}

func extractData(f *os.File, opt StegOption) []byte {
	// return []byte{0x00, 0xFF, 0xAA}
	var bits = make([]uint8, 0)
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	var pix = getPixels(img, opt)
	for i := range len(pix) {
		r, g, b, a := pix[i].RGBA()
		var dat uint32
		switch opt.channel {
		case "r":
			dat = r
		case "g":
			dat = g
		case "b":
			dat = b
		case "a":
			dat = a
		}

		for j := range len(opt.bits) {
			var bs = opt.bits[j]
			if opt.order == "msb" {
				bs = 16 - bs // most significant bit first
			}
			bits = append(bits, uint8((dat>>(bs-1))&1))
		}
	}
}
