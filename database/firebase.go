package database

import (
	"assignment-2/internal"
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
)

var (
	ctx    context.Context
	client *firestore.Client
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
	key, exists := os.LookupEnv("FIREBASE_KEY")
	if !exists {
		log.Fatal("FIREBASE_KEY environment variable not set")
	}
	opt := option.WithCredentialsFile(key)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Println("Couldn't establish connection to firebase: " + err.Error())
		os.Exit(1)
	}
	// Instantiate client
	client, err = app.Firestore(ctx)
	if err != nil {
		log.Println("Couldn't establish connection to the database: " + err.Error())
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

// Appends the webhook to the collection of other webhooks in firebase.
func AddWebhookToCollection(webhook internal.Webhook, collectionName string) error {

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

// GetWebhook returns the webhook requested by its ID, or an error if any.
func GetWebhook(collectionName string, webhookID string) (internal.Webhook, error) {
	docReference := client.Doc(collectionName + "/" + webhookID)
	documentSnapshot, err := docReference.Get(ctx)
	if err != nil {
		return internal.Webhook{}, err
	}

	var hook internal.Webhook
	if err := documentSnapshot.DataTo(&hook); err != nil {
		return internal.Webhook{}, err
	}

	return hook, nil
}

// DeleteTheWebhook finds the specified document's id in the collection specified, and deletes it.
func DeleteTheWebhook(collectionName, documentID string) (error, int) {

	// reference to the webhook document
	docReference := client.Doc(collectionName + "/" + documentID)

	// delete the webhook
	if _, err := docReference.Delete(ctx); err != nil {
		return err, http.StatusFailedDependency
	}
	return nil, 0
}

// GetAllWebhooks displays all the registered webhooks in the collection from firebase.
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
		err = doc.DataTo(&data)
		if err != nil {
			http.Error(w, "Error when sending data to document in getAllWebhooks :"+err.Error(), http.StatusInternalServerError)
		}

		documents = append(documents, data)
	}
	return documents, nil
}
