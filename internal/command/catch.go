package command

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/necroskillz/pokedex/internal/pokedex"
	pokemonapi "github.com/necroskillz/pokedex/internal/pokemon-api"
)

type CatchCommand struct {
	api     *pokemonapi.PokemonApi
	pokedex *pokedex.Pokedex
}

func NewCatchCommand(pokedex *pokedex.Pokedex, api *pokemonapi.PokemonApi) *CatchCommand {
	return &CatchCommand{
		api,
		pokedex,
	}
}

func (c *CatchCommand) Execute(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("please provide a pokemon name")
	}

	pokemonName := args[0]
	pokemon, err := c.api.GetPokemon(pokemonName)
	if err != nil {
		return fmt.Errorf("failed to get pokemon info: %w", err)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	catchChance := calculateCatchChance(pokemon.BaseExperience)
	caught := rand.Float64() <= catchChance

	if caught {
		types := []string{}
		for _, t := range pokemon.Types {
			types = append(types, t.Type.Name)
		}

		stats := []pokedex.PokemonStat{}
		for _, s := range pokemon.Stats {
			stats = append(stats, pokedex.PokemonStat{
				Stat:  s.Stat.Name,
				Value: s.BaseStat,
			})
		}

		c.pokedex.AddPokemon(pokedex.Pokemon{
			Name:     pokemon.Name,
			CaughtAt: time.Now(),
			Height:   pokemon.Height,
			Weight:   pokemon.Weight,
			Types:    types,
			Stats:    stats,
		})
		fmt.Printf("%s was caught!\n", pokemon.Name)
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func calculateCatchChance(baseExp int) float64 {
	if baseExp <= 0 {
		return 1.0 // 100% chance for zero base experience
	}

	// Using exponential decay function:
	// chance = min_chance + (max_chance - min_chance) * e^(-k * baseExp)
	// where:
	// - min_chance = 0.05 (5% minimum catch rate)
	// - max_chance = 1.0 (100% at zero)
	// - k = 0.005 (decay rate tuned to give reasonable values)

	const (
		minChance = 0.05
		maxChance = 1.0
		k         = 0.005 // tuned to give ~15% at baseExp 340
	)

	return minChance + (maxChance-minChance)*math.Exp(-k*float64(baseExp))
}
