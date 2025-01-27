package command

import (
	"fmt"

	"github.com/necroskillz/pokedex/internal/pokedex"
)

type InspectCommand struct {
	pokedex *pokedex.Pokedex
}

func NewInspectCommand(pokedex *pokedex.Pokedex) *InspectCommand {
	return &InspectCommand{
		pokedex: pokedex,
	}
}

func (c *InspectCommand) Execute(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("please provide a pokemon name")
	}

	pokemonName := args[0]
	pokemon := c.pokedex.GetByName(pokemonName)

	if pokemon == nil {
		fmt.Println("You have not caught that pokemon yet")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  - %s: %d\n", stat.Stat, stat.Value)
	}

	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t)
	}

	fmt.Printf("Caught at: %v\n", pokemon.CaughtAt.Format("1/2/2006 3:04:05 PM"))

	return nil
}
