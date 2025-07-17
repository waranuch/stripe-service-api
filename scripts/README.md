# Scripts Directory

This directory contains scripts for testing and managing the Stripe Service API.

## Available Scripts

### 1. `create_test_data.go`
**Purpose**: Creates comprehensive test data by calling the service API endpoints

**Usage**:
```bash
# Make sure the service is running first
go run main.go

# In another terminal, run the script
go run scripts/create_test_data.go
```

**Or use the Makefile**:
```bash
make test-data
```

**What it does**:
- ✅ Creates a test customer
- ✅ Creates a test product
- ✅ Creates a test price
- ✅ Creates a test payment intent
- ✅ Creates a test subscription
- ✅ Lists all customers to verify

### 2. `test_api.sh`
**Purpose**: Comprehensive API testing using curl commands

**Usage**:
```bash
# Make sure the service is running first
go run main.go

# In another terminal, run the script
./scripts/test_api.sh
```

**Or use the Makefile**:
```bash
make test-api
```

**What it does**:
- 🏥 Health check
- 👤 Customer operations (create, get, list)
- 📦 Product operations (create)
- 💰 Price operations (create)
- 💳 Payment intent operations (create)
- 📋 Subscription operations (create)
- ❌ Error handling tests
- 📊 Summary of created test data

**Requirements**:
- `curl` (usually pre-installed)
- `jq` (optional, for JSON formatting)
  - macOS: `brew install jq`
  - Ubuntu: `apt-get install jq`

## Prerequisites

### 1. Service Must Be Running
Before running any script, start the service:

```bash
# Option 1: With your real Stripe key
STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key_here go run main.go

# Option 2: With test key (will show errors but test the structure)
make start-dev
```

### 2. Environment Variables
The service requires a Stripe secret key:

```bash
export STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key_here
```

## Example Workflow

1. **Start the service**:
   ```bash
   STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key_here go run main.go
   ```

2. **Create test data**:
   ```bash
   make test-data
   ```

3. **Test all endpoints**:
   ```bash
   make test-api
   ```

4. **Verify with manual curl**:
   ```bash
   curl http://localhost:8080/api/v1/health
   curl http://localhost:8080/api/v1/customers
   ```

## Output Examples

### create_test_data.go Output:
```
🚀 Creating test data for Stripe Service...

1. Creating test customer...
✅ Created customer: cus_123456789

2. Creating test product...
✅ Created product: prod_123456789

3. Creating test price...
✅ Created price: price_123456789

4. Creating test payment intent...
✅ Created payment intent: pi_123456789

5. Creating test subscription...
✅ Created subscription: sub_123456789

🎉 Test data created successfully!
```

### test_api.sh Output:
```
🚀 Testing Stripe Service API...
✅ Service is running

1. Testing Health Check
{
  "status": "healthy",
  "service": "stripe-service"
}

2. Creating Test Customer
{
  "id": "cus_123456789",
  "email": "test@example.com",
  "name": "Test Customer"
}
✅ Customer created with ID: cus_123456789

...
```

## Troubleshooting

### Service Not Running
```
❌ Service is not running. Please start it first with: go run main.go
```
**Solution**: Start the service in another terminal

### Invalid Stripe Key
```
❌ Failed to create customer: HTTP 401: Invalid API key
```
**Solution**: Set a valid Stripe secret key

### Missing jq
```
⚠️  jq is not installed. JSON responses will not be formatted.
```
**Solution**: Install jq or ignore (scripts will still work)

## Notes

- These scripts use **test data** and won't affect real Stripe accounts
- The scripts are designed to be **idempotent** - safe to run multiple times
- All test data includes metadata marking it as created by scripts
- The service must be running on `localhost:8080` for scripts to work 