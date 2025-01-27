package reflink

import (
	"at/test"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8080/api/reflink"

func TestReflinkAPI(t *testing.T) {
	tests := []test.TestResult{
		//testCreateLink(),
		testGetLink(2),
		testUpdateLink(2),
		testClickAdd(2),
		testRegistrationAdd(2),
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

/*
	func testCreateLink() TestResult {
		endpoint := fmt.Sprintf("%s", baseURL)
		payload := map[string]interface{}{
			"worker_id":     6,
			"referral_link": "test_ref",
		}
		return SendRequest(http.MethodPost, endpoint, payload)
	}
*/
func testGetLink(workerID int) test.TestResult {
	endpoint := fmt.Sprintf("%s/%d", baseURL, workerID)
	return SendRequest(http.MethodGet, endpoint, nil)
}

func testUpdateLink(workerID int) test.TestResult {
	endpoint := fmt.Sprintf("%s/%d", baseURL, workerID)
	payload := map[string]interface{}{
		"tag": "updated_referral",
	}
	return SendRequest(http.MethodPatch, endpoint, payload)
}

func testClickAdd(workerID int) test.TestResult {
	endpoint := fmt.Sprintf("%s/clicks/%d", baseURL, workerID)
	payload := map[string]interface{}{
		"count": 5,
	}
	return SendRequest(http.MethodPatch, endpoint, payload)
}

func testRegistrationAdd(workerID int) test.TestResult {
	endpoint := fmt.Sprintf("%s/registrations/%d", baseURL, workerID)
	payload := map[string]interface{}{
		"count": 3,
	}
	return SendRequest(http.MethodPatch, endpoint, payload)
}

func SendRequest(method, endpoint string, payload map[string]interface{}) test.TestResult {
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return test.TestResult{
				Method:   method,
				Endpoint: endpoint,
				Status:   "FAILED",
				Error:    fmt.Sprintf("Error encoding payload: %v", err),
			}
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return test.TestResult{
			Method:   method,
			Endpoint: endpoint,
			Status:   "FAILED",
			Error:    fmt.Sprintf("Error creating request: %v", err),
		}
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return test.TestResult{
			Method:   method,
			Endpoint: endpoint,
			Status:   "FAILED",
			Error:    fmt.Sprintf("Request error: %v", err),
		}
	}
	defer resp.Body.Close()

	status := "SUCCESS"
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		status = "FAILED"
		bodyBytes, _ := io.ReadAll(resp.Body)
		return test.TestResult{
			Method:   method,
			Endpoint: endpoint,
			Status:   status,
			Error:    string(bodyBytes),
		}
	}

	return test.TestResult{
		Method:   method,
		Endpoint: endpoint,
		Status:   status,
		Error:    "",
	}
}
