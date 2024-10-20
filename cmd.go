package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommand struct {
	name     string
	desc     string
	callback func(c *config) error
}

func getCmds() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:     "help",
			desc:     "Displays a help message",
			callback: cmdHelp,
		},
		"exit": {
			name:     "exit",
			desc:     "Exit the pokedex",
			callback: cmdExit,
		},
		"map": {
			name:     "map",
			desc:     "display the names of 20 location areas in the pokemon world",
			callback: cmdMap,
		},
	}
}

//func cmd(input string) (func(c *config) error, error) {
//	cliCmd, ok := getCmds()[input]
//	if !ok {
//		return nil, fmt.Errorf("Not a valid cmd")
//	}
//
//	return cliCmd.callback, nil
//}

func cmdHelp(c *config) error {
	fmt.Printf("\nWelcome to the Pokedex!\n\n")
	fmt.Printf("Usage:\n\n")
	for _, cmd := range getCmds() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.desc)
	}
	fmt.Println()
	return nil
}

func cmdExit(c *config) error {
	os.Exit(0)
	return nil
}

func cmdMap(c *config) error {
	var url string
	if c.next == "" {
		url = BaseApi + "location-area/"
	} else {
		url = c.next
	}
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}
	locAreas := NewLocationArea()

	err = json.Unmarshal(body, &locAreas)
	if err != nil {
		return err
	}

	c.previous = locAreas.Previous
	c.next = locAreas.Next
	for _, result := range locAreas.Results {
		fmt.Printf("%v\n", result.Name)
	}
	return nil
}
