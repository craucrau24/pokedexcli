package pokemon

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type pokemonInnerStatDTO struct {
	Name string `json:"name"`
}

type pokemonStatDTO struct {
	BaseStat int                 `json:"base_stat"`
	Inner    pokemonInnerStatDTO `json:"stat"`
}

type pokemonInnerTypeDTO struct {
	Name string `json:"name"`
}

type pokemonTypeDTO struct {
	Inner pokemonInnerTypeDTO `json:"type"`
}

type pokemonDTO struct {
	Name           string           `json:"name"`
	BaseExperience int              `json:"base_experience"`
	Height         int              `json:"height"`
	Weight         int              `json:"weight"`
	Stats          []pokemonStatDTO `json:"stats"`
	Types          []pokemonTypeDTO `json:"types"`
}

type PokemonStat struct {
	Name     string
	BaseStat int
}

type Pokemon struct {
	Name           string
	BaseExperience int
	Height         int
	Weight         int

	Stats []PokemonStat
	Types []string
}

func JsonToPokemon(data []byte) (Pokemon, error) {
	var pokeDTO pokemonDTO
	if err := json.Unmarshal(data, &pokeDTO); err != nil {
		return Pokemon{}, fmt.Errorf("error decoding json: %w", err)
	}
	var stats []PokemonStat
	for _, stat := range pokeDTO.Stats {
		stats = append(stats, PokemonStat{
			Name:     stat.Inner.Name,
			BaseStat: stat.BaseStat,
		})
	}

	var types []string
	for _, typ := range pokeDTO.Types {
		types = append(types, typ.Inner.Name)
	}

	pokemon := Pokemon{
		Name:           pokeDTO.Name,
		BaseExperience: pokeDTO.BaseExperience,
		Height:         pokeDTO.Height,
		Weight:         pokeDTO.Weight,
		Stats:          stats,
		Types:          types,
	}
	return pokemon, nil
}

type Pokedex struct {
	caught map[string]Pokemon
}

func NewPokedex() Pokedex {
	return Pokedex{caught: make(map[string]Pokemon)}
}

func (p *Pokedex) TryCatch(pokemon Pokemon) bool {
	fmt.Printf("%s has %d base xp\n", pokemon.Name, pokemon.BaseExperience)
	if pokemon.BaseExperience <= 0 {
		fmt.Println("erroneous base experience")
		return false
	}
	roll := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(pokemon.BaseExperience)
	if roll <= 25 {
		p.caught[pokemon.Name] = pokemon
		return true
	}
	return false
}
