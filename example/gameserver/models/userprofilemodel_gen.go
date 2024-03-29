// Code generated by goctl. DO NOT EDIT.
package models

import (
	"gorm.io/gorm"
	"time"
)

const (
	UserProfileTName = "user.user_profile"
)

type (
	UserProfile struct {
		Id           int64     `gorm:"column:id"`
		Nick         string    `gorm:"column:nick"`
		Pwd          string    `gorm:"column:pwd"`
		Email        string    `gorm:"column:email"`
		CreateTime   time.Time `gorm:"column:create_time"`
		UpdateTime   time.Time `gorm:"column:update_time"`
		UserId       string    `gorm:"column:user_id"`
		VipTime      time.Time `gorm:"column:vip_time"`
		UserStatus   int64     `gorm:"column:user_status"` // 0正常 1禁用
		RegisterFrom int64     `gorm:"column:register_from"`
	}
)

func (m *UserProfile) TableName() string {
	return UserProfileTName
}
func (m *UserProfile) Create(db *gorm.DB) error {
	// m.CreateTime = time.Now()
	// m.UpdateTime = time.Now()
	return db.Table(m.TableName()).Create(m).Error
}

func (m *UserProfile) FindByPrimary(db *gorm.DB, primary int64) error {
	return IgnoreRecordNotFound(db.Table(m.TableName()).Where(" id = ?", primary).Find(m).Error)
}

func (m *UserProfile) UpdateByPrimary(db *gorm.DB, primary int64) error {
	return db.Table(m.TableName()).Where("id = ?", primary).Updates(m).Error
}

func (m *UserProfile) UpdateFieldsByPrimary(db *gorm.DB, primary int64, fields map[string]interface{}) error {
	return db.Table(m.TableName()).Where("id = ?", primary).Updates(fields).Error
}
func (m *UserProfile) DeleteByPrimary(db *gorm.DB, primary int64) error {
	return db.Table(m.TableName()).Where("id = ?", primary).Delete(m).Error
}

type UserProfileList []UserProfile

func (l *UserProfileList) FindByPrimarys(db *gorm.DB, primarys []int64) (err error) {
	if len(primarys) == 0 {
		return
	}
	err = db.Table(UserProfileTName).Where(" id in (?)", primarys).Find(l).Error
	return
}

func (l *UserProfileList) FindByPage(db *gorm.DB, page int, size int) (total int64, err error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	db = db.Table(UserProfileTName)
	//conditions
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Offset((page - 1) * size).Limit(size).Find(&l).Error
	return
}

func (l *UserProfileList) Create(db *gorm.DB, batchSize int) error {
	return db.CreateInBatches(l, batchSize).Error
}
func (m *UserProfile) FindByEmail(db *gorm.DB, key string) error {
	return IgnoreRecordNotFound(db.Table(m.TableName()).Where(" email = ?", key).Find(m).Error)
}
func (m *UserProfile) FindByUserId(db *gorm.DB, key string) error {
	return IgnoreRecordNotFound(db.Table(m.TableName()).Where(" user_id = ?", key).Find(m).Error)
}
