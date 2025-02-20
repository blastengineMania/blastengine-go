package blastengine

import (
	"fmt"
	"os"
	"testing"
)

func getClient() Client {
	return Initialize(os.Getenv("API_KEY"), os.Getenv("USER_ID"))
}

func TestTransactionSetFrom(t *testing.T) {
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

func TestTransactionSetSubject(t *testing.T) {
	client := getClient()
	transaction := client.NewTransaction()
	subject := "Test subject"
	transaction.SetSubject(subject)

	if transaction.Subject != subject {
		t.Errorf("Expected Subject to be %s, but got %s", subject, transaction.Subject)
	}
}

func TestTransactionSetTo(t *testing.T) {
	client := getClient()
	transaction := client.NewTransaction()
	to := "user@example.jp"
	transaction.SetTo(to)

	if transaction.To != to {
		t.Errorf("Expected To to be %s, but got %s", to, transaction.To)
	}
}

func TestTransactionAddCc(t *testing.T) {
	client := getClient()
	transaction := client.NewTransaction()
	cc := []string{"cc1@example.com", "cc2@example.com"}
	for _, v := range cc {
		transaction.AddCc(v)
	}

	if len(transaction.Cc) != len(cc) {
		t.Errorf("Expected Cc length to be %d, but got %d", len(cc), len(transaction.Cc))
	}

	for i, v := range cc {
		if transaction.Cc[i] != v {
			t.Errorf("Expected Cc[%d] to be %s, but got %s", i, v, transaction.Cc[i])
		}
	}
}

func TestTransactionAddBcc(t *testing.T) {
	client := getClient()
	transaction := client.NewTransaction()
	bcc := []string{"bcc1@example.com", "bcc2@example.com"}
	for _, v := range bcc {
		transaction.AddBcc(v)
	}

	if len(transaction.Bcc) != len(bcc) {
		t.Errorf("Expected Bcc length to be %d, but got %d", len(bcc), len(transaction.Bcc))
	}

	for i, v := range bcc {
		if transaction.Bcc[i] != v {
			t.Errorf("Expected Bcc[%d] to be %s, but got %s", i, v, transaction.Bcc[i])
		}
	}
}

func TestTransactionSetInsertCode(t *testing.T) {
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

func TestTransactionSetEncode(t *testing.T) {
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
	fmt.Println(defaultTransaction.Encode)
}

func TestTransactionSetTextPart(t *testing.T) {
	client := getClient()
	transaction := client.NewTransaction()
	textPart := "This is a text part"
	transaction.SetTextPart(textPart)

	if transaction.TextPart != textPart {
		t.Errorf("Expected TextPart to be %s, but got %s", textPart, transaction.TextPart)
	}
}

func TestTransactionSetHtmlPart(t *testing.T) {
	client := getClient()

	transaction := client.NewTransaction()
	htmlPart := "<p>This is an HTML part</p>"
	transaction.SetHtmlPart(htmlPart)

	if transaction.HtmlPart != htmlPart {
		t.Errorf("Expected HtmlPart to be %s, but got %s", htmlPart, transaction.HtmlPart)
	}
}

func TestTransactionSend(t *testing.T) {
	// Mock HTTP server
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
	err = transaction.Get()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if transaction.DeliveryId == 0 {
		t.Errorf("Expected DeliveryId to be not 0, but got 0")
	}
	if transaction.Status != "SENT" {
		t.Errorf("Expected Status to be sent, but got %s", transaction.Status)
	}
	if transaction.DeliveryType != "TRANSACTION" {
		t.Errorf("Expected deliveryType to be TRANSACTION, but got %s", transaction.DeliveryType)
	}
	if transaction.CreatedTime.IsZero() {
		t.Errorf("Expected createdTime to be zero, but got %v", transaction.CreatedTime)
	}
	if !transaction.ReservationTime.IsZero() {
		t.Errorf("Expected reservationTime to be not zero, but got zero")
	}
	if transaction.UpdatedTime.IsZero() {
		t.Errorf("Expected updatedTime to be not zero, but got zero")
	}
}

func TestTransactionSendMultipart(t *testing.T) {
	client := getClient()

	transaction := client.NewTransaction()
	transaction.SetFrom(os.Getenv("FROM"), "Test User")
	transaction.SetTo(os.Getenv("TO"))
	transaction.SetSubject("Test subject")
	transaction.SetTextPart("This is a text part")
	transaction.SetHtmlPart("<p>This is an HTML part</p>")
	transaction.AddAttachment("README.md")
	transaction.AddAttachment("LICENSE")
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
