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
func GetWebhook(w http.ResponseWriter, r *http.Request, webhookID string) internal.Webhook {
	var webhookInQuestion internal.Webhook
	for _, webhook := range webhooks {
		if webhook.WebhookId == webhookID {
			webhookInQuestion = webhook
			return webhookInQuestion
		}
	}
	return webhookInQuestion
}

func DeleteWebhook(w http.ResponseWriter, r *http.Request, webhookId string) {

	for i, webhook := range webhooks {
		if webhook.WebhookId == webhookId {
			webhooks = append(webhooks[:i], webhooks[i+1:]...)
			return
		}
	}

}

// Returns all webhooks for now
func getAllWebhooks(w http.ResponseWriter, r *http.Request) []internal.Webhook {
	err := json.NewEncoder(w).Encode(webhooks)
	if err != nil {
		http.Error(w, "Something went wrong when getting all webhooks "+err.Error(), http.StatusServiceUnavailable)
	}
	return webhooks
}
