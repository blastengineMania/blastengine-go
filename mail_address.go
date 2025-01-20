package blastenginego

type MailAddress struct {
	Email string
	Name  string
}

func NewMailAddress(email, name string) MailAddress {
	return MailAddress{
		Email: email,
		Name:  name,
	}
}
