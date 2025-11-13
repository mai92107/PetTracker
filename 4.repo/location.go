package repo

import (
	"context"

	"batchLog/0.core/global"
	gormTable "batchLog/0.core/gorm"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetLatestDeviceRecordByDeviceId(deviceId string) (*gormTable.DeviceLocation, error) {
	mongoDb := global.Repository.DB.MongoDb.Reading
	collection := mongoDb.Collection("pettrack")

	var result gormTable.DeviceLocation
	err := collection.FindOne(
		context.TODO(),
		bson.M{"device_id": deviceId},
		options.FindOne().SetSort(bson.D{{"recorded_at", -1}}),
	).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
