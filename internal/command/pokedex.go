package command

import (
	"fmt"

	"github.com/necroskillz/pokedex/internal/pokedex"
)

type PokedexCommand struct {
	pokedex *pokedex.Pokedex
}

func NewPokedexCommand(pokedex *pokedex.Pokedex) *PokedexCommand {
	return &PokedexCommand{
		pokedex: pokedex,
	}
}

func (c *PokedexCommand) Execute(args []string) error {
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
