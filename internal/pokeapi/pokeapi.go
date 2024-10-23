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

type LocationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}
