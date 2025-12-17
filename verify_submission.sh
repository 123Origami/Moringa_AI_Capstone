#!/bin/bash

# ======================================================
# Go Beginner Toolkit - Submission Verification Script
# Moringa AI Capstone Project
# Run this before submitting to ensure everything is ready
# ======================================================

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Banner
echo -e "${CYAN}"
echo "╔══════════════════════════════════════════════════════════╗"
echo "║    Go Beginner Toolkit - Submission Verifier             ║"
echo "║    Moringa AI Capstone Project                          ║"
echo "╚══════════════════════════════════════════════════════════╝"
echo -e "${NC}"
echo ""

# Function to print section headers
section() {
    echo ""
    echo -e "${BLUE}══════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE} $1 ${NC}"
    echo -e "${BLUE}══════════════════════════════════════════════════════════${NC}"
}

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "  ${GREEN}✅ $2${NC}"
    elif [ $1 -eq 2 ]; then
        echo -e "  ${YELLOW}⚠️  $2${NC}"
    else
        echo -e "  ${RED}❌ $2${NC}"
    fi
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check file exists and has content
file_check() {
    if [ ! -f "$1" ]; then
        echo 1  # File doesn't exist
    elif [ ! -s "$1" ]; then
        echo 2  # File exists but is empty
    else
        echo 0  # File exists and has content
    fi
}

# Start verification
section "📋 SYSTEM CHECK"

# Check Go installation
if command_exists go; then
    print_status 0 "Go installed: $(go version)"
else
    print_status 1 "Go is NOT installed"
    exit 1
fi

# Check minimum Go version
GO_VERSION=$(go version | grep -oE '[0-9]+\.[0-9]+')
REQUIRED_VERSION=1.21
if awk -v v1=$GO_VERSION -v v2=$REQUIRED_VERSION 'BEGIN {exit (v1 < v2)}'; then
    print_status 0 "Go version $GO_VERSION >= required $REQUIRED_VERSION"
else
    print_status 1 "Go version $GO_VERSION < required $REQUIRED_VERSION"
fi

section "📁 REQUIRED FILES CHECK"

# List of required files
required_files=(
    "main.go"
    "go.mod"
    "README.md"
    "toolkit-guide.md"
    "LICENSE"
    ".gitignore"
)

all_files_ok=true
for file in "${required_files[@]}"; do
    status=$(file_check "$file")
    case $status in
        0) print_status 0 "$file - exists with content" ;;
        1) print_status 1 "$file - MISSING"; all_files_ok=false ;;
        2) print_status 2 "$file - exists but EMPTY"; all_files_ok=false ;;
    esac
done

if $all_files_ok; then
    print_status 0 "All required files present"
else
    echo -e "${YELLOW}⚠️  Some files are missing or empty${NC}"
fi

section "🔧 GO MODULE CHECK"

# Check go.mod
if [ -f "go.mod" ]; then
    module_name=$(grep '^module' go.mod | head -1 | cut -d' ' -f2)
    if [ -n "$module_name" ]; then
        print_status 0 "Module: $module_name"
    else
        print_status 1 "No module name found in go.mod"
    fi
fi

section "⚙️ BUILD & COMPILATION CHECK"

# Try to build
build_output=$(go build -o /tmp/go_toolkit_build_test main.go 2>&1)
if [ $? -eq 0 ]; then
    print_status 0 "Build successful"
    rm -f /tmp/go_toolkit_build_test
else
    print_status 1 "Build failed"
    echo -e "${YELLOW}Build error output:${NC}"
    echo "$build_output"
fi

# Check for syntax errors
fmt_output=$(go fmt main.go 2>&1)
if [ $? -eq 0 ]; then
    print_status 0 "Code formatting check passed"
else
    print_status 1 "Code formatting issues"
    echo "$fmt_output"
fi

section "🧪 FUNCTIONALITY TEST"

# Test server in background
echo -e "  Testing server startup..."
timeout 5s go run main.go > /tmp/server_test.log 2>&1 &
SERVER_PID=$!
sleep 3

# Test endpoints
echo -e "  Testing endpoints..."
endpoints=("/" "/api" "/health")
all_endpoints_ok=true

for endpoint in "${endpoints[@]}"; do
    if curl -s -f "http://localhost:8080$endpoint" > /dev/null 2>&1; then
        print_status 0 "Endpoint $endpoint - reachable"
    else
        print_status 1 "Endpoint $endpoint - NOT reachable"
        all_endpoints_ok=false
    fi
