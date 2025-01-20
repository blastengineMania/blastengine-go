package blastenginego

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Transaction struct {
	From       MailAddress
	To         string
	Cc         []string
	Bcc        []string
	InsertCode map[string]string
	Subject    string
	Encode     string
	TextPart   string
	HtmlPart   string
	Client     *Client
}

func NewTransaction() *Transaction {
	return &Transaction{
		Encode: "UTF-8",
	}
}

func (t *Transaction) SetFrom(email, name string) {
	t.From = NewMailAddress(email, name)
}

func (t *Transaction) SetTo(to string) {
	t.To = to
}

func (t *Transaction) SetCc(cc []string) {
	t.Cc = cc
}

func (t *Transaction) SetBcc(bcc []string) {
	t.Bcc = bcc
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

func (t *Transaction) Send() error {
	url := "https://app.engn.jp/api/v1/deliveries/transaction"

	// Convert InsertCode map to array of key-value pairs with __キー__ format
	insertCodeArray := make([]map[string]string, 0, len(t.InsertCode))
	for key, value := range t.InsertCode {
		insertCodeArray = append(insertCodeArray, map[string]string{"key": "__" + key + "__", "value": value})
	}

	// Create a temporary struct to hold the modified InsertCode
	tempTransaction := struct {
		From       MailAddress        `json:"from"`
		To         string             `json:"to"`
		Cc         []string           `json:"cc"`
		Bcc        []string           `json:"bcc"`
		InsertCode []map[string]string `json:"insert_code"`
		Subject    string             `json:"subject"`
		Encode     string             `json:"encode"`
		TextPart   string             `json:"text_part"`
		HtmlPart   string             `json:"html_part"`
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
	jsonData, err := json.Marshal(tempTransaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.Client.generateToken())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Println("Error response:", bodyString)
		return fmt.Errorf("received non-201 response: %d", resp.StatusCode)
	}

	return nil
}
