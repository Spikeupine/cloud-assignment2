package assignment_two

import (
	"assignment-2/database"
	"assignment-2/external/router"
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

// severalIds is responsible for holding the ids of webhooks created within this testfile. Gives control to delete them
// once the tests are over, and is also checked with assert empty to ensure this.
var severalIds []string

// registerTestingID takes in the id of the webhook that's registered in the system, and appends it to the list of
// webhooks registered from the test file.
func registerTestingId(webhookId string) {
	severalIds = append(severalIds, webhookId)
}

// Returns the list with webhook ids
func getIds() []string {
	return severalIds
}

// allIdsDeleted wipes the list, because it initializes a new empty one.
func allIdsDeleted() {
	severalIds = []string{}
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

// TestWebhookRegistration runs tests on the registration of webhooks. It creates requests, and uses the requests in the
// methods to perform the different scenarios. Adds the webhooks it has registered to the list of webhooks made within
// test file, for later deletion. It is its own method for code readability, and easier reading of where the tests
// necessarily fail.
func TestWebhookRegistration(t *testing.T) {
	TestSetup()
	defer TestTearDown()

	//Here I am creating a webhook to perform test on.
	WebhookRegistration := internal.Webhook{
		Url:     "https://webhook.site/22b1fade-ac45-431c-81a6-8f68a918b7c6",
		Country: "NO",
		Event:   "REGISTER",
	}

	// marshals the webhook registration so it is as json.
	body, err := json.Marshal(WebhookRegistration)

	//Sets up the server to the endpoint.
	server := httptest.NewServer(http.HandlerFunc(handlers.NotificationsHandler))
	err = os.Setenv(router.TESTING, "true")
	if err != nil {
		t.Fatal("Cannot test when stubs cannot be used.")
	}
	defer server.Close()

	//Specified the url to use in request.
	url := server.URL

	//Initializes client.
	client := http.Client{}

	// This recorder will record http errors for example.
	rec := httptest.NewRecorder()

	//sending the http request post to testserver url notifications handler with body of new webhook to be registered.
	responseRegistration, err := client.Post(url, "Content type: application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("errer" + err.Error())
	}

	// Decodes the responsebody into webhook struct.
	err = json.NewDecoder(responseRegistration.Body).Decode(&WebhookRegistration)
	if err != nil {
		t.Errorf("Error in getting response from posting new webhook " + err.Error())
		t.Fatal()
	}

	//Uses the request from registered response to make a new webhook registration. It is put into webhooks in firebase.
	handlers.WebhookRegistration(rec, responseRegistration.Request, database.WebhookCollection)

	//initializes recorder to record the responses
	record := httptest.NewRecorder()

	//Calls method from notification handler to get the webhook. Response is recorded to record.
	handlers.GetWebhook(record, database.WebhookCollection, WebhookRegistration.WebhookId)

	//Instatiates assert, so that we can compare our results to what we expect.
	asrt := assert.New(t)

	//Asserts that we expect status code 200, and checks what the recorder actually has recorded of status code.
	asrt.Equal(http.StatusOK, record.Code)

	recordagain := httptest.NewRecorder()

	//Trying to process the wrong content
	handlers.GetWebhook(recordagain, database.WebhookCollection, "723kjk")

	//I expect the recorder should return a bad request-status, as I have written it should in the code.
	//recorder will be measured up to this.
	asrt.Equal(http.StatusBadRequest, recordagain.Code)

}

// TestGetWebhooks tests method of retrieving all the webhooks in the collection. Since this request goes to the actual
// webhooks in firebase, we first register the number for how many there are now from GET request with unspecified id,
// add a few, and then compare the number. For each webhook added, they are added to the list of webhooks created within
// test file.
func TestGetWebhooks(t *testing.T) {

	err := godotenv.Load()
	if err != nil {
		os.Exit(1)
	}
	database.FirebaseConnect()
	err = os.Setenv(router.TESTING, "true")
	if err != nil {
		t.Fatal("Cannot test when stubs cannot be used.")
	}

	record := httptest.NewRecorder()

	//Sends the recorder inn to record what response I get from running GetWebhooks method with the request made earlier
	handlers.GetWebhooks(record, database.WebhookCollection)

	asser := assert.New(t)

	//Should return OK, recorder will show.
	asser.Equal(http.StatusOK, record.Code)

	var listOfAllHooksInDatabase []internal.Webhook

	//Tries to insert the response from the body of recorder into the list of webhooks I made just over.
	err = json.NewDecoder(record.Body).Decode(&listOfAllHooksInDatabase)
	if err != nil {
		t.Errorf("Error in reading content into list of webhooks " + err.Error())
		t.Fatal()
	}

	//As it says, it is the old number of webhooks, before I insert any new ones to test.
	oldNumberOfHooks := len(listOfAllHooksInDatabase)

	//To double check, I also use the method created to count all the webhooks.
	numberOfHooks, err := database.CountWebhooks(database.WebhookCollection)

	//Here I am creating a webhook to perform test on.
	newWebhook := internal.Webhook{
		Url:     "https://webhook.site/22b1fade-ac45-431c-81a6-8f68a918b7c6",
		Country: "NO",
		Event:   "REGISTER",
	}
	//Here I am creating a webhook to perform test on.
	newerWebhook := internal.Webhook{
		Url:     "https://webhook.site/22b1fade-ac45-431c-81a6-8f68a918b7c6",
		Country: "NO",
		Event:   "REGISTER",
	}

	//Here I am creating a webhook to perform test on.
	newestWebhook := internal.Webhook{
		Url:     "https://webhook.site/22b1fade-ac45-431c-81a6-8f68a918b7c6",
		Country: "NO",
		Event:   "REGISTER",
	}

	//Adds all three webhooks I just created to a list to add them to database. Intention is to compare the numbers
	// before and after to see if GetWebhooks will provide a deterministic number.
	var webhooksToAdd []internal.Webhook
	webhooksToAdd = append(webhooksToAdd, newestWebhook, newWebhook, newerWebhook)

	//Loops over the different webhooks just created to add them to the collection.
	for _, webhook := range webhooksToAdd {
		// marshals the webhook registration so it is as json.
		body, err := json.Marshal(webhook)

		//Sets up the server to the endpoint.
		server := httptest.NewServer(http.HandlerFunc(handlers.NotificationsHandler))
		defer server.Close()

		//Specified the url to use in request.
		url := server.URL

		//Initializes client.
		client := http.Client{}

		// This recorder will record http errors for example.
		rec := httptest.NewRecorder()

		//sending the http request post to testserver url notifications handler with body of new webhook to be registered.
		responseRegistration, err := client.Post(url, "Content type: application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Errorf("errer" + err.Error())
		}

		// Decodes the responsebody into webhook struct.
		err = json.NewDecoder(responseRegistration.Body).Decode(&webhook)
		if err != nil {
			t.Errorf("Error in getting response from posting new webhook " + err.Error())
			t.Fatal()
		}

		//Uses the request from registered response to make a new webhook registration. It is put into webhooks in firebase.
		handlers.WebhookRegistration(rec, responseRegistration.Request, database.WebhookCollection)

		//Here the id of the webhook is registered in the list of webhooks created within test environment, to be deleted
		//after test.
		registerTestingId(webhook.WebhookId)

	}

	// New number for webhokks in database. Should be three more than what was.
	newNumberOfWebhooks, err := database.CountWebhooks(database.WebhookCollection)
	if err != nil {
		t.Errorf("Error counting webhooks")
	}

	asrt := assert.New(t)

	//Old number of webhooks from count method should be the same as the one where i checked length of list I created
	// with all the webhooks registered.
	asrt.Equal(oldNumberOfHooks, numberOfHooks)

	//Old number plus three should be the same as the new number.
	asrt.Equal(oldNumberOfHooks+3, newNumberOfWebhooks)

	//The status code from the recorder should be 200.
	asrt.Equal(http.StatusOK, record.Code)

}

func TestGetWebhook(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		os.Exit(1)
	}

	err = os.Setenv(router.TESTING, "true")
	if err != nil {
		t.Fatal("Cannot test when stubs cannot be used.")
	}
	database.FirebaseConnect()

	server := httptest.NewServer(http.HandlerFunc(handlers.NotificationsHandler))

	//Here I am creating a webhook to perform test on.
	WebhookRegistration := internal.Webhook{
		Url:     "https://webhook.site/22b1fade-ac45-431c-81a6-8f68a918b7c6",
		Country: "NO",
		Event:   "INVOKE",
	}

	// marshals the webhook registration so it is as json.
	body, err := json.Marshal(WebhookRegistration)

	//Specified the url to use in request.
	url := server.URL

	//Initializes client.
	client := http.Client{}

	// This recorder will record http errors for example.
	rec := httptest.NewRecorder()

	//sending the http request post to testserver url notifications handler with body of new webhook to be registered.
	responseRegistration, err := client.Post(url, "Content type: application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("errer" + err.Error())
	}

	var webhookRegistration internal.Webhook

	// Decodes the responsebody into webhook struct.
	err = json.NewDecoder(responseRegistration.Body).Decode(&webhookRegistration)
	if err != nil {
		t.Errorf("Error in getting response from posting new webhook " + err.Error())
		t.Fatal()
	}

	//Uses the request from registered response to make a new webhook registration. It is put into webhooks in firebase.
	handlers.WebhookRegistration(rec, responseRegistration.Request, database.WebhookCollection)

	//Here the id of the webhook is registered in the list of webhooks created within test environment, to be deleted
	//after test.
	registerTestingId(webhookRegistration.WebhookId)

	//initializes recorder to record the responses
	record := httptest.NewRecorder()

	//Calls method from notification handler to get the webhook. Response is recorded to record.
	handlers.GetWebhook(record, database.WebhookCollection, webhookRegistration.WebhookId)

	//Instatiates assert, so that we can compare our results to what we expect.
	asrt := assert.New(t)

	asrt.Equal(http.StatusOK, record.Code)

	//initializes recorder to record the responses
	recordNew := httptest.NewRecorder()

	//Calls method from notification handler to get the webhook. Response is recorded to record. Nonesense
	// webhook id is passed on as parameter.
	handlers.GetWebhook(recordNew, database.WebhookCollection, "nonesense123")

	// As it is a bad request, recorder should return bas request, as expected.
	asrt.Equal(http.StatusBadRequest, recordNew.Code)

	recordAnew := httptest.NewRecorder()

	//Testing what happens when the GetWebhook method is called with empty webhook id.
	handlers.GetWebhook(recordAnew, database.WebhookCollection, "")

	//Should return StatusBadRequest, checks with recorder.
	asrt.Equal(http.StatusBadRequest, recordAnew.Code)

}

// TestDeleteWebhook deletes the webhooks that has been added to the collection of webhook id's within this test file.
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

	err = os.Setenv(router.TESTING, "true")
	if err != nil {
		t.Fatal("Cannot test when stubs cannot be used.")
	}

	//Initializes client.
	client := http.Client{}

	rec := httptest.NewRecorder()

	//Sets up the server to the endpoint.
	server := httptest.NewServer(http.HandlerFunc(handlers.NotificationsHandler))
	// When finished with code, it closes.
	defer server.Close()
	// gets the id's of the webhooks to delete, and it is a list of strings.
	listOfIdsOfWebhooksToDelete := getIds()

	//Logic is: for each webhook id, delete the webhook, and test after that if you can get the webhook from collection.
	for _, id := range listOfIdsOfWebhooksToDelete {
		//Specifies which webhook to delete in url.
		url := server.URL + "/" + id

		//Method to test. rec records the response.
		handlers.DeleteWebhook(rec, database.WebhookCollection, id)

		//get-request to the url with the specified webhook is sent.
		responseGetwebhook, err := client.Get(url)
		if err != nil {
			t.Errorf("Error sending get request to notification service " + err.Error())
		}
		var testHookDelete internal.Webhook

		//Here we try to decode the response into a variable of webhook to populate it. It sbould end up empty if
		//deleted.
		err = json.NewDecoder(responseGetwebhook.Body).Decode(&testHookDelete)
		if err == nil {
			t.Errorf("Error in deleting webhook, as we can retrieve it even after deletion" + err.Error())
			t.Fatal()
		}

		asrt := assert.New(t)

		//Checks if that webhook is actually empty, as it should be.
		asrt.Empty(testHookDelete)

		//Tries out the method for GetWebhook with if of deleted webhook. Records response in rec.
		handlers.GetWebhook(rec, database.WebhookCollection, id)

		asrt.Equal(testHookDelete.WebhookId, "")
		asrt.Equal(http.StatusBadRequest, rec.Code)

	}

	//Sends the request to method, with empty webhook id.
	handlers.DeleteWebhook(rec, database.WebhookCollection, "")

	asrt := assert.New(t)

	//Asserts that it should return StatusBadRequest
	asrt.Equal(http.StatusBadRequest, rec.Code)

	//Deletes all the id's from the list.
	allIdsDeleted()

	//asserts that all the id's have ben deleted from the list, by asserting empty.
	asrt.Empty(severalIds, "Expect no more IDs in the list of registered webhooks within test file")
}
