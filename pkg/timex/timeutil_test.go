package timex

import (
	"fmt"
	"testing"
)

func TestNowHour(t *testing.T) {
	fmt.Println(UnixNowHour())
}
