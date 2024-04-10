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

	// Firebase initialisation
	ctx = context.Background()

	opt := option.WithCredentialsFile("./firebase_privatekey.json")
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

	//Gets the port from the environment. If empty, sets it to 8080 as default.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("$PORT not set. Default: " + port)
	}
	addr := ":" + port

	// Register the routes and corresponding handlers
	http.HandleFunc(internal.StatusPath, handlers.StatusHandler)
	http.HandleFunc(internal.DashboardsPath, handlers.DashboardsHandler)

	// Starts the server
	log.Println("Server starting on http://localhost:%s" + port + " ...")
	log.Printf("Firestore REST service listening on %s ...\n", addr)
	if errServ := http.ListenAndServe(addr, nil); errServ != nil {
		panic(errServ)
	}

}
