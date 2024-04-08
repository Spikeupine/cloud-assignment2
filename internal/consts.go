package internal

// API's
// Link to the countries api
const CountriesApi = "http://129.241.150.113:8080/v3.1/"

// Link to the currency api
const CurrencyApi = "http://129.241.150.113:9090/currency/"

const CountriesFields = "fields=name,cca2,capitalInfo,population,area,currencies"

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
const ApplicationJson = "Content type: "

// Example of ISO code to be used to check if services that require ISO code is available.
const IsoExample = "alpha?codes=no"
