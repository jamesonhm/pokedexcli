package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("Pokedex> ")
		scanner.Scan()
		callback, err := cmd(scanner.Text())
		if err != nil {
			fmt.Println(err)
		}
		if err := callback(); err != nil {
			fmt.Println(err)
		}
	}
}
