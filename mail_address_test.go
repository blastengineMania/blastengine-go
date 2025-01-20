package blastenginego

import "testing"

func TestNewMailAddress(t *testing.T) {
	email := "test@example.com"
	name := "Test User"
	mailAddress := NewMailAddress(email, name)

	if mailAddress.Email != email {
		t.Errorf("Expected Email to be %s, but got %s", email, mailAddress.Email)
	}

	if mailAddress.Name != name {
		t.Errorf("Expected Name to be %s, but got %s", name, mailAddress.Name)
	}

	// Test with empty name
	mailAddress = NewMailAddress(email, "")

	if mailAddress.Email != email {
		t.Errorf("Expected Email to be %s, but got %s", email, mailAddress.Email)
	}

	if mailAddress.Name != "" {
		t.Errorf("Expected Name to be empty, but got %s", mailAddress.Name)
	}
}
