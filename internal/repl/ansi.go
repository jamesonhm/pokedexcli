package repl

import (
	"fmt"
	"os"
)

const _ESC = "\033"

func csi1(n int, char byte) {
	fmt.Fprintf(os.Stdout, "%s[%d%c", _ESC, n, char)
}

func csi2(n int, m int, char byte) {
	fmt.Fprintf(os.Stdout, "%s[%d;%d%c", _ESC, n, m, char)
}

func clearRow() {
	csi1(2, 'K')
}

func clearRows(n int) {
	for i := 0; i < n; i++ {
		csi1(2, 'K')
		csi1(1, 'F')
	}
}

func moveToRowStart() {
	csi1(1, 'G')
}

func queryCursorPos() {
	csi1(6, 'n')
}

func moveCursorTo(x, y int) {
	csi2(y+1, x+1, 'H')
}
