package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080/api/v1"

type TestData struct {
	CustomerID      string
	ProductID       string
	PriceID         string
	PaymentIntentID string
	SubscriptionID  string
}

func main() {
	fmt.Println("🚀 Creating test data for Stripe Service...")

	// Check if service is running
	if !isServiceRunning() {
		fmt.Println("❌ Service is not running. Please start the service first with: go run main.go")
		return
	}

	var testData TestData

	// 1. Create a test customer
	fmt.Println("\n1. Creating test customer...")
	customerID, err := createTestCustomer()
	if err != nil {
		fmt.Printf("❌ Failed to create customer: %v\n", err)
		return
	}
	testData.CustomerID = customerID
	fmt.Printf("✅ Created customer: %s\n", customerID)

	// 2. Create a test product
	fmt.Println("\n2. Creating test product...")
	productID, err := createTestProduct()
	if err != nil {
		fmt.Printf("❌ Failed to create product: %v\n", err)
		return
	}
	testData.ProductID = productID
	fmt.Printf("✅ Created product: %s\n", productID)

	// 3. Create a test price
	fmt.Println("\n3. Creating test price...")
	priceID, err := createTestPrice(productID)
	if err != nil {
		fmt.Printf("❌ Failed to create price: %v\n", err)
		return
	}
	testData.PriceID = priceID
	fmt.Printf("✅ Created price: %s\n", priceID)

	// 4. Create a test payment intent
	fmt.Println("\n4. Creating test payment intent...")
	paymentIntentID, err := createTestPaymentIntent(customerID)
	if err != nil {
		fmt.Printf("❌ Failed to create payment intent: %v\n", err)
		return
	}
	testData.PaymentIntentID = paymentIntentID
	fmt.Printf("✅ Created payment intent: %s\n", paymentIntentID)

	// 5. Create a test subscription
	fmt.Println("\n5. Creating test subscription...")
	subscriptionID, err := createTestSubscription(customerID, priceID)
	if err != nil {
		fmt.Printf("❌ Failed to create subscription: %v\n", err)
		return
	}
	testData.SubscriptionID = subscriptionID
	fmt.Printf("✅ Created subscription: %s\n", subscriptionID)

	// 6. List customers to verify
	fmt.Println("\n6. Listing customers...")
	if err := listCustomers(); err != nil {
		fmt.Printf("❌ Failed to list customers: %v\n", err)
	}

	// Print summary
	fmt.Println("\n🎉 Test data created successfully!")
	fmt.Println("📊 Summary:")
	fmt.Printf("   Customer ID: %s\n", testData.CustomerID)
	fmt.Printf("   Product ID: %s\n", testData.ProductID)
	fmt.Printf("   Price ID: %s\n", testData.PriceID)
	fmt.Printf("   Payment Intent ID: %s\n", testData.PaymentIntentID)
	fmt.Printf("   Subscription ID: %s\n", testData.SubscriptionID)

	fmt.Println("\n💡 You can now test the API endpoints with this data!")
	fmt.Println("   Example: curl http://localhost:8080/api/v1/customers/" + testData.CustomerID)
}

func isServiceRunning() bool {
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func createTestCustomer() (string, error) {
	payload := map[string]interface{}{
		"email":       "test@example.com",
		"name":        "Test Customer",
		"phone":       "+1234567890",
		"description": "Test customer created by script",
		"metadata": map[string]string{
			"source": "test_script",
			"env":    "development",
		},
	}

	resp, err := makeRequest("POST", "/customers", payload)
	if err != nil {
		return "", err
	}

	var customer map[string]interface{}
	if err := json.Unmarshal(resp, &customer); err != nil {
		return "", err
	}

	return customer["id"].(string), nil
}

func createTestProduct() (string, error) {
	payload := map[string]interface{}{
		"name":        "Test Product",
		"description": "A test product created by script",
		"active":      true,
		"metadata": map[string]string{
			"category": "test",
			"source":   "test_script",
		},
	}

	resp, err := makeRequest("POST", "/products", payload)
	if err != nil {
		return "", err
	}

	var product map[string]interface{}
	if err := json.Unmarshal(resp, &product); err != nil {
		return "", err
	}

	return product["id"].(string), nil
}

func createTestPrice(productID string) (string, error) {
	payload := map[string]interface{}{
		"product_id":  productID,
		"unit_amount": 999, // $9.99
		"currency":    "usd",
		"type":        "one_time",
		"active":      true,
		"metadata": map[string]string{
			"source": "test_script",
		},
	}

	resp, err := makeRequest("POST", "/prices", payload)
	if err != nil {
		return "", err
	}

	var price map[string]interface{}
	if err := json.Unmarshal(resp, &price); err != nil {
		return "", err
	}

	return price["id"].(string), nil
}

func createTestPaymentIntent(customerID string) (string, error) {
	payload := map[string]interface{}{
		"amount":      2000, // $20.00
		"currency":    "usd",
		"customer_id": customerID,
		"description": "Test payment intent",
		"metadata": map[string]string{
			"source": "test_script",
		},
	}

	resp, err := makeRequest("POST", "/payment-intents", payload)
	if err != nil {
		return "", err
	}

	var paymentIntent map[string]interface{}
	if err := json.Unmarshal(resp, &paymentIntent); err != nil {
		return "", err
	}

	return paymentIntent["id"].(string), nil
}

func createTestSubscription(customerID, priceID string) (string, error) {
	payload := map[string]interface{}{
		"customer_id": customerID,
		"price_id":    priceID,
		"metadata": map[string]string{
			"source": "test_script",
			"plan":   "test_plan",
		},
	}

	resp, err := makeRequest("POST", "/subscriptions", payload)
	if err != nil {
		return "", err
	}

	var subscription map[string]interface{}
	if err := json.Unmarshal(resp, &subscription); err != nil {
		return "", err
	}

	return subscription["id"].(string), nil
}

func listCustomers() error {
	resp, err := makeRequest("GET", "/customers", nil)
	if err != nil {
		return err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(resp, &response); err != nil {
		return err
	}

	customers := response["customers"].([]interface{})
	fmt.Printf("✅ Found %d customers\n", len(customers))

	for i, customer := range customers {
		c := customer.(map[string]interface{})
		fmt.Printf("   %d. %s (%s) - %s\n", i+1, c["name"], c["id"], c["email"])
	}

	return nil
}

func makeRequest(method, endpoint string, payload interface{}) ([]byte, error) {
	var body io.Reader

	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, baseURL+endpoint, body)
	if err != nil {
		return nil, err
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
