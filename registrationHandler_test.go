package assignment_two

import (
	"assignment-2/database"
	"assignment-2/internal"
	"assignment-2/internal/handlers"
	"context"
	"testing"
)

const regId = "4f8eb1d2-2245-4953-b4fa-3b04be6505d9"
const isoCode = "no"
const country = "norway"
const currency = "EUR"

func getRegistration() (registration internal.RegisterRequest) {
	registration.Id = regId
	registration.IsoCode = isoCode
	registration.Country = country
	registration.Features = getFeatures()
	return registration
}

func getFeatures() (features internal.Features) {
	features.Area = true
	features.Capital = true
	features.Coordinates = true
	features.Precipitation = true
	features.Population = true
	features.TargetCurrencies = append(features.TargetCurrencies, currency)
	features.Temperature = true
	return features
}

func TestGetAllRegistrations(t *testing.T) {
	TestSetup()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer TestTearDown()

	registrations, err := handlers.GetAllRegistrations(ctx)
	if err != nil || len(registrations) == 0 {
		t.Fatalf("Could not get all registrations %s", err)
	}
}

func createRegistration(ctx context.Context) (internal.RegistrationsResponse, error) {
	return handlers.UploadDashboard(getRegistration(), ctx)
}

func deleteRegistration(ctx context.Context, id string) error {
	return handlers.DeleteRegistration(ctx, database.DashboardCollection, id)
}

func TestCreateAndDeleteRegistration(t *testing.T) {
	TestSetup()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer TestTearDown()
	response, err := createRegistration(ctx)
	if err != nil || response.Id != regId {
		t.Fatalf("Could not upload registration with id: %s : %v", regId, err)
	}
	err = deleteRegistration(ctx, response.Id)
	if err != nil {
		t.Fatalf("New registation not deleted %v", err)
	}
}
