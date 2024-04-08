package internal

import "fmt"

// API's
// Link to the countries api
const CountriesApi = "http://129.241.150.113:8080/v3.1/"

// Link to the currency api
const CurrencyApi = "http://129.241.150.113:9090/currency/"

const CountriesFields = "fields=name,cca2,capitalInfo,population,area,currencies"

const MeteoApi = "https://api.open-meteo.com/v1/forecast?"

const MeteoField = "latitude=%f&longitude=%f&hourly=temperature_2m,precipitation&wind_speed_unit=ms&timeformat=unixtime&forecast_days=3"

// The different paths that can be used by user.
// Path to dashboards
const DashboardsPath = "/dashboard/v1/dashboards/"

// Path to notifications
const NotificationsPath = "/dashboard/v1/notifications/"

// Path to registrations
const RegistrationsPath = "/dashboard/v1/registrations/"

// Path to status
const StatusPath = "/dashboard/v1/status/"

// The type of the content. How to present it or read it.
const ContentTypeJson = "application/json"

// const of string value Content type.
const ApplicationJson = "Content type"

// Example of ISO code to be used to check if services that require ISO code is available.
const IsoExample = "alpha?codes=no"

/*
GetMeteoUrl
lat Latitude
long Longitude
returns the MeteoURL for the given coordinates
*/
func GetMeteoUrl(lat float64, long float64) string {
	arguments := fmt.Sprintf(MeteoField, lat, long)
	return MeteoApi + arguments

}
