package blastengine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
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

func (t *Transaction) SendText() error {
	url := "https://app.engn.jp/api/v1/deliveries/transaction"

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
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Println("Error response:", bodyString)
		return fmt.Errorf("received non-201 response: %d", resp.StatusCode)
	}
	// Parse the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	// Parse the response JSON
	var response struct {
		DeliveryId int `json:"delivery_id"`
	}
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}
	t.DeliveryId = response.DeliveryId
	return nil
}

func (t *Transaction) SendMultipart() error {
	url := "https://app.engn.jp/api/v1/deliveries/transaction"

	// Create a buffer to hold the multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add the JSON data part
	partHeaders := textproto.MIMEHeader{}
	partHeaders.Set("Content-Disposition", `form-data; name="data"`)
	partHeaders.Set("Content-Type", "application/json")
	dataPart, err := writer.CreatePart(partHeaders)
	if err != nil {
		return fmt.Errorf("failed to create form field: %v", err)
	}

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
	jsonData, err := json.Marshal(tempTransaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %v", err)
	}

	_, err = dataPart.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write JSON data to form field: %v", err)
	}

	// Add the file parts
	for _, attachment := range t.Attachments {
		file, err := os.Open(attachment)
		if err != nil {
			return fmt.Errorf("failed to open attachment: %v", err)
		}
		defer file.Close()

		filePart, err := writer.CreateFormFile("file", filepath.Base(attachment))
		if err != nil {
			return fmt.Errorf("failed to create form file: %v", err)
		}

		_, err = io.Copy(filePart, file)
		if err != nil {
			return fmt.Errorf("failed to copy file content to form file: %v", err)
		}
	}

	// Close the multipart writer to set the terminating boundary
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close multipart writer: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set request headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
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
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Println("Error response:", bodyString)
		return fmt.Errorf("received non-201 response: %d", resp.StatusCode)
	}

	// Parse the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the response JSON
	var response struct {
		DeliveryId int `json:"delivery_id"`
	}
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}
	t.DeliveryId = response.DeliveryId
	return nil
}
