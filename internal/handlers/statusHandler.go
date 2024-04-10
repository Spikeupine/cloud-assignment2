package handlers

import (
	"assignment-2/internal"
	"encoding/json"
	"net/http"
	"time"
)

// Start the server timer
var serviceStartTime = time.Now()

// StatusHandler gives us the status for each API
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	var response internal.Status

	// Retrieve the status of the APIs
	response.CountriesAPI = getAPIStatus(internal.CountriesApi + internal.IsoExample)
	response.MeteoAPI = getAPIStatus(internal.MeteoApi)
	response.CurrencyAPI = getAPIStatus(internal.CurrencyApi + internal.IsoExample)
	response.Uptime = int64(time.Since(serviceStartTime).Seconds())

	w.Header().Set(internal.ApplicationJson, internal.ContentTypeJson)
	json.NewEncoder(w).Encode(response)
}

// Function for reusable API statuses
func getAPIStatus(apiURL string) int {
	// Send an HTTP GET request to the specified API URL.
	resp, err := http.Get(apiURL)
	if err != nil {
		return http.StatusInternalServerError
	}

	// Close the response body to prevent resource leaks.
	defer resp.Body.Close()

	return resp.StatusCode
}
