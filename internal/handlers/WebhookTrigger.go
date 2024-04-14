package handlers

import (
	"assignment-2/database"
	"assignment-2/internal"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// IncrementCallCount increments all webhooks that is subscribing to the countryCode
func IncrementCallCount(w http.ResponseWriter, countryCode string) {

	//Method that returns a list of all the webhooks registered to specified collection, and error if any.
	webhooks, err := database.GetAllWebhooks(w, collectionNameWebhooks)
	if err != nil {
		http.Error(w, "Error when getting a list of all the webhooks :"+err.Error(), http.StatusInternalServerError)
	}

	// Goes through all the webhooks in the list.
	for _, webhook := range webhooks {

		// increments the count of calls on each of that specified country code, so every time a country with that
		// iso code is registered to registrations endpoint, it adds one to the count.
		if countryCode == webhook.Country || countryCode == "" {
			webhook.Calls++

			//If-check to find the first time the webhook is registered. Takes it to invocation-method below.
			if webhook.Calls == 1 {
				invokeWebhook(webhook.Url, internal.Webhook{
					WebhookId: webhook.WebhookId,
					Country:   webhook.Country,
					Calls:     webhook.Calls,
				})
			}

			//Uses method to update the call count on the webhook in question.
			err := database.UpdateTheCallCount(collectionNameWebhooks, webhook.WebhookId, webhook.Calls)
			if err != nil {
				http.Error(w, "Error when updating call count in webhook trigger :"+err.Error(), http.StatusNotFound)
				return
			}
		}
	}
}

// invokeWebhook invokes a POST request to the webhook at url specified in registered webhook with the body data
func invokeWebhook(url string, data internal.Webhook) {

	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	//Posts the content as content type json to the url. In our instance, a webhook.site page.
	resp, err := http.Post(url, internal.ContentTypeJson, bytes.NewBuffer(payload))
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
