package handlers

import (
	"assignment-2/internal"
	"encoding/json"
	"net/http"
	"time"
)

var registrationId int = 1

// Handles calls to registrations path, and directs different http requests to methods made for them.
func registrationsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		postToRegistrations(w, r)
		break
	case "GET":

		break

	default:

		break
	}

}

// Holds the id values of the different dashboards created by user, as well as time and date of last edit.
type registrationsResponse struct {
	id         int
	lastChange time.Time
}

// postToRegistrations posts how the user wants the dashboard to look like. The user decides which of the
// attributes they would like to have displayed on the dashboards, and it is connected to an id, so it can
// be stored and found again. This function returns a struct with the time and date of last time the dashboard
// was edited, as well as the id of the dashboard.
func postToRegistrations(w http.ResponseWriter, r *http.Request) registrationsResponse {

	//Sets the id of the dashboard, and registers the time
	registrationResponse := registrationsResponse{
		id:         registrationId,
		lastChange: time.Now(),
	}

	type Feature struct {
		Precipitations bool `json:"precipitation"`
		Population     bool `json:"population"`
		Capital        bool `json:"capital"`
		Currencies     bool `json:"currencies"`
		Coordinates    bool `json:"coordinates"`
		Area           bool `json:"area"`
	}
	type responseToDashboard struct {
		Name     string `json:"country"`
		Iso      string `json:"isoCode"`
		Features []Feature
	}

	var responsebodycontent responseToDashboard
	w.Header().Add(internal.ApplicationJson, internal.ContentTypeJson)

	//Decodes and inserts the response into the struct.
	err := json.NewDecoder(r.Body).Decode(&responsebodycontent)
	if err != nil {
		http.Error(w, "Error decoding dashboard data: "+err.Error(), http.StatusInternalServerError)
	}

	//Todo: Write code that sets the parameters the user would like into the dashboard.

	//Todo: Write dashboard to database for storage

	registrationId = registrationId + 1
	return registrationResponse
}

// Todo: write GET method.
func getRegistrations() {

}
