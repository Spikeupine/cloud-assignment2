package handlers

import (
	"assignment-2/internal"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// RegistrationsHandlerPost Handles calls to registrations path
// and directs different http requests to methods made for them.
func RegistrationsHandlerPost(client *firestore.Client, w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := postToRegistrations(client, w, r)
		if err != nil {
			http.Error(w, "Error posting to Registrations Handler: "+err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Error: This method is not supported. Please use POST", http.StatusMethodNotAllowed)
	}
}

func postToRegistrations(client *firestore.Client, w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	var dashboard internal.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&dashboard)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
	}

	// Parse the user's feature selections and create/update the dashboard in Firestore
	response, err := CreateDashboard(ctx, client, dashboard)
	if err != nil {
		http.Error(w, "Failed to create dashboard: "+err.Error(), http.StatusInternalServerError)
		return err
	}
	fmt.Fprintf(w, "Created new dashboard with ID: %s\n", response.Id)

	return nil
}

// CreateDashboard creates a new dashboard document in Firestore.
func CreateDashboard(ctx context.Context, client *firestore.Client,
	dashboard internal.RegisterRequest) (internal.RegistrationsResponse, error) {
	id, _ := uuid.NewUUID()
	_, err := client.Collection("dashboards").Doc(id.String()).Create(ctx, dashboard)
	if err != nil {
		return internal.RegistrationsResponse{}, err
	}
	registrationResponse := internal.RegistrationsResponse{
		Id:         id,
		LastChange: time.Now(),
	}
	return registrationResponse, nil
}
