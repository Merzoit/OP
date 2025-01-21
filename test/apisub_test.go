package test

import (
	"fmt"
	"net/http"
	"testing"
)

const subscribeBaseURL = "http://localhost:8080/api/subscribe"

func TestSubscribeAPI(t *testing.T) {
	tests := []TestResult{
		testAddSubscribe(),
		testGetSubscribesByUser(2),
		testGetSubscribesBySponsor(1),
	}

	for _, result := range tests {
		if result.Status != "SUCCESS" {
			t.Errorf("FAIL: Method: %s, Endpoint: %s, Error: %s",
				result.Method, result.Endpoint, result.Error)
		} else {
			t.Logf("PASS: Method: %s, Endpoint: %s", result.Method, result.Endpoint)
		}
	}
}

func testAddSubscribe() TestResult {
	endpoint := fmt.Sprintf("%s", subscribeBaseURL)
	payload := map[string]interface{}{
		"user_id":    1,
		"sponsor_id": 2,
	}
	return sendRequest(http.MethodPost, endpoint, payload)
}

func testGetSubscribesByUser(userID int) TestResult {
	endpoint := fmt.Sprintf("%s/user/%d", subscribeBaseURL, userID)
	return sendRequest(http.MethodGet, endpoint, nil)
}

func testGetSubscribesBySponsor(sponsorID int) TestResult {
	endpoint := fmt.Sprintf("%s/sponsor/%d", subscribeBaseURL, sponsorID)
	return sendRequest(http.MethodGet, endpoint, nil)
}
