package internal

import (
	"cloud.google.com/go/firestore"
	"context"
	"net/http"
	"time"
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

// CurrencyInfo Defines a separate struct for currency information
type CurrencyInfo struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type RegistrationsResponse struct {
	Id         string
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
	Id         string    `json:"id"`
	Country    string    `json:"country"`
	IsoCode    string    `json:"isoCode"`
	LastChange time.Time `json:"lastChange"`
	Features   Features
}

type RegisterMap struct {
	Registers map[string]RegisterRequest `json:"registers"`
}

type Basics struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Client         *firestore.Client
	Ctx            context.Context
	Endpoint       string
	ID             string
}

type Status struct {
	CountriesAPI   int    `json:"countries_api"`
	MeteoAPI       int    `json:"meteo_api"`
	CurrencyAPI    int    `json:"currency_api"`
	NotificationDB int    `json:"notification_db"`
	Webhooks       int    `json:"webhooks"`
	Version        string `json:"version"`
	Uptime         int64  `json:"uptime"`
}

type PopulatedDashboard struct {
	ID            string            `json:"id"`
	Country       string            `json:"country"`
	IsoCode       string            `json:"isoCode"`
	Features      DashboardFeatures `json:"features"`
	LastRetrieval time.Time         `json:"lastRetrieval"`
}

type DashboardFeatures struct {
	Temperature      float64              `json:"temperature,omitempty"`
	Precipitation    float64              `json:"precipitation,omitempty"`
	Capital          string               `json:"capital,omitempty"`
	Coordinates      DashboardCoordinates `json:"coordinates,omitempty"`
	Population       int                  `json:"population,omitempty"`
	Area             float64              `json:"area,omitempty"`
	TargetCurrencies map[string]float64   `json:"targetCurrencies,omitempty"`
}

type DashboardCoordinates struct {
	Latitude  float32 `json:"latitude,omitempty"`
	Longitude float32 `json:"longitude,omitempty"`
}
