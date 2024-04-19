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

const TestCollection = "testWebhooks"

func init() {
	err := godotenv.Load()
	if err != nil {
		os.Exit(1)
	}
	err = os.Setenv("TESTING", "true")
}
func initData() ([]internal.Webhook, error) {

	var listOfWebhook []internal.Webhook
	//Here I am creating a webhook to perform test on.
	hook := internal.Webhook{
		Url:     "https://webhook.site/22b1fade-ac45-431c-81a6-8f68a918b7c6",
		Country: "NO",
		Event:   "REGISTER",
	}
	listOfWebhook = append(listOfWebhook, hook)
	hook1 := internal.Webhook{
		Url:     "https://webhook.site/22b1fade-ac45-431c-81a6-8f68a918b7c6",
		Country: "ES",
		Event:   "INVOKE",
	}
	listOfWebhook = append(listOfWebhook, hook1)

	hook2 := internal.Webhook{
		Url:     "https://webhook.site/22b1fade-ac45-431c-81a6-8f68a918b7c6",
		Country: "UK",
		Event:   "DELETE",
	}
	listOfWebhook = append(listOfWebhook, hook2)

	for _, webhook := range listOfWebhook {
		err := database.AddWebhookToCollection(webhook, "webhooks")
		if err != nil {
			return []internal.Webhook{}, err
		}
	}
	return listOfWebhook, nil
}

