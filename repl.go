package main

import (
	"fmt"
	"os"

	"github.com/jamesonhm/pokedexcli/internal/pokeapi"
	"golang.org/x/term"
)

type Config struct {
	next     string
	previous string
	client   *pokeapi.Client
	pokedex  map[string]pokeapi.Pokemon
}

func NewConfig(client *pokeapi.Client) *Config {
	return &Config{
		client:  client,
		pokedex: map[string]pokeapi.Pokemon{},
	}
}

func (c *Config) AddPokemon(p pokeapi.Pokemon) {
	name := p.Name
	c.pokedex[name] = p
}

func (c *Config) Previous() string {
	return c.previous
}

func (c *Config) UpdatePrev(new string) {
	c.previous = new
}

func (c *Config) Next() string {
	return c.next
}

func (c *Config) UpdateNext(new string) {
	c.next = new
}

func (c *Config) RawPrint(fmtStr string, args ...any) {
	fmt.Print("\r\x1b[K")
	fmt.Printf(fmtStr, args...)
}

type cmdHistory struct {
	history []string
	idx     int
}

func newHistory() *cmdHistory {
	return &cmdHistory{
		idx: -1,
	}
}

func (ch *cmdHistory) addCmd(cmd string) {
	ch.history = append(ch.history, cmd)
	ch.idx++
}

func (ch *cmdHistory) prevCmd() string {
	//fmt.Printf("len: %d | idx: %d", len(ch.history), ch.idx)
	if len(ch.history) > 0 && ch.idx < len(ch.history) && ch.idx >= 0 {
		prev := ch.history[ch.idx]
		//fmt.Printf("%s", prev)
		if ch.idx > 0 {
			ch.idx--
		}
		return prev
	}
	return ""
}

func (ch *cmdHistory) nextCmd() string {
	if len(ch.history) > 0 && ch.idx < len(ch.history) {
		next := ch.history[ch.idx]
		if ch.idx < len(ch.history)-1 {
			ch.idx++
		}
		return next
	}
	return ""
}

func (ch *cmdHistory) sync() {
	ch.history = ch.history[:ch.idx+1]
	ch.idx = len(ch.history) - 1
}

func runRepl(config *Config) {
	//scanner := bufio.NewScanner(os.Stdin)
	ch := newHistory()
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	for {
		//fmt.Print("\r\x1b[K")
		//fmt.Printf("Pokedex> ")
		config.RawPrint("Pokedex> ")
		input, err := readLine(ch, config)
		if err != nil {
			fmt.Println("\nError reading input:", err)
			continue
		}

		// Handle special keys
		switch input {
		case "\x03":
			config.RawPrint("Exiting...\n\r")
			return
		}

		// Normal input processing
		if input != "" {
			// add to history
			ch.addCmd(input)
			// process command
			words := cleanInput(input)
			//fmt.Printf("\nwords: %v", words)
			cmdName := words[0]
			if cmdName == "exit" {
				config.RawPrint("Exiting...\n\r")
				return
			}
			//			cliCmd, ok := getCmds()[cmdName]
			//			if !ok {
			//				//fmt.Println("Not a valid cmd")
			//				fmt.Print("\n")
			//				config.RawPrint("%s not a valid cmd\n", cmdName)
			//				continue
			//			}
			//			var args []string
			//			if len(words) > 1 {
			//				args = words[1:]
			//			}
			//			fmt.Print("\n")
			//			if err := cliCmd.callback(config, args...); err != nil {
			//				fmt.Println(err)
			//				continue
			//			}
			ch.sync()
		}
	}
}

func readLine(h *cmdHistory, c *Config) (string, error) {
	var buf []byte
	tmp := make([]byte, 1)

	for {
		n, err := os.Stdin.Read(tmp)
		if err != nil {
			return "", err
		}
		if n == 0 {
			continue
		}

		switch tmp[0] {
		// Handle escape sequences
		case '\x1b':
			// Read [ char
			os.Stdin.Read(tmp)
			if tmp[0] != '[' {
				continue
			}
			// Read actual code (A for up, B for down, etc,)
			os.Stdin.Read(tmp)
			if tmp[0] == 'A' {
				//c.RawPrint("Pokedex> " + "UPARROW")
				prev := h.prevCmd()
				//if prev != "" {
				c.RawPrint("Pokedex> " + prev)
				//fmt.Print(prev)
				buf = []byte(prev)
				continue
				//}
				//continue
			} else if tmp[0] == 'B' {
				next := h.nextCmd()
				//if next != "" {
				c.RawPrint("Pokedex> " + next)
				buf = []byte(next)
				continue
				//}
				//continue
			}
			return "\x1b[" + string(tmp[0]), nil
		// Ctrl-C
		case '\x03':
			return "\x03", nil
		//Backspace
		case '\x7f':
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
				fmt.Print("\b \b")
			}
			continue
		// Normal Sequence followed by Enter
		case '\r', '\n':
			fmt.Print("\r")
			return string(buf), nil
		}
		fmt.Print(string(tmp[0]))
		buf = append(buf, tmp[0])
	}
}

func clearLine() {
	fmt.Print("\r\x1b[K")
}

//func cleanInput(text string) []string {
//	output := strings.ToLower(text)
//	words := strings.Fields(output)
//	return words
//}
