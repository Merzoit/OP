package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TestResult struct {
	Method   string
	Endpoint string
	Status   string
	Error    string
}

func SendRequest(method, endpoint string, payload map[string]interface{}) TestResult {
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
