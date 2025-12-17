// Go Beginner Toolkit - HTTP Server Implementation
// Created as part of Moringa AI Capstone Project
// Date: December 2024

package main

// Import necessary packages from the standard library.
import (
    "encoding/json" // For encoding data to JSON.
    "fmt"           // For formatted I/O (printing).
    "log"           // For logging errors.
    "net/http"      // For building HTTP servers and clients.
    "time"          // For handling time (used in our goroutine sleep).
)

// HomeHandler handles requests to the root route ("/").
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    // Set the Content-Type header to plain text.
    w.Header().Set("Content-Type", "text/plain")
    
    // Log the incoming request
    log.Printf("📝 HomeHandler called: %s %s", r.Method, r.URL.Path)
    
    // Write the HTTP status code 200 (OK) and our message.
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Hello, World from Go!\n")
    fmt.Fprintf(w, "Server time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
    fmt.Fprintf(w, "Visit /api for JSON response\n")
}

// ApiHandler handles requests to the "/api" route.
func ApiHandler(w http.ResponseWriter, r *http.Request) {
    // Log the API request
    log.Printf("📡 API request: %s %s", r.Method, r.URL.Path)
    
    // Define a simple struct (like an object) for our JSON response.
    type Message struct {
        Text       string    `json:"message"`   // Struct tag tells the encoder to use "message" as the JSON key.
        Status     string    `json:"status"`
        Timestamp  string    `json:"timestamp"`
        Endpoint   string    `json:"endpoint"`
        Version    string    `json:"version"`
    }

    // Create an instance of the Message struct.
    response := Message{
        Text:      "Welcome to the Go API!",
        Status:    "success",
        Timestamp: time.Now().Format(time.RFC3339),
        Endpoint:  "/api",
        Version:   "1.0.0",
    }

    // Set the Content-Type header to application/json.
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    // Encode the struct to JSON and write it to the response.
    // If encoding fails, log the error.
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("❌ Error encoding JSON: %v", err)
        http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
    }
}

// HealthHandler provides a health check endpoint
func HealthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    
    healthStatus := map[string]interface{}{
        "status":    "healthy",
        "service":   "go-beginner-toolkit",
        "timestamp": time.Now().Unix(),
        "uptime":    time.Since(startTime).String(),
    }
    
    json.NewEncoder(w).Encode(healthStatus)
}

// NotFoundHandler handles 404 responses
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusNotFound)
    
    errorResponse := map[string]string{
        "error":   "Endpoint not found",
        "path":    r.URL.Path,
        "message": "Try / or /api endpoints",
    }
    
    json.NewEncoder(w).Encode(errorResponse)
    log.Printf("❌ 404 Not Found: %s", r.URL.Path)
}

// Global variable to track server start time
var startTime time.Time

func main() {
    // Record server start time
    startTime = time.Now()
    
    fmt.Println("🚀 Go Beginner Toolkit Server")
    fmt.Println("=============================")
    fmt.Println("Project: Moringa AI Capstone")
    fmt.Println("Technology: Go (Golang)")
    fmt.Println("Author: [Your Name]")
    fmt.Println("=============================")
    
    // ---- Goroutine Example: Concurrent Logging ----
    // The 'go' keyword starts this function in a new goroutine (lightweight thread).
    go func() {
        // Wait a moment for the server to start announcing.
        time.Sleep(100 * time.Millisecond)
        fmt.Println("\n📢 Server Information:")
        fmt.Println("   Port: 8080")
        fmt.Println("   Endpoints:")
        fmt.Println("     • GET /         - Welcome message")
        fmt.Println("     • GET /api      - JSON API response")
        fmt.Println("     • GET /health   - Health check")
        fmt.Println("   Press Ctrl+C to stop the server")
        fmt.Println("==========================================")
    }()
    // The main function continues immediately without waiting for the goroutine.

    // ---- Routing ----
    // Register our handler functions to specific URL patterns.
    http.HandleFunc("/", HomeHandler)
    http.HandleFunc("/api", ApiHandler)
    http.HandleFunc("/health", HealthHandler)
    
    // Custom 404 handler
    http.HandleFunc("/404", NotFoundHandler)
    
    // ---- Server Configuration ----
    port := ":8080"
    server := &http.Server{
        Addr:         port,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  30 * time.Second,
    }
    
    fmt.Printf("⏳ Starting server on http://localhost%s\n", port)
    fmt.Println("🔄 Initializing...")

    // ---- Start the Server ----
    // log.Fatal is used because ListenAndServe only returns an error.
    // If the server starts successfully, it blocks here indefinitely.
    log.Printf("✅ Server started at %s", time.Now().Format("2006-01-02 15:04:05"))
    log.Printf("🌐 Listening on http://localhost%s", port)
    
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("❌ Failed to start server: %v", err)
    }
}