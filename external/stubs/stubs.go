package stubs

import (
	"assignment-2/external"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func readStubFile(filePath string) ([]byte, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("File error: %v\n", err)
		return nil, err
	}
	return file, nil
}

func MeteoStub() (external.MeteoObject, error) {
	var meteoData external.MeteoObject
	path := filepath.Clean(external.StubDataPath + "meteoWeather.json")
	meteoFile, err := readStubFile(path)
	if err != nil {
		return meteoData, err
	}
	err = json.Unmarshal(meteoFile, &meteoData)
	if err != nil {
		log.Printf("File parse error: %v\n", err)
		return external.MeteoObject{}, err
	}
	return meteoData, nil
}

func CurrencyStub() (external.CurrencyObject, error) {
	var currency external.CurrencyObject
	path := filepath.Clean(external.StubDataPath + "currencies.json")
	currencyFile, err := readStubFile(path)
	if err != nil {
		return external.CurrencyObject{}, err
	}
	err = json.Unmarshal(currencyFile, &currency)
	if err != nil {
		log.Printf("File parse error: %v\n", err)
		return external.CurrencyObject{}, err
	}
	return currency, nil
}

func CountryStub() (external.CountriesObject, error) {
	var country external.CountriesObject
	path := filepath.Clean(external.StubDataPath + "countriesAPI.json")
	countryFile, err := readStubFile(path)
	if err != nil {
		return external.CountriesObject{}, err
	}
	err = json.Unmarshal(countryFile, &country)
	if err != nil {
		log.Printf("File parse error: %v\n", err)
		return external.CountriesObject{}, err
	}
	return country, nil
}
