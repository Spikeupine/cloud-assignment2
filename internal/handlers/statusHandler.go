package handlers

import (
	"assignment-2/database"
	"assignment-2/external"
	"assignment-2/internal"
	"encoding/json"
	"net/http"
	"time"
)

// Start the server timer
var serviceStartTime = time.Now()

// StatusHandler gives us the status for each API
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8") // This sets the header for unsupported methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", internal.ContentTypeJson)

	var response internal.Status

	// Retrieve the status of the APIs
	response.CountriesAPI = getAPIStatus(external.CountriesApi + internal.IsoExample)
	response.MeteoAPI = getAPIStatus(external.MeteoApi)
	response.CurrencyAPI = getAPIStatus(external.CurrencyApi + "nok")
	webhooksCount, err := database.CountWebhooks("webhooks")
	if err != nil {
		http.Error(w, "Failed to get webhook count", http.StatusInternalServerError)
		return
	}
	firestoreStatusCode := http.StatusOK
	webhooksCount, err = database.GetNumberOfDocuments("webhooks")
	if err != nil {
		firestoreStatusCode = http.StatusNotFound
	}
	response.NotificationDB = firestoreStatusCode
	response.Webhooks = webhooksCount
	response.Version = "v1"
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
