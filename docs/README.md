# 📚 API Documentation

This directory contains the comprehensive API documentation for the Stripe Service API.

## 📋 Files

- **`api-documentation.html`** - Interactive HTML documentation with Swagger UI
- **`../openapi.yaml`** - OpenAPI 3.0 specification file
- **`README.md`** - This file

## 🚀 Quick Start

### View Documentation

1. **Open HTML Documentation**:
   ```bash
   # Option 1: Open directly in browser
   open docs/api-documentation.html
   
   # Option 2: Serve locally with Python
   make docs-serve
   # Then visit: http://localhost:8000/docs/api-documentation.html
   ```

2. **View OpenAPI Specification**:
   ```bash
   # View the raw OpenAPI YAML
   cat openapi.yaml
   
   # Or use any OpenAPI viewer/editor
   # - Swagger Editor: https://editor.swagger.io/
   # - Insomnia: Import openapi.yaml
   # - Postman: Import openapi.yaml
   ```

### Generate Documentation

```bash
# Generate documentation files
make docs

# Serve documentation locally
make docs-serve
```

## 📖 Documentation Features

### Interactive API Explorer
- **Try It Out**: Test API endpoints directly from the documentation
- **Real Examples**: All endpoints include realistic request/response examples
- **Schema Validation**: Request/response schemas with validation rules
- **Error Handling**: Comprehensive error response documentation

### Comprehensive Coverage
- **10 API Endpoints** - All service endpoints documented
- **Data Models** - Complete schema definitions for all request/response types
- **Authentication** - API key authentication documentation
- **Error Responses** - Standard error formats and codes

### Developer-Friendly
- **Quick Start Guide** - Get up and running quickly
- **Code Examples** - Ready-to-use curl commands
- **Testing Scripts** - Links to automated testing tools
- **Resource Links** - Related documentation and guides

## 🔧 API Endpoints Overview

### Health Check
- `GET /api/v1/health` - Service health status

### Customer Management
- `POST /api/v1/customers` - Create customer
- `GET /api/v1/customers` - List customers
- `GET /api/v1/customers/{id}` - Get customer

### Payment Processing
- `POST /api/v1/payment-intents` - Create payment intent
- `POST /api/v1/payment-intents/{id}/confirm` - Confirm payment

### Product Catalog
- `POST /api/v1/products` - Create product
- `POST /api/v1/prices` - Create price

### Subscription Management
- `POST /api/v1/subscriptions` - Create subscription
- `DELETE /api/v1/subscriptions/{id}` - Cancel subscription

## 🛠️ Using the Documentation

### For Developers
1. **Start Here**: Open `api-documentation.html` in your browser
2. **Explore Endpoints**: Use the interactive interface to understand each endpoint
3. **Test APIs**: Use the "Try it out" feature to test endpoints
4. **Copy Examples**: Use the provided curl examples in your code

### For Integration
1. **Import OpenAPI**: Import `openapi.yaml` into your API client (Postman, Insomnia, etc.)
2. **Generate Code**: Use OpenAPI generators to create client libraries
3. **Validate Requests**: Use the schema definitions for request validation

### For Testing
1. **Manual Testing**: Use the interactive documentation to test endpoints
2. **Automated Testing**: Use the provided testing scripts
3. **Load Testing**: Use the API specification for load testing tools

## 📊 Documentation Statistics

- **OpenAPI Version**: 3.0.3
- **Total Endpoints**: 10
- **Request Models**: 6
- **Response Models**: 7
- **Error Responses**: 3 standard types
- **Code Examples**: 20+ curl examples
- **Test Coverage**: 80.5% (service code)

## 🔗 Related Resources

- [Project README](../README.md) - Main project documentation
- [Testing Scripts](../scripts/README.md) - API testing tools
- [Stripe API Docs](https://stripe.com/docs/api) - Official Stripe documentation
- [OpenAPI Specification](https://swagger.io/specification/) - OpenAPI standard

## 🆘 Troubleshooting

### Documentation Not Loading
1. Ensure you're serving the files over HTTP (not file://)
2. Check that `openapi.yaml` is in the correct location
3. Verify internet connection for Swagger UI CDN resources

### API Testing Issues
1. Ensure the service is running: `make start-dev`
2. Check the health endpoint: `curl localhost:8080/api/v1/health`
3. Verify your Stripe secret key is configured correctly

### OpenAPI Import Issues
1. Validate the OpenAPI spec: Use https://editor.swagger.io/
2. Check for any YAML syntax errors
3. Ensure your API client supports OpenAPI 3.0

## 🚀 Next Steps

1. **Explore the API**: Start with the interactive documentation
2. **Test Endpoints**: Use the testing scripts to create sample data
3. **Integrate**: Import the OpenAPI spec into your development tools
4. **Contribute**: Help improve the documentation by reporting issues or suggestions

Happy coding! 🎉 