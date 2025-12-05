package data

import (
	repo "batchLog/4.repo"
	"context"
)

func GetOnlineDevice(ctx context.Context) {
	repo.GetOnlineDevices(ctx)
}
