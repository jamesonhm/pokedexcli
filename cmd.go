package main

import (
	"fmt"
	"math/rand"
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
			name:     "explore <location_name>",
			desc:     "display the names of the pokemon that can be found in the named area",
			callback: cmdExplore,
		},
		"catch": {
			name:     "catch <pokemon_name>",
			desc:     "chance to catch a pokemon and add to pokedex",
			callback: cmdCatch,
		},
		"inspect": {
			name:     "inspect <pokemon_name>",
			desc:     "view stats of a pokemon you have caught",
			callback: cmdInspect,
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
	if len(args) != 1 {
		return fmt.Errorf("Provide one location name")
	}
	locDetails, err := c.client.LocationDetails(args[0])
	if err != nil {
		return err
	}

	fmt.Println("Exploring", args[0], "...")
	fmt.Println("Found Pokemon:")
	for _, encounter := range locDetails.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func cmdCatch(c *Config, args ...string) error {
	const catchLevel int = 40
	if len(args) != 1 {
		return fmt.Errorf("Provide the name of one pokemon")
	}
	name := args[0]
	p, err := c.client.Pokemon(name)
	if err != nil {
		return err
	}

	chance := rand.Intn(p.BaseExperience)
	fmt.Printf("exp: %d | chance: %d\n", p.BaseExperience, chance)
	fmt.Printf("Throwing a Pokeball at %s...\n", name)
	if chance <= catchLevel {
		fmt.Printf("%s was caught!\n", name)
		c.AddPokemon(p)
		return nil
	}
	fmt.Printf("%s escaped!\n", name)
	return nil
}

func cmdInspect(c *Config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("provide the name of a pokemon to inspect")
	}
	name := args[0]
	p, ok := c.pokedex[name]
	if !ok {
		return fmt.Errorf("you have not caught that pokemon")
	}
	fmt.Println("Name:", name)
	fmt.Println("Height:", p.Height)
	fmt.Println("Weight:", p.Weight)
	fmt.Println("Stats:")
	for _, s := range p.Stats {
		fmt.Printf("  -%s: %v\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
	return nil
}
