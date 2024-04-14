package resources

import (
	"assignment-2/external"
	"assignment-2/internal"
	"encoding/json"
	"errors"
	"net/http"
)

func RestCountryTestEndpoint(w http.ResponseWriter, r *http.Request) {
	var (
		isoCode     string
		countryName string
	)
	w.Header().Set("Content-Type", "application/json")
	stringType := r.PathValue("type")
	switch stringType {
	case "iso":
		isoCode = r.PathValue("country")
	case "name":
		countryName = r.PathValue("country")
	default:
		http.Error(w, "", http.StatusNotFound)
	}
	country, err := GetRestCountries(countryName, isoCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(country)
}

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
	link := internal.CountriesApi
	if isoCode != "" {
		link += internal.IsoExtention + isoCode
	} else if countryName != "" {
		link += internal.CountryNameExtention + countryName
	} else {
		return "", errors.New("country name and iso-code is empty")
	}
	link += "?" + internal.CountriesFields
	return link, nil
}
