package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jamesonhm/pokedexcli/internal/pokeapi"
)

type callbackFn func(c *pokeapi.Config) error

type cliCommand struct {
	name     string
	desc     string
	callback callbackFn
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
		"mapb": {
			name:     "mapb",
			desc:     "display the names of previous 20 location areas in the pokemon world",
			callback: cmdMapb,
		},
	}
}

func cmdHelp(c *pokeapi.Config) error {
	fmt.Printf("\nWelcome to the Pokedex!\n\n")
	fmt.Printf("Usage:\n\n")
	for _, cmd := range getCmds() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.desc)
	}
	fmt.Println()
	return nil
}

func cmdExit(c *pokeapi.Config) error {
	os.Exit(0)
	return nil
}

func cmdMap(c *pokeapi.Config) error {
	var url string
	if c.Next() == "" {
		url = pokeapi.BaseApi + "location-area/"
	} else {
		url = c.Next()
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
	locAreas := pokeapi.NewLocationArea()

	err = json.Unmarshal(body, &locAreas)
	if err != nil {
		return err
	}

	c.UpdatePrev(locAreas.Previous)
	c.UpdateNext(locAreas.Next)
	for _, result := range locAreas.Results {
		fmt.Printf("%v\n", result.Name)
	}
	return nil
}

func cmdMapb(c *pokeapi.Config) error {
	var url string
	if c.Previous() == "" {
		return fmt.Errorf("On first Page, no previous pages")
	}
	url = c.Previous()

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}
	locAreas := pokeapi.NewLocationArea()

	err = json.Unmarshal(body, &locAreas)
	if err != nil {
		return err
	}

	c.UpdatePrev(locAreas.Previous)
	c.UpdateNext(locAreas.Next)
	for _, result := range locAreas.Results {
		fmt.Printf("%v\n", result.Name)
	}
	return nil
}
