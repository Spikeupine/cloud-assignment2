package handlers

import (
	"assignment-2/database"
	"assignment-2/internal"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		// Parse request body
		var reg internal.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Extract ID from URL path
		id := r.PathValue("id")

		// Initialize Firestore client
		firestoreClient := database.GetClient()

		// Fetch the existing registration
		existingReg, err := fetchSingleByField(r.Context(), firestoreClient, "dashboards", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Merge existing registration with new data
		mergeRegistration(existingReg, &reg)

		// Save the updated registration
		if _, err := firestoreClient.Collection("dashboards").Doc(id).Set(r.Context(), existingReg); err != nil {
			http.Error(w, "Failed to update registration", http.StatusInternalServerError)
			return
		}

		// Respond with success message
		fmt.Fprintf(w, "Registration with ID %s updated successfully", id)

	case http.MethodDelete:
		// Extract ID from URL path
		id := r.PathValue("id")

		// Initialize Firestore client
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

func mergeRegistration(existingReg *internal.Features, newReg *internal.RegisterRequest) {
	// Update each field in the existing registration with the corresponding field from the new registration
	if newReg.Features.Temperature {
		existingReg.Temperature = true
	}
	if newReg.Features.Precipitation {
		existingReg.Precipitation = true
	}
	if newReg.Features.Capital {
		existingReg.Capital = true
	}
	if newReg.Features.Population {
		existingReg.Population = true
	}
	if newReg.Features.Area {
		existingReg.Area = true
	}
	if len(newReg.Features.TargetCurrencies) > 0 {
		existingReg.TargetCurrencies = newReg.Features.TargetCurrencies
	}

}

func fetchSingleByField(ctx context.Context, client *firestore.Client, collectionName string, documentID string) (*internal.Features, error) {
	// Create a reference to the document
	docRef := client.Collection(collectionName).Doc(documentID)

	// Get the document snapshot
	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, errors.New("Document not found")
		}
		return nil, err
	}

	// Extract the data from the document snapshot
	var dashboard internal.Features
	if err := docSnapshot.DataTo(&dashboard); err != nil {
		return nil, err
	}

	return &dashboard, nil
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
