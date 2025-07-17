# Stripe Service API - Clean & Simple

A clean, simple Golang service that provides RESTful API endpoints for Stripe payment processing.

## 🏗️ Architecture

This service follows a simple, clean architecture:

```
stripe-service/
├── main.go                    # Application entry point
├── config/                    # Configuration management
├── internal/
│   ├── models/               # Data models and request/response types
│   ├── service/              # Business logic and Stripe integration
│   └── handlers/             # HTTP request handlers
├── Dockerfile                # Container configuration
└── README.md                 # This file
```

## 🚀 Features

- **Customer Management**: Create, retrieve, and list customers
- **Payment Processing**: Create and confirm payment intents
- **Product Catalog**: Create products and prices
- **Subscriptions**: Create and manage subscriptions
- **Input Validation**: Comprehensive request validation
- **Error Handling**: Proper HTTP error responses
- **Health Checks**: Service health monitoring

## 📋 Prerequisites

- Go 1.21 or later
- Stripe account and API keys

## 🔧 Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```

## ⚙️ Configuration

Set your Stripe secret key as an environment variable:

```bash
export STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key_here
export PORT=8080
export HOST=localhost
```

Or copy the example configuration:
```bash
cp config.example .env
```

## 🏃 Running the Service

### Local Development
```bash
go run main.go
```

### Using Docker
```bash
# Build the image
docker build -t stripe-service .

# Run the container
docker run -p 8080:8080 -e STRIPE_SECRET_KEY=your_key stripe-service
```

The service will start on `http://localhost:8080`

## 📡 API Endpoints

### Health Check
- `GET /api/v1/health` - Check service health

### Customer Management
- `POST /api/v1/customers` - Create a new customer
- `GET /api/v1/customers` - List customers (with optional `limit` and `cursor` parameters)
- `GET /api/v1/customers/{id}` - Get customer by ID

### Payment Processing
- `POST /api/v1/payment-intents` - Create a payment intent
- `POST /api/v1/payment-intents/{id}/confirm` - Confirm a payment intent

### Product Management
- `POST /api/v1/products` - Create a product
- `POST /api/v1/prices` - Create a price for a product

### Subscription Management
- `POST /api/v1/subscriptions` - Create a subscription
- `DELETE /api/v1/subscriptions/{id}` - Cancel a subscription

## 📖 Interactive API Documentation

### 🚀 OpenAPI/Swagger Documentation
We provide comprehensive interactive documentation with live testing capabilities:

```bash
# Generate and view documentation
make docs

# Serve documentation locally
make docs-serve
# Then visit: http://localhost:8000/docs/api-documentation.html

# Validate OpenAPI specification
make docs-validate
```

### 📚 Documentation Features
- **Interactive API Explorer** - Test endpoints directly from the documentation
- **Complete Schema Definitions** - All request/response models documented
- **Real Examples** - Realistic curl commands and JSON examples
- **Validation Rules** - Input validation requirements clearly specified
- **Error Handling** - Comprehensive error response documentation

## 📝 API Usage Examples

### Create a Customer
```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "name": "John Doe"
  }'
```

### Create a Payment Intent
```bash
curl -X POST http://localhost:8080/api/v1/payment-intents \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 2000,
    "currency": "usd",
    "customer_id": "cus_customer_id"
  }'
```

### List Customers
```bash
curl "http://localhost:8080/api/v1/customers?limit=10"
```

## 🧪 Testing

Run the tests:
```bash
go test ./...
```

## 🔒 Security

- Store your Stripe secret key securely
- Never expose your secret key in client-side code
- Use HTTPS in production
- Validate all inputs

## 📦 Project Structure

### `/internal/models/`
Contains data models and request/response types:
- `customer.go` - Customer-related models
- `payment.go` - Payment intent models
- `product.go` - Product, price, and subscription models

### `/internal/service/`
Contains business logic:
- `stripe.go` - Stripe API integration and business logic

### `/internal/handlers/`
Contains HTTP handlers:
- `stripe.go` - HTTP request handlers with validation

### `/config/`
Contains configuration management:
- `config.go` - Environment variable handling

## 🚀 Deployment

### Docker
```bash
docker build -t stripe-service .
docker run -p 8080:8080 -e STRIPE_SECRET_KEY=your_key stripe-service
```

### Binary
```bash
go build -o stripe-service .
./stripe-service
```

## 📄 License

This project is provided as-is for educational and development purposes. 