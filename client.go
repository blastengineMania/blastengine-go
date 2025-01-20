package blastenginego

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"
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

	// Convert hash to lowercase string
	hashString := strings.ToLower(base64.StdEncoding.EncodeToString(hash[:]))

	// Base64 encode the lowercase hash string
	token := base64.StdEncoding.EncodeToString([]byte(hashString))

	return token
}
