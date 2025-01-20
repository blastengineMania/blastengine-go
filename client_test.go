package blastenginego

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
	token := client.generateToken()

	if token != expectedToken {
		t.Errorf("Expected token to be %s, but got %s", expectedToken, token)
	}
}

func TestSendTransaction(t *testing.T) {
	// Mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, but got %s", r.Method)
		}

		if r.URL.EscapedPath() != "/api/v1/deliveries/transaction" {
			t.Errorf("Expected URL path to be /api/v1/deliveries/transaction, but got %s", r.URL.EscapedPath())
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type to be application/json, but got %s", r.Header.Get("Content-Type"))
		}

		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			t.Errorf("Expected Authorization header to start with Bearer, but got %s", r.Header.Get("Authorization"))
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	client := initializeClient("testApiKey", "testUserId")

	transaction := NewTransaction()
	transaction.SetFrom("test@example.com", "Test User")
	transaction.SetTo("user@example.jp")
	transaction.SetSubject("Test subject")
	transaction.SetTextPart("This is a text part")
	transaction.SetHtmlPart("<p>This is an HTML part</p>")

	err := client.SendTransaction(transaction)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}
