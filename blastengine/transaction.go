package blastengine

type Transaction struct {
    Subject  string
    TextPart string
    HtmlPart string
}

func (t *Transaction) SetSubject(subject string) {
    t.Subject = subject
}

func (t *Transaction) SetTextPart(textPart string) {
    t.TextPart = textPart
}

func (t *Transaction) SetHtmlPart(htmlPart string) {
    t.HtmlPart = htmlPart
}

func (t *Transaction) Send() {
    // Implementation for sending the transaction email
}
