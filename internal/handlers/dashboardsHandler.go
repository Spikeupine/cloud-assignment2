package handlers

import (
	"assignment-2/internal"
	"encoding/json"
	"fmt"
	"net/http"
)

func DashboardsHandler() {

}

// TODO: change "client *http.Client" to "client *firestore.Client"
func getDashboard(id string, client *http.Client) (internal.PopulatedDashboard, error) {
	url := internal.DashboardsPath + id
	response, err := client.Get(url)

	if err != nil {
		return internal.PopulatedDashboard{}, fmt.Errorf("Failed to fetch information: %v", err)
	}
	defer response.Body.Close()

	var dashboard internal.PopulatedDashboard
	err = json.NewDecoder(response.Body).Decode(&dashboard)
	if err != nil {
		return internal.PopulatedDashboard{}, fmt.Errorf("Failed to decode information: %v", err)
	}

	return dashboard, nil
}
