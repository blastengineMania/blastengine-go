package blastengine

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// 送信先Toが1件の場合のテスト
func TestMailSendWithSingleTo(t *testing.T) {
	client := getClient()
	mail := client.NewMail()

	// 送信元の設定
	mail.SetFrom(os.Getenv("FROM"), "Test User")

	// 送信先の設定（1件）
	mail.AddTo(os.Getenv("TO"), nil)

	// 件名と本文の設定
	mail.SetSubject("Test Mail with Single To")
	mail.SetTextPart("This is a test mail with single To")

	// ListUnsubscribeの設定（nilポインタエラー回避のため）
	mail.SetListUnsubscribe(&ListUnsubscribeParams{
		Url: "https://example.com/unsubscribe",
	})

	// メール送信
	err := mail.Send(nil)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// 送信結果の確認
	if mail.DeliveryId == 0 {
		t.Errorf("Expected DeliveryId to be not 0, but got 0")
	}

	// メール情報の取得
	err = mail.Get()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// トランザクションメールとして送信されたことを確認
	if mail.DeliveryType != "BULK" {
		t.Errorf("Expected DeliveryType to be BULK, but got %s", mail.DeliveryType)
	}

	// ステータスの確認
	if mail.Status != "SENT" && mail.Status != "WAIT" {
		t.Errorf("Expected Status to be SENT, but got %s", mail.Status)
	}
}

// 送信先Toが50件以上の場合のテスト
func TestMailSendWithManyTo(t *testing.T) {
	client := getClient()
	mail := client.NewMail()

	// 送信元の設定
	mail.SetFrom(os.Getenv("FROM"), "Test User")

	// 送信先の設定（50件以上）
	baseEmail := "atsushi+%d@moongift.co.jp"
	for i := 0; i < 55; i++ {
		mail.AddTo(fmt.Sprintf(baseEmail, i), map[string]string{
			"id": fmt.Sprintf("%03d", i),
		})
	}

	// 件名と本文の設定
	mail.SetSubject("Test Mail with Many To __id__")
	mail.SetTextPart("This is a test mail with many To. Your ID is __id__")

	// ListUnsubscribeの設定（nilポインタエラー回避のため）
	mail.SetListUnsubscribe(&ListUnsubscribeParams{
		Url: "https://example.com/unsubscribe",
	})

	// メール送信
	err := mail.Send(nil)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// 送信結果の確認
	if mail.DeliveryId == 0 {
		t.Errorf("Expected DeliveryId to be not 0, but got 0")
	}

	// メール情報の取得
	err = mail.Get()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// バルクメールとして送信されたことを確認
	if mail.DeliveryType != "BULK" {
		t.Errorf("Expected DeliveryType to be BULK, but got %s", mail.DeliveryType)
	}

	// 送信先数の確認
	if mail.TotalCount != 55 {
		t.Errorf("Expected TotalCount to be 55, but got %d", mail.TotalCount)
	}
}

// 送信先Toが50件以上、添付ファイルありの場合のテスト
func TestMailSendWithManyToAndAttachment(t *testing.T) {
	client := getClient()
	mail := client.NewMail()

	// 送信元の設定
	mail.SetFrom(os.Getenv("FROM"), "Test User")

	// 送信先の設定（50件以上）
	baseEmail := "test%d@example.com"
	for i := 0; i < 55; i++ {
		mail.AddTo(fmt.Sprintf(baseEmail, i), map[string]string{
			"id": fmt.Sprintf("%03d", i),
		})
	}

	// 件名と本文の設定
	mail.SetSubject("Test Mail with Many To and Attachment __id__")
	mail.SetTextPart("This is a test mail with many To and attachment. Your ID is __id__")

	// ListUnsubscribeの設定（nilポインタエラー回避のため）
	mail.SetListUnsubscribe(&ListUnsubscribeParams{
		Url: "https://example.com/unsubscribe",
	})

	// 添付ファイルの追加
	mail.AddAttachment("README.md")
	mail.AddAttachment("LICENSE")

	// メール送信
	err := mail.Send(nil)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// 送信結果の確認
	if mail.DeliveryId == 0 {
		t.Errorf("Expected DeliveryId to be not 0, but got 0")
	}

	// メール情報の取得
	err = mail.Get()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// バルクメールとして送信されたことを確認
	if mail.DeliveryType != "BULK" {
		t.Errorf("Expected DeliveryType to be BULK, but got %s", mail.DeliveryType)
	}

	// 送信先数の確認
	if mail.TotalCount != 55 {
		t.Errorf("Expected TotalCount to be 55, but got %d", mail.TotalCount)
	}

	// 添付ファイルの確認
	if len(mail.Attaches) != 2 {
		t.Errorf("Expected 2 attachments, but got %d", len(mail.Attaches))
	}
}

