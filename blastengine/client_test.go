package blastengine

import (
    "testing"
)

func TestInitialize(t *testing.T) {
    userID := "test_user"
    apiKey := "test_key"
    client := Initialize(userID, apiKey)

    if client.UserID != userID {
        t.Errorf("Expected UserID %s, but got %s", userID, client.UserID)
    }

    if client.APIKey != apiKey {
        t.Errorf("Expected APIKey %s, but got %s", apiKey, client.APIKey)
    }

    expectedToken := "dGVzdF91c2VyZGVmYXVsdF90b2tlbl92YWx1ZQ==" // This should be the expected token value based on the userID and apiKey
    generatedToken := client.generateToken()
    if generatedToken != expectedToken {
        t.Errorf("Expected token %s, but got %s", expectedToken, generatedToken)
    }
}
