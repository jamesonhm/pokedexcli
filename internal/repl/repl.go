package repl

import (
	"fmt"
	"os"
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

// Internal Methods
func (r *Repl) clearBuffer() {
	moveCursorTo(0, 0)

	r.log("clearing buffer\n")
	clearRows()
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
			r.clearBuffer()
			//r.writeStatus()
		case 4: // Ctrl-D
			r.quit()
		case 13: // RETURN
			r.evalBuffer()
		case 27: // ESC
			r.clearBuffer()
			//r.writeStatus()
		case 127: // Backspace
			r.backspaceActiveBuffer()
		default:
			if b[0] >= 32 {
				//r.clearStatus()
				r.addBytesToBuffer([]byte{b[0]})
				//r.writeStatus()
			}
		}
	} else if n > 2 && b[0] == 27 && b[1] == 91 { // [ESCAPE, OPEN_BRACKET, ...]
		if n == 3 {
			switch b[2] {
			case 65: // Up Arrow?
				r.historyBack()
			case 66: // Down Arrow?
				r.historyForward()
			}
		}
	} else {
		r.cleanAndAddToBuffer(b)
	}
	return
}

func (r *Repl) log(format string, args ...interface{}) {
	if r.debug != nil {
		fmt.Fprintf(r.debug, format, args...)
	}
}

func (r *Repl) printPrompt() {
	moveToRowStart()
	fmt.Print(r.handler.Prompt())
}

func (r *Repl) resetBuffer() {
	r.bufferPos = 0
	r.buffer = make([]byte, 0)
	r.printPrompt()

}
