package main

const BaseApi string = "https://pokeapi.co/api/v2/"

type config struct {
	next     string
	previous string
}

func NewConfig() *config {
	return &config{}
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

func NewLocationArea() LocationArea {
	return LocationArea{}
}
