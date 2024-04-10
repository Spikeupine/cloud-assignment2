package handlers

import (
	"assignment-2/internal"
	"encoding/json"
	"net/http"
)

// DashboardsHandler handles requests related to dashboards
func DashboardsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if len(id) < 1 {
			http.Error(w, "invalid ID", http.StatusBadRequest)
			return
		}

		dashboard, err := getDashboard(id)
		if err != nil {
			http.Error(w, "error retrieving dashboard: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(dashboard)
		if err != nil {
			http.Error(w, "error encoding response: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// getDashboard retrieves dashboard information.
func getDashboard(id string) (internal.PopulatedDashboard, error) {
	var dashboard internal.PopulatedDashboard

	// Ignore binder and client for now.

	return dashboard, nil
}
