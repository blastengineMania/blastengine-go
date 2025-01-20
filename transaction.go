package blastenginego

type Transaction struct {
	From       MailAddress
	To         string
	Cc         []string
	Bcc        []string
	InsertCode map[string]string
	Subject    string
	Encode     string
	TextPart   string
	HtmlPart   string
}

func NewTransaction() *Transaction {
	return &Transaction{
		Encode: "UTF-8",
	}
}

func (t *Transaction) SetFrom(email, name string) {
	t.From = NewMailAddress(email, name)
}

func (t *Transaction) SetTo(to string) {
	t.To = to
}

func (t *Transaction) SetCc(cc []string) {
	t.Cc = cc
}

func (t *Transaction) SetBcc(bcc []string) {
	t.Bcc = bcc
}

func (t *Transaction) SetInsertCode(key string, value string) {
	if t.InsertCode == nil {
		t.InsertCode = make(map[string]string)
	}
	t.InsertCode[key] = value
}

func (t *Transaction) SetSubject(subject string) {
	t.Subject = subject
}

func (t *Transaction) SetEncode(encode string) {
	t.Encode = encode
}

func (t *Transaction) SetTextPart(textPart string) {
	t.TextPart = textPart
}

func (t *Transaction) SetHtmlPart(htmlPart string) {
	t.HtmlPart = htmlPart
}
