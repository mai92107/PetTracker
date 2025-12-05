package logafa

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logFileMu sync.Mutex
)

// 給外部呼叫（手動觸發）
func CreateLogFileNow() {
	openLogFile()
}
func StartRotateFile(ctx context.Context) {
	openLogFile()
}

// 取得當前應用的檔案名稱
func filename() string {
	now := time.Now()
	return filepath.Join("log", now.Format("2006-01-02")+".log")
}

// 開新檔 + 關舊檔
func openLogFile() {
	logFileMu.Lock()
	defer logFileMu.Unlock()

	closeOldLog()

	LogFile = getCurrentLogFile(filename)
	Debug("已開啟 log file")
}

func getCurrentLogFile(getFilename func() string) *os.File {
	filename := getFilename()

	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		Error("[LOGAFA] 建立 log 目錄失敗", "error", err)
		return nil
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		Error("[LOGAFA] 檔案開啟失敗", "error", err)
		return nil
	}
	return file
}

func closeOldLog() {
	if LogFile != nil {
		LogFile.Close()
		LogFile = nil
	}
}
