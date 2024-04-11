package handlers

import (
	"assignment-2/database"
	"assignment-2/internal"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
)

var collectionName = "dashboards"

// NotificationsHandler handles the notifications endpoint
func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		WebhookRegistration(w, r, collectionName)
	case http.MethodDelete:
		DeleteWebhook(w, r, collectionName)
	case http.MethodGet:
		GetWebhooks(w, r, collectionName)
	default:
		http.Error(w, "Method '"+r.Method+"' not supported. Currently method '"+http.MethodPost+
			"', '"+http.MethodDelete+"' and '"+http.MethodGet+"' is supported.", http.StatusNotImplemented)
	}
}

// WebhookRegistration Handles the post requests, which registers webhooks to the firebase database collection.
func WebhookRegistration(w http.ResponseWriter, r *http.Request, collectionName string) {

	// Decodes the webhook from the body of request
	webhook := internal.Webhook{}
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		http.Error(w, "Error during decoding body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if validISO, _ := regexp.MatchString("^[a-zA-Z]{3}$", webhook.Event); !validISO && webhook.Event != "" {
		http.Error(w, "Error: Invalid ISO code", http.StatusBadRequest)
		return
	}
	webhook.Event = strings.ToUpper(webhook.Event)

	// adds the webhook to the database via methods in firebase-
	webhookId, err := database.AddWebhookToCollection(webhook, collectionName)
	if err != nil {
		http.Error(w, "Error when adding webhook: "+webhook.Url+" to firebase collection "+
			collectionName+": "+err.Error(), http.StatusFailedDependency)
		return
	}

	// returns webhook id as response
	output := map[string]string{
		"webhook_id": webhookId,
	}
	w.Header().Add(internal.ApplicationJson, internal.ContentTypeJson)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, "Error during encoding: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

