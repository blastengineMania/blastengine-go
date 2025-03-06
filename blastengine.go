package blastengine

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
)

type Client struct {
	apiKey string
	userId string
}

type Response struct {
	DeliveryId int `json:"delivery_id"`
	JobId      int `json:"job_id"`
}

type Attachment struct {
	FileName string
	Content  io.Reader
}

type ListUnsubscribeParams struct {
	Email string `json:"mailto,omitempty"`
	Url   string `json:"url,omitempty"`
}

func Initialize(apiKey string, userId string) Client {
	// Initialize the client
	return Client{apiKey: apiKey, userId: userId}
}

func (c *Client) generateToken() string {
	// Concatenate userId and apiKey
	concatenated := c.userId + c.apiKey

	// Generate SHA256 hash
	hash := sha256.Sum256([]byte(concatenated))

	// Convert the hash to a lowercase hexadecimal string
	hexString := hex.EncodeToString(hash[:])

	// Base64 encode the lowercase hexadecimal string
	token := base64.URLEncoding.EncodeToString([]byte(hexString))

	return token
}

func (c *Client) NewTransaction() *Transaction {
	transaction := &Transaction{
		Encode: "UTF-8",
		Client: c,
	}
	transaction.Mail.Client = c
	return transaction
}

func (c *Client) NewBulk() *Bulk {
	bulk := &Bulk{
		Encode: "UTF-8",
		Client: c,
	}
	bulk.Mail.Client = c
	return bulk
}

func (c *Client) NewMail() *Mail {
	mail := &Mail{
		Encode: "UTF-8",
		Client: c,
	}
	mail.Client = c
	return mail
}

func (c *Client) sendRequest(method string, baseUrl string, queries url.Values, jsonData []byte, isMultipart bool, attachments []Attachment) ([]byte, error) {
	var req *http.Request
	var err error
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}
	if queries != nil {
		u.RawQuery = queries.Encode()
	}
	if isMultipart {
		var requestBody bytes.Buffer
		writer := multipart.NewWriter(&requestBody)
		partHeaders := textproto.MIMEHeader{}
		partHeaders.Set("Content-Disposition", `form-data; name="data"`)
		partHeaders.Set("Content-Type", "application/json")
		dataPart, err := writer.CreatePart(partHeaders)
		if err != nil {
			return nil, fmt.Errorf("failed to create form field: %v", err)
		}

		_, err = dataPart.Write(jsonData)
		if err != nil {
			return nil, fmt.Errorf("failed to write JSON data to form field: %v", err)
		}

		for _, attachment := range attachments {
			filePart, err := writer.CreateFormFile("file", attachment.FileName)
			if err != nil {
				return nil, fmt.Errorf("failed to create form field: %v", err)
			}
			_, err = io.Copy(filePart, attachment.Content)
			if err != nil {
				return nil, fmt.Errorf("failed to write file to form field: %v", err)
			}
		}

		err = writer.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close multipart writer: %v", err)
		}

		req, err = http.NewRequest(method, u.String(), &requestBody)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		if jsonData != nil {
			req, err = http.NewRequest(method, u.String(), bytes.NewBuffer(jsonData))
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %v", err)
			}
		} else {
			req, err = http.NewRequest(method, u.String(), nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %v", err)
			}
		}

		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Authorization", "Bearer "+c.generateToken())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Println("Error response:", bodyString)
		return nil, fmt.Errorf("received non-20x response: %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
