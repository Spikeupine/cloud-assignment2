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

var responseRegistration *http.Response
var WebhookRegistration internal.Webhook
var server httptest.Server

var rec = httptest.NewRecorder()
var client http.Client

var SeveralIds []string

func registerTestingId(webhookId string) {
	SeveralIds = append(SeveralIds, webhookId)
}

func getIds() []string {
	return SeveralIds
}

func TestMain(m *testing.M) {
	err := godotenv.Load()
	exitcode := m.Run()
	if err != nil {
		os.Exit(1)
	}
	database.FirebaseConnect()
	os.Exit(exitcode)
}

func TestWebhookRegistration(t *testing.T) {
	err := godotenv.Load()

	if err != nil {
		os.Exit(1)
	}
	database.FirebaseConnect()

	//Here I am creating a webhook to perform test on.
	WebhookRegistration := internal.Webhook{
		Url:     "https://webhook.site/22b1fade-ac45-431c-81a6-8f68a918b7c6",
		Country: "TestTestTest",
		Event:   "REGISTER",
	}

	// marshals the webhook registration so it is as json.
	body, err := json.Marshal(WebhookRegistration)

	//Sets up the server to the endpoint.
	server := httptest.NewServer(http.HandlerFunc(handlers.NotificationsHandler))
	defer server.Close()

	url := server.URL

	//Initializes client.
	client := http.Client{}

	rec := httptest.NewRecorder()

	responseRegistration, err := client.Post(url, "Content type: application/json", bytes.NewBuffer(body))
	//response, err := client.Post("https://localhost:8080/dashboards/v1/notifications/", "Content type: application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("errer" + err.Error())
	}

	err = json.NewDecoder(responseRegistration.Body).Decode(&WebhookRegistration)
	if err != nil {
		t.Errorf("Error in getting response from posting new webhook " + err.Error())
		t.Fatal()
	}

	handlers.WebhookRegistration(rec, responseRegistration.Request, "webhooks")

	newurl := server.URL + "/" + WebhookRegistration.WebhookId
	println(newurl)
	registerTestingId(WebhookRegistration.WebhookId)

	// marshals the webhook registration so it is as json.
	//body, err = json.Marshal(webhookRegistration.WebhookId)

	responsehaha := httptest.NewRequest(http.MethodGet, newurl, bytes.NewBuffer(body))

	//responseGetwebhook, err := client.Get(newurl)
	//if err != nil {
	//	t.Errorf("Error sending get request to notification service " + err.Error())
	//}

	var getHook internal.Webhook

	err = json.NewDecoder(responsehaha.Body).Decode(&getHook)
	if err != nil {
		t.Errorf("Error in getting response from posting new webhook " + err.Error())
		t.Fatal()
	}
	println(SeveralIds[0])
}

func TestDeleteWebhook(t *testing.T) {

	err := godotenv.Load()
	if err != nil {
		os.Exit(1)
	}
	database.FirebaseConnect()

	//Initializes client.
	client := http.Client{}

	rec := httptest.NewRecorder()

	//Sets up the server to the endpoint.
	server := httptest.NewServer(http.HandlerFunc(handlers.NotificationsHandler))

	defer server.Close()
	listOfIdsOfWebhooksToDelete := getIds()

	for _, id := range listOfIdsOfWebhooksToDelete {
		url := server.URL + "/" + id
		respondent, err := client.Post(url, http.MethodDelete, nil)

		handlers.DeleteWebhook(rec, respondent.Request, "webhooks", WebhookRegistration.WebhookId)

		responseGetwebhook, err := client.Get(url)

		if err != nil {
			t.Errorf("Error sending get request to notification service " + err.Error())
		}
		var testHookDelete internal.Webhook

		err = json.NewDecoder(responseGetwebhook.Body).Decode(&testHookDelete)
		if err == nil {
			t.Errorf("Error in deleting webhook, as we can retrieve it even after deletion" + err.Error())
			t.Fatal()
		}
		err = handlers.GetWebhook(rec, "webhooks", responseGetwebhook.Request, testHookDelete.WebhookId)
		if err == nil {
			t.Fatal()
		}

		// check test case results
		webhook, errorWebhook := database.GetWebhook("webhooks", testHookDelete.WebhookId)
		if errorWebhook == nil {
			t.Fatal()
		}

		asrt := assert.New(t)

		asrt.Equal(webhook, internal.Webhook{})

	}
}







