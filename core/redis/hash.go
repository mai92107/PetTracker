package redis

import (
	"batchLog/core/global"
	"batchLog/core/logafa"
)

func HSetData(key string, mapData map[string]interface{})error{
	err := global.Repository.Cache.Writing.HSet(global.Repository.Cache.CTX,key,mapData).Err()
	if err != nil{
		logafa.Error("Redis HSet 寫入失敗, key: %s, data: %+v",key,mapData)
	}
	return err
}

func HGetData(key, field string)string{
	value,err := global.Repository.Cache.Reading.HGet(global.Repository.Cache.CTX,key,field).Result()
	if err != nil{
		logafa.Error("Redis HGet 讀取失敗, key: %s, field: %s, error: %+v",key,field,err)
	}
	return value
}

func HGetAllData(key string)map[string]string{
	value,err := global.Repository.Cache.Reading.HGetAll(global.Repository.Cache.CTX,key).Result()
	if err != nil{
		logafa.Error("Redis HGetAll 讀取失敗, key: %s, error: %+v",key,err)
	}
	return value
}