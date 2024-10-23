package main

import (
	"time"

	"github.com/jamesonhm/pokedexcli/internal/pokeapi"
	"github.com/jamesonhm/pokedexcli/internal/pokecache"
)

func main() {
	cache := pokecache.NewCache(time.Second * 60)
	client := pokeapi.NewClient(time.Second*5, cache)
	config := NewConfig(client)
	runRepl(config)
}
