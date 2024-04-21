package handlers

import (
	"assignment-2/database"
	"assignment-2/internal"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
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
				return
			}
			json.NewEncoder(w).Encode(registrations)
		} else {
			document, err := fetchSingleByField(r.Context(), database.GetClient(), "dashboards", pathValue)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = json.NewEncoder(w).Encode(document)
			if err != nil {
				http.Error(w, "error encoding document"+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	case http.MethodPost:
		err := registerDashboard(w, r)
		if err != nil {
			http.Error(w, "Error posting to Registrations Handler: "+err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPut:
		// Extract ID from URL path
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "id required", http.StatusBadRequest)
			return
		}

		// Initialize Firestore client
		firestoreClient := database.GetClient()

		// Fetch the existing registration
		_, err := fetchSingleByField(r.Context(), firestoreClient, "dashboards", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Parse request body
		var updatedRegistration internal.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&updatedRegistration); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		EventWebhook(w, updatedRegistration.IsoCode, "CHANGE")

		// Makes sure the ID isn't overwritten
		updatedRegistration.Id = id
		updatedRegistration.LastChange = time.Now()

		fmt.Fprintf(w, "Registration and corresponding dashboard with ID %s updated successfully", id)

	case http.MethodDelete:
		// Extract ID from URL path
		id := r.PathValue("id")

		// Initialize Firestore client
		firestoreClient := database.GetClient()

		// Delete the registration document from Firestore based on the provided ID.
		firestore := firestoreClient.Collection("dashboards").Doc(id)
		document, err := firestore.Get(r.Context())
		if err != nil {
			http.Error(w, "Error when retrieving specified dashboard by id :"+err.Error(), http.StatusNotFound)
			return
		}
		type isoCodeForWebhook struct {
			IsoCode string `json:"isoCode"`
			Method  string
		}

		var actualIsoCodeForWebhook isoCodeForWebhook
		actualIsoCodeForWebhook.Method = "DELETE"
		if err := document.DataTo(&actualIsoCodeForWebhook); err != nil {
			http.Error(w, "Error when parsing country information for webhook check :"+err.Error(), http.StatusServiceUnavailable)
		}
		EventWebhook(w, actualIsoCodeForWebhook.IsoCode, actualIsoCodeForWebhook.Method)
		_, err = firestore.Delete(r.Context())
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

func fetchSingleByField(ctx context.Context, client *firestore.Client, collectionName string, documentID string) (internal.RegisterRequest, error) {
	// Create a reference to the document
	docRef := client.Collection(collectionName).Doc(documentID)

	// Get the document snapshot
	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return internal.RegisterRequest{}, errors.New("document not found")
		}
		return internal.RegisterRequest{}, err
	}

	// Extract the data from the document snapshot
	var dashboard internal.RegisterRequest
	if err := docSnapshot.DataTo(&dashboard); err != nil {
		return internal.RegisterRequest{}, err
	}

	return dashboard, nil
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
