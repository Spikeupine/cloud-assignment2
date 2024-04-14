package handlers

import (
	"assignment-2/database"
	"assignment-2/internal"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

var collectionName = "dashboards"

// NotificationsHandler handles the notifications endpoint
func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		WebhookRegistration(w, r, collectionName)
	case http.MethodDelete:
		pathValue := r.PathValue("id")
		if pathValue == "" {
			return
		} else {
			DeleteWebhook(w, r, collectionName, pathValue)
		}

	case http.MethodGet:
		pathValue := r.PathValue("id")
		if pathValue == "" {
			GetWebhooks(w, r, collectionName)
		} else {
			getWebhook(w, r, pathValue)
		}
	default:
		http.Error(w, "Method '"+r.Method+"' not supported. Currently method '"+http.MethodPost+
			"', '"+http.MethodDelete+"' and '"+http.MethodGet+"' is supported.", http.StatusNotImplemented)
	}
}

// WebhookRegistration Handles the post requests, which registers webhooks to the firebase database collection.
func WebhookRegistration(w http.ResponseWriter, r *http.Request, collectionName string) {

	// Decodes the webhook from the body of request
	webhook := internal.Webhook{}

	//Sets id of webhook with help from uuid.
	id, _ := uuid.NewUUID()
	webhook.WebhookId = id.String()

	//Decodes rest of the response into struct.
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		http.Error(w, "Error during decoding body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// adds the webhook to the database via methods in firebase-
	err := database.AddWebhookToCollection(webhook, collectionName)
	if err != nil {
		http.Error(w, "Error when adding webhook: "+webhook.Url+" to firebase collection "+
			collectionName+": "+err.Error(), http.StatusFailedDependency)
		return
	}

	// returns webhook id as response
	output := map[string]string{
		"webhook_id": webhook.WebhookId,
	}
	w.Header().Add(internal.ApplicationJson, internal.ContentTypeJson)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, "Error during encoding: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteWebhook handles DELETE method, and deletes the webhook with specified if from database.
func DeleteWebhook(w http.ResponseWriter, r *http.Request, collectionName string, webhookId string) {

	// deletes the webhook
	if err, sc := database.DeleteTheWebhook(collectionName, webhookId); err != nil {
		http.Error(w, "Error when deleting webhook with id '"+webhookId+"' :"+err.Error(), sc)
		return
	}
	http.Error(w, "Webhook '"+webhookId+"' was deleted.", http.StatusOK)
}

// GetWebhooks is the http.GET method that either returns all webhooks if nothing is specified, or one specific
// if there is an id inserted.
func GetWebhooks(w http.ResponseWriter, r *http.Request, collectionName string) {

	//Method from database-firebase that returns a list of struct of all the different webhooks
	webhooks := database.GetAllWebhooks()

	// encodes the resulting list to writer
	w.Header().Add(internal.ApplicationJson, internal.ContentTypeJson)
	if err := json.NewEncoder(w).Encode(webhooks); err != nil {
		http.Error(w, "Error when encoding webhooks to user: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getWebhook(w http.ResponseWriter, r *http.Request, webhookId string) {

	// get webhook from database
	webhook, err := database.GetWebhook(w, r, webhookId)
	if err != nil {
		http.Error(w, "Error when getting specified webhook "+err.Error(), http.StatusBadRequest)
	}

	// encode the resulting webhook response
	w.Header().Add("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(webhook); err != nil {
		http.Error(w, "Error during encoding: "+err.Error(), http.StatusInternalServerError)
		return
	}

}
