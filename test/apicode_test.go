package test

import (
	"fmt"
	"net/http"
	"testing"
)

const codeBaseURL = "http://localhost:8080/api/code"

func TestCodeAPI(t *testing.T) {
	// Список тестов
	tests := []TestResult{
		testCreateCode(),
		testGetCode(12345),
		testGetCodesByWorker(2),
		testDeleteCode(12345),
		testAddRequestCount(111111),
	}

	// Проверка результатов
	for _, result := range tests {
		if result.Status != "SUCCESS" {
			t.Errorf("FAIL: Method: %s, Endpoint: %s, Error: %s",
				result.Method, result.Endpoint, result.Error)
		} else {
			t.Logf("PASS: Method: %s, Endpoint: %s", result.Method, result.Endpoint)
		}
	}
}

func testCreateCode() TestResult {
	endpoint := fmt.Sprintf("%s", codeBaseURL)
	payload := map[string]interface{}{
		"access_code":        12345,
		"title":              "Test Code",
		"year":               2025,
		"description":        "This is a test description",
		"added_by_worker_id": 2,
		"request_count":      0,
	}
	return sendRequest(http.MethodPost, endpoint, payload)
}

func testGetCode(accessCode int) TestResult {
	endpoint := fmt.Sprintf("%s/%d", codeBaseURL, accessCode)
	return sendRequest(http.MethodGet, endpoint, nil)
}

func testAddRequestCount(accessCode int) TestResult {
	endpoint := fmt.Sprintf("%s/increment/%d", codeBaseURL, accessCode)
	payload := map[string]interface{}{
		"increment": 5,
	}
	return sendRequest(http.MethodPatch, endpoint, payload)
}

func testGetCodesByWorker(workerID int) TestResult {
	endpoint := fmt.Sprintf("%s/worker/%d", codeBaseURL, workerID)
	return sendRequest(http.MethodGet, endpoint, nil)
}

func testDeleteCode(accessCode int) TestResult {
	endpoint := fmt.Sprintf("%s/%d", codeBaseURL, accessCode)
	return sendRequest(http.MethodDelete, endpoint, nil)
}
