#!/bin/bash

# Go Beginner Toolkit - Test Runner
# Run all tests with different verbosity levels

set -e  # Exit on error

echo "🔧 Go Beginner Toolkit - Test Suite"
echo "======================================"

# Function to run tests with timing
run_test() {
    local test_name=$1
    local test_cmd=$2
    
    echo ""
    echo "🧪 Running: $test_name"
    echo "--------------------------------------"
    
    # Run with timing
    time {
        if eval "$test_cmd"; then
            echo "✅ $test_name PASSED"
        else
            echo "❌ $test_name FAILED"
            return 1
        fi
    }
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

echo "📊 Go Version: $(go version)"
echo "📁 Current directory: $(pwd)"
echo ""

# 1. Basic syntax check
run_test "Syntax Check" "go fmt ./..."
run_test "Vet Check" "go vet ./..."
run_test "Static Analysis" "go run golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest ./..."

# 2. Run unit tests
run_test "Unit Tests" "go test ./... -v -short"

# 3. Run benchmarks
echo ""
echo "📈 Running Benchmarks..."
go test -bench=. -benchmem ./tests/... 2>/dev/null || echo "Benchmarks not available"

# 4. Test coverage
echo ""
echo "📊 Generating Test Coverage..."
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep -E "(total|main\.go)"

# 5. Build test
echo ""
echo "🔨 Build Test..."
run_test "Build Check" "go build -o /tmp/test_build main.go && rm /tmp/test_build"

# 6. Run the server in background and test endpoints
echo ""
echo "🌐 Server Integration Test..."
timeout 5s go run main.go &
SERVER_PID=$!
sleep 2  # Wait for server to start

# Test endpoints
echo "Testing HTTP endpoints..."
if curl -s http://localhost:8080/ | grep -q "Hello, World"; then
    echo "✅ Root endpoint: PASSED"
else
    echo "❌ Root endpoint: FAILED"
fi

if curl -s http://localhost:8080/api | grep -q "Welcome to the Go API"; then
    echo "✅ API endpoint: PASSED"
else
    echo "❌ API endpoint: FAILED"
fi

# Kill the server
kill $SERVER_PID 2>/dev/null || true

echo ""
echo "======================================"
echo "🎉 All tests completed!"
echo ""
echo "📚 Next steps:"
echo "   go run main.go          # Run the server"
echo "   go test ./... -v        # Run tests with details"
echo "   go tool cover -html=coverage.out  # View coverage in browser"