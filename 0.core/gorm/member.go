package gormTable

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
    Uuid           uuid.UUID `gorm:"type:char(36);primaryKey" json:"uuid"`
    Username       string    `gorm:"type:varchar(255);unique;not null" json:"username"`
    Password       string    `gorm:"type:varchar(255);not null" json:"password"`
    Email          string    `gorm:"type:varchar(255);unique;not null" json:"email"`
    LastLoginTime  *time.Time `gorm:"type:datetime" json:"lastLoginTime"`
    PastPasswords  []PastPassword `gorm:"foreignKey:AccountUuid" json:"pastPasswords"`
    CreatedAt    time.Time `gorm:"type:timestamp;default:current_timestamp" json:"createdAt"`
}
func (a *Account)TableName()string{
    return "account"
}
type PastPassword struct {
    Id int64 `gorm:"type:bigint;autoIncrement;primaryKey" json:"id"`
    AccountUuid  uuid.UUID `gorm:"type:char(36);not null" json:"accountUuid"`
    Password     string    `gorm:"type:varchar(255);not null" json:"password"`
    CreatedAt    time.Time `gorm:"type:timestamp;default:current_timestamp" json:"createdAt"`
}
func (pp *PastPassword)TableName()string{
    return "past_password"
}

type MemberInfo struct {
    Uuid        uuid.UUID `gorm:"type:char(36);primaryKey" json:"uuid"`
    AccountUuid uuid.UUID `gorm:"type:char(36);not null" json:"accountUuid"`
    LastName    string    `gorm:"type:varchar(255)" json:"lastName"`
    FirstName   string    `gorm:"type:varchar(255)" json:"firstName"`
    NickName    string    `gorm:"type:varchar(255)" json:"nickName"`
    Email       string    `gorm:"type:varchar(255)" json:"email"`
    CreatedAt    time.Time `gorm:"type:timestamp;default:current_timestamp" json:"createdAt"`
}
func (mi *MemberInfo)TableName()string{
    return "member_info"
}

type Device struct {
    Uuid            uuid.UUID `gorm:"type:char(36);primaryKey" json:"uuid"`
    MemberInfoUuid  uuid.UUID `gorm:"type:char(36);not null" json:"memberInfoUuid"`
    DeviceId        string    `gorm:"type:varchar(255)" json:"deviceId"`
    DeviceName      string    `gorm:"type:varchar(255)" json:"deviceName"`
    CreatedAt    time.Time `gorm:"type:timestamp;default:current_timestamp" json:"createdAt"`
}
func (d *Device)TableName()string{
    return "device"
}