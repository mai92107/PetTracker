package initial

import (
	jsonModal "batchLog/0.config"
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	cron "batchLog/0.cron"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	jsoniter "github.com/json-iterator/go"
)
var (
	MariaDBSetting 	jsonModal.MariaDbConfig
	RedisDBSetting 	jsonModal.RedisDbConfig

	MosquittoBrokerSetting	jsonModal.MosquittoConfig
)

func InitAll(){
	env := initEnv()

	loadEnvFromJSON(env)
	logafaInit(env)
	
	initMachine()

	InitDeviceSequence()
	cron.CronStart()
}

func initEnv()(env string){
	flag.StringVar(&env, "env", "dev", "Environment: dev, prod, test")
	flag.Parse()
	return
}

func loadEnvFromJSON(env string){
	err := loadConfigJson(env)
	if err != nil{
		logafa.Error(" 讀取設定 json 發生異常, error: %v",err)
		return
	}

	err = loadMachineJson(env)
	if err != nil{
		logafa.Error(" 讀取機器 json 發生異常, error: %v",err)
		return
	}
}


func loadJsonFile(fileName string) (string, error) {
	wd, _ := os.Getwd()
	configFile := "0.config"
	filePath := filepath.Join(wd, configFile, fileName)
	// 讀取檔案內容為 []byte
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf(" ❌ 無法開啟 JSON 檔案: %s, error: %v", filePath, err)
	}
	return string(content), nil
}

func loadConfigJson(env string)error{
	fileName := fmt.Sprintf("config_%s.json",env)
	// 打開 JSON 檔案
	data, err := loadJsonFile(fileName)
	if err != nil {
		return nil
	}

	var config jsonModal.Config
	// 解析 JSON
	err = jsoniter.UnmarshalFromString(data, &config)
	if err != nil {
		return err
	}
	global.ConfigSetting = config
	return nil
}

func logafaInit(env string){

	switch env {
	case "dev":
		logafa.CurrentLevel = logafa.DEBUG
	case "prod":
		logafa.CurrentLevel = logafa.INFO
	case "test":
		logafa.CurrentLevel = logafa.WARN
	default:
		logafa.CurrentLevel = logafa.DEBUG
	}

	now := time.Now()

	var err error
	wd, _ := os.Getwd()
	filePath := filepath.Join(wd, "log", now.Format("2006-01-02") + ".log")
	logafa.LogFile, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("無法打開 log 檔案: %v", err))
	}
}

func loadMachineJson(env string) error{
	fileName := fmt.Sprintf("machine_%s.json",env)
	// 打開 JSON 檔案
	data, err := loadJsonFile(fileName)
	if err != nil {
		return nil
	}

	var machine jsonModal.Machine
	// 解析 JSON
	err = jsoniter.UnmarshalFromString(data, &machine)
	if err != nil {
		return fmt.Errorf("❌ 解析 JSON 失敗: %s, error: %v",fileName, err)
	}

	MariaDBSetting = machine.MariaDB
	RedisDBSetting = machine.Redis
	MosquittoBrokerSetting = machine.MosquittoBroker

	return nil
}                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         

func initMachine() {
	global.Repository = &global.Repo{
		DB: InitMariaDB(MariaDBSetting),
		Cache: InitRedis(RedisDBSetting),
	}
	global.GlobalBroker = InitMosquitto(MosquittoBrokerSetting)
}