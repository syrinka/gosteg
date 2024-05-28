package app

import (
	"os"
)

type StegOption struct {
	channel string
	bits    []uint
	order   string
	xy      string
}

func extractData(f *os.File, opt StegOption) []byte {
	return []byte{0x00, 0xFF, 0xAA}
}
