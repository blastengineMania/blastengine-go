package blastengine

import (
	"encoding/json"
	"fmt"
	"time"
)

type Attachment struct {
	DeliveryAttachID int       `json:"delivery_attach_id"`
	FileName         string    `json:"file_name"`
	Mime             string    `json:"mime"`
	Size             int       `json:"size"`
	CreatedTime      time.Time `json:"created_time"`
	UpdatedTime      time.Time `json:"updated_time"`
}

type Delivery struct {
	Client *Client

	DeliveryId      int          `json:"delivery_id"`
	From            MailAddress  `json:"from"`
	Status          string       `json:"status"`
	DeliveryTime    time.Time    `json:"delivery_time"`
	UpdatedTime     time.Time    `json:"updated_time"`
	CreatedTime     time.Time    `json:"created_time"`
	ReservationTime time.Time    `json:"reservation_time"`
	TextPart        string       `json:"text_part"`
	HtmlPart        string       `json:"html_part"`
	DeliveryType    string       `json:"delivery_type"`
	Subject         string       `json:"subject"`
	Attaches        []Attachment `json:"attaches"`
	OpenCount       int          `json:"open_count"`
	TotalCount      int          `json:"total_count"`
	SentCount       int          `json:"sent_count"`
	DropCount       int          `json:"drop_count"`
	SoftErrorCount  int          `json:"soft_error_count"`
	HardErrorCount  int          `json:"hard_error_count"`
}

func (d *Delivery) Get() error {
	if d.DeliveryId == 0 {
		return fmt.Errorf("delivery ID is required")
	}
	url := fmt.Sprintf("https://app.engn.jp/api/v1/deliveries/%d", d.DeliveryId)
	bodyBytes, err := d.Client.sendRequest("GET", url, nil, nil, false, nil)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyBytes, &d)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}
	return nil
}
