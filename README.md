# blastengine-go
Blastengine Golang SDK

## Initialization Method

To initialize the Blastengine client, you need to provide your API key and user ID. Here is an example:

```go
package main

import (
	"fmt"
	"github.com/blastengineMania/blastengine-go"
)

func main() {
	apiKey := "yourApiKey"
	userId := "yourUserId"
	client := blastengine.Initialize(apiKey, userId)
	fmt.Println("Client initialized:", client)
}
```

## Transaction Usage

To use transactions with the Blastengine client, you can follow this example:

```go
package main

import (
	"fmt"
	"github.com/blastengineMania/blastengine-go"
)

func main() {
	apiKey := "yourApiKey"
	userId := "yourUserId"
	client := blastengine.Initialize(apiKey, userId)

	transaction := client.NewTransaction()
	transaction.SetFrom("from@example.com", "Sender Name")
	transaction.SetTo("to@example.com")
	transaction.SetSubject("Test Subject")
	transaction.SetTextPart("This is a text part")
	transaction.SetHtmlPart("<p>This is an HTML part</p>")

	err := transaction.Send()
	if err != nil {
		fmt.Println("Failed to send transaction:", err)
	} else {
		fmt.Println("Transaction sent successfully")
	}
}
```

## Sending Email with Attachments

To send an email with attachments using the `Transaction` struct, you can use the `AddAttachment` method. Here is an example:

```go
package main

import (
	"fmt"
	"github.com/blastengineMania/blastengine-go"
)

func main() {
	apiKey := "yourApiKey"
	userId := "yourUserId"
	client := blastengine.Initialize(apiKey, userId)

	transaction := client.NewTransaction()
	transaction.SetFrom("from@example.com", "Sender Name")
	transaction.SetTo("to@example.com")
	transaction.SetSubject("Test Subject")
	transaction.SetTextPart("This is a text part")
	transaction.SetHtmlPart("<p>This is an HTML part</p>")
	transaction.AddAttachment("path/to/attachment1")
	transaction.AddAttachment("path/to/attachment2")

	err := transaction.Send()
	if err != nil {
		fmt.Println("Failed to send transaction:", err)
	} else {
		fmt.Println("Transaction sent successfully with attachments")
	}
}
```

## License

MIT
