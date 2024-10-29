package main

import (
	"log"
	"time"

	"github.com/jamesonhm/pokedexcli/internal/pokeapi"
	"github.com/jamesonhm/pokedexcli/internal/pokecache"
	"github.com/jamesonhm/pokedexcli/internal/repl"
)

func main() {
	cache := pokecache.NewCache(time.Second * 60)
	client := pokeapi.NewClient(time.Second*5, cache)
	config := NewConfig(client)
	//runRepl(config)
	h := &pokeHandler{
		c: config,
	}
	h.r = repl.NewRepl(h, "debug.log")

	if err := h.r.Loop(); err != nil {
		log.Fatal(err)
	}
}
