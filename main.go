package main

import (
	"time"

	"github.com/savisitor15/db-pokedex-go/internal"
)

type pokeState struct {
	globalMap *PokeMap
	globalMapCache pokecache.Cache
	baseUrl string
	previousCommand string
	fDebug bool
}

var progState pokeState = pokeState{globalMap: &PokeMap{},baseUrl: "https://pokeapi.co/api/v2/location-area/", globalMapCache: pokecache.NewCache(time.Duration(time.Second * 5))}

func main() {
	commandLoop()
}
