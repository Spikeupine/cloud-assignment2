package internal

import (
	"cloud.google.com/go/firestore"
	"context"
	"net/http"
)

// Struct for holding information about country from restcountries API.
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

type Basics struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Client         *firestore.Client
	Ctx            context.Context
	Endpoint       string
	ID             string
}

// Request for an individual dashboard identified by its ID (same as the corresponding configuration ID)
type PopulatedDashboard struct {
	Country  string `json:"country"`
	IsoCode  string `json:"isoCode"`
	Features struct {
		Temperature   float64 `json:"temperature"`
		Precipitation float64 `json:"precipitation"`
		Capital       string  `json:"capital"`
		Coordinates   struct {
			Latitude  float32 `json:"latitude"`
			Longitude float32 `json:"longitude"`
		} `json:"coordinates"`
		Population       int                `json:"population"`
		Area             float64            `json:"area"`
		TargetCurrencies map[string]float64 `json:"targetCurrencies"`
	} `json:"features"`
	LastRetrieval string `json:"lastRetrieval"`
}
