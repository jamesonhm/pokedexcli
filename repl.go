package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func runRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	config := NewConfig()
	for {
		fmt.Printf("Pokedex> ")
		scanner.Scan()

		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		cmdName := words[0]
		cliCmd, ok := getCmds()[cmdName]
		if !ok {
			fmt.Println("Not a valid cmd")
			continue
		}

		if err := cliCmd.callback(config); err != nil {
			fmt.Println(err)
			continue
		}

	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}
