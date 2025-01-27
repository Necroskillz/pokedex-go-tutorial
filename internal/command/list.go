package command

import (
	"fmt"

	"github.com/necroskillz/pokedex/internal/pokedex"
)

type ListCommand struct {
	pokedex *pokedex.Pokedex
}

func NewListCommand(pokedex *pokedex.Pokedex) *ListCommand {
	return &ListCommand{
		pokedex: pokedex,
	}
}

func (c *ListCommand) Execute(args []string) error {
	pokemon := c.pokedex.GetAll()
	if len(pokemon) == 0 {
		fmt.Println("Your Pokedex is empty!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for _, p := range pokemon {
		fmt.Printf("  - %s\n", p.Name)
	}
	return nil
}
