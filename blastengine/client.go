package blastengine

import (
    "crypto/sha256"
    "encoding/base64"
    "strings"
    "fmt"
)

type Client struct {
    UserID string
    APIKey string
}

func Initialize(userID, apiKey string) Client {
    return Client{
        UserID: userID,
        APIKey: apiKey,
    }
}

func (c *Client) generateToken() string {
    data := c.UserID + c.APIKey
    hash := sha256.Sum256([]byte(data))
    hashStr := strings.ToLower(strings.ReplaceAll(fmt.Sprintf("%x", hash), "-", ""))
    token := base64.StdEncoding.EncodeToString([]byte(hashStr))
    return strings.ReplaceAll(token, "\n", "")
}

func (c *Client) Transaction() *Transaction {
    return &Transaction{}
}
