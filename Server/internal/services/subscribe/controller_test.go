package subscribe

import (
	"at/test"
	"fmt"
	"net/http"
	"testing"
)

const subscribeBaseURL = "http://localhost:8080/api/subscribe"

func TestSubscribeAPI(t *testing.T) {
	tests := []test.TestResult{
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

func testAddSubscribe() test.TestResult {
	endpoint := fmt.Sprintf("%s/create", subscribeBaseURL)
	payload := map[string]interface{}{
		"user_id":    19,
		"sponsor_id": 2,
	}
	return test.SendRequest(http.MethodPost, endpoint, payload)
}

func testGetSubscribesByUser(userID int) test.TestResult {
	endpoint := fmt.Sprintf("%s/user/%d", subscribeBaseURL, userID)
	return test.SendRequest(http.MethodGet, endpoint, nil)
}

func testGetSubscribesBySponsor(sponsorID int) test.TestResult {
	endpoint := fmt.Sprintf("%s/sponsor/%d", subscribeBaseURL, sponsorID)
	return test.SendRequest(http.MethodGet, endpoint, nil)
}
