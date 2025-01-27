package command

import (
	"fmt"

	pokemonapi "github.com/necroskillz/pokedex/internal/pokemon-api"
)

type MapCommandHandler struct {
	NextUrl     string
	Page        int
	MapbHandler *MapbCommandHandler

	pokemonApi *pokemonapi.PokemonApi
	initialUrl string
}

func NewMapCommand(api *pokemonapi.PokemonApi) *MapCommandHandler {
	return &MapCommandHandler{
		Page:       0,
		initialUrl: api.BuildAreaUrl(0),
		pokemonApi: api,
	}
}

func (c *MapCommandHandler) Execute(args []string) error {
	url := c.initialUrl
	if c.NextUrl != "" {
		url = c.NextUrl
	}

	response, err := c.pokemonApi.GetAreas(url)
	if err != nil {
		return err
	}

	for _, area := range response.Results {
		fmt.Println(area.Name)
	}

	c.NextUrl = response.Next
	c.MapbHandler.PrevUrl = response.Previous
	c.Page++

	return nil
}

type MapbCommandHandler struct {
	PrevUrl    string
	MapHandler *MapCommandHandler

	pokemonApi *pokemonapi.PokemonApi
}

func NewMapbCommand(api *pokemonapi.PokemonApi) *MapbCommandHandler {
	return &MapbCommandHandler{
		pokemonApi: api,
	}
}

func (c *MapbCommandHandler) Execute(args []string) error {
	if c.MapHandler.Page < 2 {
		fmt.Println("you're on the first page")
		return nil
	}

	response, err := c.pokemonApi.GetAreas(c.PrevUrl)
	if err != nil {
		return err
	}

	for _, area := range response.Results {
		fmt.Println(area.Name)
	}

	c.MapHandler.NextUrl = response.Next
	c.PrevUrl = response.Previous
	c.MapHandler.Page--

	return nil
}
