package blastengine

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestBulkSetFrom(t *testing.T) {
	client := getClient()
	bulk := client.NewBulk()
	email := "test@example.com"
	name := "Test User"
	bulk.SetFrom(email, name)

	if bulk.From.Email != email {
		t.Errorf("Expected From.Email to be %s, but got %s", email, bulk.From.Email)
	}

	if bulk.From.Name != name {
		t.Errorf("Expected From.Name to be %s, but got %s", name, bulk.From.Name)
	}
}

func TestBulkSetSubject(t *testing.T) {
	client := getClient()
	bulk := client.NewBulk()
	subject := "Test subject"
	bulk.SetSubject(subject)

	if bulk.Subject != subject {
		t.Errorf("Expected Subject to be %s, but got %s", subject, bulk.Subject)
	}
}

func TestBulkSetTo(t *testing.T) {
	client := getClient()
	bulk := client.NewBulk()
	to := "user@example.jp"
	insertCode := map[string]string{
		"key": "value",
	}
	bulk.AddTo(to, insertCode)
	bulk.AddTo(to, nil)
	if bulk.To == nil {
		t.Errorf("Expected To to be not nil, but got nil")
	}
	if bulk.To[0].Email != to {
		t.Errorf("Expected To to be %s, but got %s", to, bulk.To[0].Email)
	}
	if bulk.To[0].InsertCode["key"] != "value" {
		t.Errorf("Expected InsertCode to be %s, but got %s", "value", bulk.To[0].InsertCode["key"])
	}
}

func TestBulkSetEncode(t *testing.T) {
	client := getClient()
	bulk := client.NewBulk()
	encode := "ISO-8859-1"
	bulk.SetEncode(encode)

	if bulk.Encode != encode {
		t.Errorf("Expected Encode to be %s, but got %s", encode, bulk.Encode)
	}

	// Test default value
	defaultTransaction := client.NewBulk()
	if defaultTransaction.Encode != "UTF-8" {
		t.Errorf("Expected default Encode to be UTF-8, but got %s", defaultTransaction.Encode)
	}
	fmt.Println(defaultTransaction.Encode)
}

func TestBulkSetTextPart(t *testing.T) {
	client := getClient()
	bulk := client.NewBulk()
	textPart := "This is a text part"
	bulk.SetTextPart(textPart)

	if bulk.TextPart != textPart {
		t.Errorf("Expected TextPart to be %s, but got %s", textPart, bulk.TextPart)
	}
}

func TestBulkSetHtmlPart(t *testing.T) {
	client := getClient()

	bulk := client.NewBulk()
	htmlPart := "<p>This is an HTML part</p>"
	bulk.SetHtmlPart(htmlPart)

	if bulk.HtmlPart != htmlPart {
		t.Errorf("Expected HtmlPart to be %s, but got %s", htmlPart, bulk.HtmlPart)
	}
}

func TestBulkBegin(t *testing.T) {
	client := getClient()
	bulk := client.NewBulk()
	email := "test@example.com"
	name := "Test User"
	bulk.SetFrom(email, name)
	bulk.SetSubject("Test subject")
	bulk.SetTextPart("This is a text part")
	err := bulk.Begin()
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
	if bulk.DeliveryId == 0 {
		t.Errorf("Expected DeliveryId to be not 0, but got 0")
	}
	err = bulk.Get()
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
	if bulk.DeliveryId == 0 {
		t.Errorf("Expected DeliveryId to be not 0, but got 0")
	}
	if bulk.Status != "EDIT" {
		t.Errorf("Expected Status to be sent, but got %s", bulk.Status)
	}
	if bulk.DeliveryType != "BULK" {
		t.Errorf("Expected deliveryType to be TRANSACTION, but got %s", bulk.DeliveryType)
	}
	if bulk.CreatedTime.IsZero() {
		t.Errorf("Expected createdTime to be zero, but got %v", bulk.CreatedTime)
	}
	if !bulk.ReservationTime.IsZero() {
		t.Errorf("Expected reservationTime to be not zero, but got zero")
	}
	if bulk.UpdatedTime.IsZero() {
		t.Errorf("Expected updatedTime to be not zero, but got zero")
	}
	err = bulk.Delete()
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
}

