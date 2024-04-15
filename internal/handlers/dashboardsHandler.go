package handlers

import (
	"assignment-2/database"
	"assignment-2/internal"
	"encoding/json"
	"net/http"
	"time"
)

// DashboardsHandler handles requests related to dashboards
func DashboardsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := r.PathValue("id")
		if len(id) < 1 {
			http.Error(w, "invalid ID", http.StatusBadRequest)
			return
		}

		dashboard, err := getPopulatedDashboard(id)
		if err != nil {
			http.Error(w, "error retrieving populated dashboard: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(dashboard)
		if err != nil {
			http.Error(w, "error encoding response: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func getRegisterRequestFromDatabase(id string) (internal.RegisterRequest, error) {
	var dashboard internal.RegisterRequest
	docRef := database.GetClient().Collection("dashboards").Doc(id)
	docSnapshot, err := docRef.Get(database.GetContext())
	if err != nil {
		return internal.RegisterRequest{}, err
	}
	err = docSnapshot.DataTo(&dashboard)
	if err != nil {
		return internal.RegisterRequest{}, err
	}
	return dashboard, nil
}

// Function to fetch data for each enabled feature
func getPopulatedDashboard(id string) (internal.PopulatedDashboard, error) {
	var request internal.RegisterRequest
	var populated internal.PopulatedDashboard

	// Retrieve RegisterRequest from database
	request, err := getRegisterRequestFromDatabase(id)
	if err != nil {
		return internal.PopulatedDashboard{}, err
	}

	// Start populating PopulatedDashboard
	populated.ID = request.Id
	populated.Country = request.Country
	populated.IsoCode = request.IsoCode
	populated.LastRetrieval = time.Now()

	// Populate features based on what's enabled
	if request.Features.Temperature {
		populated.Features.Temperature, err = fetchTemperature(request.Country)
		if err != nil {
			return internal.PopulatedDashboard{}, err
		}
	}
	if request.Features.Precipitation {
		populated.Features.Precipitation, err = fetchPrecipitation(request.Country)
		if err != nil {
			return internal.PopulatedDashboard{}, err
		}
	}
	if request.Features.Capital {
		populated.Features.Capital, err = fetchCapital(request.Country)
		if err != nil {
			return internal.PopulatedDashboard{}, err
		}
	}
	if request.Features.Coordinates {
		populated.Features.Coordinates.Latitude, populated.Features.Coordinates.Longitude, err = fetchCoordinates(request.Country)
		if err != nil {
			return internal.PopulatedDashboard{}, err
		}
	}
	if request.Features.Population {
		populated.Features.Population, err = fetchPopulation(request.Country)
		if err != nil {
			return internal.PopulatedDashboard{}, err
		}
	}
	if request.Features.Area {
		populated.Features.Area, err = fetchArea(request.Country)
		if err != nil {
			return internal.PopulatedDashboard{}, err
		}
	}
	if len(request.Features.TargetCurrencies) > 0 {
		populated.Features.TargetCurrencies, err = fetchTargetCurrencies(request.Features.TargetCurrencies)
		if err != nil {
			return internal.PopulatedDashboard{}, err
		}
	}

	return populated, nil
}

// Mock fetch functions (replace these with actual data fetching logic)
func fetchTemperature(country string) (float64, error) {
	// Fetch temperature data
	return 23.5, nil // Placeholder
}

func fetchPrecipitation(country string) (float64, error) {
	// Fetch precipitation data
	return 12.4, nil // Placeholder
}

func fetchCapital(country string) (string, error) {
	// Fetch capital city
	return "Capital City", nil // Placeholder
}

func fetchCoordinates(country string) (float32, float32, error) {
	// Fetch coordinates
	return 34.05, -118.25, nil // Placeholder
}

func fetchPopulation(country string) (int, error) {
	// Fetch population data
	return 5000000, nil // Placeholder
}

func fetchArea(country string) (float64, error) {
	// Fetch area data
	return 123456.78, nil // Placeholder
}

func fetchTargetCurrencies(currencies []string) (map[string]float64, error) {
	// Fetch exchange rates for given currencies
	exchangeRates := make(map[string]float64)
	for _, currency := range currencies {
		exchangeRates[currency] = 1.23 // Placeholder
	}
	return exchangeRates, nil
}
