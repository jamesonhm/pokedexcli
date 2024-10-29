package main

import (
	"fmt"
	"strings"

	"github.com/jamesonhm/pokedexcli/internal/repl"
)

type pokeHandler struct {
	r *repl.Repl
	c *Config
}

func (h *pokeHandler) Prompt() string {
	return "Pokedex> "
}

func (h *pokeHandler) Tab(buffer string) string {
	return ""
}

func (h *pokeHandler) Eval(line string) string {
	fields := cleanInput(line)

	cmdName := fields[0]
	cliCmd, ok := getCmds()[cmdName]
	if !ok {
		return fmt.Sprintf("%s not a valid cmd\n", cmdName)
	}
	var args []string
	if len(fields) > 1 {
		args = fields[1:]
	}
	msg := cliCmd.callback(h, args...)
	return msg
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}
