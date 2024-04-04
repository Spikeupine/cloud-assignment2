package internal

//Struct for holding information about country from restcountries API.
type RestCountriesStruct struct {
	Name        string                  `json:"common"`
	Iso         string                  `json:"cca2"`
	Population  int                     `json:"population"`
	Capital     []string                `json:"capital"`
	Currencies  map[string]CurrencyInfo `json:"currencies"` // Change the data type to a map
	Coordinates []float64               `json:"latlng"`
}

// Defines a separate struct for currency information
type CurrencyInfo struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}
