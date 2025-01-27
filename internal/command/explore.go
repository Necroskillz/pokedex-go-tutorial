package command

import (
	"fmt"

	pokemonapi "github.com/necroskillz/pokedex/internal/pokemon-api"
)

type ExploreCommand struct {
	api *pokemonapi.PokemonApi
}

func NewExploreCommand(api *pokemonapi.PokemonApi) *ExploreCommand {
	return &ExploreCommand{
		api,
	}
}

func (c *ExploreCommand) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a location area name")
	}

	locationArea := args[0]
	area, err := c.api.GetArea(locationArea)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", area.Name)
	fmt.Println("Found Pokemon:")
	for _, encounter := range area.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}
