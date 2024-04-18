package router

import (
	"assignment-2/external"
	"assignment-2/external/resources"
	"assignment-2/external/stubs"
	"os"
	"strings"
)

const TESTING = "TESTING"
const TRUE = "true"

func GetMeteoObject(lat float64, lng float64) (Meteo external.MeteoObject, err error) {
	if strings.ToLower(os.Getenv(TESTING)) == TRUE {
		return stubs.MeteoStub()
	} else {
		return resources.GetMeteoData(lat, lng)
	}
}

func GetCountriesObject(countryName string, isoCode string) (country external.CountriesObject, err error) {
	if strings.ToLower(os.Getenv(TESTING)) == TRUE {
		return stubs.CountryStub()
	} else {
		return resources.GetRestCountries(countryName, isoCode)
	}
}

func GetCurrencyObject(currencyCode string) (currency external.CurrencyObject, err error) {
	if strings.ToLower(os.Getenv(TESTING)) == TRUE {
		return stubs.CurrencyStub()
	} else {
		return resources.GetCurrencyData(currencyCode)
	}
}
