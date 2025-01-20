package blastenginego

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
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
