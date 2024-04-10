package handlers

import (
	"assignment-2/internal"
	"encoding/json"
	"net/http"
)

var collectionName string = "dashboards"

var webhooks = []internal.Webhook{}

// Manages registration of the different types of webhooks.
func WebhookHandler(w http.ResponseWriter, r *http.Request) {

	//Handles POST and GET
	switch r.Method {
	case http.MethodPost:
		RegisterWebhook(w, r, collectionName)
		break
	case http.MethodGet:
		GetWebhook(w, r, collectionName)
		break
	case http.MethodDelete:
		DeleteWebhook(w, r, collectionName)
		break
	default:
		http.Error(w, "Method '"+r.Method+"' not supported. Currently method '"+http.MethodPost+
			"', '"+http.MethodDelete+"' and '"+http.MethodGet+"' is supported.", http.StatusNotImplemented)
	}
}

// Registers the different webhooks. Appends the webhook to the collection of other webhooks.
func RegisterWebhook(w http.ResponseWriter, r *http.Request, collectionName string) {
	//Initializes empty struct of webhook to populate
	webhook := internal.Webhook{}

	//Populates the webhook struct with content from body of response.
	err := json.NewDecoder(r.Body).Decode(&webhook)
	if err != nil {
		http.Error(w, "Error in registration of webhook", http.StatusInternalServerError)
	}
	webhooks = append(webhooks, webhook)
}

// Gets the webhooks registered
// Todo: Write method to a mathod that just returns one webhook based on imput.
func GetWebhook(w http.ResponseWriter, r *http.Request, collectionName string) {
	//Returns all webhooks for now
	//todo: Handle the get for specific webhooks
	err := json.NewEncoder(w).Encode(webhooks)
	if err != nil {
		http.Error(w, "Something went wrong when getting all webhooks "+err.Error(), http.StatusServiceUnavailable)
	}
}
