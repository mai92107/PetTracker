package memberService

import (
	"batchLog/0.core/global"
	repo "batchLog/4.repo"
)

func AddDevice(memberId int64, deviceId, deviceName string) error {
	readingDB := global.Repository.DB.Reading
	device, err := repo.FindDeviceByDeviceId(readingDB, deviceId)
	if err != nil {
		return err
	}
	writingDB := global.Repository.DB.Writing
	err = repo.AddDevice(writingDB, memberId, device.DeviceId, deviceName)
	if err != nil {
		return err
	}
	return nil
}
