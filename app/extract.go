package app

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/sirupsen/logrus"
)

type StegOption struct {
	channel string
	bits    []uint
	order   string
	xy      string
}

func getPixels(img image.Image, opt StegOption) []color.Color {
	var bounds = img.Bounds()
	logrus.Debugf("image size: %dx%d", bounds.Dx(), bounds.Dy())
	var result = make([]color.Color, bounds.Dx()*bounds.Dy())
	var index = 0
	if opt.xy == "xy" {
		for i := range bounds.Dx() {
			for j := range bounds.Dy() {
				result[index] = img.At(i, j)
				index += 1
			}
		}
	} else {
		for i := range bounds.Dy() {
			for j := range bounds.Dx() {
				result[index] = img.At(j, i)
				index += 1
			}
		}
	}
	return result
}

func extractData(f *os.File, opt StegOption) []byte {
	var bits []uint8
	img, format, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	logrus.Debugf("image format: %s", format)

	var pix = getPixels(img, opt)
	for i := range len(pix) {
		r, g, b, a := pix[i].RGBA()
		var dat uint32
		switch opt.channel {
		case "r":
			dat = r >> 8
		case "g":
			dat = g >> 8
		case "b":
			dat = b >> 8
		case "a":
			dat = a >> 8
		}

		for j := range len(opt.bits) {
			var bs = opt.bits[j]
			if opt.order == "msb" {
				bs = 8 - bs // most significant bit first
			}
			bits = append(bits, uint8((dat>>(bs-1))&1))
		}
	}

	// convert 8 bit into 1 byte
	var bytebuf = bytes.NewBuffer([]byte{})

	var count int = 0
	var buf uint8 = 0
	for i := range len(bits) {
		buf = (buf << 1) | bits[i]
		count += 1
		if count == 8 {
			logrus.Tracef("buf = %x", buf)
			binary.Write(bytebuf, binary.LittleEndian, buf)
			count = 0
			buf = 0
		}
	}
	if count != 0 { // write the rest data
		binary.Write(bytebuf, binary.LittleEndian, buf)
	}
	return bytebuf.Bytes()
}
