package handlers

import (
	"assignment-2/database"
	"assignment-2/external/router"
	"assignment-2/internal"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

var (
	collectionNameWebhooks = database.WebhookCollection
)

// NotificationsHandler handles the notifications endpoint. It directs the http requests of different types to their
// appropriate methods. Types supported are: POST, DELETE and GET. If the method is anything else than this, an http
// error will be thrown.
func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		WebhookRegistration(w, r, collectionNameWebhooks)
	case http.MethodDelete:

		//Finds the path value based on its placeholder that was inserted in the handler in main.
		pathValue := r.PathValue(internal.Id)
		if pathValue == internal.Empty {
			return
		} else {
			DeleteWebhook(w, collectionNameWebhooks, pathValue)
		}

	case http.MethodGet:
		pathValue := r.PathValue(internal.Id)

		//If no path value, it is to show all the webhooks. See method
		if pathValue == internal.Empty {
			GetWebhooks(w, collectionNameWebhooks)
		} else {
			//If path value, show only one webhook. Path value is passed on to method, and it is the webhook id.
			GetWebhook(w, collectionNameWebhooks, pathValue)
		}
		//If none of the methods above, send an http error.
	default:
		http.Error(w, "Method '"+r.Method+"' not supported. Currently method '"+http.MethodPost+
			"', '"+http.MethodDelete+"' and '"+http.MethodGet+"' is supported.", http.StatusMethodNotAllowed)
	}
}

// WebhookRegistration Handles the post requests, which creates a webhook out of the body content if it passes
// the specified tests. It passes this webhook on to a method that writes the webhook to the firebase collection of
// webhooks.
func WebhookRegistration(w http.ResponseWriter, r *http.Request, collectionName string) {

	// Creates an empty webhook to populate.
	webhook := internal.Webhook{}

	//Sets id of webhook with help from uuid.
	id, _ := uuid.NewUUID()
	webhook.WebhookId = id.String()

	//Decodes response from response body into struct so the webhook gets populated.
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		http.Error(w, "Error during decoding body: "+err.Error(), http.StatusBadRequest)
		return
	}
	//GetCountriesObject checks if the country code written exists. When creating a webhook, you can either register
	//it to one country or all countries. FOr all countries, you don't set a country code. If you send an empty country
	//code to GetCountriesObject, it will return an error.
	country, errorCountry := router.GetCountriesObject("", webhook.Country)

	//Checks if the error from GetCountriesObjet is not empty while at the same the the iso code field of webhook is
	// not an empty string. This handles instances where invalid iso codes are registered, or misspelled.
	if errorCountry != nil && webhook.Country != "" {
		http.Error(w, "Error in country code of webhook to be registered :"+errorCountry.Error(), http.StatusBadRequest)
		return
	} else if country.Cca2 == webhook.Country || webhook.Country == "" {
		// adds the webhook to the database via methods in firebase
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
	} else {
		http.Error(w, "Error in country code of webhook. Webhook not registered.", http.StatusBadRequest)
		return
	}
}

// DeleteWebhook handles DELETE method, and deletes the webhook with specified id from database.
func DeleteWebhook(w http.ResponseWriter, collectionName string, webhookId string) {

	if webhookId != internal.Empty {
		// deletes the webhook via database method.
		err, sc := database.DeleteTheWebhook(collectionName, webhookId)
		w.Header().Add(internal.ApplicationJson, internal.ContentTypeJson)
		http.Error(w, "Successfully deleted webhook with id :"+webhookId, http.StatusOK)
		if err != nil {
			http.Error(w, "Error when deleting webhook with id '"+webhookId+"' :"+err.Error(), sc)
		}
	} else {
		http.Error(w, "Error in DeleteWebhook: Id cannot be empty ", http.StatusBadRequest)
	}
}

// GetWebhooks is the http.GET method that either returns all webhooks if no webhook id is specified, or one specific
// if there is an id inserted in the url.
func GetWebhooks(w http.ResponseWriter, collectionName string) {

	//All the webhooks are returned in a list.
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

// GetWebhook retrieves one single specified webhook by retrieving it from the database by its id, and displaying it
// to the user.
func GetWebhook(w http.ResponseWriter, collectionName string, webhookId string) {

	// Get webhook from database
	webhook, err := database.GetWebhook(collectionName, webhookId)
	if err != nil {
		http.Error(w, "Error when getting specified webhook: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Encode the resulting webhook response
	w.Header().Set(internal.ContentTypeJson, internal.ApplicationJson)
	if err := json.NewEncoder(w).Encode(webhook); err != nil {
		http.Error(w, "Error during encoding: "+err.Error(), http.StatusInternalServerError)
		return
	}
	return
}
