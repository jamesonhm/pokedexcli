package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const BaseApi string = "https://pokeapi.co/api/v2/"

type Client struct {
	baseUrl    string
	httpClient http.Client
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		baseUrl: BaseApi,
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) ListLocations(pageURL string) (LocationArea, error) {
	url := c.baseUrl + "location-area"
	if pageURL != "" {
		url = pageURL
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
