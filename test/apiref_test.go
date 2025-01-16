package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type TestResult struct {
	Method   string
	Endpoint string
	Status   string
	Error    string
}

const baseURL = "http://localhost:8080/api/reflink"

func TestReflinkAPI(t *testing.T) {
	tests := []TestResult{
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
		return sendRequest(http.MethodPost, endpoint, payload)
	}
*/
func testGetLink(workerID int) TestResult {
	endpoint := fmt.Sprintf("%s/%d", baseURL, workerID)
	return sendRequest(http.MethodGet, endpoint, nil)
}

func testUpdateLink(workerID int) TestResult {
	endpoint := fmt.Sprintf("%s/%d", baseURL, workerID)
	payload := map[string]interface{}{
		"tag": "updated_referral",
	}
	return sendRequest(http.MethodPatch, endpoint, payload)
}

func testClickAdd(workerID int) TestResult {
	endpoint := fmt.Sprintf("%s/clicks/%d", baseURL, workerID)
	payload := map[string]interface{}{
		"count": 5,
	}
	return sendRequest(http.MethodPatch, endpoint, payload)
}

func testRegistrationAdd(workerID int) TestResult {
	endpoint := fmt.Sprintf("%s/registrations/%d", baseURL, workerID)
	payload := map[string]interface{}{
		"count": 3,
	}
	return sendRequest(http.MethodPatch, endpoint, payload)
}

func sendRequest(method, endpoint string, payload map[string]interface{}) TestResult {
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return TestResult{
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
		return TestResult{
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
		return TestResult{
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
		return TestResult{
			Method:   method,
			Endpoint: endpoint,
			Status:   status,
			Error:    string(bodyBytes),
		}
	}

	return TestResult{
		Method:   method,
		Endpoint: endpoint,
		Status:   status,
		Error:    "",
	}
}
