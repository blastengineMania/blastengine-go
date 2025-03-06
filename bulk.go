package blastengine

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type BulkTo struct {
	Email      string `json:"email"`
	InsertCode map[string]string
}

type Bulk struct {
	DeliveryId      int
	From            MailAddress
	To              []BulkTo
	Subject         string
	Encode          string
	TextPart        string
	HtmlPart        string
	Attachments     []string
	ListUnsubscribe *ListUnsubscribeParams
	Client          *Client
	Mail
}

type ImportParams struct {
	IgnoreErrors bool `json:"ignore_errors,omitempty"`
	Immediate    bool `json:"immediate,omitempty"`
}

func (b *Bulk) SetDeliveryId(deliveryId int) {
	b.DeliveryId = deliveryId
	b.Mail.DeliveryId = deliveryId
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

func (b *Bulk) SetListUnsubscribe(params *ListUnsubscribeParams) {
	if params.Email != "" {
		params.Email = fmt.Sprintf("mailto:%s", params.Email)
	}
	b.ListUnsubscribe = &ListUnsubscribeParams{
		Email: params.Email,
		Url:   params.Url,
	}
}

func (b *Bulk) AddAttachment(attachment string) {
	if b.Attachments == nil {
		b.Attachments = make([]string, 0)
	}
	b.Attachments = append(b.Attachments, attachment)
}

func (b *Bulk) CreateCSVString() (string, error) {
	// ユニークなキーを収集
	keySet := make(map[string]struct{})
	for _, recipient := range b.To {
		for key := range recipient.InsertCode {
			keySet[key] = struct{}{}
		}
	}

	// キーをスライスに変換（順序は不定）
	keys := []string{"email"}
	headers := []string{"email"}
	for key := range keySet {
		headers = append(keys, "__"+key+"__")
		keys = append(keys, key)
	}

	// CSV文字列を作成
	var sb strings.Builder
	writer := csv.NewWriter(&sb)

	// ヘッダーを書き込み
	if err := writer.Write(headers); err != nil {
		return "", err
	}

	// データを書き込み
	for _, recipient := range b.To {
		row := []string{recipient.Email}
		for _, key := range keys[1:] {
			if value, exists := recipient.InsertCode[key]; exists {
				row = append(row, value)
			} else {
				row = append(row, "") // キーがない場合は空文字
			}
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	// `Flush()` は `defer` ではなく、エラー処理の影響を受けないように明示的に呼び出す
	writer.Flush()

	// `Flush()` の後にエラーを確認する
	if err := writer.Error(); err != nil {
		return "", err
	}

	return sb.String(), nil
}

func (b *Bulk) Begin() error {
	url := "https://app.engn.jp/api/v1/deliveries/bulk/begin"
	beginBulk := struct {
		From            MailAddress            `json:"from"`
		Subject         string                 `json:"subject"`
		Encode          string                 `json:"encode"`
		TextPart        string                 `json:"text_part"`
		HtmlPart        string                 `json:"html_part,omitempty"`
		ListUnsubscribe *ListUnsubscribeParams `json:"list_unsubscribe,omitempty"`
	}{
		From:            b.From,
		Subject:         b.Subject,
		Encode:          b.Encode,
		TextPart:        b.TextPart,
		HtmlPart:        b.HtmlPart,
		ListUnsubscribe: b.ListUnsubscribe,
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
	b.Mail.DeliveryId = response.DeliveryId
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
		From            MailAddress            `json:"from"`
		To              []JsonBulkTo           `json:"to"`
		Subject         string                 `json:"subject"`
		TextPart        string                 `json:"text_part"`
		HtmlPart        string                 `json:"html_part,omitempty"`
		ListUnsubscribe *ListUnsubscribeParams `json:"list_unsubscribe,omitempty"`
	}{
		From:            b.From,
		To:              bulkTos,
		Subject:         b.Subject,
		TextPart:        b.TextPart,
		HtmlPart:        b.HtmlPart,
		ListUnsubscribe: b.ListUnsubscribe,
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
	b.Mail.DeliveryId = response.DeliveryId
	return nil
}

func (b *Bulk) Import(params *ImportParams) (*Job, error) {
	if b.DeliveryId == 0 {
		return nil, fmt.Errorf("delivery ID is required")
	}
	url := fmt.Sprintf("https://app.engn.jp/api/v1/deliveries/%d/emails/import", b.DeliveryId)
	csvString, err := b.CreateCSVString()
	if err != nil {
		return nil, fmt.Errorf("failed to create CSV string: %v", err)
	}
	reader := strings.NewReader(csvString)
	attachments := []Attachment{
		{
			FileName: "import.csv",
			Content:  reader,
		},
	}
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal import params: %v", err)
	}
	responseBytes, err := b.Client.sendRequest("POST", url, nil, jsonData, true, attachments)
	if err != nil {
		return nil, err
	}
	var job Job
	err = json.Unmarshal(responseBytes, &job)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	job.Client = b.Client
	return &job, nil
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
