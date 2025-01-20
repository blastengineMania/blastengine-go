package blastenginego

import "testing"

func TestSetFrom(t *testing.T) {
	transaction := NewTransaction()
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
	transaction := NewTransaction()
	subject := "Test subject"
	transaction.SetSubject(subject)

	if transaction.Subject != subject {
		t.Errorf("Expected Subject to be %s, but got %s", subject, transaction.Subject)
	}
}

func TestSetTo(t *testing.T) {
	transaction := NewTransaction()
	to := "user@example.jp"
	transaction.SetTo(to)

	if transaction.To != to {
		t.Errorf("Expected To to be %s, but got %s", to, transaction.To)
	}
}

func TestSetCc(t *testing.T) {
	transaction := NewTransaction()
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
	transaction := NewTransaction()
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
	transaction := NewTransaction()
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
	transaction := NewTransaction()
	encode := "ISO-8859-1"
	transaction.SetEncode(encode)

	if transaction.Encode != encode {
		t.Errorf("Expected Encode to be %s, but got %s", encode, transaction.Encode)
	}

	// Test default value
	defaultTransaction := NewTransaction()
	if defaultTransaction.Encode != "UTF-8" {
		t.Errorf("Expected default Encode to be UTF-8, but got %s", defaultTransaction.Encode)
	}
}

func TestSetTextPart(t *testing.T) {
	transaction := NewTransaction()
	textPart := "This is a text part"
	transaction.SetTextPart(textPart)

	if transaction.TextPart != textPart {
		t.Errorf("Expected TextPart to be %s, but got %s", textPart, transaction.TextPart)
	}
}

func TestSetHtmlPart(t *testing.T) {
	transaction := NewTransaction()
	htmlPart := "<p>This is an HTML part</p>"
	transaction.SetHtmlPart(htmlPart)

	if transaction.HtmlPart != htmlPart {
		t.Errorf("Expected HtmlPart to be %s, but got %s", htmlPart, transaction.HtmlPart)
	}
}
