package handlers

import (
	"assignment-2/internal"
	"encoding/json"
	"net/http"
	"time"
)

// Start the server timer
var serviceStartTime = time.Now()

// StatusHandler gives us the status for each API, making it easier to see where something goes wrong
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize status codes
	countriesAPIstatus := getAPIStatus(internal.CountriesApi)
	currencyAPIstatus := getAPIStatus(internal.CurrencyApi)
	meteoAPIstatus := getAPIStatus(internal.GetMeteoUrl(0, 0))

	// Calculate uptime
	uptime := int(time.Since(serviceStartTime).Seconds())

	// Prepare and send response
	statusResponse := map[string]interface{}{
		"countries_api":   countriesAPIstatus,
		"meteo_api":       meteoAPIstatus,
		"currency_api":    currencyAPIstatus,
		"notification_db": "<http status code for *Notification database*>",
		"webhooks":        "<number of registered webhooks>",
		"version":         "v1",
		"uptime":          uptime,
	}

	w.Header().Set(internal.ApplicationJson, internal.ContentTypeJson)
	json.NewEncoder(w).Encode(statusResponse)
}

// Just a function to reduce code duplication, making it more reusable for the API statuses
func getAPIStatus(apiURL string) int {
	resp, err := http.Get(apiURL)
	if err != nil {
		return http.StatusInternalServerError
	}
	defer resp.Body.Close()

	return resp.StatusCode
}
