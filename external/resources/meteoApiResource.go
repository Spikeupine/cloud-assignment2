package resources

import (
	"assignment-2/external"
	"assignment-2/internal"
	"encoding/json"
	"fmt"
)

func GetMeteoData(latitude float64, longitude float64) (meteoObject external.MeteoObject, err error) {
	url := GetMeteoUrl(latitude, longitude)
	response, err := external.GetExternalResource(url)
	if err != nil {
		return meteoObject, err
	}
	err = json.NewDecoder(response.Body).Decode(&meteoObject)
	return meteoObject, err
}

/*
GetMeteoUrl
lat Latitude
long Longitude
returns the MeteoURL for the given coordinates
*/
func GetMeteoUrl(lat float64, long float64) string {
	arguments := fmt.Sprintf(internal.MeteoField, lat, long)
	return internal.MeteoApi + arguments
}
