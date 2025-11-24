package common

import (
	"batchLog/0.core/global"
	"time"
)

func ToUtcTime(time time.Time)time.Time{
	return time.UTC()
}

func ToLocalTime(time time.Time)time.Time{
	return time.Local()
}

func ToUtcTimeStr(time time.Time)string{
	return time.UTC().Format(global.TIME_FORMAT)
}

func ToLocalTimeStr(time time.Time)string{
	return time.Local().Format(global.TIME_FORMAT)
}