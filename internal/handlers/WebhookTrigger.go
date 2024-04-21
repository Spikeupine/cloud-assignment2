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
	Eirik Gjertsen. The webhooks part of the code was written by Mathias Iversen, and our inspiration of their
	solution is blessed by him. For project, please see link:
	https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2023-workspace/mathiaiv/assignment-2
*/

// IncrementCallCount increments all webhooks that is subscribing to the countryCode or non-specified country code.
func IncrementCallCount(w http.ResponseWriter, webhook internal.Webhook) {

	//Increments call.
	webhook.Calls++

	//Invokes webhook, which means it writes the content of the webhook to the url given in webhook.
	invokeWebhook(webhook.Url, internal.Webhook{
		WebhookId: webhook.WebhookId,
		Country:   webhook.Country,
		Calls:     webhook.Calls})

	//Uses method to update the call count on the webhook in question in its firebase collection.
	err := database.UpdateTheCallCount(collectionNameWebhooks, webhook.WebhookId, webhook.Calls)
	if err != nil {
		http.Error(w, "Error when updating call count in webhook trigger :"+err.Error(), http.StatusNotFound)
		return
	}

}

// invokeWebhook invokes a POST request to the webhook at url specified in registered webhook with the body data
func invokeWebhook(url string, data internal.Webhook) {

	//Makes the webhook content to json, and calls it payload.
	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	//Posts the content as content type json to the url. If you for example have a webhooks.site as url, you can see
	//the content there.
	resp, err := http.Post(url, internal.ContentTypeJson, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf(err.Error())
		return
	}

	//Checks the response code. If it is not 200 OK and response is not nil, it prints the responses status code.
	if resp != nil && resp.StatusCode != http.StatusOK {
		log.Printf("unexpected status code : %d", resp.StatusCode)
		return
	}

	if err := resp.Body.Close(); err != nil {
		log.Printf(err.Error())
	}
}

// EventWebhook takes in iso code from webhook and the webhook's method. It retrieves a list of all the webhooks
// registered in firebase, and loops over it. This is in order to call all the webhooks that meets the criteria.
// For example, a webhook of iso "NO" with "REGISTER" method should get its call count incremented if a new registration
// of a dashboard of "NO" is registered. EventWebhook is called upon from the places that a webhook should get action
// from, for example when a new configuration is added. It loops over all the webhooks to update the call count of them
// all. If the webhook has an empty string as iso, then all instances of its method should be registered and incremented
func EventWebhook(w http.ResponseWriter, iso string, method string) {

	//Gets the list of webhooks, or error.
	webhooks, err := database.GetAllWebhooks(w, "webhooks")
	if err != nil {
		http.Error(w, "Error when putting webhooks in list in WebhookTrigger :"+err.Error(), http.StatusInternalServerError)
	} else {
		for _, webhook := range webhooks {
			//If the country code gotten from the EventWebhook method matches the country code of the webhook ranged
			//over, and the method is the same, it should increment that call. Or it should also increment if the iso
			//code is empty.
			if webhook.Country == iso && webhook.Event == method || webhook.Country == "" {
				//if the iso code is empty in the webhook, then it is the webhook's method that decides what should
				//be incremented.
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
