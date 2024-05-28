package app

import (
	"fmt"
)

func visualize(data []byte) {
	var lines = int(len(data) / 16)
	for ln := range lines {
		var chs = ""

		fmt.Printf("%08x: ", ln*16)
		for i := range 16 {
			var index = ln*16 + i
			var bt byte
			if index >= len(data) {
				bt = 0x00
			} else {
				bt = data[index]
			}

			if bt >= 32 && bt <= 126 {
				chs = chs + string(bt)
			} else {
				chs = chs + " "
			}

			fmt.Printf("%02x ", bt)

			if i == 7 {
				fmt.Print(" ")
			}
		}
		fmt.Printf(" |%s|", chs)
		fmt.Println()
	}
}
