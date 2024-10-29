package main

import (
	"fmt"
	"math/rand"
)

type callbackFn func(h *pokeHandler, args ...string) string

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

func cmdHelp(h *pokeHandler, args ...string) string {

	msg := `Welcome to the Pokedex!
Usage:`
	for _, cmd := range getCmds() {
		//fmt.Print("\r\x1b[K")
		//fmt.Printf("%s: %s\n", cmd.name, cmd.desc)
		msg += fmt.Sprintf(`\n%s: %s\n`, cmd.name, cmd.desc)
	}
	return msg
}

func cmdExit(h *pokeHandler, args ...string) string {
	h.r.Quit()
	return ""
}

func cmdMap(h *pokeHandler, args ...string) string {
	locAreas, err := h.c.client.ListLocations(h.c.Next())
	if err != nil {
		return err.Error()
	}

	h.c.UpdatePrev(locAreas.Previous)
	h.c.UpdateNext(locAreas.Next)
	var msg string
	for _, result := range locAreas.Results {
		//fmt.Print("\r\x1b[K")
		//fmt.Printf("%v\n", result.Name)
		msg += fmt.Sprintf(`\n%v\n`, result.Name)
	}
	return msg
}

func cmdMapb(h *pokeHandler, args ...string) string {
	if h.c.Previous() == "" {
		return "On the first page"
	}

	locAreas, err := h.c.client.ListLocations(h.c.Previous())
	if err != nil {
		return err.Error()
	}

	h.c.UpdatePrev(locAreas.Previous)
	h.c.UpdateNext(locAreas.Next)
	var msg string
	for _, result := range locAreas.Results {
		//fmt.Print("\r\x1b[K")
		//fmt.Printf("%v\n", result.Name)
		msg += fmt.Sprintf(`\n%v\n`, result.Name)
	}
	return msg
}

func cmdExplore(h *pokeHandler, args ...string) string {
	if len(args) != 1 {
		return "Provide one location name"
	}
	locDetails, err := h.c.client.LocationDetails(args[0])
	if err != nil {
		return err.Error()
	}

	msg := fmt.Sprintf(`Exploring %s%s\n)
Found Pokemon:\n`, args[0], "...")
	for _, encounter := range locDetails.PokemonEncounters {
		//fmt.Print("\r\x1b[K")
		//fmt.Printf("- %s\n", encounter.Pokemon.Name)
		msg += fmt.Sprintf(`\n - %s\n`, encounter.Pokemon.Name)
	}
	return msg
}

func cmdCatch(h *pokeHandler, args ...string) string {
	const catchLevel int = 40
	if len(args) != 1 {
		return "Provide the name of one pokemon"
	}
	name := args[0]
	p, err := h.c.client.Pokemon(name)
	if err != nil {
		return err.Error()
	}

	chance := rand.Intn(p.BaseExperience)
	msg := fmt.Sprintf(`exp: %d | chance: %d
Throwing a Pokeball at %s...\n`, p.BaseExperience, chance, name)
	if chance <= catchLevel {
		msg += fmt.Sprintf(`%s was caught!\n`, name)
		h.c.AddPokemon(p)
		return msg
	}
	msg += fmt.Sprintf(`%s escaped!\n`, name)
	return msg
}

func cmdInspect(h *pokeHandler, args ...string) string {
	if len(args) != 1 {
		return "provide the name of a pokemon to inspect"
	}
	name := args[0]
	p, ok := h.c.pokedex[name]
	if !ok {
		return "you have not caught that pokemon"
	}
	msg := fmt.Sprintf(`Name: %s
Height: %d
Weight: %d
Stats:`, name, p.Height, p.Weight)
	for _, s := range p.Stats {
		msg += fmt.Sprintf(`  -%s: %v\n`, s.Stat.Name, s.BaseStat)
	}
	msg += `Types:\n`
	for _, t := range p.Types {
		msg += fmt.Sprintf(`  - %s\n`, t.Type.Name)
	}
	return msg
}

func cmdPokedex(h *pokeHandler, args ...string) string {
	if len(h.c.pokedex) == 0 {
		return "you do not have any pokemon, you are a loser"
	}
	msg := fmt.Sprintf(`Your Pokedex:\n`)
	for name := range h.c.pokedex {
		msg += fmt.Sprintf(`  - %s\n`, name)
	}
	return msg
}
