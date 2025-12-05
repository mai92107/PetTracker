package heartbeat

import (
	"batchLog/0.core/global"
	"batchLog/0.core/redis"
	"context"
	"fmt"
	"time"
)

func UpdateHeartBeat(nickname string, deviceID string) error {
	key := fmt.Sprintf("user:%s:%s", nickname, deviceID)
	heartbeat := fmt.Sprintf("%v", time.Now().UTC().UnixMilli())
	ctx := context.Background()
	err := redis.HSetFieldData(ctx, key, "heartbeat", heartbeat)
	global.Repository.Cache.Writing.Expire(global.Repository.Cache.CTX, key, 60*time.Second)
	return err
}
