package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jamesonhm/pokedexcli/internal/pokeapi"
	"github.com/jamesonhm/pokedexcli/internal/pokecache"
)

type Config struct {
	next     string
	previous string
	client   *pokeapi.Client
	cache    *pokecache.Cache
}

func NewConfig(client *pokeapi.Client, cache *pokecache.Cache) *Config {
	return &Config{
		client: client,
		cache:  cache,
	}
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
	scanner := bufio.NewScanner(os.Stdin)

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
