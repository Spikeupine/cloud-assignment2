package handlers

import (
	"assignment-2/internal"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// RegistrationsPostHandler Handles calls to registrations path
// and directs different http requests to methods made for them.
func RegistrationsPostHandler(client *firestore.Client, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := registerDashboard(client, w, r)
		if err != nil {
			http.Error(w, "Error posting to Registrations Handler: "+err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Error: This method is not supported. Please use POST", http.StatusMethodNotAllowed)
	}
}

// registerDashboard parses the json-body of the request and creates a registers a new dashboard from it
func registerDashboard(client *firestore.Client, w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	var dashboard internal.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&dashboard)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
	}

	if dashboard.IsoCode == "" && dashboard.Country == "" {
		http.Error(w, "ISO country code or country must be present: ", http.StatusBadRequest)
		return nil
	}

	// Parse the user's feature selections and create/update the dashboard in Firestore
	response, err := uploadDashboard(ctx, client, dashboard)
	if err != nil {
		http.Error(w, "Failed to create dashboard: "+err.Error(), http.StatusInternalServerError)
		return err
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Couldn't parse response from database", http.StatusInternalServerError)
	}

	return nil
}

// uploadDashboard uploads a dashboard to firestore
func uploadDashboard(ctx context.Context, client *firestore.Client,
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
