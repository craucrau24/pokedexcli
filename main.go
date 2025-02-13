package main

import (
	"bufio"
	"fmt"
	"os"
)

func inputLoop() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Pokedex > ")
	for scanner.Scan() {
		line := scanner.Text()
		words := cleanInput(line)
		fmt.Printf("Your command was: %s\n", words[0])
		fmt.Print("Pokedex > ")
	}
}

func main() {
	inputLoop()
}
