package handlers

import (
	"assignment-2/internal"
	"bytes"
	"encoding/json"
	"log"
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

// invokeWebhook invokes a POST request to the webhook at url with the body data
func invokeWebhook(url string, data structs.WebhookInvocation) {

	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf(err.Error())
		return
	}

	if resp != nil && resp.StatusCode != http.StatusOK {
		log.Printf("unexpected status code: %d", resp.StatusCode)
		return
	}

	if err := resp.Body.Close(); err != nil {
		log.Printf(err.Error())
	}
}

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

func GetWebhook(w http.ResponseWriter, r *http.Request, collectionName string) {
	//Returns all webhooks for now
	//todo: Handle the get for specific webhooks
	err := json.NewEncoder(w).Encode(webhooks)
	if err != nil {
		http.Error(w, "Something went wrong when getting all webhooks "+err.Error(), http.StatusServiceUnavailable)
	}
}
