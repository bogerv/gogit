package mailx

import (
	"fmt"
	"testing"
)

func TestAwsMailV2_Send(t *testing.T) {
	content := NewMail()
	content.From = "XXX <xxx@xxx.com>"
	content.To = "xxx@xxx.com"
	content.Subject = "Amazon SES Test"
	content.Body = "This email was sent with Amazon SES using the AWS SDK for Go."

	fmt.Println(NewAwsMailV2().Send(content))
}

func TestAwsMail_Send(t *testing.T) {
	entity := NewMail()
	entity.SmtpHost = "email-smtp.ap-southeast-1.amazonaws.com"
	entity.SmtpPort = 587
	entity.SmtpUser = "xxx" // smtp user
	entity.SmtpKey = "xxx"  // smtp key

	entity.From = "XXX <xxx@notify.xxx.com>"
	entity.To = "xxx@xxx.com"
	entity.Subject = "Amazon SES Test"
	entity.Body = "This email was sent with Amazon SES using the AWS SDK for Go."

	fmt.Println(NewAwsMail().Send(entity))
}
