package handlers

import (
	"assignment-2/database"
	"assignment-2/internal"
	"encoding/json"
	"net/http"
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

		dashboard, err := getPopulatedDashboard(id)
		if err != nil {
			http.Error(w, "error retrieving populated dashboard: "+err.Error(), http.StatusInternalServerError)
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

func getPopulatedDashboard(id string) (internal.PopulatedDashboard, error) {
	var dashboard internal.PopulatedDashboard
	docRef := database.GetClient().Collection("populatedDashboards").Doc(id)
	docSnapshot, err := docRef.Get(database.GetContext())
	if err != nil {
		return internal.PopulatedDashboard{}, err
	}
	err = docSnapshot.DataTo(&dashboard)
	if err != nil {
		return internal.PopulatedDashboard{}, err
	}
	return dashboard, nil
}
