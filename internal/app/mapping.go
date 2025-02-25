package pokedex

import (
	"fmt"
)

func isValidLocation(loc string) bool {
	for _, elm := range getGlobalMap().Results{
		if elm.Name == loc{
			return true
		}
	}
	return false
}

func getGlobalMap() PokeMap{
	if progState.GlobalMap != nil{
		return *progState.GlobalMap
	}
	newMap := PokeMap{}
	progState.GlobalMap = &newMap
	return *progState.GlobalMap
}

func getMapUrl(loc string) string{
	return fmt.Sprintf("%slocation-area/%s", progState.BaseUrl, loc)
}

func updateMap(direction int) error{
	var err error
	switch {
	case 0 == direction:
		pMap := getGlobalMap()
		if len(pMap.Results) == 0 {
			// don't have a map yet, get one
			if nMap, err := getMap[PokeMap](getMapUrl("")); err != nil {
				return err
			}else{
				progState.GlobalMap = &nMap
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
				progState.GlobalMap = &nMap
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
				progState.GlobalMap = &nMap
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


func getMap[V PokeLocation](url string) (V, error) {
	return getNetObject[V](url)
}

