package handlers

import (
	"assignment-2/internal"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"net/http"
)

func DashboardsHandler(basics internal.Basics, client *firestore.Client) {
	switch basics.Request.Method {
	case http.MethodGet:
		id := basics.ID
		if len(id) < 1 {
			http.Error(basics.ResponseWriter, "invalid ID", http.StatusBadRequest)
			return
		}

		dashboard, err := getDashboard(id, client)
		if err != nil {
			http.Error(basics.ResponseWriter, "error retrieving dashboard: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(dashboard)
		if err != nil {
			http.Error(basics.ResponseWriter, "error encoding response: "+err.Error(), http.StatusInternalServerError)
			return
		}
		basics.ResponseWriter.Header().Set("Content-Type", "application/json")
		basics.ResponseWriter.WriteHeader(http.StatusOK)
		basics.ResponseWriter.Write(resp)
	default:
		http.Error(basics.ResponseWriter, "method not allowed", http.StatusMethodNotAllowed)
	}

}

// getDashboard retrieves dashboard information.
func getDashboard(id string, client *firestore.Client) (internal.PopulatedDashboard, error) {
	var dashboard internal.PopulatedDashboard

	// Ignore binder and client for now.

	return dashboard, nil
}
