package cron

import (
	"github.com/robfig/cron/v3"
)

const TIME_LAYOUT = "15:04:05"

func CronStart(){
	c := cron.New(cron.WithSeconds())

	// 每秒鐘執行一次
	// c.AddFunc("* * * * * *", func() {
	// 	logafa.Debug("每秒執行程序, 現在時間: %+v", time.Now().Format(TIME_LAYOUT))
	// })

	// // 每分鐘執行一次（整分）
	// c.AddFunc("0 * * * * *", func() {
	// 	logafa.Info("每分執行程序, 現在時間: %+v", time.Now().Format(TIME_LAYOUT))
	// 	logafa.Info("分鐘執行完畢")
	// })

	// 每半小時執行一次
	// c.AddFunc("0 0,30 * * * *", func() {
	// 	logafa.Info("每半小時執行程序, 現在時間: %+v", time.Now().Format(TIME_LAYOUT))
	// 	persist.SaveGpsFmRedisToMaria()
	// 	logafa.Info("半小時執行完畢")
	// })

	// // 每小時執行一次（整點）
	// c.AddFunc("0 0 * * * *", func() {
	// 	logafa.Info("每小時執行程序, 現在時間: %+v", time.Now().Format(TIME_LAYOUT))		
	// 	logafa.Info("小時執行完畢")
	// })

	// // 每天執行一次（每日00:00）
	// c.AddFunc("0 0 0 * * *", func() {
	// 	fmt.Println("每天 00:00 執行", time.Now())
	// })

	c.Start()
}