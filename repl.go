package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name     string
	desc     string
	callback func() error
}

var commands = map[string]cliCommand{
	"exit": {
		name:     "exit",
		desc:     "Exit the Pokedex",
		callback: commandExit,
	},
	"help": {
		name:     "help",
		desc:     "Display a help message",
		callback: commandHelp,
	},
}

func commandDesc() []string {
	var result []string
	for _, cmd := range commands {
		result = append(result, fmt.Sprintf("%s: %s", cmd.name, cmd.desc))
	}
	return result
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")
	fmt.Println("")
	for _, line := range commandDesc() {
		fmt.Println(line)
	}

	return nil
}

func inputLoop() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Pokedex > ")
	for scanner.Scan() {
		line := scanner.Text()
		words := cleanInput(line)
		cmd, ok := commands[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			fmt.Print("Pokedex > ")
			continue
		}
		cmd.callback()
		fmt.Print("Pokedex > ")
	}
}
