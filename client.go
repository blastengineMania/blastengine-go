package blastenginego

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	apiKey string
	userId string
}

func initializeClient(apiKey string, userId string) Client {
	// Initialize the client
	c := Client{apiKey: apiKey, userId: userId}
	return c
}

func (c *Client) generateToken() string {
	// Concatenate userId and apiKey
	concatenated := c.userId + c.apiKey

	// Generate SHA256 hash
	hash := sha256.Sum256([]byte(concatenated))

	// Convert the hash to a lowercase hexadecimal string
	hexString := hex.EncodeToString(hash[:])

	// Base64 encode the lowercase hexadecimal string
	token := base64.URLEncoding.EncodeToString([]byte(hexString))

	return token
}

func (c *Client) SendTransaction(transaction *Transaction) error {
	url := "https://app.engn.jp/api/v1/deliveries/transaction"

	// Marshal the transaction struct to JSON
	jsonData, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.generateToken())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	return nil
}
