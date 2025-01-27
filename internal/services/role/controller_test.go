package role

import (
	"at/internal/services/reflink"
	"fmt"
	"net/http"
	"testing"
)

const roleBaseURL = "http://localhost:8080/api/role"

func TestRoleAPI(t *testing.T) {
	tests := []reflink.TestResult{
		testGetRole(1),
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

func testGetRole(roleID int) test.TestResult {
	endpoint := fmt.Sprintf("%s/%d", roleBaseURL, roleID)
	return test.SendRequest(http.MethodGet, endpoint, nil)
}
