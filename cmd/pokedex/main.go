package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/chzyer/readline"
	"github.com/necroskillz/pokedex/internal/command"
	"github.com/necroskillz/pokedex/internal/pokedex"
	pokemonapi "github.com/necroskillz/pokedex/internal/pokemon-api"
	"github.com/necroskillz/pokedex/internal/repl"
)

func registerCommands() map[string]command.Command {
	var commands map[string]command.Command

	api := pokemonapi.NewPokemonApi()
	pokedex := pokedex.NewPokedex()
	mapHandler := command.NewMapCommand(api)
	mapbHandler := command.NewMapbCommand(api)
	mapHandler.MapbHandler = mapbHandler
	mapbHandler.MapHandler = mapHandler

	commands = map[string]command.Command{
		"exit": {
			Name:        "exit",
			Description: "Exit the program",
			Handler:     command.NewExitCommand(),
		},
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Handler:     command.NewHelpCommand(commands),
		},
		"map": {
			Name:        "map",
			Description: "Displays location names in Pokemon World. Using the command again will display the next page.",
			Handler:     mapHandler,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Displays the previous page of location names in Pokemon World",
			Handler:     mapbHandler,
		},
		"explore": {
			Name:        "explore",
			Description: "Explore a location area to see what Pokemon can be found there",
			Handler:     command.NewExploreCommand(api),
		},
		"catch": {
			Name:        "catch",
			Description: "Try to catch a Pokemon by name",
			Handler:     command.NewCatchCommand(pokedex, api),
		},
		"inspect": {
			Name:        "inspect",
			Description: "View information about a Pokemon you have caught",
			Handler:     command.NewInspectCommand(pokedex),
		},
		"list": {
			Name:        "list",
			Description: "List all Pokemon you have caught",
			Handler:     command.NewListCommand(pokedex),
		},
	}
	return commands
}

func main() {
	commands := registerCommands()

	// Create a channel to handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Initialize readline
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "pokedex> ",
		HistoryLimit: 100,
		AutoComplete: readline.NewPrefixCompleter(
			readline.PcItem("exit"),
			readline.PcItem("help"),
			readline.PcItem("map"),
			readline.PcItem("mapb"),
			readline.PcItem("explore"),
			readline.PcItem("catch"),
			readline.PcItem("inspect"),
			readline.PcItem("list"),
		),
		EOFPrompt: "exit",
	})
	if err != nil {
		fmt.Printf("Error initializing readline: %s\n", err)
		os.Exit(1)
	}
	defer rl.Close()

	// Create a channel to handle readline errors
	errChan := make(chan error)

	// Start reading in a goroutine
	go func() {
		for {
			line, err := rl.Readline()
			if err != nil {
				errChan <- err
				return
			}

			words := repl.CleanInput(line)
			if len(words) == 0 {
				continue
			}

			command, exists := commands[words[0]]
			if !exists {
				fmt.Println("Unknown command")
				continue
			}

			err = command.Handler.Execute(words[1:])
			if err != nil {
				fmt.Printf("Error executing command: %s\n", err)
			}
		}
	}()

	// Wait for either a signal or readline error
	select {
	case <-sigChan:
		os.Exit(0)
	case err := <-errChan:
		if err != nil && err != readline.ErrInterrupt {
			fmt.Printf("Error reading input: %s\n", err)
		}
		os.Exit(0)
	}
}
