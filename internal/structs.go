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

// Webhook is used when we register webhooks for notifications endpoint
type Webhook struct {
	WebhookId string `json:"webhook_id" firebase:"webhook_id"`
	Url       string `json:"url" firebase:"url"`
	Country   string `json:"iso" firebase:"iso"`
	Event     string `json:"dashboard,omitempty" firebase:"dashboard"`
}

// WebhookResponse is used in response of the notifications endpoint
type WebhookResponse struct {
	WebhookId string `json:"webhook_id"`
	Url       string `json:"url"`
	Dashboard string `json:"dashboard"`
}

// WebhookInvocation is used when invoking a webhook
type WebhookInvocation struct {
	WebhookId string `json:"webhook_id"`
	Dashboard string `json:"dashboard"`
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
