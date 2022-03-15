package sms

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"gogit/pkg/httpx"
	"gogit/pkg/signx"
)

// Buka 天一鸿短信发送
type Buka struct {
	accessKey     string // 名称
	secretKey     string // 私钥
	applicationId string // 应用ID
}

func NewBuka(accessKey, secretKey, applicationId string) *Buka {
	return &Buka{
		accessKey:     accessKey,
		secretKey:     secretKey,
		applicationId: applicationId,
	}
}

func (b *Buka) Send(sms Sms) error {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	msg := fmt.Sprintf("%s%s%s", b.accessKey, b.secretKey, timestamp)
	log.Println(msg)
	sign := signx.MD5(msg)
	log.Println(sign)

	req := httpx.New()
	req.AddHeader("Sign", sign)
	req.AddHeader("Timestamp", timestamp)
	req.AddHeader("Api-Key", b.accessKey)

	body := map[string]string{
		"appId":    b.applicationId,
		"numbers":  sms.Numbers,
		"content":  sms.Content,
		"senderId": sms.SenderId,
	}
	bytes, err := req.Post("https://api.onbuka.com/v3/sendSms", body)
	if err != nil {
		log.Printf("sms send fail:%s\n", err)
		return err
	}
	log.Println(string(bytes))
	return nil
}
