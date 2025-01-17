# blastengine-go
Blastengine Golang SDK

## クライアントの初期化

クライアントを初期化するには、ユーザーIDとAPIキーを使用して`Initialize`関数を使用します:

```go
package main

import (
    "fmt"
    "blastengine"
)

func main() {
    userID := "your_user_id"
    apiKey := "your_api_key"
    client := blastengine.Initialize(userID, apiKey)
    fmt.Println("クライアントが初期化されました。UserID:", client.UserID)
}
```

## トランザクションの使用

トランザクションを使用するには、まずクライアントを初期化し、次に`Transaction`メソッドを使用してトランザクションを作成します。トランザクションの件名、テキスト部分、およびHTML部分を設定し、それを送信します:

```go
package main

import (
    "blastengine"
)

func main() {
    client := blastengine.Initialize("your_user_id", "your_api_key")
    transaction := client.Transaction()

    transaction.SetSubject("テストメールの件名")
    transaction.SetTextPart("テストメールの本文")
    transaction.SetHtmlPart("<h1>テストメールの本文</h1>")

    // トランザクションメールを送信
    transaction.Send()
}
```
