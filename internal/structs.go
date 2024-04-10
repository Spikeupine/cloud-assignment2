package internal

import (
	"github.com/google/uuid"
	"time"
)

// RestCountriesStruct Struct for holding information about country from restcountries API.
type RestCountriesStruct struct {
	Name        string                  `json:"common"`
	Iso         string                  `json:"cca2"`
	Population  int                     `json:"population"`
	Capital     []string                `json:"capital"`
	Currencies  map[string]CurrencyInfo `json:"currencies"` // Change the data type to a map
	Coordinates []float64               `json:"latlng"`
}

// CurrencyInfo Defines a separate struct for currency information
type CurrencyInfo struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type RegistrationsResponse struct {
	Id         uuid.UUID
	LastChange time.Time
}

type Features struct {
	Temperature      bool     `json:"temperature,omitempty"`
	Precipitation    bool     `json:"precipitation,omitempty"`
	Capital          bool     `json:"capital,omitempty"`
	Coordinates      bool     `json:"coordinates,omitempty"`
	Population       bool     `json:"population,omitempty"`
	Area             bool     `json:"area,omitempty"`
	TargetCurrencies []string `json:"targetCurrencies,omitempty"`
}

type RegisterRequest struct {
	Country  string `json:"country"`
	IsoCode  string `json:"isoCode"`
	Features Features
}
