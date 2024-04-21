package handlers

import (
	"assignment-2/database"
	"assignment-2/internal"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

/*
	The webhooks part of this assignment is inspired by last year's delivery from Mathias Iversen, Oda Steinsholt and
	Eirik Gjertsen. The webhooks part of the code was written by Mathias Iversen, and our use of these parts is blessed
	by him. For project, please see link: https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2023-workspace/mathiaiv/assignment-2
*/

// IncrementCallCount increments all webhooks that is subscribing to the countryCode
func IncrementCallCount(w http.ResponseWriter, webhook internal.Webhook) {

	webhook.Calls++

	invokeWebhook(webhook.Url, internal.Webhook{
		WebhookId: webhook.WebhookId,
		Country:   webhook.Country,
		Calls:     webhook.Calls})

	//Uses method to update the call count on the webhook in question.
	err := database.UpdateTheCallCount(collectionNameWebhooks, webhook.WebhookId, webhook.Calls)
	if err != nil {
		http.Error(w, "Error when updating call count in webhook trigger :"+err.Error(), http.StatusNotFound)
		return
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

func EventWebhook(w http.ResponseWriter, iso string, method string) {
	webhooks, err := database.GetAllWebhooks(w, "webhooks")
	if err != nil {
		http.Error(w, "Error when putting webhooks in list in WebhookTrigger :"+err.Error(), http.StatusInternalServerError)
	} else {
		for _, webhook := range webhooks {
			if webhook.Country == iso && webhook.Event == method || webhook.Country == "" {
				switch webhook.Event {
				case "REGISTER":
					IncrementCallCount(w, webhook)
				case "CHANGE":
					IncrementCallCount(w, webhook)
				case "DELETE":
					IncrementCallCount(w, webhook)
				case "INVOKE":
					IncrementCallCount(w, webhook)
				}
			}
		}
	}
}
