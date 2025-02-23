package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PokeMap struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func getGlobalMap() PokeMap{
	if progState.globalMap != nil{
		return *progState.globalMap
	}
	newMap := PokeMap{}
	progState.globalMap = &newMap
	return *progState.globalMap
}

func updateMap(direction int) error{
	var err error
	switch {
	case 0 == direction:
		pMap := getGlobalMap()
		if len(pMap.Results) == 0 {
			// don't have a map yet, get one
			if nMap, err := getMap(progState.baseUrl); err != nil {
				return err
			}else{
				progState.globalMap = &nMap
			}
		}else{
			return nil
		}
	case direction < 0:
		pMap := getGlobalMap()
		if pMap.Previous == "null" || len(pMap.Previous) == 0 {
			return fmt.Errorf("you're on the first page")
		}else{
			if nMap, err := getMap(pMap.Previous); err != nil {
				return err
			}else{
				progState.globalMap = &nMap
			}
		}
	case direction > 0:
		pMap := getGlobalMap()
		if pMap.Next == "null" || len(pMap.Next) == 0 {
			return fmt.Errorf("you're on the last page")
		}else{
			if nMap, err := getMap(pMap.Next); err != nil {
				return err
			}else{
				progState.globalMap = &nMap
			}
		}
	}
	
	return err
}

func printMap(){
	pMap := getGlobalMap()
	if len(pMap.Results) > 0 {
		for _, loc := range pMap.Results{
			fmt.Println(loc.Name)
		}
	}
}

func getMap(url string) (PokeMap, error) {
	if len(url) == 0 {
		url = progState.baseUrl
	}
	var jsonPayload []byte
	// try to pull from cache
	if byteBody, ok := progState.globalMapCache.Get(url); !ok{
		if progState.fDebug{fmt.Println("getMap - Net call!")}
		res, err := http.Get(url)
		if err != nil {
			return PokeMap{}, fmt.Errorf("Error getting url: %w", err)
		}
		defer res.Body.Close()
		if val, err :=io.ReadAll(res.Body); err != nil {
			return PokeMap{}, err
		} else {
			// Add to the cache
			progState.globalMapCache.Add(url, val)
			jsonPayload = val
		}
	}else{
		jsonPayload = byteBody
		if progState.fDebug{fmt.Println("getMap - cache hit!")}
	}
	var newMap PokeMap
	if err := json.Unmarshal(jsonPayload, &newMap); err != nil {
		return PokeMap{}, fmt.Errorf("Error decoding the json: %w", err)
	}
	return newMap, nil
}

