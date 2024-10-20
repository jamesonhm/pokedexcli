package pokeapi

const BaseApi string = "https://pokeapi.co/api/v2/"

type Config struct {
	next     string
	previous string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Previous() string {
	return c.previous
}

func (c *Config) UpdatePrev(new string) {
	c.previous = new
}

func (c *Config) Next() string {
	return c.next
}

func (c *Config) UpdateNext(new string) {
	c.next = new
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
