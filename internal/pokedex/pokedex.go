package pokedex

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

type PokemonStat struct {
	Stat  string `json:"stat"`
	Value int    `json:"value"`
}

type Pokemon struct {
	Name     string        `json:"name"`
	CaughtAt time.Time     `json:"caught_at"`
	Height   int           `json:"height"`
	Weight   int           `json:"weight"`
	Types    []string      `json:"types"`
	Stats    []PokemonStat `json:"stats"`
}

type Pokedex struct {
	pokemon []Pokemon
	mu      sync.Mutex
}

func NewPokedex() *Pokedex {
	pokedex := &Pokedex{
		pokemon: []Pokemon{},
	}

	err := pokedex.load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Error loading pokedex: %s\n", err)
	}

	return pokedex
}

func (p *Pokedex) AddPokemon(pokemon Pokemon) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.pokemon = append(p.pokemon, pokemon)

	go func() {
		err := p.save()
		if err != nil {
			fmt.Printf("error saving pokedex: %s\n", err)
		}
	}()
}

func (p *Pokedex) GetAll() []Pokemon {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.pokemon
}

func (p *Pokedex) GetByName(name string) *Pokemon {
	for _, pokemon := range p.pokemon {
		if pokemon.Name == name {
			return &pokemon
		}
	}

	return nil
}

func (p *Pokedex) save() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	json, err := json.Marshal(p.pokemon)
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(os.Getenv("HOME"), ".pokedex.json"), json, 0644)
}

func (p *Pokedex) load() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	file, err := os.ReadFile(path.Join(os.Getenv("HOME"), ".pokedex.json"))
	if err != nil {
		return err
	}

	var pokemon []Pokemon
	err = json.Unmarshal(file, &pokemon)
	if err != nil {
		return err
	}

	p.pokemon = pokemon
	return nil
}
