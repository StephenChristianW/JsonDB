package UtilsTime

import (
	"time"
)

func TimeStampToString(stamp int64) string {
	var t time.Time
	switch {
	case stamp > 1e18: // 纳秒
		t = time.Unix(0, stamp)
	case stamp > 1e14: // 毫秒
		t = time.UnixMilli(stamp)
	default: // 秒
		t = time.Unix(stamp, 0)
	}
	return t.Format("2006-01-02 15:04:05.000")
}

func TimeNow() string {
	return TimeStampToString(time.Now().UnixNano())
}
