package sms

import "testing"

// AWS 短信测试
func TestAwsSms_Send(t *testing.T) {
	aws := NewAwsSms("xxx", "xxxxxxx", "ap-southeast-1")
	sms := Sms{
		Numbers: "+886xxxxxxxxx",
		Content: "【XXX】您正在[註冊]，驗證碼：355621，驗證碼僅用於XXX官網，請勿洩露。驗證碼10分鐘內有效。",
	}
	aws.Send(sms)
}

// 天一鸿短信测试
func TestBukaSms_Send(t *testing.T) {
	buka := NewBuka("xxx", "xxx", "xxx")
	sms := Sms{
		Numbers: "91xxxxxxxxx",
		Content: "Thank you for registering. Verification code: 076833 (valid for 10 mins). The verification code is for our website only.",
	}
	buka.Send(sms)
}
