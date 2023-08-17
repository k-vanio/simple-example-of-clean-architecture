package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k-vanio/simple-example-of-clean-architecture/pkg/response"
)

func TestJsonFunction(t *testing.T) {
	t.Run("Test Successful JSON Response", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		data := struct {
			Message string `json:"response"`
		}{Message: "success"}

		response.Json(recorder, http.StatusOK, data)

		response := recorder.Result()
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, response.StatusCode)
		}

		var responseBody map[string]string
		if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			t.Errorf("Error decoding response body: %v", err)
		}

		expectedMessage := "success"
		if responseBody["response"] != expectedMessage {
			t.Errorf("Expected response message %s, but got %s", expectedMessage, responseBody["response"])
		}
	})

	t.Run("Test Internal Server Error JSON Response", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		// Passing an invalid body to simulate an error during encoding
		response.Json(recorder, http.StatusInternalServerError, make(chan int))

		response := recorder.Result()
		defer response.Body.Close()

		if response.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, but got %d", http.StatusInternalServerError, response.StatusCode)
		}

		var responseBody map[string]string
		if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			t.Errorf("Error decoding response body: %v", err)
		}

		expectedMessage := "internal server error"
		if responseBody["response"] != expectedMessage {
			t.Errorf("Expected response message %s, but got %s", expectedMessage, responseBody["response"])
		}
	})
}
