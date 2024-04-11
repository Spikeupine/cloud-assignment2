package handlers

import (
	"assignment-2/database"
	"assignment-2/internal"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"net/http"
)

// RegistrationsHandler handles requests to the registration endpoint,
func RegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		pathValue := r.PathValue("id")
		if pathValue == "" {
			registrations, err := getAllRegistrations()
			if err != nil {
				http.Error(w, "error retrieving data"+err.Error(), http.StatusInternalServerError)
			}
			json.NewEncoder(w).Encode(registrations)
		} else {
			fmt.Fprintf(w, "Path for %s", pathValue)
		}
	case http.MethodPost:
		err := registerDashboard(w, r)
		if err != nil {
			http.Error(w, "Error posting to Registrations Handler: "+err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPut:

	case http.MethodDelete:

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func getAllRegistrations() ([]internal.RegisterRequest, error) {
	var registrations []internal.RegisterRequest
	client := database.GetClient().Collection("dashboards")
	documents := client.Documents(database.GetContext())
	for {
		doc, err := documents.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return []internal.RegisterRequest{}, err
		}
		var registration internal.RegisterRequest
		err = doc.DataTo(&registration)
		if err != nil {
			return []internal.RegisterRequest{}, err
		}
		if registration.Id != "" {
			registrations = append(registrations, registration)
		}
	}
	return registrations, nil
}
