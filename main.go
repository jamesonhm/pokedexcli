package main

import (
	"time"

	"github.com/jamesonhm/pokedexcli/internal/pokeapi"
	"github.com/jamesonhm/pokedexcli/internal/pokecache"
)

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// )

func main() {
	client := pokeapi.NewClient(5 * time.Second)
	cache := pokecache.NewCache(time.Second * 5)
	config := NewConfig(client, cache)
	runRepl(config)
}
