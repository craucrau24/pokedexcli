package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliRegistry struct {
	commands map[string]cliCommand
}

type cliCommand struct {
	name     string
	desc     string
	callback func() error
}

func (r *cliRegistry) init() {
	r.commands = make(map[string]cliCommand)
	r.commands["exit"] = cliCommand{
		name:     "exit",
		desc:     "Exit the Pokedex",
		callback: commandExit,
	}
	r.commands["help"] = cliCommand{
		name:     "help",
		desc:     "Display a help message",
		callback: r.commandHelp,
	}
}

func (r *cliRegistry) commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")
	fmt.Println("")

	for _, cmd := range r.commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.desc)
	}

	return nil
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func (r *cliRegistry) execute(cmd string) error {
	cliCmd, ok := r.commands[cmd]
	if !ok {
		return fmt.Errorf("unknown command")
	}
	return cliCmd.callback()
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func inputLoop() {
	registry := cliRegistry{}
	registry.init()

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Pokedex > ")
	for scanner.Scan() {
		line := scanner.Text()
		words := cleanInput(line)
		err := registry.execute(words[0])
		if err != nil {
			fmt.Println(err)
			fmt.Print("Pokedex > ")
			continue
		}
		fmt.Print("Pokedex > ")
	}
}
