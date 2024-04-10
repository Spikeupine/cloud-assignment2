package main

import (
	"assignment-2/database"
	"assignment-2/internal"
	"assignment-2/internal/handlers"
	"log"
	"net/http"
	"os"
)

// Firebase context and client used by Firestore functions throughout the program.

// Collection name in Firestore
const collection = "something"

// Message counter to produce some variation in content
var ct = 0

// This is the start point of the entire service.
func main() {
	database.FirebaseConnect()
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
		handlers.RegistrationsPostHandler(w, r)
	})
	log.Printf("Firestore REST service listening on %s ...\n", addr)
	if errServ := http.ListenAndServe(addr, nil); errServ != nil {
		panic(errServ)
	}
	defer func() {
		database.FireBaseCloseConnection()
	}()
}
