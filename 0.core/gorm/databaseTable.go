package gormTable

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
    Uuid            uuid.UUID  `gorm:"type:char(36);primaryKey" json:"uuid"`
    MemberId        int64      `gorm:"not null" json:"memberId"`
    Username        string     `gorm:"type:varchar(255);unique;not null" json:"username"`
    Password        string     `gorm:"type:varchar(255);not null" json:"password"`
    Email           string     `gorm:"type:varchar(255);unique;not null" json:"email"`
    Identity        string     `gorm:"type:varchar(50)" json:"identity"`
    LastLoginTime   time.Time  `gorm:"type:datetime" json:"lastLoginTime"`
    CreatedAt       time.Time  `gorm:"autoCreateTime" json:"createdAt"`
}
func (a *Account)TableName()string{
    return "account"
}

type PasswordHistory struct {
    Id              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
    AccountUuid     uuid.UUID `gorm:"type:char(36);not null" json:"accountUuid"`
    Password        string    `gorm:"type:varchar(255);not null" json:"password"`
    CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
}
func (pp *PasswordHistory)TableName()string{
    return "password_history"
}

type Member struct {
    Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
    LastName  string    `gorm:"type:varchar(255)" json:"lastName"`
    FirstName string    `gorm:"type:varchar(255)" json:"firstName"`
    NickName  string    `gorm:"type:varchar(255)" json:"nickName"`
    Email     string    `gorm:"type:varchar(255)" json:"email"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}
func (m *Member)TableName()string{
    return "member"
}

type MemberDevice struct {
    Id              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
    MemberId        int64     `gorm:"not null" json:"memberId"`
    DeviceId        string    `gorm:"type:char(36);not null" json:"deviceId"`
    DeviceName      string    `gorm:"type:varchar(255)" json:"deviceName"`
    CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
}
func (md *MemberDevice)TableName()string{
    return "member_device"
}

type Device struct {
    Uuid            uuid.UUID   `gorm:"type:char(36);primaryKey" json:"uuid"`
    DeviceId        string      `gorm:"type:varchar(255)" json:"deviceId"`
    DeviceType      string      `gorm:"type:varchar(50)" json:"deviceType"`
    CreateByMember  int64       `gorm:"not null" json:"memberId"`
    Remark          string      `gorm:"type:char(50)" json:"remark"`
}
func (d *Device)TableName()string{
    return "device"
}

// CREATE TABLE member (
//     id BIGINT AUTO_INCREMENT PRIMARY KEY,
//     last_name VARCHAR(255),
//     first_name VARCHAR(255),
//     nick_name VARCHAR(255),
//     email VARCHAR(255),
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );

// CREATE TABLE device (
//     uuid CHAR(36) PRIMARY KEY,
//     device_id VARCHAR(36) UNIQUE,
//     device_type VARCHAR(50),
//     create_by_member BIGINT NOT NULL,
//     remark CHAR(50),
//     CONSTRAINT fk_device_create_by_member FOREIGN KEY (create_by_member) REFERENCES member(id)
// );


// CREATE TABLE member_device (
//     id BIGINT AUTO_INCREMENT PRIMARY KEY,
//     member_id BIGINT NOT NULL,
//     device_id VARCHAR(36) NOT NULL,
//     device_name VARCHAR(255),
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     CONSTRAINT fk_member_device_member FOREIGN KEY (member_id) REFERENCES member(id) ON DELETE CASCADE,
//     CONSTRAINT fk_member_device_device FOREIGN KEY (device_id) REFERENCES device(device_id) ON DELETE CASCADE,
//     CONSTRAINT uq_member_device UNIQUE (member_id, device_id)
// );

// CREATE TABLE account (
//     uuid CHAR(36) PRIMARY KEY,
//     member_id BIGINT NOT NULL,
//     username VARCHAR(255) NOT NULL UNIQUE,
//     password VARCHAR(255) NOT NULL,
//     email VARCHAR(255) NOT NULL UNIQUE,
//     identity VARCHAR(50),
//     last_login_time DATETIME,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     CONSTRAINT fk_account_member FOREIGN KEY (member_id) REFERENCES member(id) ON DELETE CASCADE
// );

// CREATE TABLE password_history (
//     id BIGINT AUTO_INCREMENT PRIMARY KEY,
//     account_uuid CHAR(36) NOT NULL,
//     password VARCHAR(255) NOT NULL,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     CONSTRAINT fk_password_history_account FOREIGN KEY (account_uuid) REFERENCES account(uuid) ON DELETE CASCADE
// );
