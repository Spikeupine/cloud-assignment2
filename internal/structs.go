package internal

import (
	"cloud.google.com/go/firestore"
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
	Endpoint       string
	ID             string
}

// Webhook is used when we register webhooks for notifications endpoint
type Webhook struct {
	WebhookId string `json:"webhook_id,omitempty" firebase:"webhook_id,omitempty"`
	Url       string `json:"url" firebase:"url"`
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
