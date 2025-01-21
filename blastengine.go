package blastengine

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
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

type Client struct {
	apiKey string
	userId string
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
	return transaction
}

func (c *Client) sendRequest(url string, jsonData []byte, isMultipart bool, attachments []string) (int, error) {
	var req *http.Request
	var err error

	if isMultipart {
		var requestBody bytes.Buffer
		writer := multipart.NewWriter(&requestBody)

		partHeaders := textproto.MIMEHeader{}
		partHeaders.Set("Content-Disposition", `form-data; name="data"`)
		partHeaders.Set("Content-Type", "application/json")
		dataPart, err := writer.CreatePart(partHeaders)
		if err != nil {
			return 0, fmt.Errorf("failed to create form field: %v", err)
		}

		_, err = dataPart.Write(jsonData)
		if err != nil {
			return 0, fmt.Errorf("failed to write JSON data to form field: %v", err)
		}

		for _, attachment := range attachments {
			file, err := os.Open(attachment)
			if err != nil {
				return 0, fmt.Errorf("failed to open attachment: %v", err)
			}
			defer file.Close()

			filePart, err := writer.CreateFormFile("file", filepath.Base(attachment))
			if err != nil {
				return 0, fmt.Errorf("failed to create form file: %v", err)
			}

			_, err = io.Copy(filePart, file)
			if err != nil {
				return 0, fmt.Errorf("failed to copy file content to form file: %v", err)
			}
		}

		err = writer.Close()
		if err != nil {
			return 0, fmt.Errorf("failed to close multipart writer: %v", err)
		}

		req, err = http.NewRequest("POST", url, &requestBody)
		if err != nil {
			return 0, fmt.Errorf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return 0, fmt.Errorf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Authorization", "Bearer "+c.generateToken())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Println("Error response:", bodyString)
		return 0, fmt.Errorf("received non-201 response: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %v", err)
	}

	var response struct {
		DeliveryId int `json:"delivery_id"`
	}
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return response.DeliveryId, nil
}
