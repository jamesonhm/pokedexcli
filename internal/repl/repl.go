package repl

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

type Repl struct {
	handler   Handler
	history   [][]byte
	idx       int
	reader    *_StdinReader
	buffer    []byte
	bufferPos int
	viewStart int
	viewEnd   int
	promptRow int
	height    int
	width     int
	debug     *os.File
	onEnd     func()
}

func NewRepl(handler Handler, debug string) *Repl {
	r := &Repl{
		handler:   handler,
		history:   make([][]byte, 0),
		idx:       -1,
		reader:    newStdinReader(),
		buffer:    nil,
		bufferPos: 0,
		viewStart: 0,
		viewEnd:   -1,
		promptRow: -1,
		height:    0,
		width:     0,
		onEnd:     nil,
		debug:     nil,
	}

	if debug != "" {
		debug, err := os.Create(debug)
		if err != nil {
			panic(err)
		}
		r.debug = debug
	}
	return r
}

// Exported Methods

// Loop sets the terminal into raw mode
// so any further calls to fmt.Print or similar may not behave as expected
func (r *Repl) Loop() error {
	if err := r.MakeRaw(); err != nil {
		return err
	}

	r.reader.start()

	r.printPrompt()

	queryCursorPos()

	for {
		r.reader.read()

		bts := <-r.reader.bytesChan
		r.dispatch(bts)
	}
}

func (r *Repl) MakeRaw() error {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}

	r.onEnd = func() {
		term.Restore(fd, oldState)
	}

	return nil
}

func (r *Repl) UnmakeRaw() {
	r.onEnd()
	r.onEnd = nil
}

func (r *Repl) Quit() {
	r.quit()
}

// Internal Methods
func (r *Repl) addBytesToBuffer(bs []byte) {
	if r.bufferPos == len(r.buffer) {
		xBef, _ := r.cursorCoord(-1)
		r.bufferPos += len(bs)
		len_ := len(r.buffer)
		r.buffer = append(r.buffer, bs...)

		if !r.overflow() {
			needSync := false
			for _, b := range bs {
				r.writeByte(b)
				if b != '\n' && xBef == r.width-1 {
					needSync = true
				}
			}

			if needSync {
				r.syncCursor()
			}
			r.boundPromptRow()
			return
		} else {
			// reset prev changes
			r.bufferPos -= len(bs)
			r.buffer = r.buffer[0:len_]
		}
	}

	tail := r.buffer[r.bufferPos:]

	newBuffer := make([]byte, 0)
	newBuffer = append(newBuffer, r.buffer[0:r.bufferPos]...)
	newBuffer = append(newBuffer, bs...)
	newBuffer = append(newBuffer, tail...)

	newPos := r.bufferPos + len(bs)
	r.force(newBufer, newPos) // force should take into account extra long lines
}

func (r *Repl) calcHeight(buffer []byte, x0 int, w int) int {
	_, y := relCursorCoord(buffer, x0, len(buffer), w)
	return y + 1
}

func (r *Repl) clearAfterPrompt() {
	moveCursorTo(0, r.height-1)
	if r.promptRow < 0 {
		r.updatePromptRow(0)
	}
	dy := (r.height - 1 - r.promptRow)
	clearRows(dy)
}

func (r *Repl) clearBuffer() {
	moveCursorTo(0, r.height-1)

	r.log("clearing buffer\n")
	if r.promptRow < 0 {
		r.updatePromptRow(0)
	}
	dy := (r.height - 1 - r.promptRow)

	clearRows(dy)
	clearRow()
	r.resetBuffer()
}

func (r *Repl) cursorCoord(bufferPos int) (int, int) {
	w := r.width

	if bufferPos < 0 {
		bufferPos = r.bufferPos
	}

	x, y := relCursorCoord(r.buffer[r.viewStart:], r.promptLen(), bufferPos-r.viewStart, w)
	y += r.promptRow
	return x, y
}

