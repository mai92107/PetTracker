package cron

import (
	"batchLog/0.core/logafa"
	tripService "batchLog/3.service/trip"
	"context"

	"github.com/robfig/cron/v3"
)

const EXECUTOR string = "SYSTEM_CRON"

func CronStart() {
	c := cron.New(cron.WithSeconds())

	// 每秒鐘執行一次
	executeJob(c, Second, []func(context.Context){})

	// 每分鐘執行一次
	executeJob(c, Minute, []func(context.Context){
		// func(){
		// 	for i := 0; i <= 10; i++{
		// 		println(i)
		// 		time.Sleep(1000 * time.Millisecond)
		// 	}
		// },
	})
	// 每5分鐘執行一次
	executeJob(c, Five, []func(context.Context){})

	// 每10分鐘執行一次
	executeJob(c, Ten, []func(context.Context){})

	// 每15分鐘執行一次
	executeJob(c, Quarter, []func(context.Context){})

	// 每半小時執行一次
	executeJob(c, HalfHour, []func(context.Context){})

	// 每小時執行一次
	executeJob(c, Hour, []func(context.Context){
		// data.GetOnlineDevice,
		func(ctx context.Context) {
			tripService.FlushGpsFmRdsToMongo(ctx, nil, 70)
		},
	})

	// 每半天執行一次（每日00:00, 12:00）
	executeJob(c, HalfDay, []func(context.Context){})

	// 每天執行一次（每日00:00）
	executeJob(c, Day, []func(context.Context){
		func(ctx context.Context) {
			tripService.FlushTripFmMongoToMaria(ctx, 25, EXECUTOR)
		},
		logafa.StartRotateFile,
	})

	c.Start()
}
