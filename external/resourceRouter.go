package external

import (
	"assignment-2/external/resources"
	"assignment-2/external/stubs"
	"os"
	"strings"
)

const TESTING = "TESTING"
const TRUE = "true"

func GetMeteoObject(lat float64, lng float64) (Meteo MeteoObject, err error) {
	if strings.ToLower(os.Getenv(TESTING)) == TRUE {
		return stubs.MeteoStub(), nil
	} else {
		return resources.GetMeteoData(lat, lng)
	}
}

func GetCountriesObject(countryName string, isoCode string) (country CountriesObject, err error) {
	if strings.ToLower(os.Getenv(TESTING)) == TRUE {
		return stubs.CountryStub(), nil
	} else {
		return resources.GetRestCountries(countryName, isoCode)
	}
}

func GetCurrencyObject(currencyCode string) (currency CurrencyObject, err error) {
	if strings.ToLower(os.Getenv(TESTING)) == TRUE {
		return stubs.CurrencyStub(), nil
	} else {
		return resources.GetCurrencyData(currencyCode)
	}
}