// 送信先が2件、BCCありでエラーになることを確認するテスト
func TestMailSendWithMultipleToAndBcc(t *testing.T) {
	client := getClient()
	mail := client.NewMail()

	// 送信元の設定
	mail.SetFrom(os.Getenv("FROM"), "Test User")

	// 送信先の設定（2件）
	mail.AddTo("test1@example.com", nil)
	mail.AddTo("test2@example.com", nil)

	// BCCの追加
	mail.AddBcc("bcc@example.com")

	// 件名と本文の設定
	mail.SetSubject("Test Mail with Multiple To and Bcc")
	mail.SetTextPart("This is a test mail with multiple To and Bcc")

	// ListUnsubscribeの設定（nilポインタエラー回避のため）
	mail.SetListUnsubscribe(&ListUnsubscribeParams{
		Url: "https://example.com/unsubscribe",
	})

	// メール送信（エラーが発生することを期待）
	err := mail.Send(nil)

	// エラーが発生することを確認
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}

	// エラーメッセージの確認
	expectedError := "multiple to is not supported in transaction"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', but got '%s'", expectedError, err.Error())
	}
}

// 予約送信のテスト
func TestMailSendWithReservation(t *testing.T) {
	client := getClient()
	mail := client.NewMail()

	// 送信元の設定
	mail.SetFrom(os.Getenv("FROM"), "Test User")

	// 送信先の設定（複数件）
	mail.AddTo("test1@example.com", nil)
	mail.AddTo("test2@example.com", nil)

	// 件名と本文の設定
	mail.SetSubject("Test Mail with Reservation")
	mail.SetTextPart("This is a test mail with reservation")

	// ListUnsubscribeの設定（nilポインタエラー回避のため）
	mail.SetListUnsubscribe(&ListUnsubscribeParams{
		Url: "https://example.com/unsubscribe",
	})

	// 予約時間の設定（1時間後）
	reservationTime := time.Now().Add(1 * time.Hour)

	// メール送信
	err := mail.Send(&reservationTime)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// 送信結果の確認
	if mail.DeliveryId == 0 {
		t.Errorf("Expected DeliveryId to be not 0, but got 0")
	}

	// メール情報の取得
	err = mail.Get()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// バルクメールとして送信されたことを確認
	if mail.DeliveryType != "BULK" {
		t.Errorf("Expected DeliveryType to be BULK, but got %s", mail.DeliveryType)
	}

	// 予約時間が設定されていることを確認
	if mail.ReservationTime.IsZero() {
		t.Errorf("Expected ReservationTime to be not zero, but got zero")
	}

	// ステータスの確認
	if mail.Status != "RESERVE" {
		t.Errorf("Expected Status to be RESERVE, but got %s", mail.Status)
	}
}

// トランザクションメールで予約送信するとエラーになることを確認するテスト
func TestMailSendTransactionWithReservation(t *testing.T) {
	client := getClient()
	mail := client.NewMail()

	// 送信元の設定
	mail.SetFrom(os.Getenv("FROM"), "Test User")

	// 送信先の設定（1件）
	mail.AddTo(os.Getenv("TO"), nil)
	mail.AddCc("test@example.com")

	// 件名と本文の設定
	mail.SetSubject("Test Transaction Mail with Reservation")
	mail.SetTextPart("This is a test transaction mail with reservation")

	// ListUnsubscribeの設定（nilポインタエラー回避のため）
	mail.SetListUnsubscribe(&ListUnsubscribeParams{
		Url: "https://example.com/unsubscribe",
	})

	// 予約時間の設定（1時間後）
	reservationTime := time.Now().Add(1 * time.Hour)

	// メール送信（エラーが発生することを期待）
	err := mail.Send(&reservationTime)

	// エラーが発生することを確認
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}

	// エラーメッセージの確認
	expectedError := "cc and bcc are not supported in reservation time"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', but got '%s'", expectedError, err.Error())
	}
}
