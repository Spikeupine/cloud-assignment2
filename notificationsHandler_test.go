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

var WebhookRegistration internal.Webhook

//SeveralIds is responsible for holding the ids of webhooks created within this testfile. Gives control to delete them
// once the tests are over, and is also checked with assert empty to ensure this.
var SeveralIds []string

//registerTestingID takes in the id of the webhook that's registered in the system, and appends it to the list of
// webhooks registered from the test file.
func registerTestingId(webhookId string) {
	SeveralIds = append(SeveralIds, webhookId)
}

//Returns the list with webhook ids
func getIds() []string {
	return SeveralIds
}

// allIdsDeleted wipes the list, because it initializes a new empty one.
func allIdsDeleted() {
	SeveralIds = []string{}
}

// TestMain is actually how we wished to set up the test environment. HOwever, it did not work, and we therefore set
// it up in each test method. We keep this here, because we wish to show our preferred solution, and the code also
// does not take damage of it being here.
func TestMain(m *testing.M) {
	err := godotenv.Load()
	exitcode := m.Run()
	if err != nil {
		os.Exit(1)
	}
	database.FirebaseConnect()
	os.Exit(exitcode)
}

//TestWebhookRegistration runs tests on the registration of webhooks. It creates requests, and uses the requests in the
// methods to perform the different scenarios. Adds the webhooks it has registered to the list of webhooks made within
// test file, for later deletion. It is its own method for code readability, and easier reading of where the tests
// necessarily fail.
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

	responsehaha := httptest.NewRequest(http.MethodGet, newurl, bytes.NewBuffer(body))

	var getHook internal.Webhook

	err = json.NewDecoder(responsehaha.Body).Decode(&getHook)
	if err != nil {
		t.Errorf("Error in getting response from posting new webhook " + err.Error())
		t.Fatal()
	}
}

//TestDeleteWebhook deletes the webhooks that has been added to the collection of webhook id's within this test file.
// It ranges over this list, and deletes them by creating http DELETE requests, and passing this as parameter to
// DeleteWebhook method of handler. The recorded http status code is then later compared, and a get request to the
// individual webhooks is tried sent (and test is failed if it successfully manages to retrieve webhook supposed
// to be deleted). Also empties list after all is deleted from list, and asserts it empty.
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
		asrt.Equal(testHookDelete.WebhookId, "")

	}

	url := server.URL
	respondent, err := client.Post(url, http.MethodDelete, nil)

	handlers.DeleteWebhook(rec, respondent.Request, "webhooks", "")

	asrt := assert.New(t)
	asrt.Equal(http.StatusBadRequest, rec.Code)

	allIdsDeleted()

	asrt.Empty(SeveralIds, "Expect no more IDs in the list of registered webhooks within test file")
}

// TestGetWebhooks tests method of retrieving all the webhooks in the collection. Since this request goes to the actual
// webhooks in firebase, we first register the number for how many there are now from GET request with unspecified id,
//add a few, and then compare the number. For each webhook added, they are added to the list of webhooks created within
// test file.
//
func TestGetWebhooks(t *testing.T) {

}
