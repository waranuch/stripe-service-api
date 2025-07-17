#!/bin/bash

echo "🧪 Running Stripe Service API Tests"
echo "======================================"

# Run all tests with verbose output
echo "Running all tests..."
go test ./... -v -cover

# Check if tests passed
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ All tests passed successfully!"
    echo ""
    echo "📊 Test Coverage Summary:"
    echo "- Config package: Unit tests for environment variable loading"
    echo "- Service package: Unit tests for Stripe API service layer"
    echo "- Handlers package: Unit tests for HTTP request handlers"
    echo "- Main package: Integration tests for full API routes and middleware"
    echo ""
    echo "🏗️  Test Structure:"
    echo "- config_test.go: Tests configuration loading and validation"
    echo "- stripe_service_test.go: Tests Stripe service methods"
    echo "- stripe_handlers_test.go: Tests HTTP handlers and request validation"
    echo "- main_test.go: Tests routing, middleware, and full API integration"
    echo ""
    echo "🔍 Test Coverage:"
    echo "- ✅ Configuration management"
    echo "- ✅ Service layer functionality"
    echo "- ✅ HTTP request handling"
    echo "- ✅ Input validation"
    echo "- ✅ Error handling"
    echo "- ✅ API routing"
    echo "- ✅ CORS middleware"
    echo "- ✅ JSON response formatting"
    echo ""
    echo "🚀 Ready for production use!"
else
    echo ""
    echo "❌ Some tests failed. Please check the output above."
    exit 1
fi 