done

# Test JSON response
api_response=$(curl -s "http://localhost:8080/api")
if echo "$api_response" | python3 -m json.tool > /dev/null 2>&1; then
    print_status 0 "API returns valid JSON"
else
    print_status 1 "API does not return valid JSON"
    all_endpoints_ok=false
fi

# Kill test server
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null

section "📚 DOCUMENTATION CHECK"

# Check README for required sections
required_sections=("Overview" "Installation" "AI Prompt" "Common Errors" "References")
readme_content=$(cat README.md 2>/dev/null || echo "")

all_sections_found=true
for section in "${required_sections[@]}"; do
    if echo "$readme_content" | grep -qi "$section"; then
        print_status 0 "README has '$section' section"
    else
        print_status 1 "README missing '$section' section"
        all_sections_found=false
    fi
done

# Check toolkit guide for AI prompts
if grep -q "AI Prompt" toolkit-guide.md 2>/dev/null || grep -q "GenAI" toolkit-guide.md 2>/dev/null; then
    print_status 0 "AI usage documented in toolkit guide"
else
    print_status 1 "AI documentation not found in toolkit guide"
fi

section "🎯 CAPSTONE REQUIREMENTS CHECK"

# Check for AI prompt documentation
echo -e "  Checking for capstone requirements..."
requirements_met=0
requirements_total=5

# 1. Working example
if $all_endpoints_ok; then
    ((requirements_met++))
    print_status 0 "Working example verified"
else
    print_status 1 "Working example incomplete"
fi

# 2. AI prompt documentation
if [ -f "toolkit-guide.md" ] && grep -q -i "prompt" toolkit-guide.md; then
    ((requirements_met++))
    print_status 0 "AI prompts documented"
else
    print_status 1 "AI prompts not documented"
fi

# 3. Setup instructions
if echo "$readme_content" | grep -q -i "install\|setup\|prerequisite"; then
    ((requirements_met++))
    print_status 0 "Setup instructions present"
else
    print_status 1 "Setup instructions missing"
fi

# 4. Common errors
if echo "$readme_content" | grep -q -i "error\|issue\|troubleshoot\|fix"; then
    ((requirements_met++))
    print_status 0 "Common errors documented"
else
    print_status 1 "Common errors not documented"
fi

# 5. References
if echo "$readme_content" | grep -q -i "reference\|resource\|link\|doc"; then
    ((requirements_met++))
    print_status 0 "References included"
else
    print_status 1 "References missing"
fi

section "📊 FINAL VERIFICATION REPORT"

# Calculate score
score=$((requirements_met * 100 / requirements_total))

echo -e "  ${CYAN}Capstone Requirements Met:${NC} $requirements_met/$requirements_total"
echo -e "  ${CYAN}Verification Score:${NC} $score%"

echo ""
echo -e "${BLUE}══════════════════════════════════════════════════════════${NC}"

if [ $requirements_met -eq $requirements_total ] && $all_files_ok && $all_endpoints_ok; then
    echo -e "${GREEN}🎉 ALL CHECKS PASSED! Your project is ready for submission.${NC}"
    echo -e "${GREEN}✅ You have met all capstone requirements.${NC}"
    
    # Show next steps
    echo ""
    echo -e "${CYAN}📝 Next steps for submission:${NC}"
    echo "  1. Run: chmod +x verify_submission.sh (if not already done)"
    echo "  2. Review: SUBMISSION_CHECKLIST.md"
    echo "  3. Update placeholder text in main.go (your name)"
    echo "  4. Push to GitHub repository"
    echo "  5. Submit according to instructor's directions"
    
    exit 0
else
    echo -e "${YELLOW}⚠️  SOME ISSUES FOUND${NC}"
    echo ""
    echo -e "${YELLOW}🔧 Issues to fix before submission:${NC}"
    
    if ! $all_files_ok; then
        echo "  • Check all required files exist and have content"
    fi
    
    if ! $all_endpoints_ok; then
        echo "  • Fix server endpoints that are not working"
    fi
    
    if [ $requirements_met -lt $requirements_total ]; then
        echo "  • Address missing capstone requirements"
    fi
    
    echo ""
    echo -e "${YELLOW}Run these commands to troubleshoot:${NC}"
    echo "  go run main.go        # Test server manually"
    echo "  go build              # Check for build errors"
    echo "  go vet ./...          # Check for code issues"
    
    exit 1
fi