package accountService

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	repo "batchLog/4.repo"
	"context"
	"fmt"
)

func Register(ctx context.Context, ip, username, password, email, lastName, firstName, nickName string) (map[string]interface{}, error) {

	err := validateRegister(email, nickName, username, password)
	if err != nil {
		return nil, err
	}

	tx := global.Repository.DB.MariaDb.Writing.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logafa.Error("註冊失敗(panic): %v", r)
			panic(r) // 重新拋出
		}
	}()

	memberId, err := repo.CreateMember(ctx, tx, lastName, firstName, nickName, email)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	accountUuid, err := repo.CreateAccount(ctx, tx, memberId, username, password, email)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = repo.CreatePasswordHistory(ctx, tx, accountUuid, password)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf(global.COMMON_SYSTEM_ERROR)
	}

	return Login(ctx, ip, username, password)
}

func validateRegister(email, nickName, username, password string) error {
	if username == "" {
		return fmt.Errorf("使用者帳號不可為空")
	}
	if password == "" {
		return fmt.Errorf("使用者密碼不可為空")
	}
	if email == "" {
		return fmt.Errorf("電子信箱不可為空")
	}
	if nickName == "" {
		return fmt.Errorf("使用者名稱不可為空")
	}
	return nil
}
