package pokedex

import(
	"fmt"
	"io"
	"net/http"
	"encoding/json"
)



func doGetRequest(url string) ([]byte, error){
	if byteBody, ok := progState.GlobalMapCache.Get(url); !ok {
		if progState.Debug {fmt.Println("doMapRequest - Net call!")}
		res, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("Error getting url: %s err: %w", url, err)
		}
		defer res.Body.Close()
		if val, err := io.ReadAll(res.Body); err != nil {
			return nil, fmt.Errorf("Error getting json body err: %w", err)
		} else {
			// Add to cache
			progState.GlobalMapCache.Add(url, val)
			return val, nil
		}
	} else {
		return byteBody, nil
	}
}

func getNetObject[V PokeNetRequest](url string) (V, error){
	var zeroRet V
	if len(url) == 0 {
		return zeroRet, fmt.Errorf("")
	}
	jsonPayload, err := doGetRequest(url)
	if err != nil {
		return zeroRet, err
	}
	var payload V
	if err := json.Unmarshal(jsonPayload, &payload); err != nil {
		return zeroRet, fmt.Errorf("Error decoding the json: %w", err)
	}
	return payload, nil
}
