package cron

import (
	"batchLog/0.core/global"
	logafa "batchLog/0.core/logafa"
	"batchLog/0.cron/data"
	log "batchLog/0.cron/logafa"
	"batchLog/0.cron/persist"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

const TIME_LAYOUT = "15:04:05"

var wg sync.WaitGroup

type CronType int
type cronInfo struct {
	spec     string
	infoName string
}

const (
	Second CronType = iota
	Minute
	Quarter  // 每 15 分
	HalfHour // 每 30 分
	Hour
	HalfDay // 00:00, 12:00
	Day     // 00:00
)

var cronSpecs = map[CronType]cronInfo{
	Second:   {spec: "*/1 * * * * *", infoName: "每秒"},
	Minute:   {spec: "0 * * * * *", infoName: "每分鐘"},
	Quarter:  {spec: "0 */15 * * * *", infoName: "每 15 分鐘"},
	HalfHour: {spec: "0 */30 * * * *", infoName: "每 30 分鐘"},
	Hour:     {spec: "0 0 * * * *", infoName: "每小時"},
	HalfDay:  {spec: "0 0 0,12 * * *", infoName: "每半天"},
	Day:      {spec: "0 0 0 * * *", infoName: "每天"},
}

func CronStart() {
	c := cron.New(cron.WithSeconds())

	// 每秒鐘執行一次
	executeJob(c, Second, []func(){})

	// 每分鐘執行一次
	executeJob(c, Minute, []func(){
		// func(){
		// 	for i := 0; i <= 10; i++{
		// 		println(i)
		// 		time.Sleep(1000 * time.Millisecond)
		// 	}
		// },
	})

	// 每15分鐘執行一次
	executeJob(c, Quarter, []func(){
		persist.SaveGpsFmRedisToMaria,
	})

	// 每半小時執行一次
	executeJob(c, HalfHour, []func(){})

	// 每小時執行一次
	executeJob(c, Hour, []func(){
		persist.GetLastDayTripsWithDistance,
		data.GetOnlineDevice,
	})

	// 每半天執行一次（每日00:00, 12:00）
	executeJob(c, HalfDay, []func(){})

	// 每天執行一次（每日00:00）
	executeJob(c, Day, []func(){
		log.StartRotateFile,
	})

	c.Start()
}

func executeJob(c *cron.Cron, cronType CronType, jobs []func()) {
	// 沒工作就離開
	if len(jobs) == 0 {
		return
	}

	c.AddFunc(cronSpecs[cronType].spec, func() {
		start := time.Now()
		logafa.Debug("%s執行程序, 現在時間: %+v", cronSpecs[cronType].infoName, start.Format(TIME_LAYOUT))

		var localWg sync.WaitGroup
		for _, job := range jobs {
			submitJobAsync(job, &localWg)
		}
		localWg.Wait()
		duration := time.Since(start)
		logafa.Debug("%s任務執行完畢，耗時: %v", cronSpecs[cronType].infoName, duration)
	})
}

// 工人分配執行工作
func submitJobAsync(job func(), localWg *sync.WaitGroup) {
	wg.Add(1)
	<-global.NormalWorkerPool // 取得 worker
	localWg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			localWg.Done()
			global.NormalWorkerPool <- struct{}{}
		}()
		start := time.Now()
		job()
		logafa.Debug("單一任務完成，耗時: %v", time.Since(start))
	}()
}

// 優雅結束檢查未完成工作
func CheckIsCronJobsFinished() {
	// 等待所有背景任務完成
	done := make(chan struct{})
	go func() {
		wg.Wait() // 所有 Add(1) 的都要 Done()
		close(done)
	}()

	select {
	case <-done:
		logafa.Info("所有背景任務已完成，安全關閉")
	case <-time.After(30 * time.Second):
		logafa.Warn("關閉超時，強制退出")
	}
}
