package accountService

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	repo "batchLog/4.repo"
	"fmt"
)

func Register(ip, username, password, email, lastName, firstName, nickName string)(map[string]interface{},error){

	tx := global.Repository.DB.Writing.Begin()
	defer func(){
		if r := recover();r != nil{
			logafa.Error("註冊失敗")
		}
		}()

	memberId,err := repo.CreateMember(tx,lastName,firstName,nickName,email)
	if err != nil{
		return nil, err
	}

	accountUuid,err := repo.CreateAccount(tx,memberId,username,password,email)
	if err != nil{
		return nil, err
	}

	err = repo.CreatePasswordHistory(tx,accountUuid,password)
	if err != nil{
		return nil, err
	}

	if err = tx.Commit().Error;err != nil{
		tx.Rollback()
		return nil,fmt.Errorf(global.COMMON_SYSTEM_ERROR)
	}

	return Login(ip,username,password)
}