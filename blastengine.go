package blastengine

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

type Client struct {
	apiKey string
	userId string
}

func Initialize(apiKey string, userId string) Client {
	// Initialize the client
	return Client{apiKey: apiKey, userId: userId}
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

func (c *Client) NewTransaction() *Transaction {
	transaction := &Transaction{
		Encode: "UTF-8",
		Client: c,
	}
	return transaction
}
