package handlers

import (
	"assignment-2/internal"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

const TestCollection = "testWebhooks"

func TestDeleteWebhook(t *testing.T) {

	hook := internal.Webhook{
		Url:     "https://localhost:8080/dashboards/v1/notifications/",
		Country: "NO",
		Event:   "REGISTER",
	}

	// marshal a webhook registration
	body, err := json.Marshal(hook)

	// register the webhook
	// create a request
	req := httptest.NewRequest(http.MethodPost, internal.NotificationsPath, bytes.NewBuffer(body))
	req.Header.Add("content-type", "application/json")

	recorder := httptest.NewRecorder()

	// call the handler function
	WebhookRegistration(recorder, req, collectionNameWebhooks)

	req.Header.Add(internal.ApplicationJson, internal.ContentTypeJson)

	// create test recorder
	rec := httptest.NewRecorder()
	// calls the handler function
	DeleteWebhook(rec, req, TestCollection, hook.WebhookId)

	// Create a ResponseRecorder to record the response
	recorder = httptest.NewRecorder()
	// Check the response body
	expectedResponseBody := http.StatusBadRequest

	intToUse, err := strconv.Atoi(recorder.Body.String())
	if err != nil {
		log.Fatal("Error error pants on fire")
	}

	if intToUse != expectedResponseBody {
		t.Errorf("handler returned unexpected body: got %v, want %v", recorder.Body.String(), expectedResponseBody)
	}

	/*
			// calls the handler function
			WebhookRegistration(rec, request, TestCollection)

			var hook internal.Webhook

			err = json.Unmarshal(rec.Body.Bytes(), &hook)
			if err != nil {
				t.Fatal("Error unmarshalling response : " + err.Error())
			}
			var tests []internal.Webhook

			tests[0] = hook
			hook.WebhookId = ""
			tests[1] = hook

			// iterate through the test cases
			for _, testCase := range tests {
				t.Run(testCase.WebhookId, func(t *testing.T) {

					req := httptest.NewRequest(http.MethodDelete, internal.NotificationsPath+testCase.WebhookId, nil)
					req.Header.Add(internal.ApplicationJson, internal.ContentTypeJson)

					// create the recorder
					rec := httptest.NewRecorder()

					// calls the handler function
					DeleteWebhook(rec, req, TestCollection, testCase.WebhookId)

					requestDeletedWebhook := httptest.NewRequest(http.MethodGet, internal.NotificationsPath+testCase.WebhookId, nil)
					req.Header.Add(internal.ApplicationJson, internal.ContentTypeJson)

					// Asserts what is expected
					asserts := assert.New(t)
					asserts.Equal(requestDeletedWebhook.Response, rec.Code)
					asserts.Equal(requestDeletedWebhook.Header, rec.Header().Get("content-type"))
				})
			}

			// Makes sure that there are zero webhooks left
			count, err := database.CountWebhooks(TestCollection)
			assert.NoError(t, err)
			assert.Equal(t, 0, count)
		}
	*/
}
