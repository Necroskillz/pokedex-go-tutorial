package pokemonapi

import (
	"encoding/json"
	"fmt"

	"github.com/necroskillz/pokedex/internal/cache"
)

type PokemonApi struct {
	client *cache.CachedHttpClient
}

func NewPokemonApi() *PokemonApi {
	return &PokemonApi{
		client: cache.NewCachedHttpClient(),
	}
}

type PokemonApiLocationAreasResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

type PokemonApiLocationAreaResponse struct {
	Name              string             `json:"name"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon PokemonEncounterPokemon `json:"pokemon"`
}

type PokemonEncounterPokemon struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Types          []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
}

func (p *PokemonApi) BuildAreaUrl(offset int) string {
	return fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=20", offset)
}

func (p *PokemonApi) GetAreas(url string) (*PokemonApiLocationAreasResponse, error) {
	return getJson[PokemonApiLocationAreasResponse](p.client, url)
}

func (p *PokemonApi) GetArea(name string) (*PokemonApiLocationAreaResponse, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", name)

	return getJson[PokemonApiLocationAreaResponse](p.client, url)
}

func (p *PokemonApi) GetPokemon(name string) (*Pokemon, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	return getJson[Pokemon](p.client, url)
}

func getJson[T any](client *cache.CachedHttpClient, url string) (*T, error) {
	body, err := client.GetWithCache(url)
	if err != nil {
		return nil, err
	}

	var response T
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
