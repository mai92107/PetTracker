package logafa

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"batchLog/0.core/logafa"
)

var (
	logFileMu sync.Mutex
)

// 給外部呼叫（手動觸發）
func CreateLogFileNow() {
	rotateLogFile()
}
func StartRotateFile(){
	rotateLogFile()
}

// 取得當前應用的檔案名稱
func getLogFilename() string {
	now := time.Now()
	return filepath.Join("log", now.Format("2006-01-02_1504")+".log")
}

// 開新檔 + 關舊檔
func rotateLogFile() {
	logafa.Debug("%+v 分新增 log file", time.Now().Minute())
	logFileMu.Lock()
	defer logFileMu.Unlock()

	if logafa.LogFile != nil {
		logafa.LogFile.Close()
		logafa.LogFile = nil
	}

	filename := getLogFilename()

	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		logafa.Error("[LOGAFA] 建立檔案失敗")
		return
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logafa.Error("[LOGAFA] 檔案開啟失敗")
		return
	}

	logafa.LogFile = file
	logafa.Debug("%+v 分開始使用新 log file", time.Now().Minute())
	logafa.Info("[LOGAFA] Rotated to: %s", filename)
}
