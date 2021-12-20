package mailx

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

var smtpName = "xxx" // smtp name
var smtpKey = "xxx"  // smtp key

type AwsMailV2 struct {
}

func NewAwsMailV2() *AwsMailV2 {
	return &AwsMailV2{}
}

func (m *AwsMailV2) Send(entity *Mail) error {
	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		//Region: aws.String("us-east-1"),
		//Endpoint: aws.String("email-smtp.us-east-1.amazonaws.com"),
		Region: aws.String("ap-southeast-1"),
		//Endpoint: aws.String("email-smtp.ap-southeast-1.amazonaws.com"),
	})

	creds := credentials.NewStaticCredentials(smtpName, smtpKey, "") // Token is empty
	// Create an SES session.
	svc := ses.New(sess, &aws.Config{Credentials: creds})

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(entity.To),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				//Html: &ses.Content{
				//	Charset: aws.String(content.CharSet),
				//	Data:    aws.String(HtmlBody),
				//},
				Text: &ses.Content{
					Charset: aws.String(entity.CharSet),
					Data:    aws.String(entity.Body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(entity.CharSet),
				Data:    aws.String(entity.Subject),
			},
		},
		Source: aws.String(entity.From),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String("XXXX"),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			fmt.Println(err.Error())
		}
		return err
	}

	fmt.Printf("Email Sent to address: %s result:%+v\n", entity.To, result)
	return nil
}
