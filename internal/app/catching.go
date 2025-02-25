package pokedex

import (
	"fmt"
	"math/rand"
)



func catching(pok string) (error){
	fmt.Println(fmt.Sprintf("Throwing a Pokeball at %s...", pok))
	poke, err := getPokemon[PokemonSpecies](getPokemonUrl(pok))
	if err != nil {
		return err
	}
	catch := rand.Intn(256) <= poke.CaptureRate
	if catch {
		progState.PokemonCaught[poke.Name] = poke
		fmt.Println(fmt.Sprintf("%s was caught!", poke.Name))
		return nil
	} else {
		fmt.Println(fmt.Sprintf("%s escaped!", poke.Name))
		return nil
	}
}

func getPokemonUrl(pok string) string{
	return fmt.Sprintf("%spokemon-species/%s", progState.BaseUrl, pok)
}

func getPokemon[V Pokemon](url string) (V, error){
	return getNetObject[V](url)
}
