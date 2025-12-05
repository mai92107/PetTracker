package deviceService

import (
	"batchLog/0.core/global"
	repo "batchLog/4.repo"
	"context"
	"fmt"
	"time"
)

func OnlineDeviceList(ctx context.Context) ([]string, error) {
	const timeout = 2 * time.Second
	start := time.Now()
	for {
		if global.ActiveDevicesLock.TryLock() {
			break
		}
		if time.Since(start) > timeout {
			// 鎖超時
			return nil, fmt.Errorf("警告：OnlineDeviceList() 嘗試加鎖超過 2 秒，放棄。")
		}
		time.Sleep(10 * time.Millisecond)
	}
	defer global.ActiveDevicesLock.Unlock()

	deviceIds, err := repo.GetOnlineDevices(ctx)

	return deviceIds, err
}
