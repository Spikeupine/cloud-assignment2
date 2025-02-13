package main

import (
	"assignment-2/database"
	"assignment-2/internal"
	"assignment-2/internal/handlers"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file" + err.Error())
	}
	database.FirebaseConnect()
}

// This is the start point of the entire service.
func main() {
	//Gets the port from the environment. If empty, sets it to 8080 as default.
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT not set. Default: 8080")
		port = "8080"
	}
	addr := ":" + port

	// Register the routes and their corresponding resources
	http.HandleFunc(internal.StatusPath, handlers.StatusHandler)
	http.HandleFunc(internal.RegistrationsPath, handlers.RegistrationsHandler)
	http.HandleFunc(internal.RegistrationsPath+"{id}", handlers.RegistrationsHandler)
	http.HandleFunc(internal.RegistrationsPath2, handlers.RegistrationsHandler)
	http.HandleFunc(internal.DashboardsPath, handlers.DashboardsHandler)
	http.HandleFunc(internal.DashboardsPath+"{id}", handlers.DashboardsHandler)
	http.HandleFunc(internal.NotificationsPath, handlers.NotificationsHandler)
	http.HandleFunc(internal.NotificationsPath+"{id}", handlers.NotificationsHandler)

	// Starts the server
	log.Println("Starting server on http://localhost:" + port + " ...")
	log.Printf("Firestore REST service listening on %s ...\n", addr)
	if errServ := http.ListenAndServe(addr, nil); errServ != nil {
		panic(errServ)
	}
	defer func() {
		database.FireBaseCloseConnection()
	}()
}
