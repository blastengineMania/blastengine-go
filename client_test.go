package blastenginego

import "testing"

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
	token := client.generateToken()

	if token != expectedToken {
		t.Errorf("Expected token to be %s, but got %s", expectedToken, token)
	}
}
