package repo

import (
	gormTable "batchLog/0.core/gorm"
	"batchLog/0.core/logafa"
	"fmt"

	"gorm.io/gorm"
)
func FindMemberById(tx *gorm.DB, id int64) (*gormTable.Member, error) {
    var member gormTable.Member
    err := tx.First(&member, "id = ?", id).Error
	if err != nil{
		logafa.Error("無法查詢使用者資料, error: %+v",err)
		return nil, fmt.Errorf("無法查詢使用者資料")
	}
    return &member, nil
}

func FindMemberByAccountUuid(tx *gorm.DB, accountUuid string) (*gormTable.Member, error){
	var memberInfo gormTable.Member
    err := tx.First(&memberInfo, "account_uuid = ?", accountUuid).Error
	if err != nil{
		logafa.Error("無法查詢使用者資料, error: %+v",err)
		return nil, fmt.Errorf("無法查詢使用者資料")
	}
    return &memberInfo, nil
}

func CreateMember(tx *gorm.DB, lastName, firstName, nickName, email string)(int64,error){
	member := gormTable.Member{
		LastName: lastName,
		FirstName: firstName,
		NickName: nickName,
		Email: email,
	}
	err := tx.Create(&member).Error
	if err != nil {
		logafa.Error("建立使用者資料失敗, error: %+v",err)
		return 0, fmt.Errorf("建立使用者失敗")
	}
	return member.Id,nil
}

func AddDevice(tx *gorm.DB, memberId int64, deviceId, deviceName string)error{

	memberDevice := gormTable.MemberDevice{
		MemberId: memberId,
		DeviceId: deviceId,
		DeviceName: deviceName,
	}
	err := tx.Create(&memberDevice).Error
	if err != nil {
		logafa.Error("新增使用者裝置失敗, error: %+v",err)
		return fmt.Errorf("新增使用者裝置失敗")
	}
	return nil
}

func FindMemberByDeviceId(tx *gorm.DB, deviceId string)(gormTable.Member,error){
	member := gormTable.Member{}
    err := tx.Table("member_device AS md").
        Select("m.*").
        Joins("JOIN member AS m ON md.member_id = m.id").
        Where("md.device_id = ?", deviceId).
        First(&member).Error
	if err != nil {
		logafa.Error("查詢使用者失敗, error: %+v",err)
		return member,fmt.Errorf("查詢使用者失敗")
	}
	return member,nil
}