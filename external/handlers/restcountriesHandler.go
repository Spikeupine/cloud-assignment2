package handlers

import (
	"assignment-2/external"
	"assignment-2/internal"
	"encoding/json"
	"errors"
	"net/http"
)

func GetRestCountries(countryName string, isoCode string) (country external.CountriesObject, err error) {
	link, err := GetRestCountryLink(countryName, isoCode)
	if err != nil {
		return country, err
	}
	request, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return country, err
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return country, err
	}
	err = json.NewDecoder(response.Body).Decode(&country)
	if err != nil {
		return country, err
	}
	return country, err
}

func GetRestCountryLink(countryName string, isoCode string) (string, error) {
	link := internal.CountriesApi
	if isoCode == "" {
		link += internal.IsoExtention + isoCode
	} else if countryName == "" {
		link += internal.CountryNameExtention + countryName
	} else {
		return "", errors.New("country name and iso-code is empty")
	}
	link += internal.CountriesFields
	return link, nil
}
