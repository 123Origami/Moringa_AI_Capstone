package examples

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

// Mock handler for example tests
func ExampleHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status":  "ok",
        "message": "Example response",
    })
}

func TestExampleHandler(t *testing.T) {
    // Test GET request
    req := httptest.NewRequest("GET", "/example", nil)
    rr := httptest.NewRecorder()
    ExampleHandler(rr, req)

    if rr.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", rr.Code)
    }

    // Test JSON response
    var response map[string]string
    json.Unmarshal(rr.Body.Bytes(), &response)

    if response["status"] != "ok" {
        t.Errorf("Expected status 'ok', got %s", response["status"])
    }
}

func TestPostRequestWithBody(t *testing.T) {
    // Create JSON request body
    requestBody := map[string]string{
        "name":  "Test User",
        "email": "test@example.com",
    }
    jsonBody, _ := json.Marshal(requestBody)

    req := httptest.NewRequest("POST", "/api/users",
        bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()

    // For this example, we'll use a simple echo handler
    http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Just echo back what we received
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        w.Write(jsonBody)
    }).ServeHTTP(rr, req)

    if rr.Code != http.StatusCreated {
        t.Errorf("Expected status 201, got %d", rr.Code)
    }

    // Verify the echoed response
    var echoedResponse map[string]string
    json.Unmarshal(rr.Body.Bytes(), &echoedResponse)

    if echoedResponse["name"] != "Test User" {
        t.Errorf("Expected name 'Test User', got %s", echoedResponse["name"])
    }
}

// Test for proper error handling
func TestErrorResponse(t *testing.T) {
    // This handler returns an error status
    errorHandler := func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Invalid input",
        })
    }

    req := httptest.NewRequest("GET", "/error", nil)
    rr := httptest.NewRecorder()
    errorHandler(rr, req)

    if rr.Code != http.StatusBadRequest {
        t.Errorf("Expected status 400, got %d", rr.Code)
    }

    var errorResponse map[string]string
    json.Unmarshal(rr.Body.Bytes(), &errorResponse)

    if errorResponse["error"] != "Invalid input" {
        t.Errorf("Expected error 'Invalid input', got %s", errorResponse["error"])
    }
}

// Example of testing with query parameters
func TestQueryParameters(t *testing.T) {
    handler := func(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query()
        name := query.Get("name")
        age := query.Get("age")

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "name": name,
            "age":  age,
        })
    }

    req := httptest.NewRequest("GET", "/query?name=John&age=30", nil)
    rr := httptest.NewRecorder()
    handler(rr, req)

    var response map[string]string
    json.Unmarshal(rr.Body.Bytes(), &response)

    if response["name"] != "John" {
        t.Errorf("Expected name 'John', got %s", response["name"])
    }
    if response["age"] != "30" {
        t.Errorf("Expected age '30', got %s", response["age"])
    }
}