package blastengine

import (
	"encoding/json"
	"fmt"
	"time"
)

type MailAttachment struct {
	DeliveryAttachID int       `json:"delivery_attach_id"`
	FileName         string    `json:"file_name"`
	Mime             string    `json:"mime"`
	Size             int       `json:"size"`
	CreatedTime      time.Time `json:"created_time"`
	UpdatedTime      time.Time `json:"updated_time"`
}

type Mail struct {
	Client *Client

	DeliveryId      int         `json:"delivery_id"`
	From            MailAddress `json:"from"`
	Status          string      `json:"status"`
	DeliveryTime    time.Time   `json:"delivery_time"`
	UpdatedTime     time.Time   `json:"updated_time"`
	CreatedTime     time.Time   `json:"created_time"`
	ReservationTime time.Time   `json:"reservation_time"`
	TextPart        string      `json:"text_part"`
	HtmlPart        string      `json:"html_part"`
	To              []BulkTo
	Cc              []string
	Bcc             []string
	Encode          string
	ListUnsubscribe *ListUnsubscribeParams
	Attachments     []string
	DeliveryType    string           `json:"delivery_type"`
	Subject         string           `json:"subject"`
	Attaches        []MailAttachment `json:"attaches"`
	OpenCount       int              `json:"open_count"`
	TotalCount      int              `json:"total_count"`
	SentCount       int              `json:"sent_count"`
	DropCount       int              `json:"drop_count"`
	SoftErrorCount  int              `json:"soft_error_count"`
	HardErrorCount  int              `json:"hard_error_count"`
}

func (m *Mail) SetFrom(email, name string) {
	m.From = NewMailAddress(email, name)
}

func (m *Mail) AddTo(to string, insertCode map[string]string) {
	m.To = append(m.To, BulkTo{Email: to, InsertCode: insertCode})
}

func (m *Mail) AddCc(cc string) {
	if m.Cc == nil {
		m.Cc = make([]string, 0)
	}
	m.Cc = append(m.Cc, cc)
}

func (m *Mail) AddBcc(bcc string) {
	if m.Bcc == nil {
		m.Bcc = make([]string, 0)
	}
	m.Bcc = append(m.Bcc, bcc)
}

func (m *Mail) SetSubject(subject string) {
	m.Subject = subject
}

func (m *Mail) SetEncode(encode string) {
	m.Encode = encode
}

func (m *Mail) SetTextPart(textPart string) {
	m.TextPart = textPart
}

func (m *Mail) SetHtmlPart(htmlPart string) {
	m.HtmlPart = htmlPart
}

func (m *Mail) SetListUnsubscribe(params *ListUnsubscribeParams) {
	if params.Email != "" {
		params.Email = fmt.Sprintf("mailto:%s", params.Email)
	}
	m.ListUnsubscribe = &ListUnsubscribeParams{
		Email: params.Email,
		Url:   params.Url,
	}
}

func (m *Mail) AddAttachment(attachment string) {
	if m.Attachments == nil {
		m.Attachments = make([]string, 0)
	}
	m.Attachments = append(m.Attachments, attachment)
}

func (m *Mail) Get() error {
	if m.DeliveryId == 0 {
		return fmt.Errorf("delivery ID is required")
	}
	url := fmt.Sprintf("https://app.engn.jp/api/v1/deliveries/%d", m.DeliveryId)
	bodyBytes, err := m.Client.sendRequest("GET", url, nil, nil, false, nil)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyBytes, &m)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}
	return nil
}

func (m *Mail) Send(reservationTime *time.Time) error {
	if len(m.To) == 0 {
		return fmt.Errorf("to is required")
	}
	if len(m.Bcc) > 0 || len(m.Cc) > 0 {
		// Transaction
		if reservationTime != nil {
			return fmt.Errorf("cc and bcc are not supported in reservation time")
		}
		if len(m.To) > 1 {
			return fmt.Errorf("multiple to is not supported in transaction")
		}
		return m.SendTransaction()
	} else {
		// Bulk
		return m.SendBulk(reservationTime)
	}
}

func (m *Mail) SendTransaction() error {
	transaction := m.Client.NewTransaction()
	transaction.SetSubject(m.Subject)
	transaction.SetEncode(m.Encode)
	transaction.SetTextPart(m.TextPart)
	transaction.SetHtmlPart(m.HtmlPart)
	transaction.SetFrom(m.From.Email, m.From.Name)
	transaction.SetListUnsubscribe(m.ListUnsubscribe)
	for _, attachment := range m.Attachments {
		transaction.AddAttachment(attachment)
	}
	for _, to := range m.To {
		transaction.AddTo(to.Email, to.InsertCode)
	}
	for _, cc := range m.Cc {
		transaction.AddCc(cc)
	}
	for _, bcc := range m.Bcc {
		transaction.AddBcc(bcc)
	}
	err := transaction.Send()
	if err != nil {
		return err
	}
	m.DeliveryId = transaction.DeliveryId
	return nil
}

func (m *Mail) SendBulk(reservationTime *time.Time) error {
	bulk := m.Client.NewBulk()
	bulk.SetSubject(m.Subject)
	bulk.SetEncode(m.Encode)
	bulk.SetTextPart(m.TextPart)
	bulk.SetHtmlPart(m.HtmlPart)
	bulk.SetListUnsubscribe(m.ListUnsubscribe)
	bulk.SetFrom(m.From.Email, m.From.Name)
	for _, attachment := range m.Attachments {
		bulk.AddAttachment(attachment)
	}
	err := bulk.Begin()
	if err != nil {
		return err
	}
	for _, to := range m.To {
		bulk.AddTo(to.Email, to.InsertCode)
	}
	if len(m.To) >= 50 {
		job, err := bulk.Import(&ImportParams{
			IgnoreErrors: false,
			Immediate:    false,
		})
		if err != nil {
			return err
		}
		for {
			b, err := job.Finished()
			if err != nil {
				return err
			}
			if b {
				break
			}
			time.Sleep(1 * time.Second)
		}
	} else {
		err = bulk.Update()
		if err != nil {
			return err
		}
	}
	err = bulk.Send(reservationTime)
	if err != nil {
		return err
	}
	m.DeliveryId = bulk.DeliveryId
	return nil
}
