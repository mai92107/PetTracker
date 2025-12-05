package initial

import (
	jsonModal "batchLog/0.config"
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	"batchLog/0.core/model"
	cron "batchLog/0.cron"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	jsoniter "github.com/json-iterator/go"
)

var (
	MariaDBSetting jsonModal.MariaDbConfig
	MongoDBSetting jsonModal.MongoDbConfig
	RedisDBSetting jsonModal.RedisDbConfig

	MosquittoBrokerSetting jsonModal.MosquittoConfig
)

func InitAll() {
	InitLogger()

	initWorkers()

	loadEnvFromJSON()

	initMachine()

	InitDeviceSequence()
	cron.CronStart()
}
func InitLogger() {
	logafa.CreateLogFileNow()

	handler := logafa.NewLogafaHandler(&slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	slog.SetDefault(slog.New(handler))
	logafa.Debug("Logafa 初始化完成")

}
func initWorkers() {
	maxPriorWorkers := 20
	maxNormalWorkers := 50
	// 區隔工人 做 故障隔離
	// 高級勞工
	global.PriorWorkerPool = make(chan struct{}, maxPriorWorkers)
	for i := 0; i < maxPriorWorkers; i++ {
		global.PriorWorkerPool <- struct{}{}
	}
	logafa.Debug("👮🏻‍♀️高級勞工 聘請成功", "count", maxPriorWorkers)
	// 城市打工人
	global.NormalWorkerPool = make(chan struct{}, maxNormalWorkers)
	for i := 0; i < maxNormalWorkers; i++ {
		global.NormalWorkerPool <- struct{}{}
	}
	logafa.Debug("👷🏻城市打工人 聘請成功", "count", maxNormalWorkers)
}

// func initEnv() (env string) {
// 	flag.StringVar(&env, "env", "dev", "Environment: dev, prod, test")
// 	flag.Parse()
// 	return
// }

func loadEnvFromJSON() {
	err := loadConfigJson()
	if err != nil {
		logafa.Error(" 讀取設定 json 發生異常, error: %v", err)
		return
	}

	err = loadMachineJson()
	if err != nil {
		logafa.Error(" 讀取機器 json 發生異常, error: %v", err)
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

func loadConfigJson() error {
	fileName := "config.json"
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

func loadMachineJson() error {
	fileName := "machine.json"
	// 打開 JSON 檔案
	data, err := loadJsonFile(fileName)
	if err != nil {
		return nil
	}

	var machine jsonModal.Machine
	// 解析 JSON
	err = jsoniter.UnmarshalFromString(data, &machine)
	if err != nil {
		return fmt.Errorf("❌ 解析 JSON 失敗: %s, error: %v", fileName, err)
	}

	MariaDBSetting = machine.MariaDB
	MongoDBSetting = machine.MongoDB
	RedisDBSetting = machine.Redis
	MosquittoBrokerSetting = machine.MosquittoBroker

	return nil
}

func initMachine() {
	global.Repository = &model.Repo{
		DB: &model.DataBase{
			MariaDb: InitMariaDB(MariaDBSetting),
			MongoDb: InitMongoDB(MongoDBSetting),
		},
		Cache: InitRedis(RedisDBSetting),
	}
	global.GlobalBroker = InitMosquitto(MosquittoBrokerSetting)
}
