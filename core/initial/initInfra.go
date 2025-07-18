package initial

import (
	"batchLog/core/global"
	"batchLog/core/logafa"
	"batchLog/core/redis"
)

func InitDeviceSequence(){

	device_setting := redis.HGetAllData("device_setting")
	if len(device_setting) == 0 {
		redis.HSetData("device_setting",
			map[string]interface{}{
				"device_prefix":global.ConfigSetting.DevicePrefix,
				"device_sequence":global.ConfigSetting.DeviceSequence,
			})
	}
	logafa.Debug(" ✅ 成功初始化裝置設定")
}