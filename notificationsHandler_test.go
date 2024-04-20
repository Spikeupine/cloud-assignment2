package assignment_two

import (
	"assignment-2/database"
	"assignment-2/internal"
	"assignment-2/internal/handlers"
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := godotenv.Load()
	exitcode := m.Run()
	if err != nil {
		os.Exit(1)
	}
	database.FirebaseConnect()
	os.Exit(exitcode)
}
func TestDeleteWebhook(t *testing.T) {

	err := godotenv.Load()

	if err != nil {
		os.Exit(1)
	}
	database.FirebaseConnect()

	//Here I am creating a webhook to perform test on.
	hook := internal.Webhook{
		Url:     "https://webhook.site/22b1fade-ac45-431c-81a6-8f68a918b7c6",
		Country: "TestTestTest",
		Event:   "REGISTER",
	}

	// marshals the webhook registration so it is as json.
	body, err := json.Marshal(hook)

	//Sets up the server to the endpoint.
	server := httptest.NewServer(http.HandlerFunc(handlers.NotificationsHandler))
	defer server.Close()

	url := server.URL

	//Initializes client.
	client := http.Client{}

	rec := httptest.NewRecorder()

	response, err := client.Post(url, "Content type: application/json", bytes.NewBuffer(body))
	//response, err := client.Post("https://localhost:8080/dashboards/v1/notifications/", "Content type: application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("errer" + err.Error())
	}

	err = json.NewDecoder(response.Body).Decode(&hook)
	if err != nil {
		t.Errorf("Error in getting response from posting new webhook " + err.Error())
		t.Fatal()
	}

	handlers.WebhookRegistration(rec, response.Request, "webhooks")

	handlers.DeleteWebhook(rec, response.Request, "webhooks", hook.WebhookId)

	url = server.URL + "/" + hook.WebhookId

	response, err = client.Get(url)
	//response, err := client.Post("https://localhost:8080/dashboards/v1/notifications/", "Content type: application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("errer" + err.Error())
	}
	var testHookDelete internal.Webhook

	err = json.NewDecoder(response.Body).Decode(&testHookDelete)
	if err == nil {
		t.Errorf("Error in deleting webhook, as we can retrieve it even after deletion" + err.Error())
		t.Fatal()
	}
	err = handlers.GetWebhook(rec, "webhooks", response.Request, testHookDelete.WebhookId)
	if err == nil {
		t.Fatal()
	}

	// check test case results
	webhook, errorWebhook := database.GetWebhook("webhooks", hook.WebhookId)
	if errorWebhook == nil {
		t.Fatal()
	}

	asrt := assert.New(t)

	asrt.Equal(webhook, internal.Webhook{})

}
