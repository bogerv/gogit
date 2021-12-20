package mailx

import (
	"crypto/tls"
	"github.com/go-gomail/gomail"
)

type AwsMail struct {
}

func NewAwsMail() *AwsMail {
	return &AwsMail{}
}

func (m *AwsMail) Send(entity *Mail) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", entity.From)
	msg.SetHeader("To", entity.To)
	msg.SetHeader("Subject", entity.Subject)
	// Set E-Mail body. You can set text/plain or html with text/html
	msg.SetBody("text/html", entity.Body)

	client := gomail.NewDialer(entity.SmtpHost, entity.SmtpPort, entity.SmtpUser, entity.SmtpKey)
	client.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := client.DialAndSend(msg); err != nil {
		// TODO print log
		return err
	}
	return nil
}
