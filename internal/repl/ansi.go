package repl

import (
	"fmt"
	"os"
)

const _ESC = "\033"

func csi1(n int, char byte) {
	fmt.Fprintf(os.Stdout, "%s[%d%c", _ESC, n, char)
}

func moveToRowStart() {
	csi1(1, 'G')
}

func queryCursorPos() {
	csi1(6, 'n')
}
