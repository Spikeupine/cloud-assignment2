package resources

import (
	"assignment-2/external"
	"encoding/json"
	"errors"
)

func GetRestCountries(countryName string, isoCode string) (country external.CountriesObject, err error) {
	link, err := GetRestCountryLink(countryName, isoCode)
	if err != nil {
		return country, err
	}
	response, err := external.GetExternalResource(link)
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
	link := external.CountriesApi
	if isoCode != "" {
		link += external.IsoExtention + isoCode
	} else if countryName != "" {
		link += external.CountryNameExtention + countryName
	} else {
		return "", errors.New("country name and iso-code is empty")
	}
	link += "?" + external.CountriesFields
	return link, nil
}