func TestBulkUpdate(t *testing.T) {
	client := getClient()
	bulk := client.NewBulk()
	email := "test@example.com"
	name := "Test User"
	bulk.SetFrom(email, name)
	bulk.SetSubject("Test subject")
	bulk.SetTextPart("This is a text part")
	err := bulk.Begin()
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
	to := "user@example.jp"
	insertCode := map[string]string{
		"key": "value",
	}
	bulk.AddTo(to, insertCode)
	err = bulk.Update()
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
	err = bulk.Get()
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
	bulk.Delete()
}

func TestBulkSendNow(t *testing.T) {
	client := getClient()
	bulk := client.NewBulk()
	email := os.Getenv("FROM")
	name := "Test User"
	bulk.SetFrom(email, name)
	bulk.SetSubject("Test subject __key__")
	bulk.SetTextPart("This is a text part __key__")
	err := bulk.Begin()
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
	to := "atsushi+be@moongift.co.jp"
	insertCode := map[string]string{
		"key": "001",
	}
	bulk.AddTo(to, insertCode)
	err = bulk.Update()
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
	err = bulk.Send(nil)
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
	err = bulk.Get()
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
}

func TestBulkSendNowListUnsubscribe(t *testing.T) {
	client := getClient()
	bulk := client.NewBulk()
	email := os.Getenv("FROM")
	name := "Test User"
	bulk.SetFrom(email, name)
	bulk.SetSubject("Test subject __key__")
	bulk.SetTextPart("This is a text part __key__")
	bulk.SetListUnsubscribe(&ListUnsubscribeParams{
		Url: "https://example.com/unsubscribe",
	})
	err := bulk.Begin()
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
	to := "atsushi+be@moongift.co.jp"
	insertCode := map[string]string{
		"key": "001",
	}
	bulk.AddTo(to, insertCode)
	err = bulk.Update()
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
	err = bulk.Send(nil)
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
	err = bulk.Get()
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
}

func TestGenerateCSVString(t *testing.T) {
	client := getClient()
	bulk := client.NewBulk()
	to := "atsushi+be@moongift.co.jp"
	insertCode := map[string]string{
		"key": "001",
	}
	bulk.AddTo(to, insertCode)
	csvString, err := bulk.CreateCSVString()
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
	if csvString == "" {
		t.Errorf("Expected csvString to be not empty, but got empty")
	}
	fmt.Print(csvString)
}

func TestBulkImport(t *testing.T) {
	client := getClient()
	bulk := client.NewBulk()
	email := os.Getenv("FROM")
	name := "Test User"
	bulk.SetFrom(email, name)
	bulk.SetSubject("Test subject __key__")
	bulk.SetTextPart("This is a text part __key__")
	err := bulk.Begin()
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
	to := "atsushi+be@moongift.co.jp"
	insertCode := map[string]string{
		"key": "001",
	}
	bulk.AddTo(to, insertCode)
	job, err := bulk.Import(nil)
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
	for {
		b, err := job.Finished()
		if err != nil {
			t.Errorf("Expected error to be nil, but got %v", err)
		}
		if b {
			break
		}
		time.Sleep(1 * time.Second)
	}
	err = bulk.Send(nil)
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
}

func TestBulkImportError(t *testing.T) {
	client := getClient()
	bulk := client.NewBulk()
	email := os.Getenv("FROM")
	name := "Test User"
	bulk.SetFrom(email, name)
	bulk.SetSubject("Test subject __key__")
	bulk.SetTextPart("This is a text part __key__")
	err := bulk.Begin()
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
	to := "atsushi+be@moongift.co.jp"
	insertCode := map[string]string{
		"key": "001",
	}
	bulk.AddTo(to, insertCode)
	to2 := "atsushi+be@"
	insertCode2 := map[string]string{
		"key": "001",
	}
	bulk.AddTo(to2, insertCode2)
	job, err := bulk.Import(&ImportParams{
		IgnoreErrors: false,
		Immediate:    true,
	})
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
	for {
		b, err := job.Finished()
		if err != nil {
			t.Errorf("Expected error to be nil, but got %v", err)
		}
		if b {
			break
		}
		time.Sleep(1 * time.Second)
	}
	s, err := job.Download()
	if err != nil {
		t.Errorf("Expected error to be nil, but got %v", err)
	}
	fmt.Print(s)
	bulk.Delete()
}
