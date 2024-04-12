package handlers

import (
	"assignment-2/internal"
	"encoding/json"
	"net/http"
	"time"
)

// DashboardsHandler handles requests related to dashboards
func DashboardsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := r.PathValue("id")
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

	id = "1"

	// Assign mock values to the dashboard
	dashboard.Country = "Mock Country"
	dashboard.IsoCode = "MC"
	dashboard.Features.Temperature = 25.5
	dashboard.Features.Precipitation = 10.2
	dashboard.Features.Capital = "Mock Capital"
	dashboard.Features.Coordinates.Latitude = 40.7128
	dashboard.Features.Coordinates.Longitude = -74.0060
	dashboard.Features.Population = 1000000
	dashboard.Features.Area = 1234.56
	dashboard.Features.TargetCurrencies = map[string]float64{
		"USD": 1.23,
		"EUR": 0.98,
	}
	dashboard.LastRetrieval = time.Now().Format(time.RFC3339)

	// Ignore binder and client for now.

	return dashboard, nil
}
