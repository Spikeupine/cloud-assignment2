package external

// API's
// Link to the countries api
const CountriesApi = "http://129.241.150.113:8080/v3.1/"
const IsoExtention = "alpha/"
const CountryNameExtention = "name/"

// Link to the currency api
const CurrencyApi = "http://129.241.150.113:9090/currency/"

const CountriesFields = "fields=name,cca2,capitalInfo,population,area,currencies"

const MeteoApi = "https://api.open-meteo.com/v1/forecast?"

const MeteoField = "latitude=%f&longitude=%f&hourly=temperature_2m,precipitation&wind_speed_unit=ms&timeformat=unixtime&forecast_days=1"
