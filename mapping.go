package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PokeLocation interface {
    PokeMap | PokeMapLocation
}
type PokeMap struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type PokeMapLocation struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func isValidLocation(loc string) bool {
	for _, elm := range getGlobalMap().Results{
		if elm.Name == loc{
			return true
		}
	}
	return false
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
			if nMap, err := getMap[PokeMap](progState.baseUrl); err != nil {
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
			if nMap, err := getMap[PokeMap](pMap.Previous); err != nil {
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
			if nMap, err := getMap[PokeMap](pMap.Next); err != nil {
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

func doMapRequest(url string) ([]byte, error){
	if byteBody, ok := progState.globalMapCache.Get(url); !ok {
		if progState.fDebug {fmt.Println("doMapRequest - Net call!")}
		res, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("Error getting url: %s err: %w", url, err)
		}
		defer res.Body.Close()
		if val, err := io.ReadAll(res.Body); err != nil {
			return nil, fmt.Errorf("Error getting json body err: %w", err)
		} else {
			// Add to cache
			progState.globalMapCache.Add(url, val)
			return val, nil
		}
	} else {
		return byteBody, nil
	}
}


func getMap[V PokeLocation](url string) (V, error) {
	var zeroRet V
	if len(url) == 0 {
		return zeroRet, fmt.Errorf("")
	}
	jsonPayload, err := doMapRequest(url)
	if err != nil {
		return zeroRet, err
	}
	var newMap V
	if err := json.Unmarshal(jsonPayload, &newMap); err != nil {
		return zeroRet, fmt.Errorf("Error decoding the json: %w", err)
	}
	return newMap, nil
}

