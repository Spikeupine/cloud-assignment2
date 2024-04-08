package handlers

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"net/http"
	"time"
)

var registrationId int = 1

// Handles calls to registrations path, and directs different http requests to methods made for them.
func RegistrationsHandlerPost(client *firestore.Client, w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		postToRegistrations(client, w, r)
	} else {
		http.Error(w, "Error: This method is not supported. Please use POST", http.StatusMethodNotAllowed)
	}

}

// Holds the id values of the different dashboards created by user, as well as time and date of last edit.
type RegistrationsResponse struct {
	id         string
	lastChange time.Time
}

type Feature struct {
	Precipitations bool `json:"precipitation,omitempty"`
	Population     bool `json:"population,omitempty"`
	Capital        bool `json:"capital,omitempty"`
	Currencies     bool `json:"currencies,omitempty"`
	Coordinates    bool `json:"coordinates,omitempty"`
	Area           bool `json:"area,omitempty"`
}
type ResponseToDashboard struct {
	Name     string `json:"country"`
	Iso      string `json:"isoCode"`
	Features []Feature
}

// postToRegistrations posts how the user wants the dashboard to look like. The user decides which of the
// attributes they would like to have displayed on the dashboards, and it is connected to an id, so it can
// be stored and found again. This function returns a struct with the time and date of last time the dashboard
// was edited, as well as the id of the dashboard.
func postToRegistrations(client *firestore.Client, w http.ResponseWriter, r *http.Request) error {

	//Introduces the ctx variable.
	ctx := context.Background()

	//responsebodycontent is to hold the values user choses to be stored in the dashboard.
	//var responsebodycontent ResponseToDashboard
	/*
		// Sets header type.
		w.Header().Add(internal.ApplicationJson, internal.ContentTypeJson)

		//Decodes and inserts the response into the struct.
		err := json.NewDecoder(r.Body).Decode(&responsebodycontent)
		if err != nil {
			http.Error(w, "Error decoding dashboard data: "+err.Error(), http.StatusInternalServerError)
			return err
		}

	*/

	//Let's assume this is what we got from user:
	// Example usage
	dashboard := ResponseToDashboard{
		Name: "France",
		Iso:  "FR",
		Features: []Feature{
			{Precipitations: true, Population: true, Capital: true},
		},
	}

	// Parse the user's feature selections and create/update the dashboard in Firestore
	id, err := CreateDashboard(ctx, client, dashboard, w)
	if err != nil {
		http.Error(w, "Failed to create dashboard: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	// Return the ID of the dashboard to the user
	fmt.Fprintf(w, "Created new dashboard with ID: %s\n", id)

	return nil
}

// CreateDashboard creates a new dashboard document in Firestore.
func CreateDashboard(ctx context.Context, client *firestore.Client, dashboard ResponseToDashboard, w http.ResponseWriter) (RegistrationsResponse, error) {
	docRef, _, err := client.Collection("dashboards").Add(ctx, dashboard)
	if err != nil {
		return RegistrationsResponse{}, err
	}
	registrationResponse := RegistrationsResponse{
		id:         docRef.ID,
		lastChange: time.Now(),
	}
	return registrationResponse, nil
}
