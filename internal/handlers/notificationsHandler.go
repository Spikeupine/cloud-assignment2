package handlers

import (
	"assignment-2/database"
	"assignment-2/internal"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

var (
	collectionNameWebhooks = "webhooks"
)

// NotificationsHandler handles the notifications endpoint
func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		WebhookRegistration(w, r, collectionNameWebhooks)
	case http.MethodDelete:
		pathValue := r.PathValue("id")
		if pathValue == internal.Empty {
			return
		} else {
			DeleteWebhook(w, r, collectionNameWebhooks, pathValue)
		}

	case http.MethodGet:
		pathValue := r.PathValue(internal.Id)
		if pathValue == internal.Empty {
			GetWebhooks(w, r, collectionNameWebhooks)
		} else {
			GetWebhook(w, collectionNameWebhooks, r, pathValue)
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
	err, sc := database.DeleteTheWebhook(collectionName, webhookId)
	if err != nil {
		http.Error(w, "Error when deleting webhook with id '"+webhookId+"' :"+err.Error(), sc)
	}
	if sc == 0 {
		println("successfully deleted")
	}
}

// GetWebhooks is the http.GET method that either returns all webhooks if nothing is specified, or one specific
// if there is an id inserted.
func GetWebhooks(w http.ResponseWriter, r *http.Request, collectionName string) {

	webhooks, err := database.GetAllWebhooks(w, collectionName)
	if err != nil {
		http.Error(w, "Error when receiveing variable that stores all the webhooks :"+err.Error(), http.StatusInternalServerError)
		return
	}
	// encodes the resulting list to writer
	w.Header().Add(internal.ApplicationJson, internal.ContentTypeJson)
	if err := json.NewEncoder(w).Encode(webhooks); err != nil {
		http.Error(w, "Error when encoding webhooks to user: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func GetWebhook(w http.ResponseWriter, collectionName string, r *http.Request, webhookId string) error {

	// Get webhook from database
	webhook, err := database.GetWebhook(collectionName, webhookId)
	if err != nil {
		http.Error(w, "Error when getting specified webhook: "+err.Error(), http.StatusBadRequest)
		return err
	}

	// Encode the resulting webhook response
	w.Header().Set(internal.ContentTypeJson, internal.ApplicationJson)
	if err := json.NewEncoder(w).Encode(webhook); err != nil {
		http.Error(w, "Error during encoding: "+err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}
