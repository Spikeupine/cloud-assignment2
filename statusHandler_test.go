package assignment_two

import (
	"assignment-2/internal/handlers"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testServer *httptest.Server
	testClient http.Client
)

func statusSetup() {
	// Setting up the test server
	testServer = httptest.NewServer(http.HandlerFunc(handlers.StatusHandler))
	fmt.Println("Running test server on " + testServer.URL)

	// Setting up the client
	testClient = http.Client{}
}

func statusTeardown() {
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

}
