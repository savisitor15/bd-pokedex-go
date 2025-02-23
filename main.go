package main

import (
	
)

type pokeState struct {
	globalMap *PokeMap
	globalMapCache map[string]PokeMap
	baseUrl string
	previousCommand string
	fDebug bool
}

var progState pokeState = pokeState{globalMap: &PokeMap{},baseUrl: "https://pokeapi.co/api/v2/location-area/", globalMapCache: make(map[string]PokeMap)}

func main() {
	commandLoop()
}
