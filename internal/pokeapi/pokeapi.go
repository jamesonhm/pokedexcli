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

func cachedApiRequest[T any](c *Client, url string, respT *T) error {
	if data, ok := c.cache.Get(url); ok {
		if err := json.Unmarshal(data, &respT); err != nil {
			return err
		}
		return nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if string(body) == "Not Found" {
		return fmt.Errorf("resource not found for %s", url)
	}
	err = json.Unmarshal(body, &respT)
	if err != nil {
		return err
	}

	c.cache.Add(url, body)
	return nil
}

func (c *Client) ListLocations(pageURL string) (LocationArea, error) {
	url := c.baseUrl + "location-area"
	if pageURL != "" {
		url = pageURL
	}
	locAreas := LocationArea{}
	if err := cachedApiRequest(c, url, &locAreas); err != nil {
		return locAreas, err
	}
	return locAreas, nil
}

func (c *Client) LocationDetails(name string) (LocationDetail, error) {
	url := c.baseUrl + "location-area/" + name
	locDetail := LocationDetail{}
	if err := cachedApiRequest(c, url, &locDetail); err != nil {
		return locDetail, err
	}
	return locDetail, nil
}

func (c *Client) Pokemon(name string) (Pokemon, error) {
	url := c.baseUrl + "pokemon/" + name
	pokemon := Pokemon{}
	if err := cachedApiRequest(c, url, &pokemon); err != nil {
		return pokemon, err
	}
	return pokemon, nil
}
