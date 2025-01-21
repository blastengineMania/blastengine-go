package blastengine

import (
	"encoding/json"
	"fmt"
)

type Transaction struct {
	DeliveryId  int
	From        MailAddress
	To          string
	Cc          []string
	Bcc         []string
	InsertCode  map[string]string
	Subject     string
	Encode      string
	TextPart    string
	HtmlPart    string
	Attachments []string
	Client      *Client
}

func (t *Transaction) SetFrom(email, name string) {
	t.From = NewMailAddress(email, name)
}

func (t *Transaction) SetTo(to string) {
	t.To = to
}

func (t *Transaction) AddCc(cc string) {
	if t.Cc == nil {
		t.Cc = make([]string, 0)
	}
	t.Cc = append(t.Cc, cc)
}

func (t *Transaction) AddBcc(bcc string) {
	if t.Bcc == nil {
		t.Bcc = make([]string, 0)
	}
	t.Bcc = append(t.Bcc, bcc)
}

func (t *Transaction) SetInsertCode(key string, value string) {
	if t.InsertCode == nil {
		t.InsertCode = make(map[string]string)
	}
	t.InsertCode[key] = value
}

func (t *Transaction) SetSubject(subject string) {
	t.Subject = subject
}

func (t *Transaction) SetEncode(encode string) {
	t.Encode = encode
}

func (t *Transaction) SetTextPart(textPart string) {
	t.TextPart = textPart
}

func (t *Transaction) SetHtmlPart(htmlPart string) {
	t.HtmlPart = htmlPart
}

func (t *Transaction) AddAttachment(attachment string) {
	if t.Attachments == nil {
		t.Attachments = make([]string, 0)
	}
	t.Attachments = append(t.Attachments, attachment)
}

func (t *Transaction) Send() error {
	if t.Attachments != nil {
		return t.SendMultipart()
	} else {
		return t.SendText()
	}
}

func (t *Transaction) GenerateJson() ([]byte, error) {
	// Convert InsertCode map to array of key-value pairs with __キー__ format
	insertCodeArray := make([]map[string]string, 0, len(t.InsertCode))
	for key, value := range t.InsertCode {
		insertCodeArray = append(insertCodeArray, map[string]string{"key": "__" + key + "__", "value": value})
	}
	// Create a temporary struct to hold the modified InsertCode
	tempTransaction := struct {
		From       MailAddress         `json:"from"`
		To         string              `json:"to"`
		Cc         []string            `json:"cc,omitempty"`
		Bcc        []string            `json:"bcc,omitempty"`
		InsertCode []map[string]string `json:"insert_code,omitempty"`
		Subject    string              `json:"subject"`
		Encode     string              `json:"encode"`
		TextPart   string              `json:"text_part"`
		HtmlPart   string              `json:"html_part,omitempty"`
	}{
		From:       t.From,
		To:         t.To,
		Cc:         t.Cc,
		Bcc:        t.Bcc,
		InsertCode: insertCodeArray,
		Subject:    t.Subject,
		Encode:     t.Encode,
		TextPart:   t.TextPart,
		HtmlPart:   t.HtmlPart,
	}

	// Marshal the temporary struct to JSON
	return json.Marshal(tempTransaction)
}

func (t *Transaction) SendText() error {
	url := "https://app.engn.jp/api/v1/deliveries/transaction"

	jsonData, err := t.GenerateJson()
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %v", err)
	}

	// Use the sendRequest method from Client
	deliveryId, err := t.Client.sendRequest(url, jsonData, false, nil)
	if err != nil {
		return err
	}

	t.DeliveryId = deliveryId
	return nil
}

func (t *Transaction) SendMultipart() error {
	url := "https://app.engn.jp/api/v1/deliveries/transaction"

	jsonData, err := t.GenerateJson()
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %v", err)
	}

	// Use the sendRequest method from Client
	deliveryId, err := t.Client.sendRequest(url, jsonData, true, t.Attachments)
	if err != nil {
		return err
	}

	t.DeliveryId = deliveryId
	return nil
}
