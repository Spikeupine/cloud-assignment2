package handlers

import (
	"assignment-2/database"
	"os"
	"testing"
)

const TestCollection = "testWebhooks"

func TestMain(m *testing.M) {
	database.FirebaseConnect()
	code := m.Run()
	err := database.FireBaseCloseConnection
	if err != nil {
		os.Exit(1)
	}
	os.Exit(code)
}
