package initial

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	"batchLog/0.core/redis"
	"context"
)

func InitDeviceSequence() {
	ctx := context.Background()
	device_setting := redis.HGetAllData(ctx, "device_setting")
	if len(device_setting) == 0 {
		redis.HSetData(ctx, "device_setting",
			map[string]interface{}{
				"device_prefix":   global.ConfigSetting.DevicePrefix,
				"device_sequence": global.ConfigSetting.DeviceSequence,
			})
	}
	logafa.Debug(" ✅ 成功初始化裝置設定")
}
