package initial

import (
	jsonModal "batchLog/0.config"
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMariaDB(setting jsonModal.MariaDbConfig) *global.SqlDB {
	if !setting.InUse {
		return nil
	}
	// 暫時共用
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		setting.Reading.User, setting.Reading.Password,
		setting.Reading.Host, setting.Reading.Port,
		setting.Reading.Name)

	readingDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logafa.Error(" ❌ 無法連接讀取資料庫: %v", err)
	}

	writingDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logafa.Error(" ❌ 無法連接寫入資料庫: %v", err)
	}
	logafa.Debug(" ✅ MariaDB 資料庫連接成功")
	return &global.SqlDB{
		Reading: readingDb,
		Writing: writingDb,
	}
}

func InitMongoDB(setting jsonModal.MongoDbConfig) *global.NoSqlDB {
	if !setting.InUse {
		return nil
	}
	// 暫時共用
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", setting.Reading.User, setting.Reading.Password, setting.Reading.Host, setting.Reading.Port)

	clientOptions := options.Client().ApplyURI(uri)

	// 設置執行timeout時間
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(setting.Reading.TimeoutRange)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logafa.Error("無法連接 Mongodb, error: %+v", err)
		panic(err)
	}
	logafa.Debug("✅ 成功連線 MongoDB!")
		
	// 初始化index
	initMongoIndexes(client)

	return &global.NoSqlDB{
		Reading: client,
		Writing: client,
	}
}
func initMongoIndexes(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer func() {
		if r := recover(); r != nil {
			logafa.Error("初始化Mongo Index 失敗 (panic): %v", r)
			panic(r) // 重新拋出
		}
	}()
	collection := client.Database("pettrack").Collection("pettrack")

	// 1. 地理空間索引 (2dsphere) - 用於地理查詢
	geoIndexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "location", Value: "2dsphere"}},
		Options: options.Index().
			SetName("idx_location_2dsphere"),
	}

	// 2. 複合唯一索引 - 防止重複資料
	uniqueIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "device_id", Value: 1},
			{Key: "recorded_at", Value: 1},
		},
		Options: options.Index().
			SetName("idx_device_recorded_unique").
			SetUnique(true),
	}

	// 3. 時間索引 - 用於時間範圍查詢和資料清理
	timeIndexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "recorded_at", Value: -1}},
		Options: options.Index().
			SetName("idx_recorded_at_desc"),
	}

	// 4. 裝置索引 - 用於查詢特定裝置的軌跡
	deviceIndexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "device_id", Value: 1}},
		Options: options.Index().
			SetName("idx_device_id"),
	}

	// 5. TTL 索引 - 自動刪除舊資料 (可選)
	ttlIndexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "created_at", Value: 1}},
		Options: options.Index().
			SetName("idx_created_at_ttl").
			SetExpireAfterSeconds(90 * 24 * 60 * 60), // 90 天後自動刪除
	}

	// 批次建立索引
	indexNames, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		geoIndexModel,
		uniqueIndexModel,
		timeIndexModel,
		deviceIndexModel,
		ttlIndexModel, // 如果不需要自動刪除舊資料,可以移除這個
	})

	if err != nil {
		logafa.Error("建立索引失敗: %+v", err)
		return
	}

	logafa.Info("MongoDB 索引建立成功: %v", indexNames)
}
