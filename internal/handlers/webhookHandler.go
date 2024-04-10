package handlers

import (
	"assignment-2/internal"
	"encoding/json"
	"net/http"
)

var webhooks = []internal.Webhook{}

//Manages registration of the different types of webhooks.
func WebhookHandler(w http.ResponseWriter, r *http.Request) {

	//Handles POST and GET
	switch r.Method {
	case http.MethodPost:

		//Initializes empty struct of webhook to populate
		webhook := internal.Webhook{}

		//Populates the webhook struct with content from body of response.
		err := json.NewDecoder(r.Body).Decode(&webhook)
		if err != nil {
			http.Error(w, "Error in registration of webhook", http.StatusInternalServerError)
		}
		webhooks = append(webhooks, webhook)
		break

	case http.MethodGet:
		//Returns all webhooks for now
		//todo: Handle the get for specific webhooks
		err := json.NewEncoder(w).Encode(webhooks)
		if err != nil {
			http.Error(w, "Something went wrong when getting all webhooks "+err.Error(), http.StatusServiceUnavailable)
		}
		break

	}
}
