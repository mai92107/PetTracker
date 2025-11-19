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
	openLogFile()
}
func StartRotateFile(){
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

	logafa.LogFile = getCurrentLogFile(filename)
	logafa.Debug("已開啟 log file")
}

func getCurrentLogFile(getFilename func()string)*os.File{
	filename := getFilename()
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logafa.Error("[LOGAFA] 檔案開啟失敗")
		return nil
	}
	return file
}

func closeOldLog(){
	if logafa.LogFile != nil {
		logafa.LogFile.Close()
		logafa.LogFile = nil
	}
}
