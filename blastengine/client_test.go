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
}
