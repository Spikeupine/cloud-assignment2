package database

import (
	"assignment-2/internal"
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
)

var (
	ctx      context.Context
	client   *firestore.Client
	webhooks = []internal.Webhook{}
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
func AddWebhookToCollection(webhook internal.Webhook, collectionName string) error {
	webhooks = append(webhooks, webhook)

	_, err := client.Collection(collectionName).Doc(webhook.WebhookId).Create(ctx, webhook)
	if err != nil {
		return err
	}
	return nil

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

// Deletes document with id, doccumentID from the collection with the name and collectionName in the database
func DeleteTheWebhook(collectionName, documentID string) (error, int) {
	if client == nil {
		return fmt.Errorf("firebase is not initialized"), http.StatusInternalServerError
	}

	// reference to the webhook document
	docReference := client.Doc(collectionName + "/" + documentID)

	// delete the webhook
	if _, err := docReference.Delete(ctx); err != nil {
		return err, http.StatusFailedDependency
	}
	return nil, 0
}

// Returns all webhooks for now
func GetAllWebhooks() []internal.Webhook {
	return webhooks
}
