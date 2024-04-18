package resources

import (
	"assignment-2/external"
	"encoding/json"
)

func GetCurrencyData(currencyCode string) (currencyObject external.CurrencyObject, err error) {
	url := getCurrencyUrl(currencyCode)
	response, err := external.GetExternalResource(url)
	if err != nil {
		return currencyObject, err
	}
	err = json.NewDecoder(response.Body).Decode(&currencyObject)
	return currencyObject, err
}

func getCurrencyUrl(currencyCode string) string {
	return external.CurrencyApi + currencyCode
}
