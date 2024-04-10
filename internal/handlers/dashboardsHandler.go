package handlers

import (
	"assignment-2/internal"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"net/http"
)

// DashboardsHandler handles requests related to dashboards
func DashboardsHandler(basics internal.Basics, client *firestore.Client) {
	switch basics.Request.Method {
	case http.MethodGet:
		// Retrieve the ID
		id := basics.ID
		// Check if the ID is empty
		if len(id) < 1 {
			http.Error(basics.ResponseWriter, "invalid ID", http.StatusBadRequest)
			return
		}

		// Retrieve the dashboard
		dashboard, err := getDashboard(id, client)
		if err != nil {
			http.Error(basics.ResponseWriter, "error retrieving dashboard: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Encode the dashboard data as JSON
		resp, err := json.Marshal(dashboard)
		if err != nil {
			http.Error(basics.ResponseWriter, "error encoding response: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the content type and write the response
		basics.ResponseWriter.Header().Set("Content-Type", "application/json")
		basics.ResponseWriter.WriteHeader(http.StatusOK)
		basics.ResponseWriter.Write(resp)
	default:
		http.Error(basics.ResponseWriter, "method not allowed", http.StatusMethodNotAllowed)
	}

}

// getDashboard retrieves dashboard information
func getDashboard(id string, client *firestore.Client) (internal.PopulatedDashboard, error) {
	// Create an empty dashboard.
	var dashboard internal.PopulatedDashboard

	// Ignore binder and client for now
	// Currently, not using the ID and the client

	return dashboard, nil
}
