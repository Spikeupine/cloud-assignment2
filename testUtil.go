package assignment_two

import (
	"assignment-2/database"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
)

const port = "8081"
const portName = "PORT"
const testingMode = "true"
const testingName = "TESTING"

func TestSetup() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}
	database.FirebaseConnect()
	err = os.Setenv(portName, port)
	if err != nil {
		log.Fatal("Error setting port ")
	}
	err = os.Setenv(testingName, testingMode)
	if err != nil {
		log.Fatal("Error setting testing mode")
	}
}

func TestTearDown() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	database.FireBaseCloseConnection()
}
