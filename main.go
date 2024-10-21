package main

import (
	"time"

	"github.com/jamesonhm/pokedexcli/internal/pokeapi"
)

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// )

func main() {
	client := pokeapi.NewClient(5 * time.Second)
	config := NewConfig(client)
	runRepl(config)
}