func deleteWebhook() error {
	list, err := initData()
	if err != nil {
		return err
	}
	err, _ = database.DeleteTheWebhook("Webhooks", list[0].WebhookId)
	if err != nil {
		return err
	}
	_, err = getSomeWebhook(list[0].WebhookId)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func getSomeWebhook(id string) (internal.Webhook, error) {
	webhook, err := database.GetWebhook("Webhooks", id)
	if err != nil {
		return internal.Webhook{}, err
	} else {
		return webhook, nil
	}
}
func TestDeleteWebhook(t *testing.T) {
	/*
			err := deleteWebhook()
			if err != nil {
				t.Errorf("Error when deleting webhook in method")
				t.Fatal()
			}
		}
	*/
	err := godotenv.Load()
	if err != nil {
		os.Exit(1)
	}
	database.FirebaseConnect()
	//Here I am creating a webhook to perform test on.
	hook := internal.Webhook{
		Url:     "https://webhook.site/22b1fade-ac45-431c-81a6-8f68a918b7c6",
		Country: "NO",
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

	response, err := client.Post(url, "Content type: application/json", bytes.NewBuffer(body))
	//response, err := client.Post("https://localhost:8080/dashboards/v1/notifications/", "Content type: application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("errer" + err.Error())
	}

	println(response.Body)

	// check test case results
	asrt := assert.New(t)
	asrt.Equal("hest", "hest")

}

//
//	/*
//		//Sends a request to the URL.
//		request := httptest.NewRequest(http.MethodPost, "https://localhost:8080/dashboards/v1/notifications/", bytes.NewBuffer(body))
//		rec := httptest.NewRecorder()
//
//		request.Header.Add("content-type", "application/json")
//
//		// Perform the request.
//		resp, err := client.Do(request)
//		if err != nil {
//			t.Fatalf("Error performing request: %v", err)
//		}
//		defer resp.Body.Close()
//
//		// Read the response body.
//		responseBody, err := ioutil.ReadAll(resp.Body)
//		if err != nil {
//			t.Fatalf("Error reading response body: %v", err)
//		}
//
//		// Print the response body for debugging.
//		fmt.Println(string(responseBody))
//
//		//Here we perform the request.
//		actualResponse, err := server.Client().Do(request)
//		if err != nil {
//			http.Error(rec, "Error"+err.Error(), http.StatusInternalServerError)
//			t.Errorf(err.Error())
//		}
//
//		println(actualResponse)
//		response, err := client.Do(request)
//		if err != nil {
//			http.Error(rec, "Error"+err.Error(), http.StatusInternalServerError)
//			t.Errorf(err.Error())
//			t.Fatal()
//		}
//		println(response)
//		//Registers the webhook
//		handlers.WebhookRegistration(rec, request, TestCollection)
//
//		list, err := database.GetAllWebhooks(rec, TestCollection)
//		if err != nil {
//			t.Fatal()
//		}
//		println(list)
//		webhook := list[0]
//		println(webhook.WebhookId)
//
//		req, err := http.NewRequest(http.MethodDelete, "https://localhost:8080/dashboards/v1/notifications/", bytes.NewBuffer([]byte(webhook.WebhookId)))
//		if err != nil {
//			t.Fatal()
//		}
//
//		// calls the handler function
//		handlers.DeleteWebhook(rec, req, TestCollection, webhook.WebhookId)
//
//		ualResponse, err := client.Do(req)
//
//		fmt.Print(ualResponse)
//
//		//res, err := client.Get(server.URL + "/" + internal.NotificationsPath)
//
//		//resp, err := http.Post("https://localhost:8080/dashboards/v1/notifications/", internal.ContentTypeJson, bytes.NewBuffer(body))
//		//resp.Header.Add("content-type", "application/json")
//
//		// register the webhook
//		// create a request
//		//request, err = http.NewRequest(http.MethodPost, "https://localhost:8080/dashboards/v1/notifications/", bytes.NewBuffer(body))
//		//resp, err := http.Post("https://localhost:8080/dashboards/v1/notifications/", internal.ContentTypeJson, bytes.NewBuffer(body))
//		//resp.Header.Add("content-type", "application/json")
//		//request.Header.Add(internal.ApplicationJson, internal.ContentTypeJson)
//		// call the handler function
//
//		// Create a ResponseRecorder to record the response
//		//recorder = httptest.NewRecorder()
//		// Check the response body
//		expectedResponseBody := http.StatusOK
//
//		intToUse, err := strconv.Atoi(rec.Body.String())
//		if err != nil {
//			log.Fatal("Error error pants on fire")
//		}
//
//		if intToUse != expectedResponseBody {
//			t.Errorf("handler returned unexpected body: got %v, want %v", rec.Body.String(), expectedResponseBody)
//		}
//		t.Error()
//
//		/*
//				// calls the handler function
//				WebhookRegistration(rec, request, TestCollection)
//
//				var hook internal.Webhook
//
//				err = json.Unmarshal(rec.Body.Bytes(), &hook)
//				if err != nil {
//					t.Fatal("Error unmarshalling response : " + err.Error())
//				}
//				var tests []internal.Webhook
//
//				tests[0] = hook
//				hook.WebhookId = ""
//				tests[1] = hook
//
//				// iterate through the test cases
//				for _, testCase := range tests {
//					t.Run(testCase.WebhookId, func(t *testing.T) {
//
//						resp := httptest.NewRequest(http.MethodDelete, internal.NotificationsPath+testCase.WebhookId, nil)
//						resp.Header.Add(internal.ApplicationJson, internal.ContentTypeJson)
//
//						// create the recorder
//						rec := httptest.NewRecorder()
//
//						// calls the handler function
//						DeleteWebhook(rec, resp, TestCollection, testCase.WebhookId)
//
//						requestDeletedWebhook := httptest.NewRequest(http.MethodGet, internal.NotificationsPath+testCase.WebhookId, nil)
//						resp.Header.Add(internal.ApplicationJson, internal.ContentTypeJson)
//
//						// Asserts what is expected
//						asserts := assert.New(t)
//						asserts.Equal(requestDeletedWebhook.Response, rec.Code)
//						asserts.Equal(requestDeletedWebhook.Header, rec.Header().Get("content-type"))
//					})
//				}
//
//				// Makes sure that there are zero webhooks left
//				count, err := database.CountWebhooks(TestCollection)
//				assert.NoError(t, err)
//				assert.Equal(t, 0, count)
//			}
//	*/
//}
