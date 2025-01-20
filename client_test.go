package blastenginego

import (
	"testing"
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

func TestInitializeClient(t *testing.T) {
	apiKey := "testApiKey"
	userId := "testUserId"
	client := initializeClient(apiKey, userId)

	if client.apiKey != apiKey {
		t.Errorf("Expected apiKey to be %s, but got %s", apiKey, client.apiKey)
	}

	if client.userId != userId {
		t.Errorf("Expected userId to be %s, but got %s", userId, client.userId)
	}
}

func TestGenerateToken(t *testing.T) {
	apiKey := "testApiKey"
	userId := "testUserId"
	client := initializeClient(apiKey, userId)

	expectedToken := "NGY4YjlhNzE0OWYzMTFiNDE5OTJhMmJlYTQxMDlkMmE4MmY1MTNhZWVjNWVhZDRiOGFkNzgxYzZmZmY3MTZjNg=="
	generatedToken := client.generateToken()

	if generatedToken != expectedToken {
		t.Errorf("Expected token to be %s, but got %s", expectedToken, generatedToken)
	}
}
