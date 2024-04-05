package handlers

import (
	"assignment-2/internal"
	"encoding/json"
	"net/http"
)

func HandleRestcountriesapi(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		//Introduces an empty struct of country.
		var country internal.RestCountriesStruct

		//The country is set using the getRestCountries method.
		country, _ = getRestCountries(w)

		err := printToUserInterfaceRestCountries(w, country)
		if err != nil {
			http.Error(w, "error when sending rest countries information to user.", http.StatusInternalServerError)
		}

	}
}

// Gets the country related to requested iso code. Return is a struct with information about that country, and an
// error that is empty if nothing went wrong. Else, a fitting error code with an empty struct.
func getRestCountries(w http.ResponseWriter) (internal.RestCountriesStruct, error) {

	//Sends a http-request with specified url. The response of the request is tied to the response-variable.
	//an error is created if something went wrong, and that is stored inside err-variable.
	request, err := http.NewRequest(http.MethodGet, internal.CountriesApi+internal.IsoExample, nil)

	//Handles error if error is not empty.
	if err != nil {

		//Specifies that the error is related to an unavailable service.
		//todo: Write check for valid ISO code.
		http.Error(w, "Failed to get response from country api "+err.Error(), http.StatusServiceUnavailable)

		//Returns the error, with an empty struct.
		return internal.RestCountriesStruct{}, err
	}

	//Initializes client.
	client := &http.Client{}

	//Response from the webservice is created. Request is inserted into client, that performs the request.
	response, err := client.Do(request)

	//Handles error in response if any.
	if err != nil {
		http.Error(w, "error when sending request URL into system.", http.StatusInternalServerError)
		return internal.RestCountriesStruct{}, err
	}

	//Introduces rest api struct to populate with the response.
	var restApiVariables []internal.RestCountriesStruct

	//Decodes the response into restApiVariables struct.
	err = json.NewDecoder(response.Body).Decode(&restApiVariables)
	if err != nil {

		//If unsuccessfull, handles appropriate error, and returns empty struct and error.
		http.Error(w, "error when decoding the rest country api response. "+err.Error(), http.StatusInternalServerError)
		return internal.RestCountriesStruct{}, err
	}

	//Returns the rest api variables as struct, and error if any.
	return restApiVariables[0], err
}

func printToUserInterfaceRestCountries(w http.ResponseWriter, country internal.RestCountriesStruct) error {
	// Converts the result to JSON format
	responseBody, err := json.Marshal(country)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	//Sets header type.
	w.Header().Add(internal.ApplicationJson, internal.ContentTypeJson)

	//Encodes the JSON content.
	err = json.NewEncoder(w).Encode(responseBody)
	if err != nil {

		//If unsuccessfull, handles appropriate error, and returns error.
		http.Error(w, "error when encoding the rest country api response. "+err.Error(), http.StatusInternalServerError)
		return err
	}
	w.Write(responseBody)

	return err
}
