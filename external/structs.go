package external

type MeteoObject struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Time      Weather `json:"hourly"`
}

type Weather struct {
	Time          []int     `json:"time"`
	Temperature   []float64 `json:"temperature_2m"`
	Precipitation []float64 `json:"precipitation"`
}

type CountriesObject struct {
	Location struct {
		Latlng []float64 `json:"latlng"`
	} `json:"capitalInfo"`
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
	Cca2       string              `json:"cca2"`
	Currencies map[string]Currency `json:"currencies"`
	Area       float64             `json:"area"`
	Population int                 `json:"population"`
	Capital    []string            `json:"capital"`
}

type Currency struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}
