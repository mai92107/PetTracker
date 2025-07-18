package repo

import (
	"batchLog/core/global"
	gormTable "batchLog/core/gorm"
	"batchLog/core/logafa"
	"fmt"

	"github.com/google/uuid"
)

// ✅ 定義 Repository 介面
type MemberInfoRepository interface {
    FindByUuid(uuid string) (*gormTable.MemberInfo, error)
	CreateMember(accountUuid uuid.UUID, lastName, firstName, nickName, email string)(uuid.UUID,error)
}

// ✅ 實作 Repository
type MemberInfoRepositoryImpl struct {
	DB    *global.DB
	Cache *global.Cache
}

func NewMemberInfoRepository(db *global.DB, cache *global.Cache) *MemberInfoRepositoryImpl {
	return &MemberInfoRepositoryImpl{
		DB: 		db,
		Cache: 		cache,
	}
}

func (r *MemberInfoRepositoryImpl) FindByUuid(uuid string) (*gormTable.MemberInfo, error) {
    var memberInfo gormTable.MemberInfo
    err := r.DB.Reading.First(&memberInfo, "uuid = ?", uuid).Error
    return &memberInfo, err
}

func (r *MemberInfoRepositoryImpl)CreateMember(accountUuid uuid.UUID, lastName, firstName, nickName, email string)(uuid.UUID,error){
	member := gormTable.MemberInfo{
		Uuid: uuid.New(),
		AccountUuid: accountUuid,
		LastName: lastName,
		FirstName: firstName,
		NickName: nickName,
		Email: email,
	}
	err := r.DB.Writing.Table("member_info").Create(&member).Error
	if err != nil {
		logafa.Error("建立使用者資料失敗, error: %+v",err)
		return uuid.Nil, fmt.Errorf("建立使用者失敗")
	}
	return member.Uuid,nil
}
