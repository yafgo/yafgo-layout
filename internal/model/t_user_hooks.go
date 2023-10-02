package model

import (
	"yafgo/yafgo-layout/pkg/hash"

	"gorm.io/gorm"
)

// BeforeSave GORM 的模型钩子，在创建和更新模型前调用
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {

	if u.Password != "" && !hash.BcryptIsHashed(u.Password) {
		u.Password, err = hash.BcryptHash(u.Password)
	}

	return
}
