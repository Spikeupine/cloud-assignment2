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

	//Gets all webhooks stored in webhooks, a list of webhook, or error.
	webhooks, err := database.GetAllWebhooks(w, collectionNameWebhooks)
	if err != nil {
		http.Error(w, "Error when getting a list of all the webhooks :"+err.Error(), http.StatusInternalServerError)
	}

	// Goes through all the webhooks in the list.
	for _, webhook := range webhooks {
		// increments the count of calls on each of that specified country code
		if countryCode == webhook.Country || countryCode == "" {
			webhook.Calls++
			if webhook.Calls == 1 {
				invokeWebhook(webhook.Url, internal.Webhook{
					WebhookId: webhook.WebhookId,
					Country:   webhook.Country,
					Calls:     webhook.Calls,
				})
			}
			err := database.UpdateTheCallCount(collectionNameWebhooks, webhook.WebhookId, webhook.Calls)
			if err != nil {
				http.Error(w, "Error when updating call count in webhook trigger :"+err.Error(), http.StatusNotFound)
				return
			}
		}
	}
}

// invokeWebhook invokes a POST request to the webhook at url with the body data
func invokeWebhook(url string, data internal.Webhook) {

	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf(err.Error())
		return
	}

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
