package mailx

type IMail interface {
	Send(entity Mail) error
}

type Mail struct {
	From     string // from address
	FromName string // from name
	SmtpHost string // smtp host
	SmtpPort int    // smtp port
	SmtpUser string // smtp user
	SmtpKey  string // smtp key
	Pwd      string // mail pwd
	To       string // to address
	Subject  string // mail subject
	Body     string // mail body
	CharSet  string // char set
}

func NewMail() *Mail {
	content := &Mail{
		CharSet: "UTF-8",
	}
	return content
}
