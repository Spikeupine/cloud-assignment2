package main

import (
	"assignment-2/api"
	"assignment-2/internal"
	"log"
	"net/http"
	"os"
)

// This is the start point of the entire service.
func main() {

	//Gets the port from the environment. If empty, sets it to 8080 as default.
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT not set. Default: 8080")
		port = "8080"
	}
	// Starts the server
	log.Println("Starting server on port " + port + " ...")
	http.HandleFunc(internal.DashboardsPath, api.HandleRestcountriesapi)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
