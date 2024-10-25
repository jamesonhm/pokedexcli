package main

import (
	//	"bufio"
	"fmt"
	"os"
	"strings"

	//"atomicgo.dev/keyboard"
	//"atomicgo.dev/keyboard/keys"
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

func runRepl(config *Config) {
	//scanner := bufio.NewScanner(os.Stdin)
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	for {
		fmt.Print("\r\x1b[K")
		fmt.Printf("Pokedex> ")
		input, err := readLine()
		if err != nil {
			fmt.Println("\nError reading input:", err)
			continue
		}

		// Handle special keys
		switch {
		case input == "\x03":
			fmt.Print("\r\x1b[K")
			fmt.Println("Exiting...\r")
			return
		case input == "\x1b[A":
			// check history, if not empty clear and reprint
			fmt.Println("UP PRESSED")
			clearLine()
			fmt.Print("Pokedex> " + "history here")
			continue
		case input == "\x1b[B":
			// check history, if not empty, clear and reprint
			fmt.Println("DOWN PRESSED")
			continue
		}

		// Normal input processing
		if input != "" {
			// add to history
			fmt.Printf("\n%s\n", "Entered:"+input)
			// process command
		}
		//scanner.Scan()
	}
}

//words := cleanInput(scanner.Text())
//if len(words) == 0 {
//	continue
//}

//cmdName := words[0]
//cliCmd, ok := getCmds()[cmdName]
//if !ok {
//	fmt.Println("Not a valid cmd")
//	continue
//}
//var args []string
//if len(words) > 1 {
//	args = words[1:]
//}

//if err := cliCmd.callback(config, args...); err != nil {
//	fmt.Println(err)
//	continue
//}

//	}
//}

func readLine() (string, error) {
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
		//if tmp[0] == '\x1b' {
		case '\x1b':
			// Read [ char
			os.Stdin.Read(tmp)
			if tmp[0] != '[' {
				continue
			}
			// Read actual code (A for up, B for down, etc,)
			os.Stdin.Read(tmp)
			return "\x1b[" + string(tmp[0]), nil
		//}
		// Ctrl-C
		//if tmp[0] == '\x03' {
		case '\x03':
			return "\x03", nil
		//}
		//if tmp[0] == '\x7f' {
		case '\x7f':
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
				fmt.Print("\b \b")
			}
			continue
		//}

		// Normal characters
		//if tmp[0] == '\r' || tmp[0] == '\n' {
		case '\r', '\n':
			fmt.Print("\r")
			return string(buf), nil
			//}
		}
		fmt.Print(string(tmp[0]))
		buf = append(buf, tmp[0])
		//redraw(buf)
	}
}

func redraw(input []byte) {
	fmt.Print("\r\x1b[K")
	fmt.Printf("Pokedex> %s", input)
}

func clearLine() {
	fmt.Print("\r\x1b[K")
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}
