package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jamesonhm/pokedexcli/internal/pokecache"
)

const BaseApi string = "https://pokeapi.co/api/v2/"

type Client struct {
	baseUrl    string
	httpClient http.Client
	cache      *pokecache.Cache
}

func NewClient(timeout time.Duration, cache *pokecache.Cache) *Client {
	return &Client{
		baseUrl: BaseApi,
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: cache,
	}
}

func (c *Client) ListLocations(pageURL string) (LocationArea, error) {
	url := c.baseUrl + "location-area"
	if pageURL != "" {
		url = pageURL
	}
	fmt.Println("-URL: ", url)

	if data, ok := c.cache.Get(url); ok {
		locAreas := LocationArea{}
		err := json.Unmarshal(data, &locAreas)
		if err != nil {
			return LocationArea{}, err
		}

		fmt.Println("-from cache request")
		return locAreas, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationArea{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationArea{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationArea{}, err
	}

	locAreas := LocationArea{}
	err = json.Unmarshal(body, &locAreas)
	if err != nil {
		return LocationArea{}, err
	}

	c.cache.Add(url, body)
	fmt.Println("-from get request")
	return locAreas, nil
}

func (c *Client) LocationDetails(name string) (LocationDetail, error) {
	url := c.baseUrl + "location-area/" + name
	fmt.Println("-URL: ", url)

	if data, ok := c.cache.Get(url); ok {
		locDetail := LocationDetail{}
		err := json.Unmarshal(data, &locDetail)
		if err != nil {
			return LocationDetail{}, err
		}

		fmt.Println("-from cache request")
		return locDetail, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationDetail{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationDetail{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationDetail{}, err
	}

	locDetail := LocationDetail{}
	err = json.Unmarshal(body, &locDetail)
	if err != nil {
		return LocationDetail{}, err
	}

	c.cache.Add(url, body)
	fmt.Println("-from get request")
	return locDetail, nil

}

type LocationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

type LocationDetail struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
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
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}
