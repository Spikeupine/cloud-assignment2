package handlers

import (
	"assignment-2/database"
	"assignment-2/internal"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
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
			registrations, err := GetAllRegistrations(r.Context())
			if err != nil {
				http.Error(w, "error retrieving data"+err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(registrations)
		} else {
			document, err := GetSingleRegistration(r.Context(), database.DashboardCollection, pathValue)
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
		var dashboard internal.RegisterRequest

		err := json.NewDecoder(r.Body).Decode(&dashboard)
		if err != nil {
			http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if dashboard.IsoCode == "" && dashboard.Country == "" {
			http.Error(w, "ISO country code or country must be present: ", http.StatusBadRequest)
			return
		}
		SetRegistrationValues(&dashboard)

		// Parse the user's feature selections and create/update the dashboard in Firestore
		response, err := UploadDashboard(dashboard, r.Context())
		if err != nil {
			http.Error(w, "Failed to create dashboard: "+err.Error(), http.StatusInternalServerError)
		}
		EventWebhook(w, dashboard.IsoCode, "REGISTER")

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Couldn't parse response from database", http.StatusInternalServerError)
		}
	case http.MethodPut:
		// Extract ID from URL path
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "id required", http.StatusBadRequest)
			return
		}

		// Fetch the existing registration
		_, err := GetSingleRegistration(r.Context(), database.DashboardCollection, id)
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
		// Delete the registration document from Firestore based on the provided ID.
		documentRef := database.GetDocumentRef(database.DashboardCollection, id)
		document, err := documentRef.Get(r.Context())
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
		err = DeleteRegistration(r.Context(), database.DashboardCollection, id)
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

func SetRegistrationValues(dashboard *internal.RegisterRequest) {
	id, _ := uuid.NewUUID()
	dashboard.Id = id.String()
	dashboard.LastChange = time.Now()
}

func GetSingleRegistration(ctx context.Context, collectionName string, documentID string) (internal.RegisterRequest, error) {
	// Create a reference to the document
	docRef := database.GetDocumentRef(collectionName, documentID)

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

func GetAllRegistrations(ctx context.Context) ([]internal.RegisterRequest, error) {
	var registrations []internal.RegisterRequest
	collectionRef := database.GetCollectionRef(database.DashboardCollection)
	documents := collectionRef.Documents(ctx)
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
func UploadDashboard(dashboard internal.RegisterRequest, ctx context.Context) (internal.RegistrationsResponse, error) {
	_, err := database.GetDocumentRef(database.DashboardCollection, dashboard.Id).Create(ctx, dashboard)
	if err != nil {
		return internal.RegistrationsResponse{}, err
	}
	registrationResponse := internal.RegistrationsResponse{
		Id:         dashboard.Id,
		LastChange: dashboard.LastChange,
	}
	return registrationResponse, nil
}

func DeleteRegistration(ctx context.Context, collectionName string, id string) error {
	_, err := database.GetDocumentRef(collectionName, id).Delete(ctx)
	return err
}
