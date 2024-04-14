package database

import (
	"assignment-2/internal"
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/iterator"
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

// UpdateTheCallCount of the webhook document by id and collection name. Takes in count to pass on to database.
func UpdateTheCallCount(collectionName, docId string, callCount int) error {

	_, err := client.Collection(collectionName).Doc(docId).Update(ctx, []firestore.Update{
		{
			Path:  "Calls",
			Value: callCount,
		},
	})
	return err
}

// Gets the webhook requested by its ID
func GetWebhook(w http.ResponseWriter, r *http.Request, webhookID string) (internal.Webhook, error) {
	var hook internal.Webhook
	documentContent, err := client.Doc("webhooks/" + webhookID).Get(ctx)
	if err != nil {
		return hook, err
	}

	if err := documentContent.DataTo(&hook); err != nil {
		return hook, err
	}

	return hook, nil
}

// Deletes document with id, doccumentID from the collection with the name and collectionName in the database
func DeleteTheWebhook(collectionName, documentID string) (error, int) {

	// reference to the webhook document
	docReference := client.Doc(collectionName + "/" + documentID)

	// delete the webhook
	if _, err := docReference.Delete(ctx); err != nil {
		return err, http.StatusFailedDependency
	}
	return nil, 0
}

// GetAllWebhooks displays all the registered webhooks in the collection.
func GetAllWebhooks(w http.ResponseWriter, collectionName string) ([]internal.Webhook, error) {
	//Creates map to sture the different documents in
	var documents []internal.Webhook

	//Iterator going through all the documents
	iter := client.Collection(collectionName).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break //Finished going through the documents.
		}
		if err != nil {
			http.Error(w, "Error iterating through webhooks :"+err.Error(), http.StatusInternalServerError)
		}

		var data internal.Webhook
		if err := doc.DataTo(&data); err != nil {
			return nil, fmt.Errorf("error decoding document data: %v", err)
		}

		documents = append(documents, data)
	}
	return documents, nil
}
