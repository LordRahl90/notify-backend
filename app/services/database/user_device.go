package database

import "github.com/jinzhu/gorm"

//UserDevice struct for managing user device and token
type UserDevice struct {
	gorm.Model
	User       User `json:"user_id" gorm:"foreignkey:UserRefer"`
	UserRefer  uint
	DeviceName string `json:"device_name"`
	Token      string `json:"token"`
}
