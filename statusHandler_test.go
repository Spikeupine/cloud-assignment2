package assignment_two

import (
	"assignment-2/database"
	"assignment-2/internal"
	"assignment-2/internal/handlers"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testServer *httptest.Server
	testClient http.Client
)

func statusSetup() {
	godotenv.Load()
	database.FirebaseConnect()
	// Setting up the test server
	testServer = httptest.NewServer(http.HandlerFunc(handlers.StatusHandler))
	fmt.Println("Running test server on " + testServer.URL)

	// Setting up the client
	testClient = http.Client{}
}

func statusTeardown() {
	database.FireBaseCloseConnection()
	// Closing the test server
	testServer.Close()
	fmt.Println("Closing test server on " + testServer.URL)

	// Closing the client
	testClient.CloseIdleConnections()
}

func TestStatusHandler(t *testing.T) {
	// Run setup at the start of the test and teardown at the end
	statusSetup()
	defer statusTeardown()

	testCases := []struct {
		testName            string
		httpMethod          string
		expectedContentType string
		expectedStatusCode  int
	}{
		{
			testName:            "GET request",
			httpMethod:          http.MethodGet,
			expectedContentType: internal.ContentTypeJson,
			expectedStatusCode:  http.StatusOK,
		},
		{
			testName:            "Unsupported request method",
			httpMethod:          http.MethodPost,
			expectedContentType: "text/plain; charset=utf-8",
			expectedStatusCode:  http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest(tc.httpMethod, testServer.URL, nil)
			if err != nil {
				t.Fatal("Error when creating request: " + err.Error())
			}

			// Issue request
			res, err := testClient.Do(req)
			if err != nil {
				t.Fatal("Error when issuing request: " + err.Error())
			}
			defer res.Body.Close()

			// Check test case results
			assert.Equal(t, tc.expectedStatusCode, res.StatusCode, "Status code should match")
			contentType := res.Header.Get("Content-Type")
			assert.Equal(t, tc.expectedContentType, contentType, "Content type should match")

			if res.StatusCode == http.StatusOK {
				var response internal.Status
				if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
					t.Fatal("Error decoding JSON response: " + err.Error())
				}

				// Check if the response has correct fields
				assert.NotEmpty(t, response.CountriesAPI, "Countries API status should be set")
				assert.NotEmpty(t, response.MeteoAPI, "Meteo API status should be set")
				assert.NotEmpty(t, response.CurrencyAPI, "Currency API status should be set")
				assert.Equal(t, "v1", response.Version)
				assert.Greater(t, response.Uptime, int64(0), "Uptime should be greater than 0")
			}
		})
	}
}
