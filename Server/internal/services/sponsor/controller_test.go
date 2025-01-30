package sponsor

import (
	"at/internal/services/reflink"
	"fmt"
	"net/http"
	"testing"
)

const sponsorBaseURL = "http://localhost:8080/api/sponsor"

func TestSponsorAPI(t *testing.T) {
	tests := []reflink.TestResult{
		testCreateSponsor(),
		testGetSponsor(2),
		testGetSponsors(),
		testDeleteSponsor(1),
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

func testCreateSponsor() reflink.TestResult {
	endpoint := fmt.Sprintf("%s", sponsorBaseURL)
	payload := map[string]interface{}{
		"telegram_link": "https://t.me/test",
		"price_per_sub": 10.5,
		"name":          "Test Sponsor",
	}
	return test.SendRequest(http.MethodPost, endpoint, payload)
}

func testGetSponsor(id int) test.TestResult {
	endpoint := fmt.Sprintf("%s/%d", sponsorBaseURL, id)
	return test.SendRequest(http.MethodGet, endpoint, nil)
}

func testGetSponsors() test.TestResult {
	endpoint := fmt.Sprintf("%s/all", sponsorBaseURL)
	return test.SendRequest(http.MethodGet, endpoint, nil)
}

func testDeleteSponsor(id int) test.TestResult {
	endpoint := fmt.Sprintf("%s/%d", sponsorBaseURL, id)
	return test.SendRequest(http.MethodDelete, endpoint, nil)
}
