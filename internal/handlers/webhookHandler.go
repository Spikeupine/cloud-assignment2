package handlers

import (
	"assignment-2/internal"
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
