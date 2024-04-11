package handlers

import (
	"assignment-2/database"
	"fmt"
	"net/http"
)

// RegistrationsHandler handles requests to the registration endpoint,
func RegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		pathValue := r.PathValue("id")
		if pathValue == "" {
			fmt.Fprintf(w, "All paths")
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
		// Extract ID from URL path
		id := r.PathValue("id")

		// Initialize Firestore client (replace with your Firestore initialization code)
		firestoreClient := database.GetClient()

		// Delete the registration document from Firestore based on the provided ID.
		_, err := firestoreClient.Collection("dashboards").Doc(id).Delete(r.Context())
		if err != nil {
			http.Error(w, "Error deleting registration: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with a success message.
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Registration deleted successfully")

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
