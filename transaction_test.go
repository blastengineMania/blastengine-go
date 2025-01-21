package blastengine

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func getClient() Client {
	return initialize(os.Getenv("API_KEY"), os.Getenv("USER_ID"))
}

func TestSetFrom(t *testing.T) {
	client := getClient()
	transaction := client.NewTransaction()
	email := "test@example.com"
	name := "Test User"
	transaction.SetFrom(email, name)

	if transaction.From.Email != email {
		t.Errorf("Expected From.Email to be %s, but got %s", email, transaction.From.Email)
	}

	if transaction.From.Name != name {
		t.Errorf("Expected From.Name to be %s, but got %s", name, transaction.From.Name)
	}
}

func TestSetSubject(t *testing.T) {
	client := getClient()
	transaction := client.NewTransaction()
	subject := "Test subject"
	transaction.SetSubject(subject)

	if transaction.Subject != subject {
		t.Errorf("Expected Subject to be %s, but got %s", subject, transaction.Subject)
	}
}

func TestSetTo(t *testing.T) {
	client := getClient()
	transaction := client.NewTransaction()
	to := "user@example.jp"
	transaction.SetTo(to)

	if transaction.To != to {
		t.Errorf("Expected To to be %s, but got %s", to, transaction.To)
	}
}

func TestSetCc(t *testing.T) {
	client := getClient()
	transaction := client.NewTransaction()
	cc := []string{"cc1@example.com", "cc2@example.com"}
	transaction.SetCc(cc)

	if len(transaction.Cc) != len(cc) {
		t.Errorf("Expected Cc length to be %d, but got %d", len(cc), len(transaction.Cc))
	}

	for i, v := range cc {
		if transaction.Cc[i] != v {
			t.Errorf("Expected Cc[%d] to be %s, but got %s", i, v, transaction.Cc[i])
		}
	}
}

func TestSetBcc(t *testing.T) {
	client := getClient()
	transaction := client.NewTransaction()
	bcc := []string{"bcc1@example.com", "bcc2@example.com"}
	transaction.SetBcc(bcc)

	if len(transaction.Bcc) != len(bcc) {
		t.Errorf("Expected Bcc length to be %d, but got %d", len(bcc), len(transaction.Bcc))
	}

	for i, v := range bcc {
		if transaction.Bcc[i] != v {
			t.Errorf("Expected Bcc[%d] to be %s, but got %s", i, v, transaction.Bcc[i])
		}
	}
}

func TestSetInsertCode(t *testing.T) {
	client := getClient()
	transaction := client.NewTransaction()
	transaction.SetInsertCode("code1", "value1")
	transaction.SetInsertCode("code2", "value2")

	if len(transaction.InsertCode) != 2 {
		t.Errorf("Expected InsertCode length to be %d, but got %d", 2, len(transaction.InsertCode))
	}

	if transaction.InsertCode["code1"] != "value1" {
		t.Errorf("Expected InsertCode[code1] to be %s, but got %s", "value1", transaction.InsertCode["code1"])
	}

	if transaction.InsertCode["code2"] != "value2" {
		t.Errorf("Expected InsertCode[code2] to be %s, but got %s", "value2", transaction.InsertCode["code2"])
	}
}

func TestSetEncode(t *testing.T) {
	client := getClient()
	transaction := client.NewTransaction()
	encode := "ISO-8859-1"
	transaction.SetEncode(encode)

	if transaction.Encode != encode {
		t.Errorf("Expected Encode to be %s, but got %s", encode, transaction.Encode)
	}

	// Test default value
	defaultTransaction := client.NewTransaction()
	if defaultTransaction.Encode != "UTF-8" {
		t.Errorf("Expected default Encode to be UTF-8, but got %s", defaultTransaction.Encode)
	}
}

func TestSetTextPart(t *testing.T) {
	client := getClient()
	transaction := client.NewTransaction()
	textPart := "This is a text part"
	transaction.SetTextPart(textPart)

	if transaction.TextPart != textPart {
		t.Errorf("Expected TextPart to be %s, but got %s", textPart, transaction.TextPart)
	}
}

func TestSetHtmlPart(t *testing.T) {
	client := getClient()

	transaction := client.NewTransaction()
	htmlPart := "<p>This is an HTML part</p>"
	transaction.SetHtmlPart(htmlPart)

	if transaction.HtmlPart != htmlPart {
		t.Errorf("Expected HtmlPart to be %s, but got %s", htmlPart, transaction.HtmlPart)
	}
}

func TestSend(t *testing.T) {
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

		w.WriteHeader(http.StatusCreated)
	}))
	defer mockServer.Close()
	client := getClient()

	transaction := client.NewTransaction()
	transaction.SetFrom(os.Getenv("FROM"), "Test User")
	transaction.SetTo(os.Getenv("TO"))
	transaction.SetSubject("Test subject")
	transaction.SetTextPart("This is a text part")
	transaction.SetHtmlPart("<p>This is an HTML part</p>")
	transaction.Client = &client

	err := transaction.Send()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	// Check delivery id is up to zero
	if transaction.DeliveryId == 0 {
		t.Errorf("Expected DeliveryId to be 0, but got %d", transaction.DeliveryId)
	}
}
