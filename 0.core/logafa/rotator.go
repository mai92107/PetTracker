package logafa

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	logFileMu sync.Mutex
	cronJob   *cron.Cron
)

// 給外部呼叫（手動觸發）
func CreateLogFileNow() error {
	return RotateLogFile()
}

// 取得當前應用的檔案名稱
func getLogFilename() string {
	now := time.Now()
	return filepath.Join("log", now.Format("2006-01-02_1504")+".log")
}

// 開新檔 + 關舊檔
func RotateLogFile() error {
	Info("%+v分開始切換檔案", time.Now().Minute())
	logFileMu.Lock()
	defer logFileMu.Unlock()

	if LogFile != nil {
		LogFile.Close()
		LogFile = nil
	}

	filename := getLogFilename()

	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	LogFile = file
	Info("%+v分切換檔案完成", time.Now().Minute())
	fmt.Printf("[LOGAFA] Rotated to: %s\n", filename)
	return nil
}

// 程式結束時關閉
func ShutdownRotator() {
	if cronJob != nil {
		cronJob.Stop()
	}
	logFileMu.Lock()
	if LogFile != nil {
		LogFile.Close()
		LogFile = nil
	}
	logFileMu.Unlock()
}
