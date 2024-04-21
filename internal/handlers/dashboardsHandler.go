package handlers

import (
	"assignment-2/database"
	"assignment-2/external"
	"assignment-2/external/router"
	"assignment-2/internal"
	"context"
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
		//Registers calls to webhook if registered.
		EventWebhook(w, dashboard.IsoCode, "INVOKE")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(dashboard)
		if err != nil {
			http.Error(w, "error encoding response: "+err.Error(), http.StatusInternalServerError)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func getRegisterRequestFromDatabase(id string, ctx context.Context) (internal.RegisterRequest, error) {
	var dashboard internal.RegisterRequest
	docRef := database.GetDocumentRef(database.DashboardCollection, id)
	docSnapshot, err := docRef.Get(ctx)
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
func getPopulatedDashboard(id string, ctx context.Context) (populated internal.PopulatedDashboard, err error) {
	var register internal.RegisterRequest
	// Retrieve RegisterRequest from database
	register, err = getRegisterRequestFromDatabase(id, ctx)
	if err != nil {
		return internal.PopulatedDashboard{}, err
	}
	err = populateDashboardFeatures(&populated, register)

	return populated, err
}

func populateDashboardFeatures(dashboard *internal.PopulatedDashboard, registry internal.RegisterRequest) (err error) {
	// Start populating PopulatedDashboard
	dashboard.ID = registry.Id
	dashboard.Country = registry.Country
	dashboard.IsoCode = registry.IsoCode
	dashboard.LastRetrieval = time.Now()
	dashboard.Features.TargetCurrencies = make(map[string]float64)
	countryInfo, err := router.GetCountriesObject(registry.Country, registry.IsoCode)
	if err != nil {
		return err
	}
	if dashboard.IsoCode == "" {
		dashboard.IsoCode = countryInfo.Cca2
	}
	if dashboard.Country == "" {
		dashboard.Country = countryInfo.Name.Common
	}
	currency := getFirstCurrency(countryInfo)
	currencyInfo, err := router.GetCurrencyObject(currency)
	if err != nil {
		return err
	}
	coordinates := extractCoordinates(countryInfo)
	meteoData, err := router.GetMeteoObject(coordinates.Latitude, coordinates.Longitude)
	if err != nil {
		return err
	}
	meteoWeather := meteoData.Weather

	// Populate features based on what's enabled
	if registry.Features.Temperature && len(meteoWeather.Temperature) > 0 {
		dashboard.Features.Temperature = meteoWeather.Temperature[0]
	}
	if registry.Features.Precipitation && len(meteoWeather.Precipitation) > 0 {
		dashboard.Features.Precipitation = meteoWeather.Precipitation[0]
	}
	if registry.Features.Capital && len(countryInfo.Capital) > 0 {
		dashboard.Features.Capital = countryInfo.Capital[0]
	}
	if registry.Features.Coordinates {
	}
	if registry.Features.Population {
		dashboard.Features.Population = countryInfo.Population
	}
	if registry.Features.Area {
		dashboard.Features.Area = countryInfo.Area
	}
	if len(registry.Features.TargetCurrencies) > 0 {
		for _, targetCurrency := range registry.Features.TargetCurrencies {
			rates := currencyInfo.Rates[targetCurrency]
			dashboard.Features.TargetCurrencies[targetCurrency] = rates
		}
	}
	return nil
}

func getFirstCurrency(countryObject external.CountriesObject) string {
	for currency, _ := range countryObject.Currencies {
		return currency
	}
	return ""
}

func extractCoordinates(countryInfo external.CountriesObject) internal.DashboardCoordinates {
	var coordinates internal.DashboardCoordinates
	coordinates.Longitude = countryInfo.Location.Latlng[0]
	coordinates.Latitude = countryInfo.Location.Latlng[1]
	return coordinates
}
