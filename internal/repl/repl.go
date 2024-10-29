package repl

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

type Repl struct {
	handler   Handler
	history   [][]byte
	idx       int
	reader    *_StdinReader
	buffer    []byte
	bufferPos int
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

// Internal Methods
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
				fmt.Fprintf(os.Stdout, "default, b[0] >= 32: %v\n", b)
				//r.addBytesToBuffer([]byte{b[0]})
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

func (r *Repl) printPrompt() {
	moveToRowStart()
	fmt.Print(r.handler.Prompt())
}

func (r *Repl) quit() {
	r.clearAfterPrompt()
	fmt.Print("\n\r")
	moveToRowStart()
	r.UnmakeRaw()
	os.Exit(0)
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
