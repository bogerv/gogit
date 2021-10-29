package timex

import "time"

const OneHourSec = 3600

func UnixNowHour() int64 {
	now := time.Now().Unix()
	return now - (now % OneHourSec)
}
