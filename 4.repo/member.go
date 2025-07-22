package repo

import (
	gormTable "batchLog/0.core/gorm"
	"batchLog/0.core/logafa"
	"fmt"

	"gorm.io/gorm"

	"github.com/google/uuid"
)
func FindMemberByUuid(tx *gorm.DB, uuid string) (*gormTable.MemberInfo, error) {
    var memberInfo gormTable.MemberInfo
    err := tx.First(&memberInfo, "uuid = ?", uuid).Error
	if err != nil{
		logafa.Error("無法查詢使用者資料, error: %+v",err)
		return nil, fmt.Errorf("無法查詢使用者資料")
	}
    return &memberInfo, nil
}

func FindMemberByAccountUuid(tx *gorm.DB, accountUuid string) (*gormTable.MemberInfo, error){
	var memberInfo gormTable.MemberInfo
    err := tx.First(&memberInfo, "account_uuid = ?", accountUuid).Error
	if err != nil{
		logafa.Error("無法查詢使用者資料, error: %+v",err)
		return nil, fmt.Errorf("無法查詢使用者資料")
	}
    return &memberInfo, nil
}

func CreateMember(tx *gorm.DB, accountUuid uuid.UUID, lastName, firstName, nickName, email string)(uuid.UUID,error){
	member := gormTable.MemberInfo{
		Uuid: uuid.New(),
		AccountUuid: accountUuid,
		LastName: lastName,
		FirstName: firstName,
		NickName: nickName,
		Email: email,
	}
	err := tx.Table("member_info").Create(&member).Error
	if err != nil {
		logafa.Error("建立使用者資料失敗, error: %+v",err)
		return uuid.Nil, fmt.Errorf("建立使用者失敗")
	}
	return member.Uuid,nil
}