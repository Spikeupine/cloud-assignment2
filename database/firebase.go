package database

import (
	"assignment-2/internal"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
)

var (
	ctx            context.Context
	client         *firestore.Client
	collectionName string = "dashboards"
	webhooks              = []internal.Webhook{}
)

func GetClient() *firestore.Client {
	return client
}

func GetContext() context.Context {
	return ctx
}

// FirebaseConnect establishes the connection to firebase
func FirebaseConnect() {
	ctx = context.Background()
	pathToCredentials := "./firebase_privatekey.json"
	opt := option.WithCredentialsFile(pathToCredentials)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Println("Couldn't establish connection to firebase: " + err.Error())
		os.Exit(1)
	}
	// Instantiate client
	client, err = app.Firestore(ctx)
	if err != nil {
		log.Println("Couldn't establish connection to the database" + err.Error())
		os.Exit(1)
	}
}

// FireBaseCloseConnection closes the connection to firebase
func FireBaseCloseConnection() {
	errClose := client.Close()
	if errClose != nil {
		log.Fatal("Closing of the Firebase client failed. Error:", errClose)
	}
}

// Manages registration of the different types of webhooks.

// Registers the different webhooks. Appends the webhook to the collection of other webhooks.
func AddWebhookToCollection(webhook internal.Webhook, collectionName string) (string, error) {
	webhooks = append(webhooks, webhook)

	// Adds the webhook to the collection in firestore database
	docRef, _, err := client.Collection(collectionName).Add(ctx, webhook)
	if err != nil {
		return "", err
	}
	return docRef.ID, nil

}

// Gets the webhook requested by its ID
func GetWebhook(w http.ResponseWriter, r *http.Request, webhookID string) internal.Webhook {
	var webhookInQuestion internal.Webhook
	for _, webhook := range webhooks {
		if webhook.WebhookId == webhookID {
			webhookInQuestion = webhook
			return webhookInQuestion
		}
	}
	return webhookInQuestion
}

func DeleteWebhook(w http.ResponseWriter, r *http.Request, webhookId string) {

	for i, webhook := range webhooks {
		if webhook.WebhookId == webhookId {
			webhooks = append(webhooks[:i], webhooks[i+1:]...)
			return
		}
	}

}

// Returns all webhooks for now
func getAllWebhooks(w http.ResponseWriter, r *http.Request) []internal.Webhook {
	err := json.NewEncoder(w).Encode(webhooks)
	if err != nil {
		http.Error(w, "Something went wrong when getting all webhooks "+err.Error(), http.StatusServiceUnavailable)
	}
	return webhooks
}
