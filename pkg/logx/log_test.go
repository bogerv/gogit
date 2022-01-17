package logx

import "testing"

func TestMain(m *testing.M) {
	defer Flush()
	Init("./logs/sms.log")
	m.Run()
}

func TestDebug(t *testing.T) {
	Debug("test")
	Debug("test1")
	Debug("test2")
	Debug("test3")
}

func TestInfo(t *testing.T) {
	Info("test")
}

func TestWarn(t *testing.T) {
	for i := 0; i < 10000; i++ {
		Warn("test")
	}
}

func BenchmarkWarn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Warn("test")
	}
}

func TestError(t *testing.T) {
	Error("test")
}
