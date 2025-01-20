package blastenginego

import (
	"crypto/sha256"
	"encoding/base64"
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

	// Base64 encode the hash string
	token := base64.StdEncoding.EncodeToString(hash[:])

	return token
}
