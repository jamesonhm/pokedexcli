package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func runRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("Pokedex> ")
		scanner.Scan()

		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		cmdName := words[0]
		callback, err := cmd(cmdName)

		if err != nil {
			fmt.Println(err)
			continue
		}

		if err := callback(); err != nil {
			fmt.Println(err)
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}
