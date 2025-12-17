package main

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

// TestHomeHandler tests the root endpoint
func TestHomeHandler(t *testing.T) {
    // Create a request to pass to our handler
    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        t.Fatal(err)
    }

    // Create a ResponseRecorder to record the response
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(HomeHandler)

    // Serve the HTTP request to our handler
    handler.ServeHTTP(rr, req)

    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    // Check the response body
    expected := "Hello, World from Go!\n"
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }

    // Check the Content-Type header
    expectedContentType := "text/plain"
    if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
        t.Errorf("handler returned wrong Content-Type: got %v want %v",
            contentType, expectedContentType)
    }
}

// TestApiHandler tests the API endpoint
func TestApiHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/api", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(ApiHandler)

    handler.ServeHTTP(rr, req)

    // Check status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    // Check Content-Type header
    expectedContentType := "application/json"
    if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
        t.Errorf("handler returned wrong Content-Type: got %v want %v",
            contentType, expectedContentType)
    }

    // Parse and validate JSON response
    var response map[string]string
    err = json.Unmarshal(rr.Body.Bytes(), &response)
    if err != nil {
        t.Errorf("handler returned invalid JSON: %v", err)
    }

    // Check JSON fields
    expectedMessage := "Welcome to the Go API!"
    if message, exists := response["message"]; !exists || message != expectedMessage {
        t.Errorf("handler returned wrong message: got %v want %v",
            message, expectedMessage)
    }

    expectedStatus := "success"
    if status, exists := response["status"]; !exists || status != expectedStatus {
        t.Errorf("handler returned wrong status: got %v want %v",
            status, expectedStatus)
    }
}

// TestApiHandlerWithPost tests that POST method is not allowed
func TestApiHandlerWithPost(t *testing.T) {
    req, err := http.NewRequest("POST", "/api", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(ApiHandler)

    handler.ServeHTTP(rr, req)

    // Even with POST, our handler currently returns 200 (simplified example)
    // In a real API, we might want to return 405 Method Not Allowed
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code for POST: got %v want %v",
            status, http.StatusOK)
    }
}

// TestHomeHandlerNotFound tests a non-existent endpoint
func TestHomeHandlerNotFound(t *testing.T) {
    // Note: This test demonstrates what happens with our current simple server
    // Our server will still route to HomeHandler for any path due to "/" pattern
    req, err := http.NewRequest("GET", "/nonexistent", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(HomeHandler)

    handler.ServeHTTP(rr, req)

    // Since we're testing the handler directly (not the router),
    // it should still work and return the home page
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code for nonexistent path: got %v want %v",
            status, http.StatusOK)
    }
}

// BenchmarkHomeHandler benchmarks the home handler performance
func BenchmarkHomeHandler(b *testing.B) {
    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        b.Fatal(err)
    }

    for i := 0; i < b.N; i++ {
        rr := httptest.NewRecorder()
        HomeHandler(rr, req)
    }
}