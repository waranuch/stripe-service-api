#!/bin/bash

# Test API Script for Stripe Service
# This script tests all available API endpoints

BASE_URL="http://localhost:8080/api/v1"
CUSTOMER_ID=""
PRODUCT_ID=""
PRICE_ID=""
PAYMENT_INTENT_ID=""
SUBSCRIPTION_ID=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "🚀 Testing Stripe Service API..."

# Function to check if service is running
check_service() {
    echo -e "${BLUE}Checking if service is running...${NC}"
    if curl -s "$BASE_URL/health" > /dev/null; then
        echo -e "${GREEN}✅ Service is running${NC}"
    else
        echo -e "${RED}❌ Service is not running. Please start it first with: go run main.go${NC}"
        exit 1
    fi
}

# Function to test health endpoint
test_health() {
    echo -e "\n${BLUE}1. Testing Health Check${NC}"
    curl -s "$BASE_URL/health" | jq '.'
}

# Function to create test customer
create_customer() {
    echo -e "\n${BLUE}2. Creating Test Customer${NC}"
    RESPONSE=$(curl -s -X POST "$BASE_URL/customers" \
        -H "Content-Type: application/json" \
        -d '{
            "email": "test@example.com",
            "name": "Test Customer",
            "phone": "+1234567890",
            "description": "Test customer created by script",
            "metadata": {
                "source": "test_script",
                "env": "development"
            }
        }')
    
    echo "$RESPONSE" | jq '.'
    CUSTOMER_ID=$(echo "$RESPONSE" | jq -r '.id')
    echo -e "${GREEN}✅ Customer created with ID: $CUSTOMER_ID${NC}"
}

# Function to get customer
get_customer() {
    echo -e "\n${BLUE}3. Getting Customer${NC}"
    curl -s "$BASE_URL/customers/$CUSTOMER_ID" | jq '.'
}

# Function to list customers
list_customers() {
    echo -e "\n${BLUE}4. Listing Customers${NC}"
    curl -s "$BASE_URL/customers" | jq '.'
}

# Function to create test product
create_product() {
    echo -e "\n${BLUE}5. Creating Test Product${NC}"
    RESPONSE=$(curl -s -X POST "$BASE_URL/products" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Test Product",
            "description": "A test product created by script",
            "active": true,
            "metadata": {
                "category": "test",
                "source": "test_script"
            }
        }')
    
    echo "$RESPONSE" | jq '.'
    PRODUCT_ID=$(echo "$RESPONSE" | jq -r '.id')
    echo -e "${GREEN}✅ Product created with ID: $PRODUCT_ID${NC}"
}

# Function to create test price
create_price() {
    echo -e "\n${BLUE}6. Creating Test Price${NC}"
    RESPONSE=$(curl -s -X POST "$BASE_URL/prices" \
        -H "Content-Type: application/json" \
        -d "{
            \"product_id\": \"$PRODUCT_ID\",
            \"unit_amount\": 999,
            \"currency\": \"usd\",
            \"type\": \"one_time\",
            \"active\": true,
            \"metadata\": {
                \"source\": \"test_script\"
            }
        }")
    
    echo "$RESPONSE" | jq '.'
    PRICE_ID=$(echo "$RESPONSE" | jq -r '.id')
    echo -e "${GREEN}✅ Price created with ID: $PRICE_ID${NC}"
}

# Function to create test payment intent
create_payment_intent() {
    echo -e "\n${BLUE}7. Creating Test Payment Intent${NC}"
    RESPONSE=$(curl -s -X POST "$BASE_URL/payment-intents" \
        -H "Content-Type: application/json" \
        -d "{
            \"amount\": 2000,
            \"currency\": \"usd\",
            \"customer_id\": \"$CUSTOMER_ID\",
            \"description\": \"Test payment intent\",
            \"metadata\": {
                \"source\": \"test_script\"
            }
        }")
    
    echo "$RESPONSE" | jq '.'
    PAYMENT_INTENT_ID=$(echo "$RESPONSE" | jq -r '.id')
    echo -e "${GREEN}✅ Payment Intent created with ID: $PAYMENT_INTENT_ID${NC}"
}

# Function to create test subscription
create_subscription() {
    echo -e "\n${BLUE}8. Creating Test Subscription${NC}"
    RESPONSE=$(curl -s -X POST "$BASE_URL/subscriptions" \
        -H "Content-Type: application/json" \
        -d "{
            \"customer_id\": \"$CUSTOMER_ID\",
            \"price_id\": \"$PRICE_ID\",
            \"metadata\": {
                \"source\": \"test_script\",
                \"plan\": \"test_plan\"
            }
        }")
    
    echo "$RESPONSE" | jq '.'
    SUBSCRIPTION_ID=$(echo "$RESPONSE" | jq -r '.id')
    echo -e "${GREEN}✅ Subscription created with ID: $SUBSCRIPTION_ID${NC}"
}

# Function to test error handling
test_error_handling() {
    echo -e "\n${BLUE}9. Testing Error Handling${NC}"
    
    echo -e "${YELLOW}Testing invalid customer creation (missing email):${NC}"
    curl -s -X POST "$BASE_URL/customers" \
        -H "Content-Type: application/json" \
        -d '{"name": "Test Customer"}' | jq '.'
    
    echo -e "\n${YELLOW}Testing invalid payment intent (missing amount):${NC}"
    curl -s -X POST "$BASE_URL/payment-intents" \
        -H "Content-Type: application/json" \
        -d '{"currency": "usd"}' | jq '.'
    
    echo -e "\n${YELLOW}Testing non-existent customer:${NC}"
    curl -s "$BASE_URL/customers/cus_nonexistent" | jq '.'
}

# Function to print summary
print_summary() {
    echo -e "\n${GREEN}🎉 API Testing Complete!${NC}"
    echo -e "${BLUE}📊 Summary of Created Test Data:${NC}"
    echo -e "   Customer ID: ${GREEN}$CUSTOMER_ID${NC}"
    echo -e "   Product ID: ${GREEN}$PRODUCT_ID${NC}"
    echo -e "   Price ID: ${GREEN}$PRICE_ID${NC}"
    echo -e "   Payment Intent ID: ${GREEN}$PAYMENT_INTENT_ID${NC}"
    echo -e "   Subscription ID: ${GREEN}$SUBSCRIPTION_ID${NC}"
    
    echo -e "\n${BLUE}💡 You can now test individual endpoints:${NC}"
    echo -e "   Health: curl $BASE_URL/health"
    echo -e "   Get Customer: curl $BASE_URL/customers/$CUSTOMER_ID"
    echo -e "   List Customers: curl $BASE_URL/customers"
    echo -e "   Cancel Subscription: curl -X DELETE $BASE_URL/subscriptions/$SUBSCRIPTION_ID"
}

# Main execution
main() {
    check_service
    test_health
    create_customer
    get_customer
    list_customers
    create_product
    create_price
    create_payment_intent
    create_subscription
    test_error_handling
    print_summary
}

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo -e "${YELLOW}⚠️  jq is not installed. JSON responses will not be formatted.${NC}"
    echo -e "${YELLOW}   Install jq for better output: brew install jq (macOS) or apt-get install jq (Ubuntu)${NC}"
fi

# Run the main function
main 