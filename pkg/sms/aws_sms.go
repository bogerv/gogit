package sms

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type AwsSms struct {
	accessKey string // 名称
	secretKey string // 私钥
	region    string // 区域
}

func NewAwsSms(accessKey, secretKey, region string) *AwsSms {
	return &AwsSms{
		accessKey: accessKey,
		secretKey: secretKey,
		region:    region,
	}
}

func (a *AwsSms) Send(sms Sms) error {
	creds := credentials.NewStaticCredentials(a.accessKey, a.secretKey, "")
	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(a.region),
	})
	if err != nil {
		return err
	}
	client := sns.New(sess)

	pin := &sns.PublishInput{}
	pin.SetMessage(sms.Content)
	pin.SetPhoneNumber(sms.Numbers)
	result, err := client.Publish(pin)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %s", result.String())
	return nil
}
