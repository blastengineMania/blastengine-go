package blastengine

import (
	"encoding/json"
	"fmt"
	"time"
)

type BulkTo struct {
	Email      string `json:"email"`
	InsertCode map[string]string
}

type Bulk struct {
	DeliveryId  int
	From        MailAddress
	To          []BulkTo
	Subject     string
	Encode      string
	TextPart    string
	HtmlPart    string
	Attachments []string
	Client      *Client
	Delivery
}

func (b *Bulk) SetDeliveryId(deliveryId int) {
	b.DeliveryId = deliveryId
	b.Delivery.DeliveryId = deliveryId
}

func (b *Bulk) SetFrom(email, name string) {
	b.From = NewMailAddress(email, name)
}

func (b *Bulk) AddTo(to string, insertCode map[string]string) {
	b.To = append(b.To, BulkTo{Email: to, InsertCode: insertCode})
}

func (b *Bulk) SetSubject(subject string) {
	b.Subject = subject
}

func (b *Bulk) SetEncode(encode string) {
	b.Encode = encode
}

func (b *Bulk) SetTextPart(textPart string) {
	b.TextPart = textPart
}

func (b *Bulk) SetHtmlPart(htmlPart string) {
	b.HtmlPart = htmlPart
}

func (b *Bulk) AddAttachment(attachment string) {
	if b.Attachments == nil {
		b.Attachments = make([]string, 0)
	}
	b.Attachments = append(b.Attachments, attachment)
}

func (b *Bulk) Begin() error {
	url := "https://app.engn.jp/api/v1/deliveries/bulk/begin"
	beginBulk := struct {
		From     MailAddress `json:"from"`
		Subject  string      `json:"subject"`
		Encode   string      `json:"encode"`
		TextPart string      `json:"text_part"`
		HtmlPart string      `json:"html_part,omitempty"`
	}{
		From:     b.From,
		Subject:  b.Subject,
		Encode:   b.Encode,
		TextPart: b.TextPart,
		HtmlPart: b.HtmlPart,
	}
	jsonData, err := json.Marshal(beginBulk)
	if err != nil {
		return fmt.Errorf("failed to marshal beginBulk: %v", err)
	}
	bodyBytes, err := b.Client.sendRequest("POST", url, nil, jsonData, false, nil)
	if err != nil {
		return err
	}
	var response struct {
		DeliveryId int `json:"delivery_id"`
	}
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}
	b.DeliveryId = response.DeliveryId
	b.Delivery.DeliveryId = response.DeliveryId
	return nil
}

func (b *Bulk) Update() error {
	if b.DeliveryId == 0 {
		return fmt.Errorf("delivery ID is required")
	}
	url := fmt.Sprintf("https://app.engn.jp/api/v1/deliveries/bulk/update/%d", b.DeliveryId)
	type JsonBulkTo struct {
		Email      string              `json:"email"`
		InsertCode []map[string]string `json:"insert_code,omitempty"`
	}
	bulkTos := make([]JsonBulkTo, 0, len(b.To))
	for _, to := range b.To {
		insertCodeArray := make([]map[string]string, 0, len(to.InsertCode))
		for key, value := range to.InsertCode {
			insertCodeArray = append(insertCodeArray, map[string]string{"key": "__" + key + "__", "value": value})
		}
		bulkTos = append(bulkTos, JsonBulkTo{Email: to.Email, InsertCode: insertCodeArray})
	}
	beginBulk := struct {
		From     MailAddress  `json:"from"`
		To       []JsonBulkTo `json:"to"`
		Subject  string       `json:"subject"`
		TextPart string       `json:"text_part"`
		HtmlPart string       `json:"html_part,omitempty"`
	}{
		From:     b.From,
		To:       bulkTos,
		Subject:  b.Subject,
		TextPart: b.TextPart,
		HtmlPart: b.HtmlPart,
	}
	jsonData, err := json.Marshal(beginBulk)
	if err != nil {
		return fmt.Errorf("failed to marshal beginBulk: %v", err)
	}
	bodyBytes, err := b.Client.sendRequest("PUT", url, nil, jsonData, false, nil)
	if err != nil {
		return err
	}
	var response struct {
		DeliveryId int `json:"delivery_id"`
	}
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}
	b.DeliveryId = response.DeliveryId
	b.Delivery.DeliveryId = response.DeliveryId
	return nil
}

func (b *Bulk) Calcel() error {
	if b.DeliveryId == 0 {
		return fmt.Errorf("delivery ID is required")
	}
	url := fmt.Sprintf("https://app.engn.jp/api/v1/deliveries/%d/cancel", b.DeliveryId)
	_, err := b.Client.sendRequest("PATCH", url, nil, nil, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bulk) Delete() error {
	if b.DeliveryId == 0 {
		return fmt.Errorf("delivery ID is required")
	}
	url := fmt.Sprintf("https://app.engn.jp/api/v1/deliveries/%d", b.DeliveryId)
	_, err := b.Client.sendRequest("DELETE", url, nil, nil, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bulk) Send(reservationTime *time.Time) error {
	if reservationTime != nil {
		return b.SendReservation(reservationTime)
	} else {
		return b.SendNow()
	}
}

func (b *Bulk) SendNow() error {
	url := fmt.Sprintf("https://app.engn.jp/api/v1/deliveries/bulk/commit/%d/immediate", b.DeliveryId)
	// Use the sendRequest method from Client
	_, err := b.Client.sendRequest("PATCH", url, nil, nil, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bulk) SendReservation(reservationTime *time.Time) error {
	url := fmt.Sprintf("https://app.engn.jp/api/v1/deliveries/bulk/commit/%d", b.DeliveryId)
	reservationTimeJson := struct {
		ReservationTime string `json:"reservation_time"`
	}{
		ReservationTime: reservationTime.Format("2006-01-02T15:04:05Z07:00"),
	}
	jsonData, err := json.Marshal(reservationTimeJson)
	if err != nil {
		return fmt.Errorf("failed to marshal beginBulk: %v", err)
	}
	_, err = b.Client.sendRequest("PATCH", url, nil, jsonData, false, nil)
	if err != nil {
		return err
	}
	return nil
}
