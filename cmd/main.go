package main

import (
	"assignment-2/internal"
	"assignment-2/internal/handlers"
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
)

// Firebase context and client used by Firestore functions throughout the program.
var ctx context.Context
var client *firestore.Client

// Collection name in Firestore
const collection = "something"

// Message counter to produce some variation in content
var ct = 0

// This is the start point of the entire service.
func main() {
	firebaseConnect()
	//Gets the port from the environment. If empty, sets it to 8080 as default.
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT not set. Default: 8080")
		port = "8080"
	}
	addr := ":" + port

	// Starts the server
	log.Println("Starting server on port " + port + " ...")
	http.HandleFunc(internal.DashboardsPath, handlers.HandleRestcountriesapi)
	http.HandleFunc(internal.RegistrationsPath, func(w http.ResponseWriter, r *http.Request) {
		handlers.RegistrationsHandlerPost(client, w, r)
	})
	log.Printf("Firestore REST service listening on %s ...\n", addr)
	if errServ := http.ListenAndServe(addr, nil); errServ != nil {
		panic(errServ)
	}

}

func firebaseConnect() {
	// Firebase initialisation
	ctx = context.Background()

	pathToCredentials := "./firebase_privatekey.json"
	opt := option.WithCredentialsFile(pathToCredentials)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Println(err)
		return
	}

	// Instantiate client
	client, err = app.Firestore(ctx)

	// Alternative setup, directly through Firestore (without initial reference to Firebase); but requires Project ID; useful if multiple projects are used
	// client, err := firestore.NewClient(ctx, projectID)

	// Check whether there is an error when connecting to Firestore
	if err != nil {
		log.Println(err)
		return
	}

	// Close down client at the end of the function
	defer func() {
		errClose := client.Close()
		if errClose != nil {
			log.Fatal("Closing of the Firebase client failed. Error:", errClose)
		}
	}()
}
