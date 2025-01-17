package blastengine

import (
    "testing"
)

func TestTransaction(t *testing.T) {
    client := Initialize("USER_ID", "API_KEY")
    transaction := client.Transaction()

    transaction.SetSubject("Test email subject")
    transaction.SetTextPart("Test email body")
    transaction.SetHtmlPart("<h1>Test email body</h1>")

    if transaction.Subject != "Test email subject" {
        t.Errorf("Expected Subject %s, but got %s", "Test email subject", transaction.Subject)
    }

    if transaction.TextPart != "Test email body" {
        t.Errorf("Expected TextPart %s, but got %s", "Test email body", transaction.TextPart)
    }

    if transaction.HtmlPart != "<h1>Test email body</h1>" {
        t.Errorf("Expected HtmlPart %s, but got %s", "<h1>Test email body</h1>", transaction.HtmlPart)
    }

    // Assuming Send method has some way to verify the email was sent correctly
    // This part of the test will need to be implemented based on the actual Send method implementation
    transaction.Send()
}
