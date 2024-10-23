package main

import (
	"fmt"
	"os"
)

type callbackFn func(c *Config, args ...string) error

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
		"explore": {
			name:     "explore",
			desc:     "display the names of the pokemon that can be found in the named area",
			callback: cmdExplore,
		},
	}
}

func cmdHelp(c *Config, args ...string) error {
	fmt.Printf("\nWelcome to the Pokedex!\n\n")
	fmt.Printf("Usage:\n\n")
	for _, cmd := range getCmds() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.desc)
	}
	fmt.Println()
	return nil
}

func cmdExit(c *Config, args ...string) error {
	os.Exit(0)
	return nil
}

func cmdMap(c *Config, args ...string) error {
	locAreas, err := c.client.ListLocations(c.Next())
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

func cmdMapb(c *Config, args ...string) error {
	if c.Previous() == "" {
		return fmt.Errorf("On the first page")
	}

	locAreas, err := c.client.ListLocations(c.Previous())
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

func cmdExplore(c *Config, args ...string) error {
	locDetails, err := c.client.LocationDetails(args[0])
	if err != nil {
		return err
	}

	for _, encounter := range locDetails.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}
	return nil
}
