package blastengine

type MailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func NewMailAddress(email, name string) MailAddress {
	return MailAddress{
		Email: email,
		Name:  name,
	}
}
