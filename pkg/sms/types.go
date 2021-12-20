package sms

type ISms interface {
	Send(entity Sms) error
}

type Sms struct {
	Numbers  string // 接收号码
	Content  string // 消息内容
	SenderId string // 消息内容
}
