package examples

import (
    "fmt"
    "testing"
)

// Example test showing basic Go testing patterns
func TestAddition(t *testing.T) {
    result := 2 + 2
    expected := 4
    if result != expected {
        t.Errorf("Addition failed: got %d, want %d", result, expected)
    }
}

func TestStringConcatenation(t *testing.T) {
    result := "Hello" + " " + "Go"
    expected := "Hello Go"
    if result != expected {
        t.Errorf("String concatenation failed: got %s, want %s", result, expected)
    }
}

// Table-driven test example
func TestTableDrivenAddition(t *testing.T) {
    testCases := []struct {
        name     string
        a        int
        b        int
        expected int
    }{
        {"Positive numbers", 2, 3, 5},
        {"Zero", 0, 5, 5},
        {"Negative numbers", -2, -3, -5},
        {"Mixed", -2, 5, 3},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := tc.a + tc.b
            if result != tc.expected {
                t.Errorf("%s: %d + %d = %d, want %d",
                    tc.name, tc.a, tc.b, result, tc.expected)
            }
        })
    }
}

// Example showing subtests
func TestSubTests(t *testing.T) {
    t.Run("Addition", func(t *testing.T) {
        if 1+2 != 3 {
            t.Error("1+2 should be 3")
        }
    })

    t.Run("Multiplication", func(t *testing.T) {
        if 2*3 != 6 {
            t.Error("2*3 should be 6")
        }
    })
}

// Helper function example
func add(a, b int) int {
    return a + b
}

func TestWithHelper(t *testing.T) {
    result := add(10, 20)
    if result != 30 {
        t.Errorf("add(10, 20) = %d; want 30", result)
    }
}