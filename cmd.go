package main

import "fmt"

type cliCommand struct {
	name     string
	desc     string
	callback func() error
}

func cmds() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:     "help",
			desc:     "Displays a help message",
			callback: cmdHelp,
		},
		"exit": {
			name:     "exit",
			desc:     "Exit the pokedex",
			callback: nil,
		},
	}
}

func cmd(input string) (func() error, error) {
	cliCmd, ok := cmds()[input]
	if !ok {
		return nil, fmt.Errorf("Not a valid cmd")
	}

	return cliCmd.callback, nil
}

func cmdHelp() error {
	fmt.Printf("\nWelcome to the Pokedex!\n\n")
	fmt.Printf("Usage:\n\n")
	for _, cmd := range cmds() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.desc)
	}
	fmt.Println()
	return nil
}