func (r *Repl) dispatch(b []byte) {
	n := len(b)

	r.log("keypress: %v\n", b)

	if n == 1 {
		switch b[0] {
		case 0:
			return
		case 3: // Ctrl-C
			fmt.Fprintf(os.Stdout, "Ctrl C: %v\n", b)
			r.clearBuffer()
			//r.writeStatus()
		case 4: // Ctrl-D
			fmt.Fprintf(os.Stdout, "Ctrl D: %v\n", b)
			r.quit()
		case 13: // RETURN
			fmt.Fprintf(os.Stdout, "Return: %v\n", b)
			r.evalBuffer()
		case 27: // ESC
			fmt.Fprintf(os.Stdout, "ESC: %v\n", b)
			r.clearBuffer()
			//r.writeStatus()
		case 127: // Backspace
			//r.backspaceActiveBuffer()
			fmt.Fprintf(os.Stdout, "backspace: %v\n", b)
		default:
			if b[0] >= 32 {
				//r.clearStatus()
				//fmt.Fprintf(os.Stdout, "default, b[0] >= 32: %v\n", b)
				r.addBytesToBuffer([]byte{b[0]})
				//r.writeStatus()
			}
		}
	} else if n > 2 && b[0] == 27 && b[1] == 91 { // [ESCAPE, OPEN_BRACKET, ...]
		if n == 3 {
			switch b[2] {
			case 65: // Up Arrow?
				fmt.Fprintf(os.Stdout, "UPARROW: %v\n", b)
				//r.historyBack()
			case 66: // Down Arrow?
				fmt.Fprintf(os.Stdout, "DOWNARROW: %v\n", b)
				//r.historyForward()
			}
		}
	} else {
		fmt.Fprintf(os.Stdout, "Final else condition: %v\n", b)
		//r.cleanAndAddToBuffer(b)
	}
	return
}

func (r *Repl) evalBuffer() {
	//r.clearStatus()
	r.newLine()
	out := r.handler.Eval(strings.TrimSpace(string(r.buffer)))

	if len(out) > 0 {
		outLines := strings.Split(out, "\n")
		for _, outline := range outLines {
			fmt.Print(outline)
			r.newLine()
		}
	}

	//r.appendToHistory(r.buffer)
	r.idx = -1
	r.resetBuffer()
}

func (r *Repl) log(format string, args ...interface{}) {
	if r.debug != nil {
		fmt.Fprintf(r.debug, format, args...)
	}
}

func (r *Repl) newLine() {
	fmt.Fprintf(os.Stdout, "\n\r")
}

func (r *Repl) overflow() bool {
	b := r.calcHeight() > r.height
	if !b {
		r.viewStart = 0
		r.viewEnd = -1
	}
	return b
}

func (r *Repl) printPrompt() {
	moveToRowStart()
	fmt.Print(r.handler.Prompt())
}

func (r *Repl) promptLen() int {
	return len(r.handler.Prompt())
}

func (r *Repl) quit() {
	r.clearAfterPrompt()
	fmt.Print("\n\r")
	moveToRowStart()
	r.UnmakeRaw()
	os.Exit(0)
}

func relCursorCoord(buffer []byte, x0 int, bufferPos int, w int) (int, int) {
	x := x0
	y := 0
	for j, c := range buffer {
		if j >= bufferPos {
			break
		} else if c == '\n' {
			x = 0
			y += 1
		} else {
			x += 1
		}
		if x == w {
			x = 0
			y += 1
		}
	}
	return x, y
}

func (r *Repl) resetBuffer() {
	r.bufferPos = 0
	r.buffer = make([]byte, 0)
	r.printPrompt()

}

func (r *Repl) updatePromptRow(row int) {
	if row >= r.height {
		row = r.height - 1
	} else if row < 0 {
		row = 0
	}
	r.promptRow = row
	r.log("prompt row %d/%d\n", r.promptRow, r.height-1)
}

func (r *Repl) writeByte(b byte) {
	if b == '\n' {
		r.newLine()
	} else {
		fmt.Fprintf(os.Stdout, "%c", b)
	}
}
