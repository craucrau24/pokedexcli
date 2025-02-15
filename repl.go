package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/craucrau24/pokedexcli/internal/pokecache"
)

type config struct {
	endPoint string
	next     string
	previous string
}

type cliRegistry struct {
	commands map[string]cliCommand
	cfg      config
	cache    *pokecache.Cache
}

type cliCommand struct {
	name     string
	desc     string
	callback func(cfg *config) error
}

func (r *cliRegistry) init() {
	r.cache = pokecache.NewCache(5 * time.Second)
	r.commands = make(map[string]cliCommand)
	r.addCommand(cliCommand{
		name:     "exit",
		desc:     "Exit the Pokedex",
		callback: commandExit,
	})
	r.addCommand(cliCommand{
		name:     "help",
		desc:     "Display a help message",
		callback: r.commandHelp,
	})
	r.addCommand(cliCommand{
		name:     "map",
		desc:     "Retrieve area locations. Subsequent calls paginate.",
		callback: r.commandMap,
	})
	r.addCommand(cliCommand{
		name:     "mapb",
		desc:     "Retrieve area locations. Backward pagination.",
		callback: r.commandMapb,
	})
}

func (r *cliRegistry) addCommand(cmd cliCommand) {
	r.commands[cmd.name] = cmd
}

func (r *cliRegistry) commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")
	fmt.Println("")

	for _, cmd := range r.commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.desc)
	}

	return nil
}

type MapItemSummary struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type MapResultsDTO struct {
	Count    int              `json:"count"`
	Next     string           `json:"next"`
	Previous string           `json:"previous"`
	Items    []MapItemSummary `json:"results"`
}

func (c *cliRegistry) commandMap(cfg *config) error {
	const endPoint = "location-area"
	var url string
	if endPoint != cfg.endPoint {
		url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	} else {
		url = cfg.next
	}
	data, ok := c.cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("error sending map request: %w", err)
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}
		c.cache.Add(url, data)
	}

	var results MapResultsDTO
	if err := json.Unmarshal(data, &results); err != nil {
		return fmt.Errorf("error decoding map response: %w", err)
	}
	cfg.endPoint = endPoint
	cfg.next = results.Next
	cfg.previous = results.Previous
	for _, item := range results.Items {
		fmt.Println(item.Name)
	}

	return nil
}

func (c *cliRegistry) commandMapb(cfg *config) error {
	const endPoint = "location-area"
	var url string
	if endPoint != cfg.endPoint {
		return fmt.Errorf("mapb command should be used after map command")
	}

	if cfg.previous == "" {
		return fmt.Errorf("already on first page")
	}
	url = cfg.previous

	data, ok := c.cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("error sending map request: %w", err)
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}
		c.cache.Add(url, data)
	}

	var results MapResultsDTO
	if err := json.Unmarshal(data, &results); err != nil {
		return fmt.Errorf("error decoding map response: %w", err)
	}
	cfg.endPoint = endPoint
	cfg.next = results.Next
	cfg.previous = results.Previous
	for _, item := range results.Items {
		fmt.Println(item.Name)
	}

	return nil
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func (r *cliRegistry) execute(cmd string) error {
	cliCmd, ok := r.commands[cmd]
	if !ok {
		return fmt.Errorf("unknown command")
	}
	return cliCmd.callback(&r.cfg)
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
