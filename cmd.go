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
		"pokedex": {
			name:     "pokedex",
			desc:     "view list of caught pokemon",
			callback: cmdPokedex,
		},
	}
}

func cmdHelp(c *Config, args ...string) error {
	c.RawPrint("\nWelcome to the Pokedex!\n\n")
	c.RawPrint("Usage:\n\n")
	for _, cmd := range getCmds() {
		//fmt.Print("\r\x1b[K")
		//fmt.Printf("%s: %s\n", cmd.name, cmd.desc)
		c.RawPrint("%s: %s\n", cmd.name, cmd.desc)
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
		//fmt.Print("\r\x1b[K")
		//fmt.Printf("%v\n", result.Name)
		c.RawPrint("%v\n", result.Name)
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
		//fmt.Print("\r\x1b[K")
		//fmt.Printf("%v\n", result.Name)
		c.RawPrint("%v\n", result.Name)
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

	c.RawPrint("Exploring %s%s\n", args[0], "...")
	c.RawPrint("Found Pokemon:\n")
	for _, encounter := range locDetails.PokemonEncounters {
		//fmt.Print("\r\x1b[K")
		//fmt.Printf("- %s\n", encounter.Pokemon.Name)
		c.RawPrint("- %s\n", encounter.Pokemon.Name)
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
	c.RawPrint("Throwing a Pokeball at %s...\n", name)
	if chance <= catchLevel {
		c.RawPrint("%s was caught!\n", name)
		c.AddPokemon(p)
		return nil
	}
	c.RawPrint("%s escaped!\n", name)
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
	c.RawPrint("Name: %s\n", name)
	c.RawPrint("Height: %d\n", p.Height)
	c.RawPrint("Weight: %d\n", p.Weight)
	c.RawPrint("Stats:\n")
	for _, s := range p.Stats {
		c.RawPrint("  -%s: %v\n", s.Stat.Name, s.BaseStat)
	}
	c.RawPrint("Types:\n")
	for _, t := range p.Types {
		c.RawPrint("  - %s\n", t.Type.Name)
	}
	return nil
}

func cmdPokedex(c *Config, args ...string) error {
	if len(c.pokedex) == 0 {
		return fmt.Errorf("you do not have any pokemon, you are a loser")
	}
	c.RawPrint("Your Pokedex:\n")
	for name := range c.pokedex {
		c.RawPrint("  - %s\n", name)
	}
	return nil
}
