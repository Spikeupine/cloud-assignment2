package handlers

import (
	"assignment-2/database"
	"assignment-2/internal"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// registerDashboard parses the json-body of the request and creates a registers a new dashboard from it
func registerDashboard(w http.ResponseWriter, r *http.Request) error {
	var dashboard internal.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&dashboard)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
	}

	if dashboard.IsoCode == "" && dashboard.Country == "" {
		http.Error(w, "ISO country code or country must be present: ", http.StatusBadRequest)
		return nil
	}
	id, _ := uuid.NewUUID()
	dashboard.Id = id.String()
	dashboard.LastChange = time.Now()
	// Parse the user's feature selections and create/update the dashboard in Firestore
	response, err := uploadDashboard(dashboard)
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
func uploadDashboard(dashboard internal.RegisterRequest) (internal.RegistrationsResponse, error) {
	client := database.GetClient()
	_, err := client.Collection("dashboards").Doc(dashboard.Id).Create(database.GetContext(), dashboard)
	if err != nil {
		return internal.RegistrationsResponse{}, err
	}
	registrationResponse := internal.RegistrationsResponse{
		Id:         dashboard.Id,
		LastChange: dashboard.LastChange,
	}
	return registrationResponse, nil
}
