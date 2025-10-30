package cron

import (
	"batchLog/0.core/logafa"
	"batchLog/0.cron/persist"
	"time"

	"github.com/robfig/cron/v3"
)

const TIME_LAYOUT = "15:04:05"

func CronStart() {
	c := cron.New(cron.WithSeconds())

	// // 每秒鐘執行一次
	// c.AddFunc("* * * * * *", func() {
	// 	logafa.Debug("每秒執行程序, 現在時間: %+v", time.Now().Format(TIME_LAYOUT))
	// 	logafa.Debug("%v秒執行完畢", time.Now().Second())
	// })

	// // 每分鐘執行一次
	// c.AddFunc("0 * * * * *", func() {
	// 	logafa.Debug("每分執行程序, 現在時間: %+v", time.Now().Format(TIME_LAYOUT))
	// 	logafa.Debug("%v分執行完畢", time.Now().Minute())
	// })

	// 每15分鐘執行一次
	c.AddFunc("0 */15 * * * *", func() {
		logafa.Debug("每15分鍾執行程序, 現在時間: %+v", time.Now().Format(TIME_LAYOUT))
		persist.SaveGpsFmRedisToMaria()
		logafa.Debug("%v分執行完畢", time.Now().Minute())
	})

	// // 每半小時執行一次
	// c.AddFunc("0 0,30 * * * *", func() {
	// 	logafa.Debug("每半小時執行程序, 現在時間: %+v", time.Now().Format(TIME_LAYOUT))
	// 	logafa.Debug("%v點%v執行完畢", time.Now().Hour(), time.Now().Minute())
	// })

	// // 每小時執行一次（整點）
	// c.AddFunc("0 0 * * * *", func() {
	// 	logafa.Debug("每整點執行程序, 現在時間: %+v", time.Now().Format(TIME_LAYOUT))
	// 	logafa.Debug("%v點執行完畢", time.Now().Hour())
	// })

	// 每半天執行一次（每日00:00, 12:00）
	c.AddFunc("0 0 0,12 * * *", func() {
		logafa.Debug("每半天執行程序, 現在時間: %+v", time.Now().Format(TIME_LAYOUT))
		logafa.RotateLogFile()
		logafa.Debug("每日執行完畢")
	})

	// // 每天執行一次（每日00:00）
	// c.AddFunc("0 0 0 * * *", func() {
	// 	logafa.Debug("每天 00:00 執行程序, 現在時間: %+v", time.Now().Format(TIME_LAYOUT))
	// 	logafa.Debug("每日執行完畢")
	// })

	c.Start()
}
