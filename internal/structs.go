package internal

import (
	"cloud.google.com/go/firestore"
	"context"
	"net/http"
	"time"
)

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
	WebhookId string `json:"webhook_id,omitempty" firebase:"webhook_id"`
	Url       string `json:"url" firebase:"url"`
	Country   string `json:"iso" firebase:"iso"`
	Event     string `json:"event,omitempty" firebase:"event"`
	Calls     int    `json:"calls" firebase:"calls"`
}

type InvokeWebhook struct {
	WebhookId string `json:"webhook_id,omitempty" firebase:"webhook_id"`
	Url       string `json:"url" firebase:"url"`
	Country   string `json:"iso" firebase:"iso"`
	Event     string `json:"event,omitempty" firebase:"event"`
	Calls     int    `json:"calls" firebase:"calls"`
	Time      time.Time
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
	Temperature      float64               `json:"temperature,omitempty"`
	Precipitation    float64               `json:"precipitation,omitempty"`
	Capital          string                `json:"capital,omitempty"`
	Coordinates      *DashboardCoordinates `json:"coordinates,omitempty"`
	Population       int                   `json:"population,omitempty"`
	Area             float64               `json:"area,omitempty"`
	TargetCurrencies map[string]float64    `json:"targetCurrencies,omitempty"`
}

type DashboardCoordinates struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}
